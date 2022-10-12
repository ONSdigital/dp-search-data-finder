package event_test

import (
	"context"
	"sync"
	"testing"

	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/dp-kafka/v3/kafkatest"
	"github.com/ONSdigital/dp-search-data-finder/config"
	"github.com/ONSdigital/dp-search-data-finder/event"
	"github.com/ONSdigital/dp-search-data-finder/event/mock"
	"github.com/ONSdigital/dp-search-data-finder/models"
	"github.com/ONSdigital/dp-search-data-finder/schema"
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

// kafkaStubConsumer mock which exposes Channels function returning empty channels
// to be used on tests that are not supposed to receive any kafka message
// var kafkaStubConsumer = &kafkatest.IConsumerGroupMock{
//	ChannelsFunc: func() *kafka.ConsumerGroupChannels {
//		return &kafka.ConsumerGroupChannels{}
//	},
//}

func TestConsume(t *testing.T) {
	Convey("Given kafka consumer and event handler mocks", t, func() {
		cgChannels := &kafka.ConsumerGroupChannels{Upstream: make(chan kafka.Message, 2)}
		mockConsumer := &kafkatest.IConsumerGroupMock{
			ChannelsFunc: func() *kafka.ConsumerGroupChannels { return cgChannels },
		}

		handlerWg := &sync.WaitGroup{}
		mockEventHandler := &mock.HandlerMock{
			HandleFunc: func(ctx context.Context, event *models.ReindexRequested) {
				defer handlerWg.Done()
			},
		}

		Convey("And a kafka message with the valid schema being sent to the Upstream channel", func() {
			message, err := kafkatest.NewMessage(marshal(testEvent), 0)
			if err != nil {
				t.Errorf("failed to create new message - err: %v", err)
			}

			mockConsumer.Channels().Upstream <- message
			Convey("When consume message is called", func() {
				handlerWg.Add(1)
				event.Consume(testCtx, mockConsumer, mockEventHandler, &config.Config{KafkaConfig: config.KafkaConfig{NumWorkers: 1}})
				handlerWg.Wait()

				Convey("An event is sent to the mockEventHandler ", func() {
					So(len(mockEventHandler.HandleCalls()), ShouldEqual, 1)
					So(*mockEventHandler.HandleCalls()[0].ReindexRequested, ShouldResemble, testEvent)
				})

				Convey("The message is committed and the consumer is released", func() {
					<-message.UpstreamDone()
					So(len(message.CommitCalls()), ShouldEqual, 1)
					So(len(message.ReleaseCalls()), ShouldEqual, 1)
				})
			})
		})

		Convey("And two kafka messages, one with a valid schema and one with an invalid schema", func() {
			validMessage, err := kafkatest.NewMessage(marshal(testEvent), 1)
			if err != nil {
				t.Errorf("failed to create valid kafka message - err: %v", err)
			}
			invalidMessage, err := kafkatest.NewMessage([]byte("invalid schema"), 0)
			if err != nil {
				t.Errorf("failed to create invalid kafka message - err: %v", err)
			}
			mockConsumer.Channels().Upstream <- invalidMessage
			mockConsumer.Channels().Upstream <- validMessage
			Convey("When consume messages is called", func() {
				handlerWg.Add(1)
				event.Consume(testCtx, mockConsumer, mockEventHandler, &config.Config{KafkaConfig: config.KafkaConfig{NumWorkers: 1}})
				handlerWg.Wait()

				Convey("Only the valid event is sent to the mockEventHandler ", func() {
					So(len(mockEventHandler.HandleCalls()), ShouldEqual, 1)
					So(*mockEventHandler.HandleCalls()[0].ReindexRequested, ShouldResemble, testEvent)
				})

				Convey("Only the valid message is committed, but the consumer is released for both messages", func() {
					<-validMessage.UpstreamDone()
					<-invalidMessage.UpstreamDone()
					So(len(validMessage.CommitCalls()), ShouldEqual, 1)
					So(len(invalidMessage.CommitCalls()), ShouldEqual, 1)
					So(len(validMessage.ReleaseCalls()), ShouldEqual, 1)
					So(len(invalidMessage.ReleaseCalls()), ShouldEqual, 1)
				})
			})
		})

		Convey("With a failing handler and a kafka message with the valid schema being sent to the Upstream channel", func() {
			message, err := kafkatest.NewMessage(marshal(testEvent), 0)
			if err != nil {
				t.Errorf("failed to create new kafka message - err: %v", err)
			}
			mockConsumer.Channels().Upstream <- message
			Convey("When consume message is called", func() {
				handlerWg.Add(1)
				event.Consume(testCtx, mockConsumer, mockEventHandler, &config.Config{KafkaConfig: config.KafkaConfig{NumWorkers: 1}})
				handlerWg.Wait()
				Convey("An event is sent to the mockEventHandler ", func() {
					So(len(mockEventHandler.HandleCalls()), ShouldEqual, 1)
					So(*mockEventHandler.HandleCalls()[0].ReindexRequested, ShouldResemble, testEvent)
				})
				Convey("The message is committed and the consumer is released", func() {
					<-message.UpstreamDone()
					So(len(message.CommitCalls()), ShouldEqual, 1)
					So(len(message.ReleaseCalls()), ShouldEqual, 1)
				})
			})
		})
	})
}

// marshal helper method to marshal a event into a []byte
func marshal(irEvent models.ReindexRequested) []byte {
	bytes, err := schema.ReindexRequestedEvent.Marshal(irEvent)
	So(err, ShouldBeNil)
	return bytes
}
