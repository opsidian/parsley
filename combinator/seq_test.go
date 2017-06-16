package combinator_test

import (
	"testing"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/ast/builder"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
	"github.com/opsidian/parsley/test"
	"github.com/stretchr/testify/assert"
)

func TestSeqShouldPanicIfNoParserWasGiven(t *testing.T) {
	r := test.NewReader(0, 1, false, false)
	ctx := parser.EmptyLeftRecCtx()
	assert.Panics(t, func() { combinator.Seq(nil).Parse(ctx, r) })
	assert.Panics(t, func() { combinator.SeqTry(nil, 1).Parse(ctx, r) })
}

func TestSeqShouldHandleOnlyOneParser(t *testing.T) {
	r := test.NewReader(0, 1, false, false)
	expectedRS := parser.NewResult(ast.NewTerminalNode("CHAR", r.Cursor(), 'a'), r.Clone()).AsSet()
	expectedCP := data.NewIntSet(1)
	p1 := parser.Func(func(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		return expectedCP, expectedRS
	})

	nodeBuilder := ast.NodeBuilderFunc(func(nodes []ast.Node) ast.Node { return nodes[0] })

	cp, rs := combinator.Seq(nodeBuilder, p1).Parse(parser.EmptyLeftRecCtx(), r)
	assert.Equal(t, expectedCP, cp)
	assert.Equal(t, expectedRS, rs)

	cp, rs = combinator.SeqTry(nodeBuilder, 0, p1).Parse(parser.EmptyLeftRecCtx(), r)
	assert.Equal(t, expectedCP, cp)
	assert.Equal(t, expectedRS, rs)
}

func TestSeqShouldCombineParserResults(t *testing.T) {
	parser.Stat.Reset()
	r := test.NewReader(0, 1, false, false)

	p1 := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		return parser.NoCurtailingParsers(), parser.NewResultSet(
			parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "a"), test.NewReader(1, 1, false, true)),
			parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "ab"), test.NewReader(2, 1, false, true)),
		)
	})
	p2First := true
	p2 := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		if p2First {
			p2First = false
			return parser.NoCurtailingParsers(), parser.NewResultSet(
				parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(2), "b"), test.NewReader(3, 1, false, true)),
				parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(2), "c"), test.NewReader(4, 1, false, true)),
			)
		} else {
			return parser.NoCurtailingParsers(), parser.NewResultSet(
				parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(2), "d"), test.NewReader(5, 1, false, true)),
			)
		}
	})
	nodeBuilder := ast.NodeBuilderFunc(func(nodes []ast.Node) ast.Node {
		var res string
		for _, node := range nodes {
			val, _ := node.Value()
			res += val.(string)
		}
		first := nodes[0].(ast.TerminalNode)
		return ast.NewTerminalNode("STR", first.Pos(), res)
	})

	_, rs := combinator.Seq(nodeBuilder, p1, p2).Parse(parser.EmptyLeftRecCtx(), r)
	assert.EqualValues(t, parser.NewResultSet(
		parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "ab"), test.NewReader(3, 1, false, true)),
		parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "ac"), test.NewReader(4, 1, false, true)),
		parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "abd"), test.NewReader(5, 1, false, true)),
	), rs)

	assert.EqualValues(t, 3, parser.Stat.GetSumCallCount())

	p2First = true
	_, rs = combinator.SeqTry(nodeBuilder, 0, p1, p2).Parse(parser.EmptyLeftRecCtx(), r)
	assert.EqualValues(t, parser.NewResultSet(
		parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "ab"), test.NewReader(3, 1, false, true)),
		parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "ac"), test.NewReader(4, 1, false, true)),
		parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "abd"), test.NewReader(5, 1, false, true)),
	), rs)

	assert.EqualValues(t, 6, parser.Stat.GetSumCallCount())
}

func TestSeqShouldHandleNilResults(t *testing.T) {
	r := test.NewReader(0, 1, false, false)

	p1 := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		return parser.NoCurtailingParsers(), parser.NewResult(ast.NewTerminalNode("CHAR", test.NewPosition(1), 'x'), r).AsSet()
	})

	p2 := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		return parser.NoCurtailingParsers(), nil
	})

	cp, rs := combinator.Seq(builder.Nil(), p1, p2).Parse(parser.EmptyLeftRecCtx(), r)
	assert.Equal(t, parser.NoCurtailingParsers(), cp)
	assert.Empty(t, rs)

	cp, rs = combinator.SeqTry(builder.Nil(), 2, p1, p2).Parse(parser.EmptyLeftRecCtx(), r)
	assert.Equal(t, parser.NoCurtailingParsers(), cp)
	assert.Empty(t, rs)
}

func TestSeqTryShouldMatchLongestSequence(t *testing.T) {
	r := test.NewReader(0, 1, false, false)

	res := parser.NewResult(ast.NewTerminalNode("CHAR", test.NewPosition(1), 'x'), r)

	p1 := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		return parser.NoCurtailingParsers(), res.AsSet()
	})

	p2 := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		return parser.NoCurtailingParsers(), nil
	})

	cp, rs := combinator.SeqTry(builder.All("TEST", nil), 1, p1, p2).Parse(parser.EmptyLeftRecCtx(), r)
	assert.Equal(t, parser.NoCurtailingParsers(), cp)
	assert.Equal(t, parser.NewResult(ast.NewNonTerminalNode("TEST", []ast.Node{ast.NewTerminalNode("CHAR", test.NewPosition(1), 'x')}, nil), r.Clone()).AsSet(), rs)
}

func TestSeqShouldMergeCurtailReasonsIfEmptyResult(t *testing.T) {
	r := test.NewReader(0, 1, false, false)

	p1 := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		return data.NewIntSet(0, 1), parser.NewResult(nil, r).AsSet()
	})

	p2 := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		return data.NewIntSet(1, 2), nil
	})

	cp, _ := combinator.Seq(builder.Nil(), p1, p2).Parse(parser.EmptyLeftRecCtx(), r)
	assert.EqualValues(t, data.NewIntSet(0, 1, 2), cp)

	cp, _ = combinator.SeqTry(builder.Nil(), 0, p1, p2).Parse(parser.EmptyLeftRecCtx(), r)
	assert.EqualValues(t, data.NewIntSet(0, 1, 2), cp)
}

func TestSeqShouldStopIfEOFTokenReached(t *testing.T) {
	r := test.NewReader(0, 2, false, false)

	p1 := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		return parser.NoCurtailingParsers(), parser.NewResult(ast.NewTerminalNode("CHAR", test.NewPosition(1), 'a'), r).AsSet()
	})

	p2 := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		return parser.NoCurtailingParsers(), parser.NewResultSet(
			parser.NewResult(ast.NewTerminalNode(ast.EOF, test.NewPosition(2), nil), test.NewReader(2, 0, false, true)),
			parser.NewResult(ast.NewTerminalNode("CHAR", test.NewPosition(1), 'b'), test.NewReader(1, 1, false, true)),
		)
	})

	nodeBuilder := ast.NodeBuilderFunc(func(nodes []ast.Node) ast.Node {
		return nodes[0]
	})

	_, rs := combinator.Seq(nodeBuilder, p1, p2).Parse(parser.EmptyLeftRecCtx(), r)
	assert.EqualValues(t,
		parser.NewResult(ast.NewTerminalNode("CHAR", test.NewPosition(1), 'a'), test.NewReader(2, 0, false, true)).AsSet(),
		rs,
	)

	_, rs = combinator.SeqTry(nodeBuilder, 0, p1, p2).Parse(parser.EmptyLeftRecCtx(), r)
	assert.EqualValues(t,
		parser.NewResult(ast.NewTerminalNode("CHAR", test.NewPosition(1), 'a'), test.NewReader(2, 0, false, true)).AsSet(),
		rs,
	)
}
