package event_test

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/health"
	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-mocking/httpmocks"
	dphttp "github.com/ONSdigital/dp-net/v2/http"
	"github.com/ONSdigital/dp-search-data-finder/event"
	. "github.com/smartystreets/goconvey/convey"
)

const testHost = "http://localhost:8080"

func TestReindexRequestedHandler_Handle(t *testing.T) {
	Convey("Given a successful event handler, when Handle is triggered", t, func() {
		documentContent, err := os.ReadFile("./mock/publishedindex.json")
		So(err, ShouldBeNil)
		body := httpmocks.NewReadCloserMock(documentContent, nil)
		response := httpmocks.NewResponseMock(body, http.StatusOK)

		httpClient := newMockHTTPClient(response, nil)
		zebedeeClient := newZebedeeClient(httpClient)

		eventHandler := &event.ReindexRequestedHandler{ZebedeeClient: zebedeeClient}
		err = eventHandler.Handle(testCtx, &testEvent)
		So(err, ShouldBeNil)
	})
}

func newMockHTTPClient(r *http.Response, err error) *dphttp.ClienterMock {
	return &dphttp.ClienterMock{
		SetPathsWithNoRetriesFunc: func(paths []string) {
			return
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
