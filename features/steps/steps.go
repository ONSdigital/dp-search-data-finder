package steps

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-kafka/v3/kafkatest"
	"github.com/ONSdigital/dp-search-data-finder/models"
	"github.com/ONSdigital/dp-search-data-finder/schema"
	"github.com/cucumber/godog"
	"github.com/rdumont/assistdog"
	"github.com/stretchr/testify/assert"
)

// HealthCheckTest represents a test healthcheck struct that mimics the real healthcheck struct
type HealthCheckTest struct {
	Status    string                  `json:"status"`
	Version   healthcheck.VersionInfo `json:"version"`
	Uptime    time.Duration           `json:"uptime"`
	StartTime time.Time               `json:"start_time"`
	Checks    []*Check                `json:"checks"`
}

// Check represents a health status of a registered app that mimics the real check struct
// As the component test needs to access fields that are not exported in the real struct
type Check struct {
	Name        string     `json:"name"`
	Status      string     `json:"status"`
	StatusCode  int        `json:"status_code"`
	Message     string     `json:"message"`
	LastChecked *time.Time `json:"last_checked"`
	LastSuccess *time.Time `json:"last_success"`
	LastFailure *time.Time `json:"last_failure"`
}

func (c *Component) RegisterSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I wait "([^"]*)" milliseconds`, delayTimeByMilliSeconds)
	ctx.Step(`^all of the downstream services are healthy$`, c.allOfTheDownstreamServicesAreHealthy)
	ctx.Step(`^I should receive the following health JSON response:$`, c.iShouldReceiveTheFollowingHealthJSONResponse)

	ctx.Step(`^these reindex-requested events are consumed:$`, c.theseReindexrequestedEventsAreConsumed)
	ctx.Step(`^I should receive a reindex-requested response$`, c.iShouldReceiveAReindexrequestedResponse)
	ctx.Step(`^the URLs of zebedee and dataset documents are retrieved successfully$`, c.nothingHappens)

	ctx.Step(`^get published-index request to zebedee is successful$`, c.getPublishedIndexRequestToZebedeeIsSuccessful)
	ctx.Step(`^get published-index request to zebedee is unsuccessful$`, c.getPublishedIndexRequestToZebedeeIsUnsuccessful)
	ctx.Step(`^the URLs of zebedee documents are not retrieved$`, c.nothingHappens)

	ctx.Step(`^get requests to dataset-api is successful$`, c.getRequestsToDatasetapiIsSuccessful)
	ctx.Step(`^get requests to dataset-api is unsuccessful$`, c.getRequestsToDatasetapiIsUnsuccessful)
	ctx.Step(`^the URLs of dataset documents are not retrieved$`, c.nothingHappens)
}

// delayTimeByMilliSeconds pauses the goroutine for the given seconds
func delayTimeByMilliSeconds(milliseconds string) (err error) {
	sec, err := strconv.Atoi(milliseconds)
	if err != nil {
		return
	}

	time.Sleep(time.Duration(sec) * time.Millisecond)
	return
}

func (c *Component) allOfTheDownstreamServicesAreHealthy() (err error) {
	c.fakeAPIRouter.setJSONResponseForGetHealth(200)
	err = c.fakeKafkaConsumer.Mock.Checker(context.Background(), healthcheck.NewCheckState("topic-test"))

	return
}

func (c *Component) getPublishedIndexRequestToZebedeeIsSuccessful() {
	c.fakeAPIRouter.setJSONResponseForGetPublishIndex(200)
}

func (c *Component) getPublishedIndexRequestToZebedeeIsUnsuccessful() {
	c.fakeAPIRouter.setJSONResponseForGetPublishIndex(500)
}

func (c *Component) getRequestsToDatasetapiIsSuccessful() {
	c.fakeAPIRouter.setJSONResponseForGetDatasets(200)
}

func (c *Component) getRequestsToDatasetapiIsUnsuccessful() {
	c.fakeAPIRouter.setJSONResponseForGetDatasets(500)
}

func (c *Component) iShouldReceiveTheFollowingHealthJSONResponse(expectedResponse *godog.DocString) error {
	var healthResponse, expectedHealth HealthCheckTest

	responseBody, err := io.ReadAll(c.APIFeature.HttpResponse.Body)
	if err != nil {
		return fmt.Errorf("failed to read response of search data finder component - error: %v", err)
	}

	err = json.Unmarshal(responseBody, &healthResponse)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response of search data finder component - error: %v", err)
	}

	err = json.Unmarshal([]byte(expectedResponse.Content), &expectedHealth)
	if err != nil {
		return fmt.Errorf("failed to unmarshal expected health response - error: %v", err)
	}

	c.validateHealthCheckResponse(healthResponse, expectedHealth)

	return c.errorFeature.StepError()
}

func (c *Component) nothingHappens() {}

func (c *Component) validateHealthCheckResponse(healthResponse, expectedResponse HealthCheckTest) {
	maxExpectedStartTime := c.startTime.Add((c.cfg.HealthCheckInterval + 1) * time.Second)

	assert.Equal(&c.errorFeature, expectedResponse.Status, healthResponse.Status)
	assert.True(&c.errorFeature, healthResponse.StartTime.After(c.startTime))
	assert.True(&c.errorFeature, healthResponse.StartTime.Before(maxExpectedStartTime))
	assert.Greater(&c.errorFeature, healthResponse.Uptime.Seconds(), float64(0))

	c.validateHealthVersion(healthResponse.Version, expectedResponse.Version, maxExpectedStartTime)
	for i, checkResponse := range healthResponse.Checks {
		c.validateHealthCheck(checkResponse, expectedResponse.Checks[i])
	}
}

func (c *Component) validateHealthVersion(versionResponse, expectedVersion healthcheck.VersionInfo, maxExpectedStartTime time.Time) {
	assert.True(&c.errorFeature, versionResponse.BuildTime.Before(maxExpectedStartTime))
	assert.Equal(&c.errorFeature, expectedVersion.GitCommit, versionResponse.GitCommit)
	assert.Equal(&c.errorFeature, expectedVersion.Language, versionResponse.Language)
	assert.NotEmpty(&c.errorFeature, versionResponse.LanguageVersion)
	assert.Equal(&c.errorFeature, expectedVersion.Version, versionResponse.Version)
}

func (c *Component) validateHealthCheck(checkResponse, expectedCheck *Check) {
	maxExpectedHealthCheckTime := c.startTime.Add((c.cfg.HealthCheckInterval + c.cfg.HealthCheckCriticalTimeout + 1) * time.Second)
	assert.Equal(&c.errorFeature, expectedCheck.Name, checkResponse.Name)
	assert.Equal(&c.errorFeature, expectedCheck.Status, checkResponse.Status)
	assert.Equal(&c.errorFeature, expectedCheck.StatusCode, checkResponse.StatusCode)
	assert.Equal(&c.errorFeature, expectedCheck.Message, checkResponse.Message)
	assert.True(&c.errorFeature, checkResponse.LastChecked.Before(maxExpectedHealthCheckTime))
	assert.True(&c.errorFeature, checkResponse.LastChecked.After(c.startTime))

	if expectedCheck.StatusCode == 200 {
		assert.True(&c.errorFeature, checkResponse.LastSuccess.Before(maxExpectedHealthCheckTime))
		assert.True(&c.errorFeature, checkResponse.LastSuccess.After(c.startTime))
		assert.Equal(&c.errorFeature, expectedCheck.LastFailure, checkResponse.LastFailure)
	}
}

func (c *Component) iShouldReceiveAReindexrequestedResponse() error {
	// TODO add assert

	return c.errorFeature.StepError()
}

func (c *Component) theseReindexrequestedEventsAreConsumed(table *godog.Table) error {
	c.fakeAPIRouter.setJSONResponseForGetPublishIndex(200)
	observationEvents, err := c.convertToReindexRequestedEvents(table)
	if err != nil {
		return err
	}

	// consume extracted observations
	for _, e := range observationEvents {
		if err := c.sendToConsumer(e); err != nil {
			return err
		}
	}

	time.Sleep(300 * time.Millisecond)

	return nil
}

func (c *Component) convertToReindexRequestedEvents(table *godog.Table) ([]*models.ReindexRequested, error) {
	assist := assistdog.NewDefault()
	events, err := assist.CreateSlice(&models.ReindexRequested{}, table)
	if err != nil {
		return nil, err
	}
	return events.([]*models.ReindexRequested), nil
}

func (c *Component) sendToConsumer(e *models.ReindexRequested) error {
	bytes, err := schema.ReindexRequestedEvent.Marshal(e)
	if err != nil {
		return err
	}

	newMessage, err := kafkatest.NewMessage(bytes, 0)
	if err != nil {
		return err
	}

	c.fakeKafkaConsumer.Mock.Channels().Upstream <- newMessage
	return nil
}
