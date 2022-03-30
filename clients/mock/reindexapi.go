// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"context"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-search-data-finder/clients"
	"github.com/ONSdigital/dp-search-reindex-api/models"
	"github.com/ONSdigital/dp-search-reindex-api/sdk"
	"sync"
)

var (
	lockSearchReindexClientMockChecker        sync.RWMutex
	lockSearchReindexClientMockPostTasksCount sync.RWMutex
)

// Ensure, that SearchReindexClientMock does implement SearchReindexClient.
// If this is not the case, regenerate this file with moq.
var _ clients.SearchReindexClient = &SearchReindexClientMock{}

// SearchReindexClientMock is a mock implementation of clients.SearchReindexClient.
//
//     func TestSomethingThatUsesSearchReindexClient(t *testing.T) {
//
//         // make and configure a mocked clients.SearchReindexClient
//         mockedSearchReindexClient := &SearchReindexClientMock{
//             CheckerFunc: func(in1 context.Context, in2 *healthcheck.CheckState) error {
// 	               panic("mock out the Checker method")
//             },
//             PostTasksCountFunc: func(ctx context.Context, headers sdk.Headers, jobID string) (models.Task, error) {
// 	               panic("mock out the PostTasksCount method")
//             },
//         }
//
//         // use mockedSearchReindexClient in code that requires clients.SearchReindexClient
//         // and then make assertions.
//
//     }
type SearchReindexClientMock struct {
	// CheckerFunc mocks the Checker method.
	CheckerFunc func(in1 context.Context, in2 *healthcheck.CheckState) error

	// PostTasksCountFunc mocks the PostTasksCount method.
	PostTasksCountFunc func(ctx context.Context, headers sdk.Headers, jobID string) (models.Task, error)

	// calls tracks calls to the methods.
	calls struct {
		// Checker holds details about calls to the Checker method.
		Checker []struct {
			// In1 is the in1 argument value.
			In1 context.Context
			// In2 is the in2 argument value.
			In2 *healthcheck.CheckState
		}
		// PostTasksCount holds details about calls to the PostTasksCount method.
		PostTasksCount []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Headers is the headers argument value.
			Headers sdk.Headers
			// JobID is the jobID argument value.
			JobID string
		}
	}
}

// Checker calls CheckerFunc.
func (mock *SearchReindexClientMock) Checker(in1 context.Context, in2 *healthcheck.CheckState) error {
	if mock.CheckerFunc == nil {
		panic("SearchReindexClientMock.CheckerFunc: method is nil but SearchReindexClient.Checker was just called")
	}
	callInfo := struct {
		In1 context.Context
		In2 *healthcheck.CheckState
	}{
		In1: in1,
		In2: in2,
	}
	lockSearchReindexClientMockChecker.Lock()
	mock.calls.Checker = append(mock.calls.Checker, callInfo)
	lockSearchReindexClientMockChecker.Unlock()
	return mock.CheckerFunc(in1, in2)
}

// CheckerCalls gets all the calls that were made to Checker.
// Check the length with:
//     len(mockedSearchReindexClient.CheckerCalls())
func (mock *SearchReindexClientMock) CheckerCalls() []struct {
	In1 context.Context
	In2 *healthcheck.CheckState
} {
	var calls []struct {
		In1 context.Context
		In2 *healthcheck.CheckState
	}
	lockSearchReindexClientMockChecker.RLock()
	calls = mock.calls.Checker
	lockSearchReindexClientMockChecker.RUnlock()
	return calls
}

// PostTasksCount calls PostTasksCountFunc.
func (mock *SearchReindexClientMock) PostTasksCount(ctx context.Context, headers sdk.Headers, jobID string) (models.Task, error) {
	if mock.PostTasksCountFunc == nil {
		panic("SearchReindexClientMock.PostTasksCountFunc: method is nil but SearchReindexClient.PostTasksCount was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		Headers sdk.Headers
		JobID   string
	}{
		Ctx:     ctx,
		Headers: headers,
		JobID:   jobID,
	}
	lockSearchReindexClientMockPostTasksCount.Lock()
	mock.calls.PostTasksCount = append(mock.calls.PostTasksCount, callInfo)
	lockSearchReindexClientMockPostTasksCount.Unlock()
	return mock.PostTasksCountFunc(ctx, headers, jobID)
}

// PostTasksCountCalls gets all the calls that were made to PostTasksCount.
// Check the length with:
//     len(mockedSearchReindexClient.PostTasksCountCalls())
func (mock *SearchReindexClientMock) PostTasksCountCalls() []struct {
	Ctx     context.Context
	Headers sdk.Headers
	JobID   string
} {
	var calls []struct {
		Ctx     context.Context
		Headers sdk.Headers
		JobID   string
	}
	lockSearchReindexClientMockPostTasksCount.RLock()
	calls = mock.calls.PostTasksCount
	lockSearchReindexClientMockPostTasksCount.RUnlock()
	return calls
}
