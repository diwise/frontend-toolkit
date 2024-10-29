// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"github.com/diwise/frontend-toolkit"
	"sync"
)

// BundleMock is a mock implementation of frontendtoolkit.Bundle.
//
//	func TestSomethingThatUsesBundle(t *testing.T) {
//
//		// make and configure a mocked frontendtoolkit.Bundle
//		mockedBundle := &BundleMock{
//			ForFunc: func(acceptLanguage string) frontendtoolkit.Localizer {
//				panic("mock out the For method")
//			},
//		}
//
//		// use mockedBundle in code that requires frontendtoolkit.Bundle
//		// and then make assertions.
//
//	}
type BundleMock struct {
	// ForFunc mocks the For method.
	ForFunc func(acceptLanguage string) frontendtoolkit.Localizer

	// calls tracks calls to the methods.
	calls struct {
		// For holds details about calls to the For method.
		For []struct {
			// AcceptLanguage is the acceptLanguage argument value.
			AcceptLanguage string
		}
	}
	lockFor sync.RWMutex
}

// For calls ForFunc.
func (mock *BundleMock) For(acceptLanguage string) frontendtoolkit.Localizer {
	callInfo := struct {
		AcceptLanguage string
	}{
		AcceptLanguage: acceptLanguage,
	}
	mock.lockFor.Lock()
	mock.calls.For = append(mock.calls.For, callInfo)
	mock.lockFor.Unlock()
	if mock.ForFunc == nil {
		var (
			localizerOut frontendtoolkit.Localizer
		)
		return localizerOut
	}
	return mock.ForFunc(acceptLanguage)
}

// ForCalls gets all the calls that were made to For.
// Check the length with:
//
//	len(mockedBundle.ForCalls())
func (mock *BundleMock) ForCalls() []struct {
	AcceptLanguage string
} {
	var calls []struct {
		AcceptLanguage string
	}
	mock.lockFor.RLock()
	calls = mock.calls.For
	mock.lockFor.RUnlock()
	return calls
}
