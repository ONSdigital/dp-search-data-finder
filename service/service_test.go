package service_test

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/health"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/dp-kafka/v3/kafkatest"
	"github.com/ONSdigital/dp-search-data-finder/clients"
	clientMock "github.com/ONSdigital/dp-search-data-finder/clients/mock"
	"github.com/ONSdigital/dp-search-data-finder/config"
	"github.com/ONSdigital/dp-search-data-finder/service"
	serviceMock "github.com/ONSdigital/dp-search-data-finder/service/mock"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	ctx           = context.Background()
	testBuildTime = "BuildTime"
	testGitCommit = "GitCommit"
	testVersion   = "Version"
)

var (
	errKafkaConsumer = errors.New("Kafka consumer error")
	errHealthcheck   = errors.New("healthCheck error")
)

var funcDoGetKafkaConsumerErr = func(ctx context.Context, kafkaCfg *config.KafkaConfig) (kafka.IConsumerGroup, error) {
	return nil, errKafkaConsumer
}

var funcDoGetHealthcheckErr = func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
	return nil, errHealthcheck
}

var funcDoGetHTTPServerNil = func(bindAddr string, router http.Handler) service.HTTPServer {
	return nil
}

func TestRun(t *testing.T) {
	testCfg, err := config.Get()
	if err != nil {
		t.Errorf("failed to retrieve default configuration, error is: %v", err)
	}

	Convey("Having a set of mocked dependencies", t, func() {
		consumerMock := &kafkatest.IConsumerGroupMock{
			CheckerFunc:  func(ctx context.Context, state *healthcheck.CheckState) error { return nil },
			ChannelsFunc: func() *kafka.ConsumerGroupChannels { return &kafka.ConsumerGroupChannels{} },
			StartFunc:    func() error { return nil },
		}

		producerMock := &kafkatest.IProducerMock{
			CheckerFunc:  func(ctx context.Context, state *healthcheck.CheckState) error { return nil },
			ChannelsFunc: func() *kafka.ProducerChannels { return &kafka.ProducerChannels{} },
		}

		hcMock := &serviceMock.HealthCheckerMock{
			AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
			StartFunc:    func(ctx context.Context) {},
		}

		serverWg := &sync.WaitGroup{}
		serverMock := &serviceMock.HTTPServerMock{
			ListenAndServeFunc: func() error {
				serverWg.Done()
				return nil
			},
		}

		zebedeeMock := &clientMock.ZebedeeClientMock{
			CheckerFunc: func(ctx context.Context, state *healthcheck.CheckState) error { return nil },
		}

		datasetAPIMock := &clientMock.DatasetAPIClientMock{
			CheckerFunc: func(ctx context.Context, state *healthcheck.CheckState) error { return nil },
		}

		funcDoGetKafkaConsumerOk := func(ctx context.Context, kafkaCfg *config.KafkaConfig) (kafka.IConsumerGroup, error) {
			return consumerMock, nil
		}

		funcDoGetKafkaProducerOk := func(ctx context.Context, config *config.Config) (kafka.IProducer, error) {
			return producerMock, nil
		}

		funcDoGetHealthcheckOk := func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
			return hcMock, nil
		}

		funcDoGetHealthClientOk := func(name, url string) *health.Client {
			return &health.Client{}
		}

		funcDoGetHTTPServer := func(bindAddr string, router http.Handler) service.HTTPServer {
			return serverMock
		}

		funcDoGetZebedeeOk := func(cfg *config.Config) clients.ZebedeeClient {
			return zebedeeMock
		}

		funcDoGetDatasetAPIOk := func(hcCli *health.Client) clients.DatasetAPIClient {
			return datasetAPIMock
		}

		Convey("Given that initialising Kafka consumer returns an error", func() {
			initMock := &serviceMock.InitialiserMock{
				DoGetHealthClientFunc:     funcDoGetHealthClientOk,
				DoGetHTTPServerFunc:       funcDoGetHTTPServerNil,
				DoGetKafkaConsumerFunc:    funcDoGetKafkaConsumerErr,
				DoGetZebedeeClientFunc:    funcDoGetZebedeeOk,
				DoGetDatasetAPIClientFunc: funcDoGetDatasetAPIOk,
			}
			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			_, err := service.Run(ctx, testCfg, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)

			Convey("Then service Run fails with the same error and the flag is not set", func() {
				So(err, ShouldResemble, errKafkaConsumer)
				So(svcList.KafkaConsumer, ShouldBeFalse)
				So(svcList.HealthCheck, ShouldBeFalse)
				So(svcList.ZebedeeCli, ShouldBeTrue)
				So(svcList.DatasetAPICli, ShouldBeTrue)
			})
		})

		Convey("Given that initialising healthcheck returns an error", func() {
			initMock := &serviceMock.InitialiserMock{
				DoGetHTTPServerFunc:                        funcDoGetHTTPServerNil,
				DoGetHealthCheckFunc:                       funcDoGetHealthcheckErr,
				DoGetHealthClientFunc:                      funcDoGetHealthClientOk,
				DoGetKafkaConsumerFunc:                     funcDoGetKafkaConsumerOk,
				DoGetKafkaProducerFunc:                     funcDoGetKafkaProducerOk,
				DoGetKafkaProducerForReindexTaskCountsFunc: funcDoGetKafkaProducerOk,
				DoGetZebedeeClientFunc:                     funcDoGetZebedeeOk,
				DoGetDatasetAPIClientFunc:                  funcDoGetDatasetAPIOk,
			}
			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			_, err := service.Run(ctx, testCfg, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)

			Convey("Then service Run fails with the same error and the flag is not set", func() {
				So(err, ShouldResemble, errHealthcheck)
				So(svcList.KafkaConsumer, ShouldBeTrue)
				So(svcList.HealthCheck, ShouldBeFalse)
				So(svcList.ZebedeeCli, ShouldBeTrue)
				So(svcList.DatasetAPICli, ShouldBeTrue)
			})
		})

		Convey("Given that all dependencies are successfully initialised", func() {
			initMock := &serviceMock.InitialiserMock{
				DoGetHTTPServerFunc:                        funcDoGetHTTPServer,
				DoGetHealthCheckFunc:                       funcDoGetHealthcheckOk,
				DoGetHealthClientFunc:                      funcDoGetHealthClientOk,
				DoGetKafkaConsumerFunc:                     funcDoGetKafkaConsumerOk,
				DoGetKafkaProducerFunc:                     funcDoGetKafkaProducerOk,
				DoGetKafkaProducerForReindexTaskCountsFunc: funcDoGetKafkaProducerOk,
				DoGetZebedeeClientFunc:                     funcDoGetZebedeeOk,
				DoGetDatasetAPIClientFunc:                  funcDoGetDatasetAPIOk,
			}
			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			serverWg.Add(1)
			_, err := service.Run(ctx, testCfg, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)

			Convey("Then service Run succeeds and all the flags are set", func() {
				So(err, ShouldBeNil)
				So(svcList.KafkaConsumer, ShouldBeTrue)
				So(svcList.HealthCheck, ShouldBeTrue)
				So(svcList.ZebedeeCli, ShouldBeTrue)
				So(svcList.DatasetAPICli, ShouldBeTrue)
			})

			Convey("The checkers are registered and the healthcheck and http server started", func() {
				So(len(hcMock.AddCheckCalls()), ShouldEqual, 4)
				So(hcMock.AddCheckCalls()[0].Name, ShouldResemble, "API router")
				So(hcMock.AddCheckCalls()[1].Name, ShouldResemble, "Kafka consumer")
				So(hcMock.AddCheckCalls()[2].Name, ShouldResemble, "Kafka content updated producer")
				So(hcMock.AddCheckCalls()[3].Name, ShouldResemble, "Kafka reindex task counts producer")
				So(len(initMock.DoGetHTTPServerCalls()), ShouldEqual, 1)
				So(initMock.DoGetHTTPServerCalls()[0].BindAddr, ShouldEqual, "localhost:28000")
				So(len(hcMock.StartCalls()), ShouldEqual, 1)
				serverWg.Wait() // Wait for HTTP server go-routine to finish
				So(len(serverMock.ListenAndServeCalls()), ShouldEqual, 1)
			})
		})

		Convey("Given that Checkers cannot be registered", func() {
			errAddheckFail := errors.New("Error(s) registering checkers for healthcheck")
			hcMockAddFail := &serviceMock.HealthCheckerMock{
				AddCheckFunc: func(name string, checker healthcheck.Checker) error { return errAddheckFail },
				StartFunc:    func(ctx context.Context) {},
			}

			initMock := &serviceMock.InitialiserMock{
				DoGetHTTPServerFunc: funcDoGetHTTPServerNil,
				DoGetHealthCheckFunc: func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
					return hcMockAddFail, nil
				},
				DoGetHealthClientFunc:                      funcDoGetHealthClientOk,
				DoGetKafkaConsumerFunc:                     funcDoGetKafkaConsumerOk,
				DoGetKafkaProducerFunc:                     funcDoGetKafkaProducerOk,
				DoGetKafkaProducerForReindexTaskCountsFunc: funcDoGetKafkaProducerOk,
				DoGetZebedeeClientFunc:                     funcDoGetZebedeeOk,
				DoGetDatasetAPIClientFunc:                  funcDoGetDatasetAPIOk,
			}
			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			_, err := service.Run(ctx, testCfg, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)

			Convey("Then service Run fails, but all checks try to register", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldResemble, fmt.Sprintf("unable to register checkers: %s", errAddheckFail.Error()))
				So(svcList.HealthCheck, ShouldBeTrue)
				So(svcList.KafkaConsumer, ShouldBeTrue)
				So(svcList.ZebedeeCli, ShouldBeTrue)
				So(svcList.DatasetAPICli, ShouldBeTrue)
				So(len(hcMockAddFail.AddCheckCalls()), ShouldEqual, 4)
				So(hcMockAddFail.AddCheckCalls()[0].Name, ShouldResemble, "API router")
				So(hcMockAddFail.AddCheckCalls()[1].Name, ShouldResemble, "Kafka consumer")
				So(hcMockAddFail.AddCheckCalls()[2].Name, ShouldResemble, "Kafka content updated producer")
				So(hcMockAddFail.AddCheckCalls()[3].Name, ShouldResemble, "Kafka reindex task counts producer")
			})
		})
	})
}

func TestClose(t *testing.T) {
	testCfg, err := config.Get()
	if err != nil {
		t.Errorf("failed to retrieve default configuration, error is: %v", err)
	}

	datasetAPIMock := &clientMock.DatasetAPIClientMock{
		CheckerFunc: func(ctx context.Context, state *healthcheck.CheckState) error { return nil },
	}

	funcDoGetDatasetAPIOk := func(hcCli *health.Client) clients.DatasetAPIClient {
		return datasetAPIMock
	}

	Convey("Having a correctly initialised service", t, func() {
		hcStopped := false

		consumerMock := &kafkatest.IConsumerGroupMock{
			StopFunc:     func() error { return nil },
			CloseFunc:    func(ctx context.Context, optFuncs ...kafka.OptFunc) error { return nil },
			CheckerFunc:  func(ctx context.Context, state *healthcheck.CheckState) error { return nil },
			ChannelsFunc: func() *kafka.ConsumerGroupChannels { return &kafka.ConsumerGroupChannels{} },
			StartFunc:    func() error { return nil },
		}

		producerMock := &kafkatest.IProducerMock{
			CheckerFunc:  func(ctx context.Context, state *healthcheck.CheckState) error { return nil },
			ChannelsFunc: func() *kafka.ProducerChannels { return &kafka.ProducerChannels{} },
		}

		// healthcheck Stop does not depend on any other service being closed/stopped
		hcMock := &serviceMock.HealthCheckerMock{
			AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
			StartFunc:    func(ctx context.Context) {},
			StopFunc:     func() { hcStopped = true },
		}

		// server Shutdown will fail if healthcheck is not stopped
		serverMock := &serviceMock.HTTPServerMock{
			ListenAndServeFunc: func() error { return nil },
			ShutdownFunc: func(ctx context.Context) error {
				if !hcStopped {
					return errors.New("Server stopped before healthcheck")
				}
				return nil
			},
		}

		zebedeeMock := &clientMock.ZebedeeClientMock{
			CheckerFunc: func(ctx context.Context, state *healthcheck.CheckState) error { return nil },
		}

		Convey("Closing the service results in all the dependencies being closed in the expected order", func() {
			initMock := &serviceMock.InitialiserMock{
				DoGetHTTPServerFunc: func(bindAddr string, router http.Handler) service.HTTPServer { return serverMock },
				DoGetHealthCheckFunc: func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
					return hcMock, nil
				},
				DoGetHealthClientFunc: func(name, url string) *health.Client { return &health.Client{} },
				DoGetKafkaConsumerFunc: func(ctx context.Context, kafkaCfg *config.KafkaConfig) (kafka.IConsumerGroup, error) {
					return consumerMock, nil
				},
				DoGetKafkaProducerFunc: func(ctx context.Context, config *config.Config) (kafka.IProducer, error) {
					return producerMock, nil
				},
				DoGetKafkaProducerForReindexTaskCountsFunc: func(ctx context.Context, config *config.Config) (kafka.IProducer, error) {
					return producerMock, nil
				},
				DoGetZebedeeClientFunc:    func(cfg *config.Config) clients.ZebedeeClient { return zebedeeMock },
				DoGetDatasetAPIClientFunc: funcDoGetDatasetAPIOk,
			}

			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			svc, err := service.Run(ctx, testCfg, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)
			So(err, ShouldBeNil)

			err = svc.Close(context.Background())
			So(err, ShouldBeNil)
			So(len(consumerMock.StopCalls()), ShouldEqual, 1)
			So(len(hcMock.StopCalls()), ShouldEqual, 1)
			So(len(consumerMock.CloseCalls()), ShouldEqual, 1)
			So(len(serverMock.ShutdownCalls()), ShouldEqual, 1)
		})
		Convey("If services fail to stop, the Close operation tries to close all dependencies and returns an error", func() {
			failingserverMock := &serviceMock.HTTPServerMock{
				ListenAndServeFunc: func() error { return nil },
				ShutdownFunc: func(ctx context.Context) error {
					return errors.New("Failed to stop http server")
				},
			}

			initMock := &serviceMock.InitialiserMock{
				DoGetHTTPServerFunc: func(bindAddr string, router http.Handler) service.HTTPServer { return failingserverMock },
				DoGetHealthCheckFunc: func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
					return hcMock, nil
				},
				DoGetHealthClientFunc: func(name, url string) *health.Client { return &health.Client{} },
				DoGetKafkaConsumerFunc: func(ctx context.Context, kafkaCfg *config.KafkaConfig) (kafka.IConsumerGroup, error) {
					return consumerMock, nil
				},
				DoGetKafkaProducerFunc: func(ctx context.Context, config *config.Config) (kafka.IProducer, error) {
					return producerMock, nil
				},
				DoGetKafkaProducerForReindexTaskCountsFunc: func(ctx context.Context, config *config.Config) (kafka.IProducer, error) {
					return producerMock, nil
				},
				DoGetZebedeeClientFunc:    func(cfg *config.Config) clients.ZebedeeClient { return zebedeeMock },
				DoGetDatasetAPIClientFunc: funcDoGetDatasetAPIOk,
			}

			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			svc, err := service.Run(ctx, testCfg, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)
			So(err, ShouldBeNil)

			err = svc.Close(context.Background())
			So(err, ShouldNotBeNil)
			So(len(hcMock.StopCalls()), ShouldEqual, 1)
			So(len(failingserverMock.ShutdownCalls()), ShouldEqual, 1)
		})
	})
}
