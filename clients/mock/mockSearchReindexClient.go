// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"context"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-search-data-finder/clients"
	searchReindex "github.com/ONSdigital/dp-search-reindex-api/models"
	searchReindexSDK "github.com/ONSdigital/dp-search-reindex-api/sdk"
	"sync"
)

// Ensure, that SearchReindexClientMock does implement clients.SearchReindexClient.
// If this is not the case, regenerate this file with moq.
var _ clients.SearchReindexClient = &SearchReindexClientMock{}

// SearchReindexClientMock is a mock implementation of clients.SearchReindexClient.
//
// 	func TestSomethingThatUsesSearchReindexClient(t *testing.T) {
//
// 		// make and configure a mocked clients.SearchReindexClient
// 		mockedSearchReindexClient := &SearchReindexClientMock{
// 			CheckerFunc: func(contextMoqParam context.Context, checkState *healthcheck.CheckState) error {
// 				panic("mock out the Checker method")
// 			},
// 			PatchJobFunc: func(contextMoqParam context.Context, headers searchReindexSDK.Headers, s string, patchOpsList searchReindexSDK.PatchOpsList) error {
// 				panic("mock out the PatchJob method")
// 			},
// 			PostJobFunc: func(contextMoqParam context.Context, headers searchReindexSDK.Headers) (searchReindex.Job, error) {
// 				panic("mock out the PostJob method")
// 			},
// 		}
//
// 		// use mockedSearchReindexClient in code that requires clients.SearchReindexClient
// 		// and then make assertions.
//
// 	}
type SearchReindexClientMock struct {
	// CheckerFunc mocks the Checker method.
	CheckerFunc func(contextMoqParam context.Context, checkState *healthcheck.CheckState) error

	// PatchJobFunc mocks the PatchJob method.
	PatchJobFunc func(contextMoqParam context.Context, headers searchReindexSDK.Headers, s string, patchOpsList searchReindexSDK.PatchOpsList) error

	// PostJobFunc mocks the PostJob method.
	PostJobFunc func(contextMoqParam context.Context, headers searchReindexSDK.Headers) (searchReindex.Job, error)

	// calls tracks calls to the methods.
	calls struct {
		// Checker holds details about calls to the Checker method.
		Checker []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// CheckState is the checkState argument value.
			CheckState *healthcheck.CheckState
		}
		// PatchJob holds details about calls to the PatchJob method.
		PatchJob []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// Headers is the headers argument value.
			Headers searchReindexSDK.Headers
			// S is the s argument value.
			S string
			// PatchOpsList is the patchOpsList argument value.
			PatchOpsList searchReindexSDK.PatchOpsList
		}
		// PostJob holds details about calls to the PostJob method.
		PostJob []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// Headers is the headers argument value.
			Headers searchReindexSDK.Headers
		}
	}
	lockChecker  sync.RWMutex
	lockPatchJob sync.RWMutex
	lockPostJob  sync.RWMutex
}

// Checker calls CheckerFunc.
func (mock *SearchReindexClientMock) Checker(contextMoqParam context.Context, checkState *healthcheck.CheckState) error {
	if mock.CheckerFunc == nil {
		panic("SearchReindexClientMock.CheckerFunc: method is nil but SearchReindexClient.Checker was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
		CheckState      *healthcheck.CheckState
	}{
		ContextMoqParam: contextMoqParam,
		CheckState:      checkState,
	}
	mock.lockChecker.Lock()
	mock.calls.Checker = append(mock.calls.Checker, callInfo)
	mock.lockChecker.Unlock()
	return mock.CheckerFunc(contextMoqParam, checkState)
}

// CheckerCalls gets all the calls that were made to Checker.
// Check the length with:
//     len(mockedSearchReindexClient.CheckerCalls())
func (mock *SearchReindexClientMock) CheckerCalls() []struct {
	ContextMoqParam context.Context
	CheckState      *healthcheck.CheckState
} {
	var calls []struct {
		ContextMoqParam context.Context
		CheckState      *healthcheck.CheckState
	}
	mock.lockChecker.RLock()
	calls = mock.calls.Checker
	mock.lockChecker.RUnlock()
	return calls
}

// PatchJob calls PatchJobFunc.
func (mock *SearchReindexClientMock) PatchJob(contextMoqParam context.Context, headers searchReindexSDK.Headers, s string, patchOpsList searchReindexSDK.PatchOpsList) error {
	if mock.PatchJobFunc == nil {
		panic("SearchReindexClientMock.PatchJobFunc: method is nil but SearchReindexClient.PatchJob was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
		Headers         searchReindexSDK.Headers
		S               string
		PatchOpsList    searchReindexSDK.PatchOpsList
	}{
		ContextMoqParam: contextMoqParam,
		Headers:         headers,
		S:               s,
		PatchOpsList:    patchOpsList,
	}
	mock.lockPatchJob.Lock()
	mock.calls.PatchJob = append(mock.calls.PatchJob, callInfo)
	mock.lockPatchJob.Unlock()
	return mock.PatchJobFunc(contextMoqParam, headers, s, patchOpsList)
}

// PatchJobCalls gets all the calls that were made to PatchJob.
// Check the length with:
//     len(mockedSearchReindexClient.PatchJobCalls())
func (mock *SearchReindexClientMock) PatchJobCalls() []struct {
	ContextMoqParam context.Context
	Headers         searchReindexSDK.Headers
	S               string
	PatchOpsList    searchReindexSDK.PatchOpsList
} {
	var calls []struct {
		ContextMoqParam context.Context
		Headers         searchReindexSDK.Headers
		S               string
		PatchOpsList    searchReindexSDK.PatchOpsList
	}
	mock.lockPatchJob.RLock()
	calls = mock.calls.PatchJob
	mock.lockPatchJob.RUnlock()
	return calls
}

// PostJob calls PostJobFunc.
func (mock *SearchReindexClientMock) PostJob(contextMoqParam context.Context, headers searchReindexSDK.Headers) (searchReindex.Job, error) {
	if mock.PostJobFunc == nil {
		panic("SearchReindexClientMock.PostJobFunc: method is nil but SearchReindexClient.PostJob was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
		Headers         searchReindexSDK.Headers
	}{
		ContextMoqParam: contextMoqParam,
		Headers:         headers,
	}
	mock.lockPostJob.Lock()
	mock.calls.PostJob = append(mock.calls.PostJob, callInfo)
	mock.lockPostJob.Unlock()
	return mock.PostJobFunc(contextMoqParam, headers)
}

// PostJobCalls gets all the calls that were made to PostJob.
// Check the length with:
//     len(mockedSearchReindexClient.PostJobCalls())
func (mock *SearchReindexClientMock) PostJobCalls() []struct {
	ContextMoqParam context.Context
	Headers         searchReindexSDK.Headers
} {
	var calls []struct {
		ContextMoqParam context.Context
		Headers         searchReindexSDK.Headers
	}
	mock.lockPostJob.RLock()
	calls = mock.calls.PostJob
	mock.lockPostJob.RUnlock()
	return calls
}
