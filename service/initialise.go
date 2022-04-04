package service

import (
	"context"
	"net/http"

	apihealthcheck "github.com/ONSdigital/dp-api-clients-go/v2/health"
	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dpkafka "github.com/ONSdigital/dp-kafka/v3"
	dphttp "github.com/ONSdigital/dp-net/v2/http"
	"github.com/ONSdigital/dp-search-data-finder/clients"
	"github.com/ONSdigital/dp-search-data-finder/config"
	dpsearchreindex "github.com/ONSdigital/dp-search-reindex-api/sdk/v1"
)

// ExternalServiceList holds the initialiser and initialisation state of external services.
type ExternalServiceList struct {
	HealthCheck         bool
	KafkaConsumer       bool
	Init                Initialiser
	ZebedeeCli          bool
	SearchReindexAPICli bool
}

// NewServiceList creates a new service list with the provided initialiser
func NewServiceList(initialiser Initialiser) *ExternalServiceList {
	return &ExternalServiceList{
		HealthCheck:         false,
		KafkaConsumer:       false,
		Init:                initialiser,
		ZebedeeCli:          false,
		SearchReindexAPICli: false,
	}
}

// Init implements the Initialiser interface to initialise dependencies
type Init struct{}

// GetHTTPServer creates an http server and sets the Server flag to true
func (e *ExternalServiceList) GetHTTPServer(bindAddr string, router http.Handler) HTTPServer {
	s := e.Init.DoGetHTTPServer(bindAddr, router)
	return s
}

// GetKafkaConsumer creates a Kafka consumer and sets the consumer flag to true
func (e *ExternalServiceList) GetKafkaConsumer(ctx context.Context, cfg *config.Config) (dpkafka.IConsumerGroup, error) {
	consumer, err := e.Init.DoGetKafkaConsumer(ctx, &cfg.KafkaConfig)
	if err != nil {
		return nil, err
	}
	e.KafkaConsumer = true
	return consumer, nil
}

// GetHealthCheck creates a healthcheck with versionInfo and sets teh HealthCheck flag to true
func (e *ExternalServiceList) GetHealthCheck(cfg *config.Config, buildTime, gitCommit, version string) (HealthChecker, error) {
	hc, err := e.Init.DoGetHealthCheck(cfg, buildTime, gitCommit, version)
	if err != nil {
		return nil, err
	}
	e.HealthCheck = true
	return hc, nil
}

// DoGetHTTPServer creates an HTTP Server with the provided bind address and router
func (e *Init) DoGetHTTPServer(bindAddr string, router http.Handler) HTTPServer {
	s := dphttp.NewServer(bindAddr, router)
	s.HandleOSSignals = false
	return s
}

// GetZebedee return zebedee client
func (e *ExternalServiceList) GetZebedee(cfg *config.Config) clients.ZebedeeClient {
	zebedeeClient := e.Init.DoGetZebedeeClient(cfg)
	e.ZebedeeCli = true
	return zebedeeClient
}

// DoGetZebedeeClient gets and initialises the Zebedee Client
func (e *Init) DoGetZebedeeClient(cfg *config.Config) clients.ZebedeeClient {
	zebedeeClient := zebedee.New(cfg.ZebedeeURL)
	return zebedeeClient
}

// GetZebedee return searchreindexapi client
func (e *ExternalServiceList) GetSearchReindexAPI(cfg *config.Config, httpClient dphttp.Clienter) clients.SearchReindexClient {
	searchReindexClient := e.Init.DoGetSearchReindexClient(cfg, httpClient)
	e.SearchReindexAPICli = true
	return searchReindexClient
}

// DoGetSearchReindexClient gets and initialises the SearchReindex Client
func (e *Init) DoGetSearchReindexClient(cfg *config.Config, httpClient dphttp.Clienter) clients.SearchReindexClient {
	healthClient := apihealthcheck.NewClientWithClienter("dp-search-data-finder", cfg.SearchReindexURL, httpClient)
	searchReindexClient := dpsearchreindex.NewClientWithHealthcheck(cfg.ServiceAuthToken, healthClient)
	return searchReindexClient
}

// DoGetKafkaConsumer returns a Kafka Consumer group
func (e *Init) DoGetKafkaConsumer(ctx context.Context, kafkaCfg *config.KafkaConfig) (dpkafka.IConsumerGroup, error) {
	kafkaOffset := dpkafka.OffsetNewest
	if kafkaCfg.OffsetOldest {
		kafkaOffset = dpkafka.OffsetOldest
	}
	cgConfig := &dpkafka.ConsumerGroupConfig{
		KafkaVersion: &kafkaCfg.Version,
		Offset:       &kafkaOffset,
		Topic:        kafkaCfg.ReindexRequestedTopic,
		GroupName:    kafkaCfg.ReindexRequestedGroup,
		BrokerAddrs:  kafkaCfg.Brokers,
	}
	if kafkaCfg.SecProtocol == config.KafkaTLSProtocolFlag {
		cgConfig.SecurityConfig = dpkafka.GetSecurityConfig(
			kafkaCfg.SecCACerts,
			kafkaCfg.SecClientCert,
			kafkaCfg.SecClientKey,
			kafkaCfg.SecSkipVerify,
		)
	}
	kafkaConsumer, err := dpkafka.NewConsumerGroup(
		ctx,
		cgConfig,
	)
	if err != nil {
		return nil, err
	}

	return kafkaConsumer, nil
}

// DoGetHealthCheck creates a healthcheck with versionInfo
func (e *Init) DoGetHealthCheck(cfg *config.Config, buildTime, gitCommit, version string) (HealthChecker, error) {
	versionInfo, err := healthcheck.NewVersionInfo(buildTime, gitCommit, version)
	if err != nil {
		return nil, err
	}
	hc := healthcheck.New(versionInfo, cfg.HealthCheckCriticalTimeout, cfg.HealthCheckInterval)
	return &hc, nil
}
