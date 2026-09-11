package tools

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSessionManager_AddGet(t *testing.T) {
	sm := NewSessionManager()
	session := &ProcessSession{
		ID:        "test-1",
		Command:   "echo hello",
		Status:    "running",
		StartTime: 1000,
	}

	sm.Add(session)

	got, err := sm.Get("test-1")
	require.NoError(t, err)
	require.Equal(t, "test-1", got.ID)
}

func TestSessionManager_Remove(t *testing.T) {
	sm := NewSessionManager()
	session := &ProcessSession{
		ID:        "test-1",
		Command:   "echo hello",
		Status:    "running",
		StartTime: 1000,
	}
	sm.Add(session)
	sm.Remove("test-1")

	_, err := sm.Get("test-1")
	require.ErrorIs(t, err, ErrSessionNotFound)
}

func TestSessionManager_List(t *testing.T) {
	sm := NewSessionManager()
	sm.Add(&ProcessSession{
		ID:        "test-1",
		Command:   "echo hello",
		Status:    "running",
		StartTime: 1000,
	})
	sm.Add(&ProcessSession{
		ID:        "test-2",
		Command:   "echo world",
		Status:    "running",
		StartTime: 1001,
	})
	sm.Add(&ProcessSession{
		ID:        "test-3",
		Command:   "echo done",
		Status:    "done",
		StartTime: 1002,
	})

	sessions := sm.List()
	require.Len(t, sessions, 3)

	ids := make(map[string]bool)
	for _, s := range sessions {
		ids[s.ID] = true
	}
	require.True(t, ids["test-1"])
	require.True(t, ids["test-2"])
	require.True(t, ids["test-3"])
}

func TestProcessSession_IsDone(t *testing.T) {
	session := &ProcessSession{Status: "running"}
	require.False(t, session.IsDone())

	session.Status = "done"
	require.True(t, session.IsDone())

	session.Status = "exited"
	require.True(t, session.IsDone())
}

func TestProcessSession_ToSessionInfo(t *testing.T) {
	session := &ProcessSession{
		ID:        "test-1",
		PID:       12345,
		Command:   "echo hello",
		Status:    "running",
		StartTime: 1000,
	}

	info := session.ToSessionInfo()
	require.Equal(t, "test-1", info.ID)
	require.Equal(t, "echo hello", info.Command)
	require.Equal(t, "running", info.Status)
	require.Equal(t, 12345, info.PID)
	require.Equal(t, int64(1000), info.StartedAt)
}

func TestSessionManager_CleanupUsesCompletedAt(t *testing.T) {
	sm := NewSessionManager()
	now := time.Now().Unix()

	oldDone := &ProcessSession{
		ID:          "old-done",
		Status:      "done",
		StartTime:   now - int64(2*time.Hour.Seconds()),
		CompletedAt: now - int64(2*time.Hour.Seconds()),
	}
	recentlyFinished := &ProcessSession{
		ID:          "recent-done",
		Status:      "done",
		StartTime:   now - int64(2*time.Hour.Seconds()),
		CompletedAt: now - 60,
	}
	stillRunning := &ProcessSession{
		ID:        "running",
		Status:    "running",
		StartTime: now - int64(2*time.Hour.Seconds()),
	}

	sm.Add(oldDone)
	sm.Add(recentlyFinished)
	sm.Add(stillRunning)
	sm.cleanupOldSessions()

	if _, err := sm.Get("old-done"); err == nil {
		t.Fatal("expected old completed session to be cleaned up")
	}
	if _, err := sm.Get("recent-done"); err != nil {
		t.Fatalf("recently finished session should be kept: %v", err)
	}
	if _, err := sm.Get("running"); err != nil {
		t.Fatalf("running session should be kept: %v", err)
	}
}
