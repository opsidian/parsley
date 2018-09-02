package text

import (
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
)

// LeftTrim skips the whitespaces before it tries to match the given parser
func LeftTrim(p parsley.Parser, wsMode WsMode) *parser.NamedFunc {
	return parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet) {
		return p.Parse(ctx, leftRecCtx, ctx.Reader().(*Reader).SkipWhitespaces(pos, wsMode))
	}).WithName(p.Name)
}

// RightTrim reads and skips the whitespaces after any parser matches and updates the reader position
func RightTrim(p parsley.Parser, wsMode WsMode) *parser.NamedFunc {
	return parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet) {
		tr := ctx.Reader().(*Reader)
		res, cp := p.Parse(ctx, leftRecCtx, pos)
		if res != nil {
			res = ast.SetReaderPos(res, func(pos parsley.Pos) parsley.Pos { return tr.SkipWhitespaces(pos, wsMode) })
		}

		if err := ctx.Error(); err != nil && int(err.Pos()) >= int(pos) {
			errPos := tr.SkipWhitespaces(err.Pos(), wsMode)
			if errPos > err.Pos() {
				ctx.SetError(errPos, err.Cause())
			}
		}
		return res, cp
	}).WithName(p.Name)
}

// Trim removes all whitespaces before and after the result token
func Trim(p parsley.Parser) *parser.NamedFunc {
	return RightTrim(LeftTrim(p, WsSpacesNl), WsSpacesNl)
}
