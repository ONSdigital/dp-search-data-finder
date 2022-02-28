package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

const KafkaTLSProtocolFlag = "TLS"

// Config represents service configuration for dp-search-data-finder
type Config struct {
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	OutputFilePath             string        `envconfig:"OUTPUT_FILE_PATH"`
	KafkaConfig                KafkaConfig
}

// KafkaConfig contains the config required to connect to Kafka
type KafkaConfig struct {
	Brokers               []string `envconfig:"KAFKA_ADDR"`
	Version               string   `envconfig:"KAFKA_VERSION"`
	OffsetOldest          bool     `envconfig:"KAFKA_OFFSET_OLDEST"`
	SecProtocol           string   `envconfig:"KAFKA_SEC_PROTO"`
	SecCACerts            string   `envconfig:"KAFKA_SEC_CA_CERTS"`
	SecClientKey          string   `envconfig:"KAFKA_SEC_CLIENT_KEY"    json:"-"`
	SecClientCert         string   `envconfig:"KAFKA_SEC_CLIENT_CERT"`
	SecSkipVerify         bool     `envconfig:"KAFKA_SEC_SKIP_VERIFY"`
	NumWorkers            int      `envconfig:"KAFKA_NUM_WORKERS"`
	ReindexRequestedTopic string   `envconfig:"KAFKA_REINDEX_REQUESTED_TOPIC"`
	ContentUpdatedTopic   string   `envconfig:"KAFKA_CONTENT_UPDATED_TOPIC"`
}

var cfg *Config

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		BindAddr:                   "localhost:28000",
		GracefulShutdownTimeout:    5 * time.Second,
		HealthCheckInterval:        30 * time.Second,
		HealthCheckCriticalTimeout: 90 * time.Second,
		OutputFilePath:             "/tmp/helloworld.txt",
		KafkaConfig: KafkaConfig{
			Brokers:               []string{"localhost:9092"},
			Version:               "1.0.2",
			OffsetOldest:          true,
			NumWorkers:            1,
			ReindexRequestedTopic: "reindex-requested",
			ContentUpdatedTopic:   "content-updated",
		},
	}

	return cfg, envconfig.Process("", cfg)
}