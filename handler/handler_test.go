package handler_test

import (
	"context"
	"fmt"
	"strconv"
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
			{
				URI: "http://www.ons.gov.uk",
			},
		},
		Offset:     0,
		TotalCount: 1,
	}

	getPublishedIndexFunc = func(ctx context.Context, publishedIndexRequestParams *zebedeeClient.PublishedIndexRequestParams) (zebedeeClient.PublishedIndex, error) {
		return mockZebedeePublishedIndexResponse, nil
	}

	generateFakeDatasetID = func(num int) string {
		return fmt.Sprintf("RM%s", strconv.Itoa(num))
	}

	generateMockDatasetAPIResponse = func(numberOfDatasets int) datasetClient.List {
		items := make([]datasetClient.Dataset, numberOfDatasets)

		for i := 0; i < numberOfDatasets; i++ {
			newID := generateFakeDatasetID(i + 1)
			newItem := datasetClient.Dataset{
				ID: newID,
				Current: &datasetClient.DatasetDetails{
					ID: newID,
				},
			}
			items = append(items, newItem)
		}

		mockDatasetAPIResponse := datasetClient.List{
			Count: numberOfDatasets,
			Limit: 0,
			Items: items,
		}

		return mockDatasetAPIResponse
	}

	errDatasetAPI = errors.New("dataset-api test error")
	getDatasetsOk = func(_ context.Context, _ string, _ string, _ string, q *datasetClient.QueryParams) (datasetClient.List, error) {
		if q.Offset == 0 {
			return generateMockDatasetAPIResponse(2), nil
		}
		return datasetClient.List{}, nil
	}

	generateMockDatasetEditionAPIResponse = func(id string) []datasetClient.EditionsDetails {
		return []datasetClient.EditionsDetails{
			{
				ID: id,
				Current: datasetClient.Edition{
					ID: id,
					Links: datasetClient.Links{
						LatestVersion: datasetClient.Link{
							ID: id,
						},
					},
				},
			},
		}
	}

	getFullEditionsDetailsOk = func(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, datasetID string) ([]datasetClient.EditionsDetails, error) {
		return generateMockDatasetEditionAPIResponse(datasetID), nil
	}

	generateMockDatasetVersionAPIResponse = func(id string) datasetClient.Metadata {
		return datasetClient.Metadata{
			DatasetLinks: datasetClient.Links{
				LatestVersion: datasetClient.Link{
					URL: fmt.Sprintf("http://www.ons.gov.uk/datasets/%s", id),
				},
			},
		}
	}

	getVersionMetadataOk = func(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, id string, edition string, version string) (datasetClient.Metadata, error) {
		return generateMockDatasetVersionAPIResponse(id), nil
	}
	getDatasetsError = func(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, q *datasetClient.QueryParams) (datasetClient.List, error) {
		return datasetClient.List{}, errDatasetAPI
	}

	generateExpectedDatasetContentUpdatedEvent = func(id int) models.ContentUpdated {
		datasetID := generateFakeDatasetID(id)

		return models.ContentUpdated{
			URI:         fmt.Sprintf("http://www.ons.gov.uk/datasets/%s", datasetID),
			JobID:       testEvent.JobID,
			TraceID:     testEvent.TraceID,
			SearchIndex: testEvent.SearchIndex,
			DataType:    "datasets",
		}
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
		JobID:       testEvent.JobID,
		TraceID:     testEvent.TraceID,
		SearchIndex: testEvent.SearchIndex,
		DataType:    "legacy",
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

			// Due to parallel calls, the events can end up in any order
			events := make([]models.ContentUpdated, len(contentUpdaterMock.ContentUpdateCalls()))

			for i, call := range contentUpdaterMock.ContentUpdateCalls() {
				events[i] = call.Event
			}

			Convey("Then a legacy event should be called to be produced", func() {
				So(events, ShouldHaveLength, 3)
				So(events, ShouldContain, expectedLegacyContentUpdatedEvent)
			})

			Convey("Then a datasets event should be called to be produced", func() {
				So(events, ShouldHaveLength, 3)
				So(events, ShouldContain, generateExpectedDatasetContentUpdatedEvent(1))
				So(events, ShouldContain, generateExpectedDatasetContentUpdatedEvent(2))
			})
		})
	})
}
