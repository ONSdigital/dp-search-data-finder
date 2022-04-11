package handler_test

import (
	"context"
	"testing"

	zebedeeclient "github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	clientMock "github.com/ONSdigital/dp-search-data-finder/clients/mock"
	"github.com/ONSdigital/dp-search-data-finder/handler"
	"github.com/ONSdigital/dp-search-data-finder/models"
	searchReindexSDK "github.com/ONSdigital/dp-search-reindex-api/sdk"
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
	getPublishedIndexEmpty = func(ctx context.Context, publishedIndexRequestParams *zebedeeclient.PublishedIndexRequestParams) (zebedeeclient.PublishedIndex, error) {
		return zebedeeclient.PublishedIndex{}, errZebedee
	}

	getPublishedIndexFunc = func(ctx context.Context, publishedIndexRequestParams *zebedeeclient.PublishedIndexRequestParams) (zebedeeclient.PublishedIndex, error) {
		return zebedeeclient.PublishedIndex{}, nil
	}

	errSearchReindex = errors.New("search reindex test error")
	patchJobFails    = func(context.Context, searchReindexSDK.Headers, string, []searchReindexSDK.PatchOperation) (string, error) {
		return "", errSearchReindex
	}

	patchJobFunc = func(context.Context, searchReindexSDK.Headers, string, []searchReindexSDK.PatchOperation) (string, error) {
		return `"56b6890f1321590998d5fd8d293b620581ff3531"`, nil
	}
)

func TestReindexRequestedHandler_Handle(t *testing.T) {
	t.Parallel()
	Convey("Given an event handler working successfully, and an event containing a URI", t, func() {
		var zebedeeMock = &clientMock.ZebedeeClientMock{GetPublishedIndexFunc: getPublishedIndexFunc}
		var searchReindexMock = &clientMock.SearchReindexClientMock{PatchJobFunc: patchJobFunc}
		eventHandler := &handler.ReindexRequestedHandler{ZebedeeCli: zebedeeMock,
			SearchReindexCli: searchReindexMock}

		Convey("When given a valid event", func() {
			err := eventHandler.Handle(testCtx, &testEvent)

			Convey("Then Zebedee is called 1 time with no error ", func() {
				So(err, ShouldBeNil)

				So(zebedeeMock.GetPublishedIndexCalls(), ShouldNotBeEmpty)
				So(zebedeeMock.GetPublishedIndexCalls(), ShouldHaveLength, 1)
			})

			Convey("And Search Reindex is called 1 time with no error ", func() {
				So(err, ShouldBeNil)

				So(searchReindexMock.PatchJobCalls(), ShouldNotBeEmpty)
				So(searchReindexMock.PatchJobCalls(), ShouldHaveLength, 1)
			})
		})
	})
	Convey("Given an event handler not working successfully with Zebedee, and an event containing a jobId", t, func() {
		var zebedeeMockInError = &clientMock.ZebedeeClientMock{GetPublishedIndexFunc: getPublishedIndexEmpty}
		var searchReindexMock = &clientMock.SearchReindexClientMock{PatchJobFunc: patchJobFunc}
		eventHandler := &handler.ReindexRequestedHandler{ZebedeeCli: zebedeeMockInError, SearchReindexCli: searchReindexMock}

		Convey("When given a valid event", func() {
			err := eventHandler.Handle(testCtx, &testEvent)

			Convey("Then Zebedee is called 1 time with the expected error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, errZebedee.Error())

				So(zebedeeMockInError.GetPublishedIndexCalls(), ShouldNotBeEmpty)
				So(zebedeeMockInError.GetPublishedIndexCalls(), ShouldHaveLength, 1)
			})

			Convey("And so the handler returns and Search Reindex is not called", func() {
				So(searchReindexMock.PatchJobCalls(), ShouldBeEmpty)
				So(searchReindexMock.PatchJobCalls(), ShouldHaveLength, 0)
			})
		})
	})
	Convey("Given an event handler not working successfully with Search Reindex, and an event containing a jobId", t, func() {
		var zebedeeMock = &clientMock.ZebedeeClientMock{GetPublishedIndexFunc: getPublishedIndexFunc}
		var searchReindexMockInError = &clientMock.SearchReindexClientMock{PatchJobFunc: patchJobFails}
		eventHandler := &handler.ReindexRequestedHandler{ZebedeeCli: zebedeeMock, SearchReindexCli: searchReindexMockInError}

		Convey("When given a valid event", func() {
			err := eventHandler.Handle(testCtx, &testEvent)

			Convey("Then Zebedee is called 1 time", func() {
				So(zebedeeMock.GetPublishedIndexCalls(), ShouldNotBeEmpty)
				So(zebedeeMock.GetPublishedIndexCalls(), ShouldHaveLength, 1)
			})

			Convey("And Search Reindex is called 1 time with the expected error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, errSearchReindex.Error())

				So(searchReindexMockInError.PatchJobCalls(), ShouldNotBeEmpty)
				So(searchReindexMockInError.PatchJobCalls(), ShouldHaveLength, 1)
			})
		})
	})
}
