// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"context"
	datasetclient "github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-search-data-finder/clients"
	"sync"
)

// Ensure, that DatasetAPIClientMock does implement clients.DatasetAPIClient.
// If this is not the case, regenerate this file with moq.
var _ clients.DatasetAPIClient = &DatasetAPIClientMock{}

// DatasetAPIClientMock is a mock implementation of clients.DatasetAPIClient.
//
// 	func TestSomethingThatUsesDatasetAPIClient(t *testing.T) {
//
// 		// make and configure a mocked clients.DatasetAPIClient
// 		mockedDatasetAPIClient := &DatasetAPIClientMock{
// 			CheckerFunc: func(contextMoqParam context.Context, checkState *healthcheck.CheckState) error {
// 				panic("mock out the Checker method")
// 			},
// 			GetDatasetsFunc: func(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, q *datasetclient.QueryParams) (datasetclient.List, error) {
// 				panic("mock out the GetDatasets method")
// 			},
// 			GetFullEditionsDetailsFunc: func(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, datasetID string) ([]datasetclient.EditionsDetails, error) {
// 				panic("mock out the GetFullEditionsDetails method")
// 			},
// 			GetVersionMetadataFunc: func(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, id string, edition string, version string) (datasetclient.Metadata, error) {
// 				panic("mock out the GetVersionMetadata method")
// 			},
// 		}
//
// 		// use mockedDatasetAPIClient in code that requires clients.DatasetAPIClient
// 		// and then make assertions.
//
// 	}
type DatasetAPIClientMock struct {
	// CheckerFunc mocks the Checker method.
	CheckerFunc func(contextMoqParam context.Context, checkState *healthcheck.CheckState) error

	// GetDatasetsFunc mocks the GetDatasets method.
	GetDatasetsFunc func(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, q *datasetclient.QueryParams) (datasetclient.List, error)

	// GetFullEditionsDetailsFunc mocks the GetFullEditionsDetails method.
	GetFullEditionsDetailsFunc func(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, datasetID string) ([]datasetclient.EditionsDetails, error)

	// GetVersionMetadataFunc mocks the GetVersionMetadata method.
	GetVersionMetadataFunc func(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, id string, edition string, version string) (datasetclient.Metadata, error)

	// calls tracks calls to the methods.
	calls struct {
		// Checker holds details about calls to the Checker method.
		Checker []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// CheckState is the checkState argument value.
			CheckState *healthcheck.CheckState
		}
		// GetDatasets holds details about calls to the GetDatasets method.
		GetDatasets []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// UserAuthToken is the userAuthToken argument value.
			UserAuthToken string
			// ServiceAuthToken is the serviceAuthToken argument value.
			ServiceAuthToken string
			// CollectionID is the collectionID argument value.
			CollectionID string
			// Q is the q argument value.
			Q *datasetclient.QueryParams
		}
		// GetFullEditionsDetails holds details about calls to the GetFullEditionsDetails method.
		GetFullEditionsDetails []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// UserAuthToken is the userAuthToken argument value.
			UserAuthToken string
			// ServiceAuthToken is the serviceAuthToken argument value.
			ServiceAuthToken string
			// CollectionID is the collectionID argument value.
			CollectionID string
			// DatasetID is the datasetID argument value.
			DatasetID string
		}
		// GetVersionMetadata holds details about calls to the GetVersionMetadata method.
		GetVersionMetadata []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// UserAuthToken is the userAuthToken argument value.
			UserAuthToken string
			// ServiceAuthToken is the serviceAuthToken argument value.
			ServiceAuthToken string
			// CollectionID is the collectionID argument value.
			CollectionID string
			// ID is the id argument value.
			ID string
			// Edition is the edition argument value.
			Edition string
			// Version is the version argument value.
			Version string
		}
	}
	lockChecker                sync.RWMutex
	lockGetDatasets            sync.RWMutex
	lockGetFullEditionsDetails sync.RWMutex
	lockGetVersionMetadata     sync.RWMutex
}

// Checker calls CheckerFunc.
func (mock *DatasetAPIClientMock) Checker(contextMoqParam context.Context, checkState *healthcheck.CheckState) error {
	if mock.CheckerFunc == nil {
		panic("DatasetAPIClientMock.CheckerFunc: method is nil but DatasetAPIClient.Checker was just called")
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
//     len(mockedDatasetAPIClient.CheckerCalls())
func (mock *DatasetAPIClientMock) CheckerCalls() []struct {
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

// GetDatasets calls GetDatasetsFunc.
func (mock *DatasetAPIClientMock) GetDatasets(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, q *datasetclient.QueryParams) (datasetclient.List, error) {
	if mock.GetDatasetsFunc == nil {
		panic("DatasetAPIClientMock.GetDatasetsFunc: method is nil but DatasetAPIClient.GetDatasets was just called")
	}
	callInfo := struct {
		Ctx              context.Context
		UserAuthToken    string
		ServiceAuthToken string
		CollectionID     string
		Q                *datasetclient.QueryParams
	}{
		Ctx:              ctx,
		UserAuthToken:    userAuthToken,
		ServiceAuthToken: serviceAuthToken,
		CollectionID:     collectionID,
		Q:                q,
	}
	mock.lockGetDatasets.Lock()
	mock.calls.GetDatasets = append(mock.calls.GetDatasets, callInfo)
	mock.lockGetDatasets.Unlock()
	return mock.GetDatasetsFunc(ctx, userAuthToken, serviceAuthToken, collectionID, q)
}

// GetDatasetsCalls gets all the calls that were made to GetDatasets.
// Check the length with:
//     len(mockedDatasetAPIClient.GetDatasetsCalls())
func (mock *DatasetAPIClientMock) GetDatasetsCalls() []struct {
	Ctx              context.Context
	UserAuthToken    string
	ServiceAuthToken string
	CollectionID     string
	Q                *datasetclient.QueryParams
} {
	var calls []struct {
		Ctx              context.Context
		UserAuthToken    string
		ServiceAuthToken string
		CollectionID     string
		Q                *datasetclient.QueryParams
	}
	mock.lockGetDatasets.RLock()
	calls = mock.calls.GetDatasets
	mock.lockGetDatasets.RUnlock()
	return calls
}

// GetFullEditionsDetails calls GetFullEditionsDetailsFunc.
func (mock *DatasetAPIClientMock) GetFullEditionsDetails(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, datasetID string) ([]datasetclient.EditionsDetails, error) {
	if mock.GetFullEditionsDetailsFunc == nil {
		panic("DatasetAPIClientMock.GetFullEditionsDetailsFunc: method is nil but DatasetAPIClient.GetFullEditionsDetails was just called")
	}
	callInfo := struct {
		Ctx              context.Context
		UserAuthToken    string
		ServiceAuthToken string
		CollectionID     string
		DatasetID        string
	}{
		Ctx:              ctx,
		UserAuthToken:    userAuthToken,
		ServiceAuthToken: serviceAuthToken,
		CollectionID:     collectionID,
		DatasetID:        datasetID,
	}
	mock.lockGetFullEditionsDetails.Lock()
	mock.calls.GetFullEditionsDetails = append(mock.calls.GetFullEditionsDetails, callInfo)
	mock.lockGetFullEditionsDetails.Unlock()
	return mock.GetFullEditionsDetailsFunc(ctx, userAuthToken, serviceAuthToken, collectionID, datasetID)
}

// GetFullEditionsDetailsCalls gets all the calls that were made to GetFullEditionsDetails.
// Check the length with:
//     len(mockedDatasetAPIClient.GetFullEditionsDetailsCalls())
func (mock *DatasetAPIClientMock) GetFullEditionsDetailsCalls() []struct {
	Ctx              context.Context
	UserAuthToken    string
	ServiceAuthToken string
	CollectionID     string
	DatasetID        string
} {
	var calls []struct {
		Ctx              context.Context
		UserAuthToken    string
		ServiceAuthToken string
		CollectionID     string
		DatasetID        string
	}
	mock.lockGetFullEditionsDetails.RLock()
	calls = mock.calls.GetFullEditionsDetails
	mock.lockGetFullEditionsDetails.RUnlock()
	return calls
}

// GetVersionMetadata calls GetVersionMetadataFunc.
func (mock *DatasetAPIClientMock) GetVersionMetadata(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, id string, edition string, version string) (datasetclient.Metadata, error) {
	if mock.GetVersionMetadataFunc == nil {
		panic("DatasetAPIClientMock.GetVersionMetadataFunc: method is nil but DatasetAPIClient.GetVersionMetadata was just called")
	}
	callInfo := struct {
		Ctx              context.Context
		UserAuthToken    string
		ServiceAuthToken string
		CollectionID     string
		ID               string
		Edition          string
		Version          string
	}{
		Ctx:              ctx,
		UserAuthToken:    userAuthToken,
		ServiceAuthToken: serviceAuthToken,
		CollectionID:     collectionID,
		ID:               id,
		Edition:          edition,
		Version:          version,
	}
	mock.lockGetVersionMetadata.Lock()
	mock.calls.GetVersionMetadata = append(mock.calls.GetVersionMetadata, callInfo)
	mock.lockGetVersionMetadata.Unlock()
	return mock.GetVersionMetadataFunc(ctx, userAuthToken, serviceAuthToken, collectionID, id, edition, version)
}

// GetVersionMetadataCalls gets all the calls that were made to GetVersionMetadata.
// Check the length with:
//     len(mockedDatasetAPIClient.GetVersionMetadataCalls())
func (mock *DatasetAPIClientMock) GetVersionMetadataCalls() []struct {
	Ctx              context.Context
	UserAuthToken    string
	ServiceAuthToken string
	CollectionID     string
	ID               string
	Edition          string
	Version          string
} {
	var calls []struct {
		Ctx              context.Context
		UserAuthToken    string
		ServiceAuthToken string
		CollectionID     string
		ID               string
		Edition          string
		Version          string
	}
	mock.lockGetVersionMetadata.RLock()
	calls = mock.calls.GetVersionMetadata
	mock.lockGetVersionMetadata.RUnlock()
	return calls
}