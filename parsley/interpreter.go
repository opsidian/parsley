// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsley

// Interpreter defines an interface to evaluate the given nonterminal node
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Interpreter
type Interpreter interface {
	Eval(userCtx interface{}, node NonTerminalNode) (interface{}, Error)
}

// StaticCheckerInterpreter defines an interpreter which is also a static checker
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . StaticCheckerInterpreter
type StaticCheckerInterpreter interface {
	Interpreter
	StaticChecker
}

// NodeTransformerInterpreter defines an interpreter which is also a node transformer
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . NodeTransformerInterpreter
type NodeTransformerInterpreter interface {
	Interpreter
	NodeTransformer
}
