package ast_test

import (
	"errors"
	"testing"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/test"
	"github.com/stretchr/testify/assert"
)

func TestTerminalNode(t *testing.T) {
	pos := test.NewPosition(1)
	node := ast.NewTerminalNode("TOKEN", pos, "VALUE")
	assert.Equal(t, "TOKEN", node.Token())
	assert.Equal(t, pos, node.Pos())
	actualVal, actualErr := node.Value()
	assert.Equal(t, "VALUE", actualVal)
	assert.Nil(t, actualErr)

	assert.Equal(t, "T{VALUE, Pos{1}}", node.String())

	node = ast.NewTerminalNode("TOKEN", pos, nil)
	assert.Equal(t, "T{TOKEN, Pos{1}}", node.String())
}

func TestNonTerminalNode(t *testing.T) {
	expectedValue := 3
	expectedErr := errors.New("E")
	var actualValues []interface{}
	interpreterFunc := ast.InterpreterFunc(func(values []interface{}) (interface{}, error) {
		actualValues = values
		return expectedValue, expectedErr
	})

	nodes := []ast.Node{
		ast.NewTerminalNode("1", test.NewPosition(0), 1),
		ast.NewTerminalNode("2", test.NewPosition(2), 2),
	}

	node := ast.NewNonTerminalNode("+", nodes, interpreterFunc)
	assert.Equal(t, "+", node.Token())
	assert.Equal(t, nodes[0].Pos(), node.Pos())
	assert.Equal(t, nodes, node.Children())
	actualVal, actualErr := node.Value()
	assert.Equal(t, []interface{}{1, 2}, actualValues)
	assert.Equal(t, expectedValue, actualVal)
	assert.Equal(t, expectedErr, actualErr)

	assert.Equal(t, "NT{+, [T{1, Pos{0}} T{2, Pos{2}}]}", node.String())
}

func TestNonTerminalNodeShouldPanicWithoutChildren(t *testing.T) {
	assert.Panics(t, func() { ast.NewNonTerminalNode("X", nil, getInterpreterFunc(nil, nil)) })
}

func TestNonTerminalNodeShouldPanicWithoutInterpreter(t *testing.T) {
	randomNode := ast.NewTerminalNode("X", test.NewPosition(1), getInterpreterFunc(nil, nil))
	assert.Panics(t, func() { ast.NewNonTerminalNode("X", []ast.Node{randomNode}, nil) })
}

func TestNonTerminalNodeValueShouldReturnErrorIfChildHasError(t *testing.T) {
	expectedErr := errors.New("E")
	randomNode := ast.NewTerminalNode("X", test.NewPosition(1), getInterpreterFunc(nil, nil))
	badChild := ast.NewNonTerminalNode("BAD", []ast.Node{randomNode}, ast.InterpreterFunc(func(values []interface{}) (interface{}, error) {
		return nil, expectedErr
	}))
	node := ast.NewNonTerminalNode("X", []ast.Node{badChild}, getInterpreterFunc(nil, nil))
	actualVal, actualErr := node.Value()
	assert.Nil(t, actualVal)
	assert.Equal(t, expectedErr, actualErr)
}

func TestNonTerminalNodeValueShouldHandleNilNodes(t *testing.T) {
	var actualValues []interface{}
	interpreterFunc := ast.InterpreterFunc(func(values []interface{}) (interface{}, error) {
		actualValues = values
		return nil, nil
	})

	nodes := []ast.Node{
		ast.NewTerminalNode("1", test.NewPosition(0), 1),
		nil,
		ast.NewTerminalNode("2", test.NewPosition(2), 2),
	}

	node := ast.NewNonTerminalNode("+", nodes, interpreterFunc)
	node.Value()
	assert.Equal(t, []interface{}{1, nil, 2}, actualValues)
}

func TestNonTerminalNodeValueShouldIgnoreEmptyNodes(t *testing.T) {
	var actualValues []interface{}
	interpreterFunc := ast.InterpreterFunc(func(values []interface{}) (interface{}, error) {
		actualValues = values
		return nil, nil
	})

	nodes := []ast.Node{
		ast.NewTerminalNode("1", test.NewPosition(0), 1),
		ast.NewTerminalNode(ast.EMPTY, test.NewPosition(0), nil),
		ast.NewTerminalNode("2", test.NewPosition(2), 2),
	}

	node := ast.NewNonTerminalNode("+", nodes, interpreterFunc)
	node.Value()
	assert.Equal(t, []interface{}{1, 2}, actualValues)
}