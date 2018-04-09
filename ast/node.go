// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ast

import (
	"fmt"

	"github.com/opsidian/parsley/parsley"
)

// NIL is the nil token
const NIL = "NIL"

// EOF is the end of file token
const EOF = "EOF"

// TerminalNode is a leaf node in the AST
type TerminalNode struct {
	token     string
	value     interface{}
	pos       parsley.Pos
	readerPos int
}

// NewTerminalNode creates a new TerminalNode instance
func NewTerminalNode(token string, value interface{}, pos parsley.Pos, readerPos int) *TerminalNode {
	return &TerminalNode{
		token:     token,
		value:     value,
		pos:       pos,
		readerPos: readerPos,
	}
}

// Token returns with the node token
func (t *TerminalNode) Token() string {
	return t.token
}

// Value returns with the value of the node
func (t *TerminalNode) Value(ctx interface{}) (interface{}, parsley.Error) {
	return t.value, nil
}

// Pos returns the position
func (t *TerminalNode) Pos() parsley.Pos {
	return t.pos
}

// ReaderPos returns the position of the first character immediately after this node
func (t *TerminalNode) ReaderPos() int {
	return t.readerPos
}

func (t *TerminalNode) SetReaderPos(f func(int) int) {
	t.readerPos = f(t.readerPos)
}

// String returns with a string representation of the node
func (t *TerminalNode) String() string {
	if t.value == nil {
		return fmt.Sprintf("T{%s, %d, %d}", t.token, t.pos, t.readerPos)
	}
	return fmt.Sprintf("T{%v, %d, %d}", t.value, t.pos, t.readerPos)
}

// NonTerminalNode represents a branch node in the AST
type NonTerminalNode struct {
	token       string
	children    []parsley.Node
	pos         parsley.Pos
	readerPos   int
	interpreter parsley.Interpreter
}

// NewNonTerminalNode creates a new NonTerminalNode instance
func NewNonTerminalNode(token string, children []parsley.Node) *NonTerminalNode {
	if len(children) == 0 {
		panic("a NonTerminalNode should always have children")
	}
	node := &NonTerminalNode{
		token:    token,
		children: children,
	}

	for _, child := range children {
		if child.Pos() != parsley.NilPos {
			node.pos = child.Pos()
			break
		}
	}
	node.readerPos = children[len(children)-1].ReaderPos()

	return node
}

// Bind adds an interpreter function to the node
func (n *NonTerminalNode) Bind(interpreter parsley.Interpreter) *NonTerminalNode {
	n.interpreter = interpreter
	return n
}

// Token returns with the node token
func (n *NonTerminalNode) Token() string {
	return n.token
}

// Value returns with the value of the node
func (n *NonTerminalNode) Value(ctx interface{}) (interface{}, parsley.Error) {
	if n.interpreter == nil {
		return nil, nil
	}
	return n.interpreter.Eval(ctx, n.children)
}

// Pos returns the position
func (n *NonTerminalNode) Pos() parsley.Pos {
	return n.pos
}

// Children returns with the children
func (n *NonTerminalNode) Children() []parsley.Node {
	return n.children
}

// ReaderPos returns the position of the first character immediately after this node
func (n *NonTerminalNode) ReaderPos() int {
	return n.readerPos
}

// SetReaderPos amends the reader position using the given function
func (n *NonTerminalNode) SetReaderPos(f func(int) int) {
	n.readerPos = f(n.readerPos)
}

// String returns with a string representation of the node
func (n *NonTerminalNode) String() string {
	return fmt.Sprintf("NT{%s, %s, %d, %d}", n.token, n.children, n.pos, n.readerPos)
}

// NodeList contains a list of nodes, should be used when a parser returns with multiple results
type NodeList []parsley.Node

// Token returns with NODE_LIST
func (nl NodeList) Token() string {
	return "NODE_LIST"
}

// Value returns with the value of the first result
func (nl NodeList) Value(ctx interface{}) (interface{}, parsley.Error) {
	if len(nl) > 0 {
		return nl[0].Value(ctx)
	}
	return nil, nil
}

// Pos returns the value of the first pos (all nodes should have the same position)
func (nl NodeList) Pos() parsley.Pos {
	if len(nl) > 0 {
		return nl[0].Pos()
	}
	return parsley.NilPos
}

// ReaderPos should not be called on a NodeList
func (nl NodeList) ReaderPos() int {
	panic("ReaderPos() should not be called on NodeList")
}

// Append appends a new node to the list
func (nl *NodeList) Append(node parsley.Node) {
	switch v := node.(type) {
	case NodeList:
		for _, node := range v {
			nl.Append(node)
		}
	case EmptyNode:
		for _, node := range *nl {
			if node == v {
				return
			}
		}
		*nl = append(*nl, v)
	default:
		*nl = append(*nl, v)
	}
}

// Walk runs the given function on all nodes
func (nl NodeList) Walk(f func(i int, n parsley.Node) bool) {
	for i, node := range nl {
		if f(i, node) {
			break
		}
	}
}

// EmptyNode represents an empty node
type EmptyNode int

// Token returns with EMPTY
func (e EmptyNode) Token() string {
	return NIL
}

// Value returns with nil
func (e EmptyNode) Value(ctx interface{}) (interface{}, parsley.Error) {
	return nil, nil
}

// Pos returns with NilPosition
func (e EmptyNode) Pos() parsley.Pos {
	return parsley.NilPos
}

// ReaderPos returns the reader position
func (e EmptyNode) ReaderPos() int {
	return int(e)
}

// String returns with a string representation of the node
func (e EmptyNode) String() string {
	return "NIL"
}

// AppendNode appends
func AppendNode(n1, n2 parsley.Node) parsley.Node {
	if n1 == nil {
		return n2
	}
	if n2 == nil {
		return n1
	}
	switch n := n1.(type) {
	case NodeList:
		n.Append(n2)
		return n
	default:
		nl := NodeList([]parsley.Node{n1})
		nl.Append(n2)
		return nl
	}
}

// Walkable is a generic interface to allow to apply a function on the node
type Walkable interface {
	Walk(f func(i int, n parsley.Node) bool)
}

// WalkNode applies the given function to the node
func WalkNode(node parsley.Node, f func(i int, n parsley.Node) bool) {
	switch n := node.(type) {
	case Walkable:
		n.Walk(f)
	default:
		f(0, node)
	}
}

// ReaderPosSetter allows to change the reader position on a node
type ReaderPosSetter interface {
	SetReaderPos(f func(int) int)
}

// SetReaderPos sets the reader position on a node
func SetReaderPos(node parsley.Node, f func(int) int) parsley.Node {
	switch n := node.(type) {
	case ReaderPosSetter:
		n.SetReaderPos(f)
	case EmptyNode:
		return EmptyNode(f(int(n)))
	case NodeList:
		for i, item := range n {
			n[i] = SetReaderPos(item, f)
		}
	default:
		panic("invalid node type for SetReaderPos(), you may need to implement the ast.ReaderPosSetter interface")
	}
	return node
}
