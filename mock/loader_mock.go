// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"github.com/diwise/frontend-toolkit"
	"sync"
)

// LoaderMock is a mock implementation of frontendtoolkit.Loader.
//
//	func TestSomethingThatUsesLoader(t *testing.T) {
//
//		// make and configure a mocked frontendtoolkit.Loader
//		mockedLoader := &LoaderMock{
//			LoadFunc: func(name string) frontendtoolkit.Asset {
//				panic("mock out the Load method")
//			},
//			LoadFromSha256Func: func(sha string) (frontendtoolkit.Asset, error) {
//				panic("mock out the LoadFromSha256 method")
//			},
//		}
//
//		// use mockedLoader in code that requires frontendtoolkit.Loader
//		// and then make assertions.
//
//	}
type LoaderMock struct {
	// LoadFunc mocks the Load method.
	LoadFunc func(name string) frontendtoolkit.Asset

	// LoadFromSha256Func mocks the LoadFromSha256 method.
	LoadFromSha256Func func(sha string) (frontendtoolkit.Asset, error)

	// calls tracks calls to the methods.
	calls struct {
		// Load holds details about calls to the Load method.
		Load []struct {
			// Name is the name argument value.
			Name string
		}
		// LoadFromSha256 holds details about calls to the LoadFromSha256 method.
		LoadFromSha256 []struct {
			// Sha is the sha argument value.
			Sha string
		}
	}
	lockLoad           sync.RWMutex
	lockLoadFromSha256 sync.RWMutex
}

// Load calls LoadFunc.
func (mock *LoaderMock) Load(name string) frontendtoolkit.Asset {
	callInfo := struct {
		Name string
	}{
		Name: name,
	}
	mock.lockLoad.Lock()
	mock.calls.Load = append(mock.calls.Load, callInfo)
	mock.lockLoad.Unlock()
	if mock.LoadFunc == nil {
		var (
			assetOut frontendtoolkit.Asset
		)
		return assetOut
	}
	return mock.LoadFunc(name)
}

// LoadCalls gets all the calls that were made to Load.
// Check the length with:
//
//	len(mockedLoader.LoadCalls())
func (mock *LoaderMock) LoadCalls() []struct {
	Name string
} {
	var calls []struct {
		Name string
	}
	mock.lockLoad.RLock()
	calls = mock.calls.Load
	mock.lockLoad.RUnlock()
	return calls
}

// LoadFromSha256 calls LoadFromSha256Func.
func (mock *LoaderMock) LoadFromSha256(sha string) (frontendtoolkit.Asset, error) {
	callInfo := struct {
		Sha string
	}{
		Sha: sha,
	}
	mock.lockLoadFromSha256.Lock()
	mock.calls.LoadFromSha256 = append(mock.calls.LoadFromSha256, callInfo)
	mock.lockLoadFromSha256.Unlock()
	if mock.LoadFromSha256Func == nil {
		var (
			assetOut frontendtoolkit.Asset
			errOut   error
		)
		return assetOut, errOut
	}
	return mock.LoadFromSha256Func(sha)
}

// LoadFromSha256Calls gets all the calls that were made to LoadFromSha256.
// Check the length with:
//
//	len(mockedLoader.LoadFromSha256Calls())
func (mock *LoaderMock) LoadFromSha256Calls() []struct {
	Sha string
} {
	var calls []struct {
		Sha string
	}
	mock.lockLoadFromSha256.RLock()
	calls = mock.calls.LoadFromSha256
	mock.lockLoadFromSha256.RUnlock()
	return calls
}
