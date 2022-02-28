package event_test

import (
	"os"
	"testing"

	"github.com/ONSdigital/dp-search-data-finder/config"
	"github.com/ONSdigital/dp-search-data-finder/event"
	. "github.com/smartystreets/goconvey/convey"
)

func TestReindexRequestedHandler_Handle(t *testing.T) {

	Convey("Given a successful event handler, when Handle is triggered", t, func() {
		eventHandler := &event.ReindexRequestedHandler{}
		filePath := "/tmp/helloworld.txt"
		os.Remove(filePath)
		err := eventHandler.Handle(testCtx, &config.Config{OutputFilePath: filePath}, &testEvent)
		So(err, ShouldBeNil)
	})

	Convey("handler returns an error when cannot write to file", t, func() {
		eventHandler := &event.ReindexRequestedHandler{}
		filePath := ""
		err := eventHandler.Handle(testCtx, &config.Config{OutputFilePath: filePath}, &testEvent)
		So(err, ShouldNotBeNil)
	})
}
