package steps

import (
	"context"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ONSdigital/dp-kafka/v3/kafkatest"
	"github.com/ONSdigital/dp-search-data-finder/event"
	"github.com/ONSdigital/dp-search-data-finder/schema"
	"github.com/ONSdigital/dp-search-data-finder/service"
	"github.com/cucumber/godog"
	"github.com/rdumont/assistdog"
	"github.com/stretchr/testify/assert"
)

func (c *Component) RegisterSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^these reindex-requested events are consumed:$`, c.theseReindexrequestedEventsAreConsumed)
	ctx.Step(`^I should receive a reindex-requested response$`, c.iShouldReceiveAReindexrequestedResponse)
}

func (c *Component) iShouldReceiveAReindexrequestedResponse() error {
	content, err := ioutil.ReadFile(c.cfg.OutputFilePath)
	if err != nil {
		return err
	}

	assert.Equal(c, "Hello there!", string(content))

	return c.StepError()
}

func (c *Component) theseReindexrequestedEventsAreConsumed(table *godog.Table) error {

	observationEvents, err := c.convertToReindexRequestedEvents(table)
	if err != nil {
		return err
	}

	signals := registerInterrupt()

	// run application in separate goroutine
	go func() {
		c.svc, err = service.Run(context.Background(), c.serviceList, "", "", "", c.errorChan)
	}()

	// consume extracted observations
	for _, e := range observationEvents {
		if err := c.sendToConsumer(e); err != nil {
			return err
		}
	}

	time.Sleep(300 * time.Millisecond)

	// kill application
	signals <- os.Interrupt

	return nil
}

func (c *Component) convertToReindexRequestedEvents(table *godog.Table) ([]*event.ReindexRequested, error) {
	assist := assistdog.NewDefault()
	events, err := assist.CreateSlice(&event.ReindexRequested{}, table)
	if err != nil {
		return nil, err
	}
	return events.([]*event.ReindexRequested), nil
}

func (c *Component) sendToConsumer(e *event.ReindexRequested) error {
	bytes, err := schema.ReindexRequestedEvent.Marshal(e)
	if err != nil {
		return err
	}

	c.KafkaConsumer.Channels().Upstream <- kafkatest.NewMessage(bytes, 0)
	return nil

}

func registerInterrupt() chan os.Signal {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	return signals
}