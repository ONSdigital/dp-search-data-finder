package event_test

import (
	"context"
	"testing"
	"time"

	dpkafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/dp-kafka/v3/kafkatest"
	"github.com/ONSdigital/dp-search-data-finder/config"
	"github.com/ONSdigital/dp-search-data-finder/event"
	"github.com/ONSdigital/dp-search-data-finder/event/mock"
	"github.com/ONSdigital/dp-search-data-finder/models"
	"github.com/ONSdigital/dp-search-data-finder/schema"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	someURI          = "test-uri"
	someDataType     = "test-datatype"
	someCollectionID = "test-collectionID"
	someJobID        = "test-Jobid"
	someTask         = "test-task"
	testCount        = 10
	someTraceID      = "w34234dgdge335g3333"
	someSearchIndex  = "test-searchindex"
)

var (
	ctx = context.Background()

	expectedContentUpdatedEvent = models.ContentUpdated{
		URI:          someURI,
		DataType:     someDataType,
		CollectionID: someCollectionID,
		JobID:        someJobID,
		TraceID:      someTraceID,
		SearchIndex:  someSearchIndex,
	}

	expectedReindexTaskCountsEvent = models.ReindexTaskCounts{
		JobID:               someJobID,
		Task:                someTask,
		ExtractionCompleted: false,
		Count:               testCount,
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

		Convey("When ContentUpdate is called on the event producer with EnablePublishContentUpdatedTopic enabled", func() {
			err := contentUpdatedProducer.ContentUpdate(ctx, &config.Config{EnablePublishContentUpdatedTopic: true}, expectedContentUpdatedEvent)
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

func TestProducer_ReindexTaskCounts(t *testing.T) {
	Convey("Given ReindexTaskCounts has been configured correctly", t, func() {
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
				return schema.ReindexTaskCounts.Marshal(s)
			},
		}

		// event is message
		reindexTaskCountsProducer := event.ReindexTaskCountsProducer{
			Producer:   kafkaProducerMock,
			Marshaller: marshallerMock,
		}

		Convey("When ReindexTaskCounts is called on the event producer with EnablePublishReindexTaskCountsTopic enabled", func() {
			err := reindexTaskCountsProducer.TaskCounts(ctx, &config.Config{EnableReindexTaskCounts: true}, expectedReindexTaskCountsEvent)
			So(err, ShouldBeNil)

			var avroBytes []byte
			var testTimeout = time.Second * 5
			select {
			case avroBytes = <-pChannels.Output:
				t.Log("avro byte sent to producer output")
			case <-time.After(testTimeout):
				t.Fatalf("failing test due to timing out after %v seconds", testTimeout)
			}

			Convey("Then the expected bytes are sent to producer.output", func() {
				var actual models.ReindexTaskCounts
				err = schema.ReindexTaskCounts.Unmarshal(avroBytes, &actual)
				So(err, ShouldBeNil)
				So(expectedReindexTaskCountsEvent, ShouldResemble, actual)
			})
		})
	})
}
