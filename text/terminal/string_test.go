package terminal_test

import (
	"fmt"
	"testing"

	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringShouldMatch(t *testing.T) {
	type TC struct {
		input    string
		expected string
		cursor   int
	}
	testCases := []TC{
		TC{`""`, "", 2},
		TC{`"a"`, "a", 3},
		TC{`"a" "b"`, "a", 3},
		TC{`"abcd"`, "abcd", 6},
		TC{`"'"`, "'", 3},
		TC{`"\a\b\f\n\r\t\v"`, "\a\b\f\n\r\t\v", 16},
		TC{`"\x67"`, "\x67", 6},
		TC{`"\uAB12"`, "\uAB12", 8},
		TC{`"\U0001F355"`, "\U0001F355", 12},
		TC{"``", ``, 2},
		TC{"`a`", `a`, 3},
		TC{"`a` `b`", `a`, 3},
		TC{"`abcd`", `abcd`, 6},
		TC{"`'`", `'`, 3},
		TC{"`" + `\a\b\f\n\r\t\v` + "`", `\a\b\f\n\r\t\v`, 16},
		TC{"`" + `\x67` + "`", `\x67`, 6},
		TC{"`" + `\uAB12` + "`", `\uAB12`, 8},
		TC{"`" + `\U0001F355` + "`", `\U0001F355`, 12},
	}
	for _, tc := range testCases {
		r := text.NewReader([]byte(tc.input), true)
		_, res := terminal.String().Parse(parser.EmptyLeftRecCtx(), r)
		require.NotNil(t, res, fmt.Sprintf("Failed to parse: %s", tc.input))
		actual, _ := res[0].Node().Value()
		assert.Equal(t, tc.expected, actual)
		assert.Equal(t, tc.cursor, res[0].Reader().Cursor().Pos())
	}
}

func TestStringShouldNotMatch(t *testing.T) {
	type TC struct {
		input string
	}
	testCases := []TC{
		TC{``},
		TC{`"`},
		TC{`'`},
		TC{"`"},
		TC{"''"},
		TC{"'a'"},
		TC{"5"},
		TC{"a"},
	}
	for _, tc := range testCases {
		r := text.NewReader([]byte(tc.input), true)
		_, res := terminal.String().Parse(parser.EmptyLeftRecCtx(), r)
		require.Nil(t, res, fmt.Sprintf("Should fail to parse: %s", tc.input))
	}
}