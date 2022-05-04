package handler_test

import (
	"context"
	"testing"

	zebedeeClient "github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	clientMock "github.com/ONSdigital/dp-search-data-finder/clients/mock"
	"github.com/ONSdigital/dp-search-data-finder/handler"
	"github.com/ONSdigital/dp-search-data-finder/models"
	searchReindexModels "github.com/ONSdigital/dp-search-reindex-api/models"
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
	getPublishedIndexEmpty = func(ctx context.Context, publishedIndexRequestParams *zebedeeClient.PublishedIndexRequestParams) (zebedeeClient.PublishedIndex, error) {
		return zebedeeClient.PublishedIndex{}, errZebedee
	}

	getPublishedIndexFunc = func(ctx context.Context, publishedIndexRequestParams *zebedeeClient.PublishedIndexRequestParams) (zebedeeClient.PublishedIndex, error) {
		return zebedeeClient.PublishedIndex{}, nil
	}

	errSearchReindexAPI = errors.New("search reindex api test error")
	patchJobFails       = func(context.Context, searchReindexSDK.Headers, string, []searchReindexSDK.PatchOperation) (*searchReindexSDK.RespHeaders, error) {
		return &searchReindexSDK.RespHeaders{}, errSearchReindexAPI
	}

	patchJobFunc = func(context.Context, searchReindexSDK.Headers, string, []searchReindexSDK.PatchOperation) (*searchReindexSDK.RespHeaders, error) {
		respHeaders := searchReindexSDK.RespHeaders{
			ETag: `"56b6890f1321590998d5fd8d293b620581ff3531"`,
		}
		return &respHeaders, nil
	}

	postTaskEmpty = func(ctx context.Context, headers searchReindexSDK.Headers, jobID string, taskToCreate searchReindexModels.TaskToCreate) (*searchReindexSDK.RespHeaders, *searchReindexModels.Task, error) {
		return &searchReindexSDK.RespHeaders{}, &searchReindexModels.Task{}, errSearchReindexAPI
	}

	postTaskFunc = func(ctx context.Context, headers searchReindexSDK.Headers, jobID string, taskToCreate searchReindexModels.TaskToCreate) (*searchReindexSDK.RespHeaders, *searchReindexModels.Task, error) {
		return &searchReindexSDK.RespHeaders{}, &searchReindexModels.Task{}, nil
	}
)

func TestReindexRequestedHandler_Handle(t *testing.T) {
	t.Parallel()
	Convey("Given an event handler working successfully, and an event containing a URI", t, func() {
		var zebedeeMock = &clientMock.ZebedeeClientMock{GetPublishedIndexFunc: getPublishedIndexFunc}
		var searchReindexMock = &clientMock.SearchReindexClientMock{PatchJobFunc: patchJobFunc, PostTaskFunc: postTaskFunc}
		eventHandler := &handler.ReindexRequestedHandler{ZebedeeCli: zebedeeMock, SearchReindexCli: searchReindexMock}

		Convey("When given a valid event", func() {
			err := eventHandler.Handle(testCtx, &testEvent)

			Convey("Then error is nil ", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then Zebedee is called 1 time", func() {
				So(zebedeeMock.GetPublishedIndexCalls(), ShouldNotBeEmpty)
				So(zebedeeMock.GetPublishedIndexCalls(), ShouldHaveLength, 1)
			})

			Convey("And Search Reindex is called 2 times in total", func() {
				So(searchReindexMock.PatchJobCalls(), ShouldNotBeEmpty)
				So(searchReindexMock.PatchJobCalls(), ShouldHaveLength, 1)
				So(searchReindexMock.PostTaskCalls(), ShouldNotBeEmpty)
				So(searchReindexMock.PostTaskCalls(), ShouldHaveLength, 1)
			})
		})
	})
	Convey("Given an event handler not working successfully with Zebedee, and an event containing a jobId", t, func() {
		var zebedeeMockInError = &clientMock.ZebedeeClientMock{GetPublishedIndexFunc: getPublishedIndexEmpty}
		var searchReindexMock = &clientMock.SearchReindexClientMock{PatchJobFunc: patchJobFunc, PostTaskFunc: postTaskEmpty}
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
				So(searchReindexMock.PostTaskCalls(), ShouldBeEmpty)
				So(searchReindexMock.PostTaskCalls(), ShouldHaveLength, 0)
			})
		})
	})
}

func TestGetPayloads(t *testing.T) {
	t.Parallel()
	Convey("Given an empty task name", t, func() {
		testTaskNames := map[string]string{}
		var zebedeeMock = &clientMock.ZebedeeClientMock{GetPublishedIndexFunc: getPublishedIndexFunc}
		var SearchReindexMock = &clientMock.SearchReindexClientMock{PostTaskFunc: postTaskFunc}
		eventHandler := &handler.ReindexRequestedHandler{ZebedeeCli: zebedeeMock, SearchReindexCli: SearchReindexMock}

		Convey("When preparing the payload", func() {
			payload, err := eventHandler.GetPayload(testCtx, 10, testTaskNames)

			Convey("Then it should throws an error", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("And the payload is empty", func() {
				So(payload, ShouldBeEmpty)
			})
		})
	})
	Convey("Given a valid task name", t, func() {
		testTaskNames := map[string]string{
			"zebedee": "zebedee",
		}
		var zebedeeMock = &clientMock.ZebedeeClientMock{GetPublishedIndexFunc: getPublishedIndexFunc}
		var SearchReindexMock = &clientMock.SearchReindexClientMock{PostTaskFunc: postTaskFunc}
		eventHandler := &handler.ReindexRequestedHandler{ZebedeeCli: zebedeeMock, SearchReindexCli: SearchReindexMock}

		Convey("When preparing payloads", func() {
			payload, err := eventHandler.GetPayload(testCtx, 10, testTaskNames)

			Convey("Then it should not throws error", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the payload is not empty", func() {
				So(payload, ShouldNotBeEmpty)
			})
		})
	})
	Convey("Given an invalid task name", t, func() {
		testTaskNames := map[string]string{
			"invalidTask": "invalid",
		}
		var zebedeeMock = &clientMock.ZebedeeClientMock{GetPublishedIndexFunc: getPublishedIndexFunc}
		var searchReindexMock = &clientMock.SearchReindexClientMock{PostTaskFunc: postTaskFunc}
		eventHandler := &handler.ReindexRequestedHandler{ZebedeeCli: zebedeeMock, SearchReindexCli: searchReindexMock}

		Convey("When preparing payloads", func() {
			payload, err := eventHandler.GetPayload(testCtx, 10, testTaskNames)

			Convey("Then it should throws error", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("Then it throws error", func() {
				So(payload, ShouldBeEmpty)
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
				So(err.Error(), ShouldEqual, errSearchReindexAPI.Error())

				So(searchReindexMockInError.PatchJobCalls(), ShouldNotBeEmpty)
				So(searchReindexMockInError.PatchJobCalls(), ShouldHaveLength, 1)
			})
		})
	})
}
