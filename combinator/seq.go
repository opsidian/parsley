// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package combinator

import (
	"fmt"

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
		nodes:             []parsley.Node{},
	}
	return p.Parse(ctx, leftRecCtx, pos)
}

// Name returns with a new parser function which overrides the returned error
// if its position is the same as the reader's position
// The error will be: "was expecting <name>"
func (s *Sequence) Name(name string) parser.Func {
	return parser.ReturnError(s, fmt.Errorf("was expecting %s", name))
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
				nodesCopy := make([]parsley.Node, depth)
				copy(nodesCopy[0:depth], s.nodes[0:depth])
				s.result = ast.AppendNode(s.result, ast.NewNonTerminalNode(s.token, nodesCopy, s.interpreter))
				if s.nodes[depth-1] != nil && s.nodes[depth-1].Token() == ast.EOF {
					return true
				}
			} else { // It's an empty result
				s.result = ast.AppendNode(s.result, ast.NewEmptyNonTerminalNode(s.token, pos, s.interpreter))
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
	return newSeq(len(parsers), parsers...)
}

// SeqTry tries to apply all parsers after each other and returns with all combinations of the results.
// It needs to match the first parser at least
func SeqTry(parsers ...parsley.Parser) *Sequence {
	return newSeq(1, parsers...)
}

func newSeq(min int, parsers ...parsley.Parser) *Sequence {
	lookup := func(i int) parsley.Parser {
		if i < len(parsers) {
			return parsers[i]
		}
		return nil
	}
	l := len(parsers)
	lenCheck := func(len int) bool {
		return len >= min && len <= l
	}
	return Seq("SEQ", lookup, lenCheck)
}
