package steps

import (
	"context"
	"net/http"
	"os"

	componenttest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/dp-kafka/v3/kafkatest"
	dphttp "github.com/ONSdigital/dp-net/v2/http"
	"github.com/ONSdigital/dp-search-data-finder/config"
	"github.com/ONSdigital/dp-search-data-finder/service"
	"github.com/ONSdigital/dp-search-data-finder/service/mock"
)

type Component struct {
	componenttest.ErrorFeature
	serviceList   *service.ExternalServiceList
	KafkaConsumer kafka.IConsumerGroup
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

	cfg, err := config.Get()
	if err != nil {
		return nil
	}

	c.cfg = cfg

	initMock := &mock.InitialiserMock{
		DoGetKafkaConsumerFunc: c.DoGetConsumer,
		DoGetHealthCheckFunc:   c.DoGetHealthCheck,
		DoGetHTTPServerFunc:    c.DoGetHTTPServer,
	}

	c.serviceList = service.NewServiceList(initMock)

	return c
}

func (c *Component) Close() {
	os.Remove(c.cfg.OutputFilePath)
}

func (c *Component) Reset() {
	os.Remove(c.cfg.OutputFilePath)
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

func funcCheck(ctx context.Context, state *healthcheck.CheckState) error {
	return nil
}
