package parser

import (
	"regexp"

	"github.com/opsidian/parsec/ast"
	"github.com/opsidian/parsec/reader"
)

// Empty always matches and returns with an empty result
func Empty() Func {
	return Func(func(leftRecCtx IntMap, r *reader.Reader) *ParserResult {
		return NewParserResult(nil, NewResult(nil, r))
	})
}

// End matches the end of the input
func End() Func {
	return Func(func(leftRecCtx IntMap, r *reader.Reader) *ParserResult {
		if r.IsEOF() {
			return NewParserResult(nil, NewResult(ast.NewTerminalNode(reader.EOF, r.Cursor(), nil), r))
		}
		return nil
	})
}

// Rune matches one specific character
func Rune(char rune, token string) Func {
	return Func(func(leftRecCtx IntMap, r *reader.Reader) *ParserResult {
		if matches, pos := r.ReadMatch("^" + regexp.QuoteMeta(string(char))); matches != nil {
			return NewParserResult(nil, NewResult(ast.NewTerminalNode(token, pos, char), r))
		}
		return nil
	})
}
