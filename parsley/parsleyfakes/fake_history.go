// Code generated by counterfeiter. DO NOT EDIT.
package parsleyfakes

import (
	"sync"

	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parsley"
)

type FakeHistory struct {
	SaveResultStub        func(parserIndex int, pos int, result *parsley.Result)
	saveResultMutex       sync.RWMutex
	saveResultArgsForCall []struct {
		parserIndex int
		pos         int
		result      *parsley.Result
	}
	GetResultStub        func(parserIndex int, pos int, leftRecCtx data.IntMap) (*parsley.Result, bool)
	getResultMutex       sync.RWMutex
	getResultArgsForCall []struct {
		parserIndex int
		pos         int
		leftRecCtx  data.IntMap
	}
	getResultReturns struct {
		result1 *parsley.Result
		result2 bool
	}
	getResultReturnsOnCall map[int]struct {
		result1 *parsley.Result
		result2 bool
	}
	RegisterCallStub        func()
	registerCallMutex       sync.RWMutex
	registerCallArgsForCall []struct{}
	CallCountStub           func() int
	callCountMutex          sync.RWMutex
	callCountArgsForCall    []struct{}
	callCountReturns        struct {
		result1 int
	}
	callCountReturnsOnCall map[int]struct {
		result1 int
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeHistory) SaveResult(parserIndex int, pos int, result *parsley.Result) {
	fake.saveResultMutex.Lock()
	fake.saveResultArgsForCall = append(fake.saveResultArgsForCall, struct {
		parserIndex int
		pos         int
		result      *parsley.Result
	}{parserIndex, pos, result})
	fake.recordInvocation("SaveResult", []interface{}{parserIndex, pos, result})
	fake.saveResultMutex.Unlock()
	if fake.SaveResultStub != nil {
		fake.SaveResultStub(parserIndex, pos, result)
	}
}

func (fake *FakeHistory) SaveResultCallCount() int {
	fake.saveResultMutex.RLock()
	defer fake.saveResultMutex.RUnlock()
	return len(fake.saveResultArgsForCall)
}

func (fake *FakeHistory) SaveResultArgsForCall(i int) (int, int, *parsley.Result) {
	fake.saveResultMutex.RLock()
	defer fake.saveResultMutex.RUnlock()
	return fake.saveResultArgsForCall[i].parserIndex, fake.saveResultArgsForCall[i].pos, fake.saveResultArgsForCall[i].result
}

func (fake *FakeHistory) GetResult(parserIndex int, pos int, leftRecCtx data.IntMap) (*parsley.Result, bool) {
	fake.getResultMutex.Lock()
	ret, specificReturn := fake.getResultReturnsOnCall[len(fake.getResultArgsForCall)]
	fake.getResultArgsForCall = append(fake.getResultArgsForCall, struct {
		parserIndex int
		pos         int
		leftRecCtx  data.IntMap
	}{parserIndex, pos, leftRecCtx})
	fake.recordInvocation("GetResult", []interface{}{parserIndex, pos, leftRecCtx})
	fake.getResultMutex.Unlock()
	if fake.GetResultStub != nil {
		return fake.GetResultStub(parserIndex, pos, leftRecCtx)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getResultReturns.result1, fake.getResultReturns.result2
}

func (fake *FakeHistory) GetResultCallCount() int {
	fake.getResultMutex.RLock()
	defer fake.getResultMutex.RUnlock()
	return len(fake.getResultArgsForCall)
}

func (fake *FakeHistory) GetResultArgsForCall(i int) (int, int, data.IntMap) {
	fake.getResultMutex.RLock()
	defer fake.getResultMutex.RUnlock()
	return fake.getResultArgsForCall[i].parserIndex, fake.getResultArgsForCall[i].pos, fake.getResultArgsForCall[i].leftRecCtx
}

func (fake *FakeHistory) GetResultReturns(result1 *parsley.Result, result2 bool) {
	fake.GetResultStub = nil
	fake.getResultReturns = struct {
		result1 *parsley.Result
		result2 bool
	}{result1, result2}
}

func (fake *FakeHistory) GetResultReturnsOnCall(i int, result1 *parsley.Result, result2 bool) {
	fake.GetResultStub = nil
	if fake.getResultReturnsOnCall == nil {
		fake.getResultReturnsOnCall = make(map[int]struct {
			result1 *parsley.Result
			result2 bool
		})
	}
	fake.getResultReturnsOnCall[i] = struct {
		result1 *parsley.Result
		result2 bool
	}{result1, result2}
}

func (fake *FakeHistory) RegisterCall() {
	fake.registerCallMutex.Lock()
	fake.registerCallArgsForCall = append(fake.registerCallArgsForCall, struct{}{})
	fake.recordInvocation("RegisterCall", []interface{}{})
	fake.registerCallMutex.Unlock()
	if fake.RegisterCallStub != nil {
		fake.RegisterCallStub()
	}
}

func (fake *FakeHistory) RegisterCallCallCount() int {
	fake.registerCallMutex.RLock()
	defer fake.registerCallMutex.RUnlock()
	return len(fake.registerCallArgsForCall)
}

func (fake *FakeHistory) CallCount() int {
	fake.callCountMutex.Lock()
	ret, specificReturn := fake.callCountReturnsOnCall[len(fake.callCountArgsForCall)]
	fake.callCountArgsForCall = append(fake.callCountArgsForCall, struct{}{})
	fake.recordInvocation("CallCount", []interface{}{})
	fake.callCountMutex.Unlock()
	if fake.CallCountStub != nil {
		return fake.CallCountStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.callCountReturns.result1
}

func (fake *FakeHistory) CallCountCallCount() int {
	fake.callCountMutex.RLock()
	defer fake.callCountMutex.RUnlock()
	return len(fake.callCountArgsForCall)
}

func (fake *FakeHistory) CallCountReturns(result1 int) {
	fake.CallCountStub = nil
	fake.callCountReturns = struct {
		result1 int
	}{result1}
}

func (fake *FakeHistory) CallCountReturnsOnCall(i int, result1 int) {
	fake.CallCountStub = nil
	if fake.callCountReturnsOnCall == nil {
		fake.callCountReturnsOnCall = make(map[int]struct {
			result1 int
		})
	}
	fake.callCountReturnsOnCall[i] = struct {
		result1 int
	}{result1}
}

func (fake *FakeHistory) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.saveResultMutex.RLock()
	defer fake.saveResultMutex.RUnlock()
	fake.getResultMutex.RLock()
	defer fake.getResultMutex.RUnlock()
	fake.registerCallMutex.RLock()
	defer fake.registerCallMutex.RUnlock()
	fake.callCountMutex.RLock()
	defer fake.callCountMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeHistory) recordInvocation(key string, args []interface{}) {
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

var _ parsley.History = new(FakeHistory)