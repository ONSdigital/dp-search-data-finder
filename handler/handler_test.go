package handler_test

import (
	"context"
	"testing"

	datasetClient "github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	zebedeeClient "github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	clientMock "github.com/ONSdigital/dp-search-data-finder/clients/mock"

	"github.com/ONSdigital/dp-search-data-finder/config"
	"github.com/ONSdigital/dp-search-data-finder/event/mock"
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

	mockZebedeePublishedIndexResponse = zebedeeClient.PublishedIndex{
		Count: 1,
		Limit: 0,
		Items: []zebedeeClient.PublishedIndexItem{
			{URI: "http://www.ons.gov.uk"},
		},
		Offset:     0,
		TotalCount: 1,
	}

	getPublishedIndexFunc = func(ctx context.Context, publishedIndexRequestParams *zebedeeClient.PublishedIndexRequestParams) (zebedeeClient.PublishedIndex, error) {
		return mockZebedeePublishedIndexResponse, nil
	}

	mockDatasetAPIResponse = datasetClient.List{
		Count: 1,
		Limit: 0,
		Items: []datasetClient.Dataset{
			{
				ID: "RM007",
				Current: &datasetClient.DatasetDetails{
					ID: "RM007",
				}},
		},
		Offset:     0,
		TotalCount: 1,
	}

	errDatasetAPI = errors.New("dataset-api test error")
	getDatasetsOk = func(_ context.Context, _ string, _ string, _ string, q *datasetClient.QueryParams) (datasetClient.List, error) {
		if q.Offset == 0 {
			return mockDatasetAPIResponse, nil
		}
		return datasetClient.List{}, nil
	}

	mockDatasetEditionAPIResponse = []datasetClient.EditionsDetails{
		{
			ID: "RM007",
			Current: datasetClient.Edition{
				ID: "1",
				Links: datasetClient.Links{
					LatestVersion: datasetClient.Link{
						ID: "1",
					},
				},
			},
		},
	}
	getFullEditionsDetailsOk = func(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, datasetID string) ([]datasetClient.EditionsDetails, error) {
		return mockDatasetEditionAPIResponse, nil
	}
	mockDatasetVersionAPIResponse = datasetClient.Metadata{
		DatasetLinks: datasetClient.Links{
			LatestVersion: datasetClient.Link{
				URL: "http://www.ons.gov.uk/datasets/RM007",
			},
		},
	}
	getVersionMetadataOk = func(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, id string, edition string, version string) (datasetClient.Metadata, error) {
		return mockDatasetVersionAPIResponse, nil
	}
	getDatasetsError = func(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, q *datasetClient.QueryParams) (datasetClient.List, error) {
		return datasetClient.List{}, errDatasetAPI
	}
)

func TestHandle(t *testing.T) {
	testCfg, err := config.Get()
	if err != nil {
		t.Errorf("failed to retrieve default configuration, error is: %v", err)
	}

	testCfg.EnablePublishContentUpdatedTopic = true

	expectedLegacyContentUpdatedEvent := models.ContentUpdated{
		URI:         "http://www.ons.gov.uk",
		JobID:       "",
		TraceID:     "",
		SearchIndex: "",
		DataType:    "legacy",
	}

	expectedDatasetContentUpdatedEvent := &models.ContentUpdated{
		URI:         "http://www.ons.gov.uk/datasets/RM007",
		JobID:       "",
		TraceID:     "",
		SearchIndex: "",
		DataType:    "datasets",
	}

	Convey("Given an event handler, working zebedee and dataset clients, kafka producer and an event containing a URI", t, func() {
		zebedeeMock := &clientMock.ZebedeeClientMock{GetPublishedIndexFunc: getPublishedIndexFunc}
		datasetAPIMock := &clientMock.DatasetAPIClientMock{
			GetDatasetsFunc:            getDatasetsOk,
			GetFullEditionsDetailsFunc: getFullEditionsDetailsOk,
			GetVersionMetadataFunc:     getVersionMetadataOk,
		}

		contentUpdaterMock := &mock.ContentUpdaterMock{
			ContentUpdateFunc: func(ctx context.Context, cfg *config.Config, event models.ContentUpdated) error {
				return nil
			},
		}

		eventHandler := &handler.ReindexRequestedHandler{Config: testCfg, ZebedeeCli: zebedeeMock, DatasetAPICli: datasetAPIMock, ContentUpdatedProducer: contentUpdaterMock}

		Convey("When given a valid event", func() {
			eventHandler.Handle(testCtx, &testEvent)

			Convey("Then Zebedee and Dataset API are called to get document urls", func() {
				So(zebedeeMock.GetPublishedIndexCalls(), ShouldNotBeEmpty)
				So(zebedeeMock.GetPublishedIndexCalls(), ShouldHaveLength, 1)
				So(datasetAPIMock.GetDatasetsCalls(), ShouldNotBeEmpty)
				So(datasetAPIMock.GetDatasetsCalls(), ShouldHaveLength, 2)
			})

			Convey("Then a legacy event should be called to be produced", func() {
				So(contentUpdaterMock.ContentUpdateCalls(), ShouldHaveLength, 2)
				So(contentUpdaterMock.ContentUpdateCalls()[0].Event.URI, ShouldEqual, expectedLegacyContentUpdatedEvent.URI)
				So(contentUpdaterMock.ContentUpdateCalls()[0].Event.DataType, ShouldEqual, expectedLegacyContentUpdatedEvent.DataType)
			})

			Convey("Then a datasets event should be called to be produced", func() {
				So(contentUpdaterMock.ContentUpdateCalls(), ShouldHaveLength, 2)
				So(contentUpdaterMock.ContentUpdateCalls()[1].Event.URI, ShouldEqual, expectedDatasetContentUpdatedEvent.URI)
				So(contentUpdaterMock.ContentUpdateCalls()[1].Event.DataType, ShouldEqual, expectedDatasetContentUpdatedEvent.DataType)
			})
		})
	})
}
