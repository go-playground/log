// Code generated by counterfeiter. DO NOT EDIT.
package fake

import (
	"sync"
)

type Client struct {
	MessageWithExtrasStub        func(level string, msg string, extras map[string]interface{})
	messageWithExtrasMutex       sync.RWMutex
	messageWithExtrasArgsForCall []struct {
		level  string
		msg    string
		extras map[string]interface{}
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Client) MessageWithExtras(level string, msg string, extras map[string]interface{}) {
	fake.messageWithExtrasMutex.Lock()
	fake.messageWithExtrasArgsForCall = append(fake.messageWithExtrasArgsForCall, struct {
		level  string
		msg    string
		extras map[string]interface{}
	}{level, msg, extras})
	fake.recordInvocation("MessageWithExtras", []interface{}{level, msg, extras})
	fake.messageWithExtrasMutex.Unlock()
	if fake.MessageWithExtrasStub != nil {
		fake.MessageWithExtrasStub(level, msg, extras)
	}
}

func (fake *Client) MessageWithExtrasCallCount() int {
	fake.messageWithExtrasMutex.RLock()
	defer fake.messageWithExtrasMutex.RUnlock()
	return len(fake.messageWithExtrasArgsForCall)
}

func (fake *Client) MessageWithExtrasArgsForCall(i int) (string, string, map[string]interface{}) {
	fake.messageWithExtrasMutex.RLock()
	defer fake.messageWithExtrasMutex.RUnlock()
	return fake.messageWithExtrasArgsForCall[i].level, fake.messageWithExtrasArgsForCall[i].msg, fake.messageWithExtrasArgsForCall[i].extras
}

func (fake *Client) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.messageWithExtrasMutex.RLock()
	defer fake.messageWithExtrasMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Client) recordInvocation(key string, args []interface{}) {
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