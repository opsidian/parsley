// Code generated by counterfeiter. DO NOT EDIT.
package parsleyfakes

import (
	"sync"

	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parsley"
)

type FakeParser struct {
	ParseStub        func(*parsley.Context, data.IntMap, parsley.Pos) (parsley.Node, data.IntSet, parsley.Error)
	parseMutex       sync.RWMutex
	parseArgsForCall []struct {
		arg1 *parsley.Context
		arg2 data.IntMap
		arg3 parsley.Pos
	}
	parseReturns struct {
		result1 parsley.Node
		result2 data.IntSet
		result3 parsley.Error
	}
	parseReturnsOnCall map[int]struct {
		result1 parsley.Node
		result2 data.IntSet
		result3 parsley.Error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeParser) Parse(arg1 *parsley.Context, arg2 data.IntMap, arg3 parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
	fake.parseMutex.Lock()
	ret, specificReturn := fake.parseReturnsOnCall[len(fake.parseArgsForCall)]
	fake.parseArgsForCall = append(fake.parseArgsForCall, struct {
		arg1 *parsley.Context
		arg2 data.IntMap
		arg3 parsley.Pos
	}{arg1, arg2, arg3})
	fake.recordInvocation("Parse", []interface{}{arg1, arg2, arg3})
	fake.parseMutex.Unlock()
	if fake.ParseStub != nil {
		return fake.ParseStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	fakeReturns := fake.parseReturns
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *FakeParser) ParseCallCount() int {
	fake.parseMutex.RLock()
	defer fake.parseMutex.RUnlock()
	return len(fake.parseArgsForCall)
}

func (fake *FakeParser) ParseCalls(stub func(*parsley.Context, data.IntMap, parsley.Pos) (parsley.Node, data.IntSet, parsley.Error)) {
	fake.parseMutex.Lock()
	defer fake.parseMutex.Unlock()
	fake.ParseStub = stub
}

func (fake *FakeParser) ParseArgsForCall(i int) (*parsley.Context, data.IntMap, parsley.Pos) {
	fake.parseMutex.RLock()
	defer fake.parseMutex.RUnlock()
	argsForCall := fake.parseArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeParser) ParseReturns(result1 parsley.Node, result2 data.IntSet, result3 parsley.Error) {
	fake.parseMutex.Lock()
	defer fake.parseMutex.Unlock()
	fake.ParseStub = nil
	fake.parseReturns = struct {
		result1 parsley.Node
		result2 data.IntSet
		result3 parsley.Error
	}{result1, result2, result3}
}

func (fake *FakeParser) ParseReturnsOnCall(i int, result1 parsley.Node, result2 data.IntSet, result3 parsley.Error) {
	fake.parseMutex.Lock()
	defer fake.parseMutex.Unlock()
	fake.ParseStub = nil
	if fake.parseReturnsOnCall == nil {
		fake.parseReturnsOnCall = make(map[int]struct {
			result1 parsley.Node
			result2 data.IntSet
			result3 parsley.Error
		})
	}
	fake.parseReturnsOnCall[i] = struct {
		result1 parsley.Node
		result2 data.IntSet
		result3 parsley.Error
	}{result1, result2, result3}
}

func (fake *FakeParser) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.parseMutex.RLock()
	defer fake.parseMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeParser) recordInvocation(key string, args []interface{}) {
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

var _ parsley.Parser = new(FakeParser)
