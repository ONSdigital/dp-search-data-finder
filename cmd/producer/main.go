package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/dp-search-data-finder/config"
	"github.com/ONSdigital/dp-search-data-finder/models"
	"github.com/ONSdigital/dp-search-data-finder/schema"
	"github.com/ONSdigital/log.go/v2/log"
)

const serviceName = "dp-search-data-finder"

func main() {
	log.Namespace = serviceName
	ctx := context.Background()

	// Get Config
	cfg, err := config.Get()
	if err != nil {
		log.Fatal(ctx, "error getting config", err)
		os.Exit(1)
	}

	// Create Kafka ContentUpdatedProducer
	pConfig := &kafka.ProducerConfig{
		BrokerAddrs:  cfg.KafkaConfig.Brokers,
		Topic:        cfg.KafkaConfig.ContentUpdatedTopic,
		KafkaVersion: &cfg.KafkaConfig.Version,
	}
	// KafkaTLSProtocol informs service to use TLS protocol for kafka
	if cfg.KafkaConfig.SecProtocol == config.KafkaTLSProtocolFlag {
		pConfig.SecurityConfig = kafka.GetSecurityConfig(
			cfg.KafkaConfig.SecCACerts,
			cfg.KafkaConfig.SecClientCert,
			cfg.KafkaConfig.SecClientKey,
			cfg.KafkaConfig.SecSkipVerify,
		)
	}

	kafkaProducer, err := kafka.NewProducer(ctx, pConfig)
	if err != nil {
		log.Fatal(ctx, "fatal error trying to create kafka producer", err, log.Data{"topic": cfg.KafkaConfig.ContentUpdatedTopic})
		os.Exit(1)
	}

	// kafka error logging go-routines
	kafkaProducer.LogErrors(ctx)

	time.Sleep(500 * time.Millisecond)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		e := scanEvent(scanner)
		log.Info(ctx, "sending content-updated event", log.Data{"contentUpdatedEvent": e})

		bytes, err := schema.ContentUpdatedEvent.Marshal(e)
		if err != nil {
			log.Fatal(ctx, "content-updated event error", err)
			os.Exit(1)
		}

		// Send bytes to Output channel, after calling Initialise just in case it is not initialised.
		if err := kafkaProducer.Initialise(ctx); err != nil {
			log.Warn(ctx, "failed to initialise kafka producer")
			return
		}
		kafkaProducer.Channels().Output <- bytes
	}
}

// scanEvent creates a ContentUpdated event according to the user input
func scanEvent(scanner *bufio.Scanner) *models.ContentUpdated {
	fmt.Println("--- [Send Kafka ContentUpdated] ---")

	fmt.Println("Please type the URI like /help")
	fmt.Printf("$ ")
	scanner.Scan()
	uri := scanner.Text()

	return &models.ContentUpdated{
		URI:          uri,
		DataType:     "legacy",
		CollectionID: "collectionID",
		JobID:        "",
		TraceID:      "054435ded",
		SearchIndex:  "ONS",
	}
}
