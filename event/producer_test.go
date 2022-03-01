package event_test

import (
	"context"
	"testing"
	"time"

	dpkafka "github.com/ONSdigital/dp-kafka/v2"
	"github.com/ONSdigital/dp-kafka/v2/kafkatest"
	"github.com/ONSdigital/dp-search-data-finder/event"
	"github.com/ONSdigital/dp-search-data-finder/event/mock"
	"github.com/ONSdigital/dp-search-data-finder/models"
	"github.com/ONSdigital/dp-search-data-finder/schema"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	someURI         = "test-uri"
	someDataType    = "test-datatype"
	someJobID       = "test-Jobid"
	someTraceID     = "w34234dgdge335g3333"
	someSearchIndex = "test-searchindex"
)

var (
	ctx = context.Background()

	expectedContentUpdatedEvent = models.ContentUpdated{
		URI:         someURI,
		DataType:    someDataType,
		JobID:       someJobID,
		TraceID:     someTraceID,
		SearchIndex: someSearchIndex,
	}
)

func TestProducer_ContentUpdated(t *testing.T) {
	Convey("Given ContentUpdatedProducer has been configured correctly", t, func() {
		pChannels := &dpkafka.ProducerChannels{
			Output: make(chan []byte, 1),
		}

		kafkaProducerMock := &kafkatest.IProducerMock{
			ChannelsFunc: func() *dpkafka.ProducerChannels {
				return pChannels
			},
		}

		marshallerMock := &mock.MarshallerMock{
			MarshalFunc: func(s interface{}) ([]byte, error) {
				return schema.ContentUpdatedEvent.Marshal(s)
			},
		}

		// event is message
		contentUpdatedProducer := event.ContentUpdatedProducer{
			Producer:   kafkaProducerMock,
			Marshaller: marshallerMock,
		}
		Convey("When ContentUpdate is called on the event producer", func() {
			err := contentUpdatedProducer.ContentUpdate(ctx, expectedContentUpdatedEvent)
			So(err, ShouldBeNil)

			var avroBytes []byte
			var testTimeout = time.Second * 5
			select {
			case avroBytes = <-pChannels.Output:
				t.Log("avro byte sent to producer output")
			case <-time.After(testTimeout):
				t.Fatalf("failing test due to timing out after %v seconds", testTimeout)
				t.FailNow()
			}

			Convey("Then the expected bytes are sent to producer.output", func() {
				var actual models.ContentUpdated
				err = schema.ContentUpdatedEvent.Unmarshal(avroBytes, &actual)
				So(err, ShouldBeNil)
				So(expectedContentUpdatedEvent, ShouldResemble, actual)
			})
		})
	})
}
