package service

import (
	"context"
	"time"

	"github.com/ONSdigital/dp-api-clients-go/v2/health"
	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/dp-search-data-finder/config"
	"github.com/ONSdigital/dp-search-data-finder/event"
	"github.com/ONSdigital/dp-search-data-finder/handler"
	"github.com/ONSdigital/dp-search-data-finder/schema"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// Service contains all the configs, server and clients to run the event handler service
type Service struct {
	server          HTTPServer
	router          *mux.Router
	serviceList     *ExternalServiceList
	healthCheck     HealthChecker
	consumer        kafka.IConsumerGroup
	shutdownTimeout time.Duration
}

// Run the service
func Run(ctx context.Context, cfg *config.Config, serviceList *ExternalServiceList, buildTime, gitCommit, version string, svcErrors chan error) (*Service, error) {
	log.Info(ctx, "running service")

	// Get HTTP Server with collectionID checkHeader middleware
	r := mux.NewRouter()
	s := serviceList.GetHTTPServer(cfg.BindAddr, r)

	// Get health client for api router
	routerHealthClient := serviceList.GetHealthClient("api-router", cfg.APIRouterURL)

	// Get the zebedee client
	zebedeeClient := serviceList.GetZebedee(cfg)

	// Get dataset-api client
	datasetAPIClient := serviceList.GetDatasetAPI(routerHealthClient)

	// Get Kafka consumer
	consumer, err := serviceList.GetKafkaConsumer(ctx, cfg)
	if err != nil {
		log.Fatal(ctx, "failed to initialise kafka consumer", err)
		return nil, err
	}

	// Get Kafka producer
	producer, err := serviceList.GetKafkaProducer(ctx, cfg)
	if err != nil {
		log.Fatal(ctx, "failed to initialise kafka producer", err)
		return nil, err
	}

	// Get Kafka producer for reindex task counts
	producerForReindexTaskCounts, err := serviceList.GetKafkaProducerForReindexTaskCounts(ctx, cfg)
	if err != nil {
		log.Fatal(ctx, "failed to initialise kafka producer", err)
		return nil, err
	}

	contentUpdatedProducer := event.ContentUpdatedProducer{
		Marshaller: schema.ContentUpdatedEvent,
		Producer:   producer,
	}

	reindexTaskCountsProducer := event.ReindexTaskCountsProducer{
		Marshaller: schema.ReindexTaskCounts,
		Producer:   producerForReindexTaskCounts,
	}

	// Event Handler for Kafka Consumer
	eventhandler := &handler.ReindexRequestedHandler{
		Config:                    cfg,
		ZebedeeCli:                zebedeeClient,
		DatasetAPICli:             datasetAPIClient,
		ContentUpdatedProducer:    contentUpdatedProducer,
		ReindexTaskCountsProducer: reindexTaskCountsProducer,
	}

	event.Consume(ctx, consumer, eventhandler, cfg)
	if consumerStartErr := consumer.Start(); consumerStartErr != nil {
		log.Fatal(ctx, "error starting the consumer", consumerStartErr)
		return nil, consumerStartErr
	}

	// Get HealthCheck
	hc, err := serviceList.GetHealthCheck(cfg, buildTime, gitCommit, version)
	if err != nil {
		log.Fatal(ctx, "could not instantiate healthcheck", err)
		return nil, err
	}

	if err := registerCheckers(ctx, hc, consumer, producer, producerForReindexTaskCounts, routerHealthClient); err != nil {
		return nil, errors.Wrap(err, "unable to register checkers")
	}

	r.StrictSlash(true).Path("/health").HandlerFunc(hc.Handler)
	hc.Start(ctx)

	// Run the http server in a new go-routine
	go func() {
		if err := s.ListenAndServe(); err != nil {
			svcErrors <- errors.Wrap(err, "failure in http listen and serve")
		}
	}()

	return &Service{
		server:          s,
		router:          r,
		serviceList:     serviceList,
		healthCheck:     hc,
		consumer:        consumer,
		shutdownTimeout: cfg.GracefulShutdownTimeout,
	}, nil
}

// Close gracefully shuts the service down in the required order, with timeout
func (svc *Service) Close(ctx context.Context) error {
	timeout := svc.shutdownTimeout
	log.Info(ctx, "commencing graceful shutdown", log.Data{"graceful_shutdown_timeout": timeout})
	ctx, cancel := context.WithTimeout(ctx, timeout)

	// track shutdown gracefully closes up
	var gracefulShutdown bool

	go func() {
		defer cancel()
		var hasShutdownError bool

		// stop healthcheck, as it depends on everything else
		if svc.serviceList.HealthCheck {
			svc.healthCheck.Stop()
		}

		// If kafka consumer exists, stop listening to it.
		// This will automatically stop the event consumer loops and no more messages will be processed.
		// The kafka consumer will be closed after the service shuts down.
		if svc.serviceList.KafkaConsumer {
			log.Info(ctx, "stopping kafka consumer listener")
			if err := svc.consumer.Stop(); err != nil {
				log.Error(ctx, "error stopping kafka consumer listener", err)
				hasShutdownError = true
			}
			log.Info(ctx, "stopped kafka consumer listener")
		}

		// stop any incoming requests before closing any outbound connections
		if err := svc.server.Shutdown(ctx); err != nil {
			log.Error(ctx, "failed to shutdown http server", err)
			hasShutdownError = true
		}

		// If kafka consumer exists, close it.
		if svc.serviceList.KafkaConsumer {
			log.Info(ctx, "closing kafka consumer")
			if err := svc.consumer.Close(ctx); err != nil {
				log.Error(ctx, "error closing kafka consumer", err)
				hasShutdownError = true
			}
			log.Info(ctx, "closed kafka consumer")
		}

		if !hasShutdownError {
			gracefulShutdown = true
		}
	}()

	// wait for shutdown success (via cancel) or failure (timeout)
	<-ctx.Done()

	if !gracefulShutdown {
		err := errors.New("failed to shutdown gracefully")
		log.Error(ctx, "failed to shutdown gracefully ", err)
		return err
	}

	log.Info(ctx, "graceful shutdown was successful")
	return nil
}

func registerCheckers(ctx context.Context, hc HealthChecker, consumer kafka.IConsumerGroup, contentUpdatedProducer, reindexTaskCounts kafka.IProducer, routerHealthClient *health.Client) error {
	hasErrors := false

	if err := hc.AddCheck("API router", routerHealthClient.Checker); err != nil {
		hasErrors = true
		log.Error(ctx, "error adding check for api-router", err)
	}
	if err := hc.AddCheck("Kafka consumer", consumer.Checker); err != nil {
		hasErrors = true
		log.Error(ctx, "error adding check for kafka", err)
	}
	if err := hc.AddCheck("Kafka content updated producer", contentUpdatedProducer.Checker); err != nil {
		hasErrors = true
		log.Error(ctx, "error adding check for kafka", err)
	}
	if err := hc.AddCheck("Kafka reindex task counts producer", reindexTaskCounts.Checker); err != nil {
		hasErrors = true
		log.Error(ctx, "error adding check for kafka", err)
	}
	if hasErrors {
		return errors.New("Error(s) registering checkers for healthcheck")
	}

	return nil
}
