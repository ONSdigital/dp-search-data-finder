package steps

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	zebedeeclient "github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	componenttest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/dp-kafka/v3/kafkatest"
	"github.com/ONSdigital/dp-search-data-finder/clients"
	"github.com/ONSdigital/dp-search-data-finder/config"
	"github.com/ONSdigital/dp-search-data-finder/service"
	"github.com/ONSdigital/dp-search-data-finder/service/mock"
	"github.com/pkg/errors"
)

const (
	gitCommitHash = "3t7e5s1t4272646ef477f8ed755"
	appVersion    = "v1.2.3"
)

type state string

var s state = "state"

type Component struct {
	APIFeature        *componenttest.APIFeature
	cfg               *config.Config
	errorFeature      componenttest.ErrorFeature
	errorChan         chan error
	fakeZebedee       *ZebedeeFeature
	fakeKafkaConsumer kafka.IConsumerGroup
	HTTPServer        *http.Server
	serviceList       *service.ExternalServiceList
	serviceRunning    bool
	startTime         time.Time
	svc               *service.Service
	zebedeeClient     clients.ZebedeeClient
}

func NewSearchDataFinderComponent() (*Component, error) {
	ctx := context.Background()

	c := &Component{
		HTTPServer: &http.Server{},
		errorChan:  make(chan error),
	}

	cfg, err := config.Get()
	if err != nil {
		return nil, err
	}

	c.cfg = cfg

	consumer := kafkatest.NewMessageConsumer(true)
	consumer.CheckerFunc = funcCheck
	consumer.StartFunc = func() error { return nil }
	c.fakeKafkaConsumer = consumer

	c.fakeZebedee = NewZebedeeFeature()
	c.cfg.ZebedeeURL = c.fakeZebedee.FakeAPI.ResolveURL("")
	c.zebedeeClient = zebedeeclient.New(c.cfg.ZebedeeURL)

	initMock := &mock.InitialiserMock{
		DoGetKafkaConsumerFunc: func(ctx context.Context, kafkaCfg *config.KafkaConfig) (kafkaConsumer kafka.IConsumerGroup, err error) {
			return c.fakeKafkaConsumer, nil
		},
		DoGetHealthCheckFunc: getHealthCheckOK,
		DoGetHTTPServerFunc:  c.getHTTPServer,
		DoGetZebedeeClientFunc: func(cfg *config.Config) clients.ZebedeeClient {
			return c.zebedeeClient
		},
	}

	// Setup API health endpoints prior to starting component
	c.fakeZebedee.setJSONResponseForGetHealth("/health", 200)

	// Setup healthcheck critical timeout and interval so tests can run faster then
	// using the existing defaults or those set in local or remote environment
	c.cfg.HealthCheckInterval = 1 * time.Second
	c.cfg.HealthCheckCriticalTimeout = 2 * time.Second

	c.serviceList = service.NewServiceList(initMock)

	c.startTime = time.Now()
	c.svc, err = service.Run(ctx, c.cfg, c.serviceList, c.startTime.GoString(), gitCommitHash, appVersion, c.errorChan)
	if err != nil {
		return nil, errors.Wrap(err, "running service failed")
	}

	c.serviceRunning = true

	return c, nil
}

// InitAPIFeature initialises the ApiFeature
func (c *Component) InitAPIFeature() *componenttest.APIFeature {
	c.APIFeature = componenttest.NewAPIFeature(c.InitialiseService)

	return c.APIFeature
}

func (c *Component) Close() error {
	if c.svc != nil && c.serviceRunning {
		c.svc.Close(context.Background())
		c.serviceRunning = false
	}

	return nil
}

func (c *Component) Reset() (*Component, error) {
	ctx := context.WithValue(context.Background(), s, "empty")
	if err := c.fakeKafkaConsumer.Checker(ctx, healthcheck.NewCheckState("topic-test")); err != nil {
		return c, err
	}

	return c, nil
}

// InitialiseService returns the http.Handler that's contained within the component.
func (c *Component) InitialiseService() (http.Handler, error) {
	return c.HTTPServer.Handler, nil
}

func getHealthCheckOK(cfg *config.Config, buildTime, gitCommit, version string) (service.HealthChecker, error) {
	componentBuildTime := strconv.Itoa(int(time.Now().Unix()))
	versionInfo, err := healthcheck.NewVersionInfo(componentBuildTime, gitCommitHash, appVersion)
	if err != nil {
		return nil, err
	}
	hc := healthcheck.New(versionInfo, cfg.HealthCheckCriticalTimeout, cfg.HealthCheckInterval)
	return &hc, nil
}

func (c *Component) getHTTPServer(bindAddr string, router http.Handler) service.HTTPServer {
	c.HTTPServer.Addr = bindAddr
	c.HTTPServer.Handler = router
	return c.HTTPServer
}

func funcCheck(ctx context.Context, state *healthcheck.CheckState) error {
	var str string
	healthState := ctx.Value(s)
	if healthState != nil {
		str = fmt.Sprintf("%v", healthState)
	}

	if str == "empty" {
		if err := state.Update("", "", 0); err != nil {
			return err
		}
	} else {
		if err := state.Update(healthcheck.StatusOK, "OK", 0); err != nil {
			return err
		}
	}

	return nil
}
