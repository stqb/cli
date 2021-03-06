// This file was generated by counterfeiter
package v2fakes

import (
	"sync"

	"code.cloudfoundry.org/cli/actor/v2action"
	"code.cloudfoundry.org/cli/command/v2"
)

type FakeSetHealthCheckActor struct {
	SetApplicationHealthCheckTypeByNameAndSpaceStub        func(name string, spaceGUID string, healthCheckType string) (v2action.Warnings, error)
	setApplicationHealthCheckTypeByNameAndSpaceMutex       sync.RWMutex
	setApplicationHealthCheckTypeByNameAndSpaceArgsForCall []struct {
		name            string
		spaceGUID       string
		healthCheckType string
	}
	setApplicationHealthCheckTypeByNameAndSpaceReturns struct {
		result1 v2action.Warnings
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeSetHealthCheckActor) SetApplicationHealthCheckTypeByNameAndSpace(name string, spaceGUID string, healthCheckType string) (v2action.Warnings, error) {
	fake.setApplicationHealthCheckTypeByNameAndSpaceMutex.Lock()
	fake.setApplicationHealthCheckTypeByNameAndSpaceArgsForCall = append(fake.setApplicationHealthCheckTypeByNameAndSpaceArgsForCall, struct {
		name            string
		spaceGUID       string
		healthCheckType string
	}{name, spaceGUID, healthCheckType})
	fake.recordInvocation("SetApplicationHealthCheckTypeByNameAndSpace", []interface{}{name, spaceGUID, healthCheckType})
	fake.setApplicationHealthCheckTypeByNameAndSpaceMutex.Unlock()
	if fake.SetApplicationHealthCheckTypeByNameAndSpaceStub != nil {
		return fake.SetApplicationHealthCheckTypeByNameAndSpaceStub(name, spaceGUID, healthCheckType)
	} else {
		return fake.setApplicationHealthCheckTypeByNameAndSpaceReturns.result1, fake.setApplicationHealthCheckTypeByNameAndSpaceReturns.result2
	}
}

func (fake *FakeSetHealthCheckActor) SetApplicationHealthCheckTypeByNameAndSpaceCallCount() int {
	fake.setApplicationHealthCheckTypeByNameAndSpaceMutex.RLock()
	defer fake.setApplicationHealthCheckTypeByNameAndSpaceMutex.RUnlock()
	return len(fake.setApplicationHealthCheckTypeByNameAndSpaceArgsForCall)
}

func (fake *FakeSetHealthCheckActor) SetApplicationHealthCheckTypeByNameAndSpaceArgsForCall(i int) (string, string, string) {
	fake.setApplicationHealthCheckTypeByNameAndSpaceMutex.RLock()
	defer fake.setApplicationHealthCheckTypeByNameAndSpaceMutex.RUnlock()
	return fake.setApplicationHealthCheckTypeByNameAndSpaceArgsForCall[i].name, fake.setApplicationHealthCheckTypeByNameAndSpaceArgsForCall[i].spaceGUID, fake.setApplicationHealthCheckTypeByNameAndSpaceArgsForCall[i].healthCheckType
}

func (fake *FakeSetHealthCheckActor) SetApplicationHealthCheckTypeByNameAndSpaceReturns(result1 v2action.Warnings, result2 error) {
	fake.SetApplicationHealthCheckTypeByNameAndSpaceStub = nil
	fake.setApplicationHealthCheckTypeByNameAndSpaceReturns = struct {
		result1 v2action.Warnings
		result2 error
	}{result1, result2}
}

func (fake *FakeSetHealthCheckActor) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.setApplicationHealthCheckTypeByNameAndSpaceMutex.RLock()
	defer fake.setApplicationHealthCheckTypeByNameAndSpaceMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeSetHealthCheckActor) recordInvocation(key string, args []interface{}) {
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

var _ v2.SetHealthCheckActor = new(FakeSetHealthCheckActor)
