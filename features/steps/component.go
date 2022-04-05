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
	searchReindexClient "github.com/ONSdigital/dp-search-reindex-api/sdk"
	searchReindex "github.com/ONSdigital/dp-search-reindex-api/sdk/v1"
)

type Component struct {
	cfg                  *config.Config
	errorFeature         componenttest.ErrorFeature
	errorChan            chan error
	fakeSearchReindexAPI *SearchReindexFeature
	KafkaConsumer        kafka.IConsumerGroup
	searchReindexClient  searchReindex.Client
	serviceList          *service.ExternalServiceList
	svc                  *service.Service
	zebedeeClient        clients.ZebedeeClient
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

	c.fakeSearchReindexAPI = NewSearchReindexFeature()
	c.cfg.SearchReindexURL = c.fakeSearchReindexAPI.FakeSearchAPI.ResolveURL("")
	c.searchReindexClient = *searchReindex.New(c.cfg.SearchReindexURL, "")

	// Setup responses from registered checkers for component
	c.fakeSearchReindexAPI.setJSONResponseForGetHealth("/health", 200)

	funcDoGetSearchReindexCli := func(cfg *config.Config) searchReindexClient.Client {
		return &c.searchReindexClient
	}

	initMock := &mock.InitialiserMock{
		DoGetKafkaConsumerFunc:       c.DoGetConsumer,
		DoGetHealthCheckFunc:         c.DoGetHealthCheck,
		DoGetHTTPServerFunc:          c.DoGetHTTPServer,
		DoGetZebedeeClientFunc:       c.DoGetZebedeeClient,
		DoGetSearchReindexClientFunc: funcDoGetSearchReindexCli,
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
