// This file was generated by counterfeiter
package fakes

import (
	"github.com/robdimsdale/wundergo"
	"sync"
)

type FakeHTTPHelper struct {
	GetStub        func(url string) ([]byte, error)
	getMutex       sync.RWMutex
	getArgsForCall []struct {
		url string
	}
	getReturns struct {
		result1 []byte
		result2 error
	}
	PostStub        func(url string, body []byte) ([]byte, error)
	postMutex       sync.RWMutex
	postArgsForCall []struct {
		url  string
		body []byte
	}
	postReturns struct {
		result1 []byte
		result2 error
	}
	PutStub        func(url string, body []byte) ([]byte, error)
	putMutex       sync.RWMutex
	putArgsForCall []struct {
		url  string
		body []byte
	}
	putReturns struct {
		result1 []byte
		result2 error
	}
	PatchStub        func(url string, body []byte) ([]byte, error)
	patchMutex       sync.RWMutex
	patchArgsForCall []struct {
		url  string
		body []byte
	}
	patchReturns struct {
		result1 []byte
		result2 error
	}
	DeleteStub        func(url string) error
	deleteMutex       sync.RWMutex
	deleteArgsForCall []struct {
		url string
	}
	deleteReturns struct {
		result1 error
	}
}

func (fake *FakeHTTPHelper) Get(url string) ([]byte, error) {
	fake.getMutex.Lock()
	defer fake.getMutex.Unlock()
	fake.getArgsForCall = append(fake.getArgsForCall, struct {
		url string
	}{url})
	if fake.GetStub != nil {
		return fake.GetStub(url)
	} else {
		return fake.getReturns.result1, fake.getReturns.result2
	}
}

func (fake *FakeHTTPHelper) GetCallCount() int {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	return len(fake.getArgsForCall)
}

func (fake *FakeHTTPHelper) GetArgsForCall(i int) string {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	return fake.getArgsForCall[i].url
}

func (fake *FakeHTTPHelper) GetReturns(result1 []byte, result2 error) {
	fake.GetStub = nil
	fake.getReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *FakeHTTPHelper) Post(url string, body []byte) ([]byte, error) {
	fake.postMutex.Lock()
	defer fake.postMutex.Unlock()
	fake.postArgsForCall = append(fake.postArgsForCall, struct {
		url  string
		body []byte
	}{url, body})
	if fake.PostStub != nil {
		return fake.PostStub(url, body)
	} else {
		return fake.postReturns.result1, fake.postReturns.result2
	}
}

func (fake *FakeHTTPHelper) PostCallCount() int {
	fake.postMutex.RLock()
	defer fake.postMutex.RUnlock()
	return len(fake.postArgsForCall)
}

func (fake *FakeHTTPHelper) PostArgsForCall(i int) (string, []byte) {
	fake.postMutex.RLock()
	defer fake.postMutex.RUnlock()
	return fake.postArgsForCall[i].url, fake.postArgsForCall[i].body
}

func (fake *FakeHTTPHelper) PostReturns(result1 []byte, result2 error) {
	fake.PostStub = nil
	fake.postReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *FakeHTTPHelper) Put(url string, body []byte) ([]byte, error) {
	fake.putMutex.Lock()
	defer fake.putMutex.Unlock()
	fake.putArgsForCall = append(fake.putArgsForCall, struct {
		url  string
		body []byte
	}{url, body})
	if fake.PutStub != nil {
		return fake.PutStub(url, body)
	} else {
		return fake.putReturns.result1, fake.putReturns.result2
	}
}

func (fake *FakeHTTPHelper) PutCallCount() int {
	fake.putMutex.RLock()
	defer fake.putMutex.RUnlock()
	return len(fake.putArgsForCall)
}

func (fake *FakeHTTPHelper) PutArgsForCall(i int) (string, []byte) {
	fake.putMutex.RLock()
	defer fake.putMutex.RUnlock()
	return fake.putArgsForCall[i].url, fake.putArgsForCall[i].body
}

func (fake *FakeHTTPHelper) PutReturns(result1 []byte, result2 error) {
	fake.PutStub = nil
	fake.putReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *FakeHTTPHelper) Patch(url string, body []byte) ([]byte, error) {
	fake.patchMutex.Lock()
	defer fake.patchMutex.Unlock()
	fake.patchArgsForCall = append(fake.patchArgsForCall, struct {
		url  string
		body []byte
	}{url, body})
	if fake.PatchStub != nil {
		return fake.PatchStub(url, body)
	} else {
		return fake.patchReturns.result1, fake.patchReturns.result2
	}
}

func (fake *FakeHTTPHelper) PatchCallCount() int {
	fake.patchMutex.RLock()
	defer fake.patchMutex.RUnlock()
	return len(fake.patchArgsForCall)
}

func (fake *FakeHTTPHelper) PatchArgsForCall(i int) (string, []byte) {
	fake.patchMutex.RLock()
	defer fake.patchMutex.RUnlock()
	return fake.patchArgsForCall[i].url, fake.patchArgsForCall[i].body
}

func (fake *FakeHTTPHelper) PatchReturns(result1 []byte, result2 error) {
	fake.PatchStub = nil
	fake.patchReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *FakeHTTPHelper) Delete(url string) error {
	fake.deleteMutex.Lock()
	defer fake.deleteMutex.Unlock()
	fake.deleteArgsForCall = append(fake.deleteArgsForCall, struct {
		url string
	}{url})
	if fake.DeleteStub != nil {
		return fake.DeleteStub(url)
	} else {
		return fake.deleteReturns.result1
	}
}

func (fake *FakeHTTPHelper) DeleteCallCount() int {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	return len(fake.deleteArgsForCall)
}

func (fake *FakeHTTPHelper) DeleteArgsForCall(i int) string {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	return fake.deleteArgsForCall[i].url
}

func (fake *FakeHTTPHelper) DeleteReturns(result1 error) {
	fake.DeleteStub = nil
	fake.deleteReturns = struct {
		result1 error
	}{result1}
}

var _ wundergo.HTTPHelper = new(FakeHTTPHelper)
