package steps

import (
	"context"
	"net/http"

	zebedeeclient "github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	componenttest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/dp-kafka/v3/kafkatest"
	dphttp "github.com/ONSdigital/dp-net/v2/http"
	"github.com/ONSdigital/dp-search-data-finder/clients"
	zebedeeClientMock "github.com/ONSdigital/dp-search-data-finder/clients/mock"
	"github.com/ONSdigital/dp-search-data-finder/config"
	"github.com/ONSdigital/dp-search-data-finder/service"
	"github.com/ONSdigital/dp-search-data-finder/service/mock"
)

type Component struct {
	componenttest.ErrorFeature
	serviceList   *service.ExternalServiceList
	KafkaConsumer kafka.IConsumerGroup
	zebedeeClient clients.ZebedeeClient
	errorChan     chan error
	svc           *service.Service
	cfg           *config.Config
}

func NewComponent() *Component {
	c := &Component{errorChan: make(chan error)}

	consumer := kafkatest.NewMessageConsumer(false)
	consumer.CheckerFunc = funcCheck
	consumer.StartFunc = func() error { return nil }
	c.KafkaConsumer = consumer

	c.zebedeeClient = &zebedeeClientMock.ZebedeeClientMock{
		CheckerFunc: func(in1 context.Context, in2 *healthcheck.CheckState) error { return nil },
		GetPublishedIndexFunc: func(ctx context.Context, publishedIndexRequestParams *zebedeeclient.PublishedIndexRequestParams) (zebedeeclient.PublishedIndex, error) {
			return zebedeeclient.PublishedIndex{}, nil
		},
	}

	cfg, err := config.Get()
	if err != nil {
		return nil
	}

	c.cfg = cfg

	initMock := &mock.InitialiserMock{
		DoGetKafkaConsumerFunc: c.DoGetConsumer,
		DoGetHealthCheckFunc:   c.DoGetHealthCheck,
		DoGetHTTPServerFunc:    c.DoGetHTTPServer,
		DoGetZebedeeClientFunc: c.DoGetZebedeeClient,
	}

	c.serviceList = service.NewServiceList(initMock)

	return c
}

func (c *Component) Close() {
}

func (c *Component) Reset() {
}

func (c *Component) DoGetHealthCheck(cfg *config.Config, buildTime, gitCommit, version string) (service.HealthChecker, error) {
	return &mock.HealthCheckerMock{
		AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
		StartFunc:    func(ctx context.Context) {},
		StopFunc:     func() {},
	}, nil
}

func (c *Component) DoGetHTTPServer(bindAddr string, router http.Handler) service.HTTPServer {
	return dphttp.NewServer(bindAddr, router)
}

func (c *Component) DoGetConsumer(ctx context.Context, kafkaCfg *config.KafkaConfig) (kafkaConsumer kafka.IConsumerGroup, err error) {
	return c.KafkaConsumer, nil
}

func (c *Component) DoGetZebedeeClient(cfg *config.Config) clients.ZebedeeClient {
	return c.zebedeeClient
}

func funcCheck(ctx context.Context, state *healthcheck.CheckState) error {
	return nil
}