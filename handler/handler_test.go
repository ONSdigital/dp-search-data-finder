package handler_test

import (
	"context"
	"testing"

	datasetClient "github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	zebedeeClient "github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	clientMock "github.com/ONSdigital/dp-search-data-finder/clients/mock"
	"github.com/ONSdigital/dp-search-data-finder/config"
	"github.com/ONSdigital/dp-search-data-finder/handler"
	"github.com/ONSdigital/dp-search-data-finder/models"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	testCtx = context.Background()

	testEvent = models.ReindexRequested{
		JobID:       "job id",
		SearchIndex: "search index",
		TraceID:     "oe433dpe446657gge",
	}

	errZebedee             = errors.New("zebedee test error")
	getPublishedIndexEmpty = func(ctx context.Context, publishedIndexRequestParams *zebedeeClient.PublishedIndexRequestParams) (zebedeeClient.PublishedIndex, error) {
		return zebedeeClient.PublishedIndex{}, errZebedee
	}
	getPublishedIndexFunc = func(ctx context.Context, publishedIndexRequestParams *zebedeeClient.PublishedIndexRequestParams) (zebedeeClient.PublishedIndex, error) {
		return zebedeeClient.PublishedIndex{}, nil
	}

	errDatasetAPI = errors.New("dataset-api test error")
	getDatasetsOk = func(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, q *datasetClient.QueryParams) (datasetClient.List, error) {
		return datasetClient.List{}, nil
	}
	getFullEditionsDetailsOk = func(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, datasetID string) ([]datasetClient.EditionsDetails, error) {
		return []datasetClient.EditionsDetails{}, nil
	}
	getVersionMetadataOk = func(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, id string, edition string, version string) (datasetClient.Metadata, error) {
		return datasetClient.Metadata{}, nil
	}
	getDatasetsError = func(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, q *datasetClient.QueryParams) (datasetClient.List, error) {
		return datasetClient.List{}, errDatasetAPI
	}
)

func TestReindexRequestedHandler_Handle(t *testing.T) {
	t.Parallel()

	testCfg, err := config.Get()
	if err != nil {
		t.Errorf("failed to retrieve default configuration, error is: %v", err)
	}

	Convey("Given an event handler working successfully, and an event containing a URI", t, func() {
		var zebedeeMock = &clientMock.ZebedeeClientMock{GetPublishedIndexFunc: getPublishedIndexFunc}
		var datasetAPIMock = &clientMock.DatasetAPIClientMock{
			GetDatasetsFunc:            getDatasetsOk,
			GetFullEditionsDetailsFunc: getFullEditionsDetailsOk,
			GetVersionMetadataFunc:     getVersionMetadataOk,
		}
		eventHandler := &handler.ReindexRequestedHandler{Config: testCfg, ZebedeeCli: zebedeeMock, DatasetAPICli: datasetAPIMock}

		Convey("When given a valid event", func() {
			eventHandler.Handle(testCtx, &testEvent)

			Convey("Then Zebedee and Dataset API are called to get document urls", func() {
				So(zebedeeMock.GetPublishedIndexCalls(), ShouldNotBeEmpty)
				So(zebedeeMock.GetPublishedIndexCalls(), ShouldHaveLength, 1)
				So(datasetAPIMock.GetDatasetsCalls(), ShouldNotBeEmpty)
				So(datasetAPIMock.GetDatasetsCalls(), ShouldHaveLength, 1)
			})
		})
	})

	Convey("Given an event handler not working successfully with Zebedee, and an event containing a jobId", t, func() {
		var zebedeeMockInError = &clientMock.ZebedeeClientMock{GetPublishedIndexFunc: getPublishedIndexEmpty}
		var datasetAPIMock = &clientMock.DatasetAPIClientMock{
			GetDatasetsFunc:            getDatasetsOk,
			GetFullEditionsDetailsFunc: getFullEditionsDetailsOk,
			GetVersionMetadataFunc:     getVersionMetadataOk,
		}
		eventHandler := &handler.ReindexRequestedHandler{Config: testCfg, ZebedeeCli: zebedeeMockInError, DatasetAPICli: datasetAPIMock}

		Convey("When given a valid event", func() {
			eventHandler.Handle(testCtx, &testEvent)

			Convey("Then Zebedee is called 1 time with the expected error which is logged", func() {
				So(zebedeeMockInError.GetPublishedIndexCalls(), ShouldNotBeEmpty)
				So(zebedeeMockInError.GetPublishedIndexCalls(), ShouldHaveLength, 1)
				So(datasetAPIMock.GetDatasetsCalls(), ShouldNotBeEmpty)
				So(datasetAPIMock.GetDatasetsCalls(), ShouldHaveLength, 1)
			})
		})
	})

	Convey("Given an event handler not working successfully with Dataset API, and an event containing a jobId", t, func() {
		var zebedeeMock = &clientMock.ZebedeeClientMock{GetPublishedIndexFunc: getPublishedIndexFunc}
		var datasetAPIMockErr = &clientMock.DatasetAPIClientMock{GetDatasetsFunc: getDatasetsError}
		eventHandler := &handler.ReindexRequestedHandler{Config: testCfg, ZebedeeCli: zebedeeMock, DatasetAPICli: datasetAPIMockErr}

		Convey("When given a valid event", func() {
			eventHandler.Handle(testCtx, &testEvent)

			Convey("Then Dataset API is called 1 time with the expected error which is logged", func() {
				So(zebedeeMock.GetPublishedIndexCalls(), ShouldNotBeEmpty)
				So(zebedeeMock.GetPublishedIndexCalls(), ShouldHaveLength, 1)
				So(datasetAPIMockErr.GetDatasetsCalls(), ShouldNotBeEmpty)
				So(datasetAPIMockErr.GetDatasetsCalls(), ShouldHaveLength, 1)
			})
		})
	})
}
