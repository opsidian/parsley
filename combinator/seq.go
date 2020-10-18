// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package combinator

import (
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
)

// Sequence is a recursive and-type combinator
type Sequence struct {
	token        string
	parserLookUp func(int) parsley.Parser
	lenCheck     func(int) bool
	interpreter  parsley.Interpreter
	customErr    error
	returnSingle bool
	returnEmpty  bool
}

// Seq tries to apply all parsers after each other and returns with all combinations of the results.
// The parsers should be generated by the parserLookUp function for the given index.
// When there is no parser for the given index then nil should be returned.
// The lenCheck function should return true if the longest possible match is valid
func Seq(token string, parserLookUp func(int) parsley.Parser, lenCheck func(int) bool) *Sequence {
	return &Sequence{
		token:        token,
		parserLookUp: parserLookUp,
		lenCheck:     lenCheck,
	}
}

// Bind binds the given interpreter
func (s *Sequence) Bind(interpreter parsley.Interpreter) *Sequence {
	s.interpreter = interpreter
	return s
}

// Parse parses the given input
func (s *Sequence) Parse(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
	p := &sequence{
		token:             s.token,
		parserLookUp:      s.parserLookUp,
		lenCheck:          s.lenCheck,
		interpreter:       s.interpreter,
		curtailingParsers: data.EmptyIntSet,
		nodes:             nil,
		returnSingle:      s.returnSingle,
		returnEmpty:       s.returnEmpty,
	}
	res, cp, err := p.Parse(ctx, leftRecCtx, pos)
	if err != nil && s.customErr != nil && err.Pos() == pos {
		err = parsley.NewError(pos, s.customErr)
	}

	return res, cp, err
}

// Name overrides the returned error if its position is the same as the reader's position
// The error will be: "was expecting <name>"
func (s *Sequence) Name(name string) *Sequence {
	s.customErr = parsley.NotFoundError(name)
	return s
}

// Token sets the result token
func (s *Sequence) Token(token string) *Sequence {
	s.token = token
	return s
}

// ReturnSingle will change the result of the parser if it returns with a non terminal node
// with a single child.
// In this case directly the child will returned.
func (s *Sequence) ReturnSingle() *Sequence {
	s.returnSingle = true
	return s
}

// ReturnEmpty will change the result of the parser if it returns with a non terminal node
// without children.
// In this case an empty node is returned instead.
func (s *Sequence) ReturnEmpty() *Sequence {
	s.returnEmpty = true
	return s
}

type sequence struct {
	token             string
	parserLookUp      func(i int) parsley.Parser
	lenCheck          func(i int) bool
	interpreter       parsley.Interpreter
	curtailingParsers data.IntSet
	result            parsley.Node
	err               parsley.Error
	nodes             []parsley.Node
	returnSingle      bool
	returnEmpty       bool
}

// Parse runs the recursive parser
func (s *sequence) Parse(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
	s.parse(0, ctx, leftRecCtx, pos, true)
	if s.result == nil {
		return nil, s.curtailingParsers, s.err
	}

	if s.err != nil {
		ctx.SetError(s.err)
	}

	return s.result, s.curtailingParsers, nil
}

func (s *sequence) parse(depth int, ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos, mergeCurtailingParsers bool) bool {
	var cp data.IntSet
	var res parsley.Node
	var err parsley.Error
	nextParser := s.parserLookUp(depth)
	if nextParser != nil {
		ctx.RegisterCall()
		res, cp, err = nextParser.Parse(ctx, leftRecCtx, pos)
		if err != nil && (s.err == nil || err.Pos() >= s.err.Pos()) {
			s.err = err
		}
	}

	if mergeCurtailingParsers {
		s.curtailingParsers = s.curtailingParsers.Union(cp)
	}

	if res != nil {
		switch rest := res.(type) {
		case ast.NodeList:
			for i, node := range rest {
				if s.parseNext(i, node, depth, ctx, leftRecCtx, pos, mergeCurtailingParsers) {
					return true
				}
			}
		default:
			if s.parseNext(0, rest, depth, ctx, leftRecCtx, pos, mergeCurtailingParsers) {
				return true
			}
		}
	}

	if res == nil {
		if s.lenCheck(depth) {
			if depth > 0 {
				if depth == 1 && s.returnSingle {
					s.result = ast.AppendNode(s.result, s.nodes[0])
				} else {
					nodesCopy := make([]parsley.Node, depth)
					copy(nodesCopy[0:depth], s.nodes[0:depth])
					s.result = ast.AppendNode(s.result, ast.NewNonTerminalNode(s.token, nodesCopy, s.interpreter))
				}
				if s.nodes[depth-1] != nil && s.nodes[depth-1].Token() == parser.EOF {
					return true
				}
			} else { // It's an empty result
				if s.returnEmpty {
					s.result = ast.AppendNode(s.result, ast.EmptyNode(pos))
				} else {
					s.result = ast.AppendNode(s.result, ast.NewEmptyNonTerminalNode(s.token, pos, s.interpreter))
				}
			}
		}
	}
	return false
}

func (s *sequence) parseNext(i int, node parsley.Node, depth int, ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos, mergeCurtailingParsers bool) bool {
	if len(s.nodes) < depth+1 {
		s.nodes = append(s.nodes, node)
	} else {
		s.nodes[depth] = node
	}
	if i > 0 || node.ReaderPos() > pos {
		leftRecCtx = data.EmptyIntMap
		mergeCurtailingParsers = false
	}
	if s.parse(depth+1, ctx, leftRecCtx, node.ReaderPos(), mergeCurtailingParsers) {
		return true
	}
	return false
}

// SeqOf tries to apply all parsers after each other and returns with all combinations of the results.
// Only matches are returned where all parsers were applied successfully.
func SeqOf(parsers ...parsley.Parser) *Sequence {
	l := len(parsers)
	lookup := func(i int) parsley.Parser {
		if i < l {
			return parsers[i]
		}
		return nil
	}
	lenCheck := func(len int) bool {
		return len == l
	}
	return Seq("SEQ", lookup, lenCheck)
}

// SeqTry tries to apply all parsers after each other and returns with all combinations of the results.
// It needs to match the first parser at least
func SeqTry(parsers ...parsley.Parser) *Sequence {
	l := len(parsers)
	lookup := func(i int) parsley.Parser {
		if i < l {
			return parsers[i]
		}
		return nil
	}
	lenCheck := func(len int) bool {
		return len > 0 && len <= l
	}
	return Seq("SEQ", lookup, lenCheck)
}

// SeqFirstOrAll tries to apply all parsers after each other and returns with all combinations of the results.
// If it can't match all parsers, but it can match the first one it will return with the result of the first one.
func SeqFirstOrAll(parsers ...parsley.Parser) *Sequence {
	l := len(parsers)
	lookup := func(i int) parsley.Parser {
		if i < l {
			return parsers[i]
		}
		return nil
	}
	lenCheck := func(len int) bool {
		return len == 1 || len == l
	}
	return Seq("SEQ", lookup, lenCheck)
}
