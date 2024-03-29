package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

const KafkaTLSProtocolFlag = "TLS"

// Config represents service configuration for dp-search-data-finder
type Config struct {
	APIRouterURL                     string        `envconfig:"API_ROUTER_URL"`
	BindAddr                         string        `envconfig:"BIND_ADDR"`
	GracefulShutdownTimeout          time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckCriticalTimeout       time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	HealthCheckInterval              time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	KafkaConfig                      KafkaConfig
	ServiceAuthToken                 string        `envconfig:"SERVICE_AUTH_TOKEN"   json:"-"`
	ZebedeeClientTimeout             time.Duration `envconfig:"ZEBEDEE_CLIENT_TIMEOUT"`
	EnablePublishContentUpdatedTopic bool          `envconfig:"ENABLE_PUBLISH_CONTENT_UPDATED_TOPIC"`
	EnableReindexTaskCounts          bool          `envconfig:"ENABLE_REINDEX_TASK_COUNTS"`
	ZebedeeRequestLimit              int           `envconfig:"ZEBEDEE_REQUEST_LIMIT"`
	ZebedeeURL                       string        `envconfig:"ZEBEDEE_URL"`
	TaskNameValues                   string        `envconfig:"TASK_NAME_VALUES"`
}

// KafkaConfig contains the config required to connect to Kafka
type KafkaConfig struct {
	Brokers                []string `envconfig:"KAFKA_ADDR"`
	ContentUpdatedTopic    string   `envconfig:"KAFKA_CONTENT_UPDATED_TOPIC"`
	ReindexTaskCountsTopic string   `envconfig:"KAFKA_REINDEX_TASK_COUNTS_TOPIC"`
	ConsumerGroup          string   `envconfig:"KAFKA_CONSUMER_GROUP"`
	NumWorkers             int      `envconfig:"KAFKA_NUM_WORKERS"`
	OffsetOldest           bool     `envconfig:"KAFKA_OFFSET_OLDEST"`
	ReindexRequestedTopic  string   `envconfig:"KAFKA_REINDEX_REQUESTED_TOPIC"`
	SecProtocol            string   `envconfig:"KAFKA_SEC_PROTO"`
	SecCACerts             string   `envconfig:"KAFKA_SEC_CA_CERTS"`
	SecClientCert          string   `envconfig:"KAFKA_SEC_CLIENT_CERT"`
	SecClientKey           string   `envconfig:"KAFKA_SEC_CLIENT_KEY"    json:"-"`
	SecSkipVerify          bool     `envconfig:"KAFKA_SEC_SKIP_VERIFY"`
	Version                string   `envconfig:"KAFKA_VERSION"`
}

var cfg *Config

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		APIRouterURL:                     "http://localhost:23200/v1",
		BindAddr:                         "localhost:28000",
		EnablePublishContentUpdatedTopic: false,
		EnableReindexTaskCounts:          false,
		GracefulShutdownTimeout:          5 * time.Second,
		HealthCheckCriticalTimeout:       90 * time.Second,
		HealthCheckInterval:              30 * time.Second,
		KafkaConfig: KafkaConfig{
			Brokers:                []string{"localhost:9092", "localhost:9093", "localhost:9094"},
			ContentUpdatedTopic:    "content-updated",
			ReindexTaskCountsTopic: "reindex-task-counts",
			ConsumerGroup:          "dp-search-data-finder",
			NumWorkers:             1,
			OffsetOldest:           true,
			ReindexRequestedTopic:  "reindex-requested",
			SecProtocol:            "",
			SecCACerts:             "",
			SecClientCert:          "",
			SecClientKey:           "",
			SecSkipVerify:          false,
			Version:                "1.0.2",
		},
		ServiceAuthToken:     "",
		ZebedeeClientTimeout: 30 * time.Second,
		ZebedeeRequestLimit:  10,
		ZebedeeURL:           "http://localhost:8082",
		TaskNameValues:       "dataset-api,zebedee",
	}

	return cfg, envconfig.Process("", cfg)
}
