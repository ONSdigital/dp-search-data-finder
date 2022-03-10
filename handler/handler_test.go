package handler_test

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/health"
	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-mocking/httpmocks"
	dphttp "github.com/ONSdigital/dp-net/v2/http"
	"github.com/ONSdigital/dp-search-data-finder/handler"
	"github.com/ONSdigital/dp-search-data-finder/models"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	testCtx = context.Background()

	testEvent = models.ReindexRequested{
		JobID:       "job id",
		SearchIndex: "search index",
		TraceID:     "trace id",
	}
)

const testHost = "http://localhost:8080"

func TestReindexRequestedHandler_Handle(t *testing.T) {
	Convey("Given a successful event handler, when Handle is triggered", t, func() {
		documentContent, err := os.ReadFile("./publishedindex.json")
		So(err, ShouldBeNil)
		body := httpmocks.NewReadCloserMock(documentContent, nil)
		response := httpmocks.NewResponseMock(body, http.StatusOK)

		httpClient := newMockHTTPClient(response, nil)
		zebedeeClient := newZebedeeClient(httpClient)

		eventHandler := &handler.ReindexRequestedHandler{ZebedeeClient: zebedeeClient}
		err = eventHandler.Handle(testCtx, &testEvent)

		Convey("Then no error is returned", func() {
			So(err, ShouldBeNil)
		})
	})

	Convey("Given an event handler containing a nil ZebedeeClient, when Handle is triggered", t, func() {
		eventHandler := &handler.ReindexRequestedHandler{ZebedeeClient: nil}
		err := eventHandler.Handle(testCtx, &testEvent)

		Convey("Then an error is returned", func() {
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "the Zebedee client in the reindex requested handler must not be nil")
		})
	})
}

func newMockHTTPClient(r *http.Response, err error) *dphttp.ClienterMock {
	return &dphttp.ClienterMock{
		SetPathsWithNoRetriesFunc: func(paths []string) {
		},
		DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
			return r, err
		},
		GetPathsWithNoRetriesFunc: func() []string {
			return []string{"/healthcheck"}
		},
	}
}

func newZebedeeClient(httpClient *dphttp.ClienterMock) *zebedee.Client {
	healthClient := health.NewClientWithClienter("", testHost, httpClient)
	zebedeeClient := zebedee.NewWithHealthClient(healthClient)
	return zebedeeClient
}
