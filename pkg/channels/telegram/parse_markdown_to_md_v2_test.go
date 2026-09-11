package telegram

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed testdata/md2_all_formats.txt
var md2AllFormats string

func Test_markdownToTelegramMarkdownV2(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "heading -> bolding",
			input:    `## HeadingH2 #`,
			expected: "*HeadingH2 \\#*",
		},
		{
			name:     "strikethrough",
			input:    "~strikethroughMD~",
			expected: "~strikethroughMD~",
		},
		{
			name:     "inline URL",
			input:    "[inline URL](http://www.example.com/)",
			expected: "[inline URL](http://www.example.com/)",
		},
		{
			name:     "all telegram formats",
			input:    md2AllFormats,
			expected: md2AllFormats,
		},
		{
			name:     "empty",
			input:    "",
			expected: "",
		},
		{
			name:     "one letter",
			input:    "o",
			expected: "o",
		},
		{
			name:     "",
			input:    "*Last update: ~10 24h*",
			expected: "*Last update: \\~10 24h*",
		},
		{
			name:     "",
			input:    "<Market Capitalization>",
			expected: "\\<Market Capitalization\\>",
		},
		{
			name:     "nested inline entities",
			input:    "*bold _italic_ still*",
			expected: "*bold _italic_ still*",
		},
		{
			name:     "sequential entities",
			input:    "*a* ~b~ `c`",
			expected: "*a* ~b~ `c`",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := markdownToTelegramMarkdownV2(tc.input)

			require.EqualValues(t, tc.expected, actual)
		})
	}
}

func Test_markdownToTelegramMarkdownV2_ManySequentialEntities(t *testing.T) {
	var b strings.Builder
	for i := 0; i < 2000; i++ {
		b.WriteString("*x*")
		b.WriteByte(' ')
	}
	out := markdownToTelegramMarkdownV2(b.String())
	require.Contains(t, out, "*x*")
	require.Greater(t, len(out), 1000)
}
