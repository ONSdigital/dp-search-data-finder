package handler_test

import (
	"context"
	"testing"

	zebedeeclient "github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	clientMock "github.com/ONSdigital/dp-search-data-finder/clients/mock"
	"github.com/ONSdigital/dp-search-data-finder/handler"
	"github.com/ONSdigital/dp-search-data-finder/models"
	searchReindexClient "github.com/ONSdigital/dp-search-reindex-api/sdk"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	searchreindexmodels "github.com/ONSdigital/dp-search-reindex-api/models"
)

var (
	testCtx = context.Background()

	testEvent = models.ReindexRequested{
		JobID:       "job id",
		SearchIndex: "search index",
		TraceID:     "oe433dpe446657gge",
	}

	errZebedee             = errors.New("zebedee test error")
	getPublishedIndexEmpty = func(ctx context.Context, publishedIndexRequestParams *zebedeeclient.PublishedIndexRequestParams) (zebedeeclient.PublishedIndex, error) {
		return zebedeeclient.PublishedIndex{}, errZebedee
	}

	getPublishedIndexFunc = func(ctx context.Context, publishedIndexRequestParams *zebedeeclient.PublishedIndexRequestParams) (zebedeeclient.PublishedIndex, error) {
		return zebedeeclient.PublishedIndex{}, nil
	}

	errSearchReindexApi = errors.New("search reindex api test error")
	postTasksCountEmpty = func(ctx context.Context, headers searchReindexClient.Headers, jobID string) (searchreindexmodels.Task, error) {
		return searchreindexmodels.Task{}, errSearchReindexApi
	}

	postTasksCountFunc = func(ctx context.Context, headers searchReindexClient.Headers, jobID string) (searchreindexmodels.Task, error) {
		return searchreindexmodels.Task{}, nil
	}
)

func TestReindexRequestedHandler_Handle(t *testing.T) {
	t.Parallel()
	Convey("Given an event handler working successfully, and an event containing a URI", t, func() {
		var zebedeeMock = &clientMock.ZebedeeClientMock{GetPublishedIndexFunc: getPublishedIndexFunc}
		var SearchReindexMock = &clientMock.SearchReindexClientMock{PostTasksCountFunc: postTasksCountFunc}
		eventHandler := &handler.ReindexRequestedHandler{ZebedeeCli: zebedeeMock, SearchReindexCli: SearchReindexMock}

		Convey("When given a valid event", func() {
			err := eventHandler.Handle(testCtx, &testEvent)

			Convey("Then error is nil ", func() {
				So(err, ShouldBeNil)

			})
			Convey("And Zebedee and Searchreindexapi is called 1 time with no error ", func() {
				So(zebedeeMock.GetPublishedIndexCalls(), ShouldNotBeEmpty)
				So(zebedeeMock.GetPublishedIndexCalls(), ShouldHaveLength, 1)
				So(SearchReindexMock.PostTasksCountCalls(), ShouldNotBeEmpty)
				So(SearchReindexMock.PostTasksCountCalls(), ShouldHaveLength, 1)
			})
		})
	})
	Convey("Given an event handler not working successfully, and an event containing a jobId", t, func() {
		var zebedeeMockInError = &clientMock.ZebedeeClientMock{GetPublishedIndexFunc: getPublishedIndexEmpty}
		var SearchReindexMock = &clientMock.SearchReindexClientMock{PostTasksCountFunc: postTasksCountEmpty}
		eventHandler := &handler.ReindexRequestedHandler{ZebedeeCli: zebedeeMockInError, SearchReindexCli: SearchReindexMock}

		Convey("When given a valid event", func() {
			err := eventHandler.Handle(testCtx, &testEvent)

			Convey("Then error is not nil ", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("And Zebedee is called 1 time and Searchreindexapi is not called", func() {

				So(err.Error(), ShouldEqual, errZebedee.Error())

				So(zebedeeMockInError.GetPublishedIndexCalls(), ShouldNotBeEmpty)
				So(zebedeeMockInError.GetPublishedIndexCalls(), ShouldHaveLength, 1)
				So(SearchReindexMock.PostTasksCountCalls(), ShouldBeEmpty)
				So(SearchReindexMock.PostTasksCountCalls(), ShouldHaveLength, 0)
			})
		})
	})
}
