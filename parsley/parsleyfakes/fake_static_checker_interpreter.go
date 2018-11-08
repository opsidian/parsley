// Code generated by counterfeiter. DO NOT EDIT.
package parsleyfakes

import (
	"sync"

	"github.com/opsidian/parsley/parsley"
)

type FakeStaticCheckerInterpreter struct {
	EvalStub        func(userCtx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error)
	evalMutex       sync.RWMutex
	evalArgsForCall []struct {
		userCtx interface{}
		node    parsley.NonTerminalNode
	}
	evalReturns struct {
		result1 interface{}
		result2 parsley.Error
	}
	evalReturnsOnCall map[int]struct {
		result1 interface{}
		result2 parsley.Error
	}
	StaticCheckStub        func(userCtx interface{}, node parsley.NonTerminalNode) (string, parsley.Error)
	staticCheckMutex       sync.RWMutex
	staticCheckArgsForCall []struct {
		userCtx interface{}
		node    parsley.NonTerminalNode
	}
	staticCheckReturns struct {
		result1 string
		result2 parsley.Error
	}
	staticCheckReturnsOnCall map[int]struct {
		result1 string
		result2 parsley.Error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeStaticCheckerInterpreter) Eval(userCtx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
	fake.evalMutex.Lock()
	ret, specificReturn := fake.evalReturnsOnCall[len(fake.evalArgsForCall)]
	fake.evalArgsForCall = append(fake.evalArgsForCall, struct {
		userCtx interface{}
		node    parsley.NonTerminalNode
	}{userCtx, node})
	fake.recordInvocation("Eval", []interface{}{userCtx, node})
	fake.evalMutex.Unlock()
	if fake.EvalStub != nil {
		return fake.EvalStub(userCtx, node)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.evalReturns.result1, fake.evalReturns.result2
}

func (fake *FakeStaticCheckerInterpreter) EvalCallCount() int {
	fake.evalMutex.RLock()
	defer fake.evalMutex.RUnlock()
	return len(fake.evalArgsForCall)
}

func (fake *FakeStaticCheckerInterpreter) EvalArgsForCall(i int) (interface{}, parsley.NonTerminalNode) {
	fake.evalMutex.RLock()
	defer fake.evalMutex.RUnlock()
	return fake.evalArgsForCall[i].userCtx, fake.evalArgsForCall[i].node
}

func (fake *FakeStaticCheckerInterpreter) EvalReturns(result1 interface{}, result2 parsley.Error) {
	fake.EvalStub = nil
	fake.evalReturns = struct {
		result1 interface{}
		result2 parsley.Error
	}{result1, result2}
}

func (fake *FakeStaticCheckerInterpreter) EvalReturnsOnCall(i int, result1 interface{}, result2 parsley.Error) {
	fake.EvalStub = nil
	if fake.evalReturnsOnCall == nil {
		fake.evalReturnsOnCall = make(map[int]struct {
			result1 interface{}
			result2 parsley.Error
		})
	}
	fake.evalReturnsOnCall[i] = struct {
		result1 interface{}
		result2 parsley.Error
	}{result1, result2}
}

func (fake *FakeStaticCheckerInterpreter) StaticCheck(userCtx interface{}, node parsley.NonTerminalNode) (string, parsley.Error) {
	fake.staticCheckMutex.Lock()
	ret, specificReturn := fake.staticCheckReturnsOnCall[len(fake.staticCheckArgsForCall)]
	fake.staticCheckArgsForCall = append(fake.staticCheckArgsForCall, struct {
		userCtx interface{}
		node    parsley.NonTerminalNode
	}{userCtx, node})
	fake.recordInvocation("StaticCheck", []interface{}{userCtx, node})
	fake.staticCheckMutex.Unlock()
	if fake.StaticCheckStub != nil {
		return fake.StaticCheckStub(userCtx, node)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.staticCheckReturns.result1, fake.staticCheckReturns.result2
}

func (fake *FakeStaticCheckerInterpreter) StaticCheckCallCount() int {
	fake.staticCheckMutex.RLock()
	defer fake.staticCheckMutex.RUnlock()
	return len(fake.staticCheckArgsForCall)
}

func (fake *FakeStaticCheckerInterpreter) StaticCheckArgsForCall(i int) (interface{}, parsley.NonTerminalNode) {
	fake.staticCheckMutex.RLock()
	defer fake.staticCheckMutex.RUnlock()
	return fake.staticCheckArgsForCall[i].userCtx, fake.staticCheckArgsForCall[i].node
}

func (fake *FakeStaticCheckerInterpreter) StaticCheckReturns(result1 string, result2 parsley.Error) {
	fake.StaticCheckStub = nil
	fake.staticCheckReturns = struct {
		result1 string
		result2 parsley.Error
	}{result1, result2}
}

func (fake *FakeStaticCheckerInterpreter) StaticCheckReturnsOnCall(i int, result1 string, result2 parsley.Error) {
	fake.StaticCheckStub = nil
	if fake.staticCheckReturnsOnCall == nil {
		fake.staticCheckReturnsOnCall = make(map[int]struct {
			result1 string
			result2 parsley.Error
		})
	}
	fake.staticCheckReturnsOnCall[i] = struct {
		result1 string
		result2 parsley.Error
	}{result1, result2}
}

func (fake *FakeStaticCheckerInterpreter) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.evalMutex.RLock()
	defer fake.evalMutex.RUnlock()
	fake.staticCheckMutex.RLock()
	defer fake.staticCheckMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeStaticCheckerInterpreter) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ parsley.StaticCheckerInterpreter = new(FakeStaticCheckerInterpreter)
