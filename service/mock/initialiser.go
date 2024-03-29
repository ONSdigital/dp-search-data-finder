// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"context"
	"github.com/ONSdigital/dp-api-clients-go/v2/health"
	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/dp-search-data-finder/clients"
	"github.com/ONSdigital/dp-search-data-finder/config"
	"github.com/ONSdigital/dp-search-data-finder/service"
	"net/http"
	"sync"
)

// Ensure, that InitialiserMock does implement service.Initialiser.
// If this is not the case, regenerate this file with moq.
var _ service.Initialiser = &InitialiserMock{}

// InitialiserMock is a mock implementation of service.Initialiser.
//
//	func TestSomethingThatUsesInitialiser(t *testing.T) {
//
//		// make and configure a mocked service.Initialiser
//		mockedInitialiser := &InitialiserMock{
//			DoGetDatasetAPIClientFunc: func(hcCli *health.Client) clients.DatasetAPIClient {
//				panic("mock out the DoGetDatasetAPIClient method")
//			},
//			DoGetHTTPServerFunc: func(bindAddr string, router http.Handler) service.HTTPServer {
//				panic("mock out the DoGetHTTPServer method")
//			},
//			DoGetHealthCheckFunc: func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
//				panic("mock out the DoGetHealthCheck method")
//			},
//			DoGetHealthClientFunc: func(name string, url string) *health.Client {
//				panic("mock out the DoGetHealthClient method")
//			},
//			DoGetKafkaConsumerFunc: func(ctx context.Context, kafkaCfg *config.KafkaConfig) (kafka.IConsumerGroup, error) {
//				panic("mock out the DoGetKafkaConsumer method")
//			},
//			DoGetKafkaProducerFunc: func(ctx context.Context, cfg *config.Config) (kafka.IProducer, error) {
//				panic("mock out the DoGetKafkaProducer method")
//			},
//			DoGetKafkaProducerForReindexTaskCountsFunc: func(ctx context.Context, cfg *config.Config) (kafka.IProducer, error) {
//				panic("mock out the DoGetKafkaProducerForReindexTaskCounts method")
//			},
//			DoGetZebedeeClientFunc: func(cfg *config.Config) clients.ZebedeeClient {
//				panic("mock out the DoGetZebedeeClient method")
//			},
//		}
//
//		// use mockedInitialiser in code that requires service.Initialiser
//		// and then make assertions.
//
//	}
type InitialiserMock struct {
	// DoGetDatasetAPIClientFunc mocks the DoGetDatasetAPIClient method.
	DoGetDatasetAPIClientFunc func(hcCli *health.Client) clients.DatasetAPIClient

	// DoGetHTTPServerFunc mocks the DoGetHTTPServer method.
	DoGetHTTPServerFunc func(bindAddr string, router http.Handler) service.HTTPServer

	// DoGetHealthCheckFunc mocks the DoGetHealthCheck method.
	DoGetHealthCheckFunc func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error)

	// DoGetHealthClientFunc mocks the DoGetHealthClient method.
	DoGetHealthClientFunc func(name string, url string) *health.Client

	// DoGetKafkaConsumerFunc mocks the DoGetKafkaConsumer method.
	DoGetKafkaConsumerFunc func(ctx context.Context, kafkaCfg *config.KafkaConfig) (kafka.IConsumerGroup, error)

	// DoGetKafkaProducerFunc mocks the DoGetKafkaProducer method.
	DoGetKafkaProducerFunc func(ctx context.Context, cfg *config.Config) (kafka.IProducer, error)

	// DoGetKafkaProducerForReindexTaskCountsFunc mocks the DoGetKafkaProducerForReindexTaskCounts method.
	DoGetKafkaProducerForReindexTaskCountsFunc func(ctx context.Context, cfg *config.Config) (kafka.IProducer, error)

	// DoGetZebedeeClientFunc mocks the DoGetZebedeeClient method.
	DoGetZebedeeClientFunc func(cfg *config.Config) clients.ZebedeeClient

	// calls tracks calls to the methods.
	calls struct {
		// DoGetDatasetAPIClient holds details about calls to the DoGetDatasetAPIClient method.
		DoGetDatasetAPIClient []struct {
			// HcCli is the hcCli argument value.
			HcCli *health.Client
		}
		// DoGetHTTPServer holds details about calls to the DoGetHTTPServer method.
		DoGetHTTPServer []struct {
			// BindAddr is the bindAddr argument value.
			BindAddr string
			// Router is the router argument value.
			Router http.Handler
		}
		// DoGetHealthCheck holds details about calls to the DoGetHealthCheck method.
		DoGetHealthCheck []struct {
			// Cfg is the cfg argument value.
			Cfg *config.Config
			// BuildTime is the buildTime argument value.
			BuildTime string
			// GitCommit is the gitCommit argument value.
			GitCommit string
			// Version is the version argument value.
			Version string
		}
		// DoGetHealthClient holds details about calls to the DoGetHealthClient method.
		DoGetHealthClient []struct {
			// Name is the name argument value.
			Name string
			// URL is the url argument value.
			URL string
		}
		// DoGetKafkaConsumer holds details about calls to the DoGetKafkaConsumer method.
		DoGetKafkaConsumer []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// KafkaCfg is the kafkaCfg argument value.
			KafkaCfg *config.KafkaConfig
		}
		// DoGetKafkaProducer holds details about calls to the DoGetKafkaProducer method.
		DoGetKafkaProducer []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Cfg is the cfg argument value.
			Cfg *config.Config
		}
		// DoGetKafkaProducerForReindexTaskCounts holds details about calls to the DoGetKafkaProducerForReindexTaskCounts method.
		DoGetKafkaProducerForReindexTaskCounts []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Cfg is the cfg argument value.
			Cfg *config.Config
		}
		// DoGetZebedeeClient holds details about calls to the DoGetZebedeeClient method.
		DoGetZebedeeClient []struct {
			// Cfg is the cfg argument value.
			Cfg *config.Config
		}
	}
	lockDoGetDatasetAPIClient                  sync.RWMutex
	lockDoGetHTTPServer                        sync.RWMutex
	lockDoGetHealthCheck                       sync.RWMutex
	lockDoGetHealthClient                      sync.RWMutex
	lockDoGetKafkaConsumer                     sync.RWMutex
	lockDoGetKafkaProducer                     sync.RWMutex
	lockDoGetKafkaProducerForReindexTaskCounts sync.RWMutex
	lockDoGetZebedeeClient                     sync.RWMutex
}

// DoGetDatasetAPIClient calls DoGetDatasetAPIClientFunc.
func (mock *InitialiserMock) DoGetDatasetAPIClient(hcCli *health.Client) clients.DatasetAPIClient {
	if mock.DoGetDatasetAPIClientFunc == nil {
		panic("InitialiserMock.DoGetDatasetAPIClientFunc: method is nil but Initialiser.DoGetDatasetAPIClient was just called")
	}
	callInfo := struct {
		HcCli *health.Client
	}{
		HcCli: hcCli,
	}
	mock.lockDoGetDatasetAPIClient.Lock()
	mock.calls.DoGetDatasetAPIClient = append(mock.calls.DoGetDatasetAPIClient, callInfo)
	mock.lockDoGetDatasetAPIClient.Unlock()
	return mock.DoGetDatasetAPIClientFunc(hcCli)
}

// DoGetDatasetAPIClientCalls gets all the calls that were made to DoGetDatasetAPIClient.
// Check the length with:
//
//	len(mockedInitialiser.DoGetDatasetAPIClientCalls())
func (mock *InitialiserMock) DoGetDatasetAPIClientCalls() []struct {
	HcCli *health.Client
} {
	var calls []struct {
		HcCli *health.Client
	}
	mock.lockDoGetDatasetAPIClient.RLock()
	calls = mock.calls.DoGetDatasetAPIClient
	mock.lockDoGetDatasetAPIClient.RUnlock()
	return calls
}

// DoGetHTTPServer calls DoGetHTTPServerFunc.
func (mock *InitialiserMock) DoGetHTTPServer(bindAddr string, router http.Handler) service.HTTPServer {
	if mock.DoGetHTTPServerFunc == nil {
		panic("InitialiserMock.DoGetHTTPServerFunc: method is nil but Initialiser.DoGetHTTPServer was just called")
	}
	callInfo := struct {
		BindAddr string
		Router   http.Handler
	}{
		BindAddr: bindAddr,
		Router:   router,
	}
	mock.lockDoGetHTTPServer.Lock()
	mock.calls.DoGetHTTPServer = append(mock.calls.DoGetHTTPServer, callInfo)
	mock.lockDoGetHTTPServer.Unlock()
	return mock.DoGetHTTPServerFunc(bindAddr, router)
}

// DoGetHTTPServerCalls gets all the calls that were made to DoGetHTTPServer.
// Check the length with:
//
//	len(mockedInitialiser.DoGetHTTPServerCalls())
func (mock *InitialiserMock) DoGetHTTPServerCalls() []struct {
	BindAddr string
	Router   http.Handler
} {
	var calls []struct {
		BindAddr string
		Router   http.Handler
	}
	mock.lockDoGetHTTPServer.RLock()
	calls = mock.calls.DoGetHTTPServer
	mock.lockDoGetHTTPServer.RUnlock()
	return calls
}

// DoGetHealthCheck calls DoGetHealthCheckFunc.
func (mock *InitialiserMock) DoGetHealthCheck(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
	if mock.DoGetHealthCheckFunc == nil {
		panic("InitialiserMock.DoGetHealthCheckFunc: method is nil but Initialiser.DoGetHealthCheck was just called")
	}
	callInfo := struct {
		Cfg       *config.Config
		BuildTime string
		GitCommit string
		Version   string
	}{
		Cfg:       cfg,
		BuildTime: buildTime,
		GitCommit: gitCommit,
		Version:   version,
	}
	mock.lockDoGetHealthCheck.Lock()
	mock.calls.DoGetHealthCheck = append(mock.calls.DoGetHealthCheck, callInfo)
	mock.lockDoGetHealthCheck.Unlock()
	return mock.DoGetHealthCheckFunc(cfg, buildTime, gitCommit, version)
}

// DoGetHealthCheckCalls gets all the calls that were made to DoGetHealthCheck.
// Check the length with:
//
//	len(mockedInitialiser.DoGetHealthCheckCalls())
func (mock *InitialiserMock) DoGetHealthCheckCalls() []struct {
	Cfg       *config.Config
	BuildTime string
	GitCommit string
	Version   string
} {
	var calls []struct {
		Cfg       *config.Config
		BuildTime string
		GitCommit string
		Version   string
	}
	mock.lockDoGetHealthCheck.RLock()
	calls = mock.calls.DoGetHealthCheck
	mock.lockDoGetHealthCheck.RUnlock()
	return calls
}

// DoGetHealthClient calls DoGetHealthClientFunc.
func (mock *InitialiserMock) DoGetHealthClient(name string, url string) *health.Client {
	if mock.DoGetHealthClientFunc == nil {
		panic("InitialiserMock.DoGetHealthClientFunc: method is nil but Initialiser.DoGetHealthClient was just called")
	}
	callInfo := struct {
		Name string
		URL  string
	}{
		Name: name,
		URL:  url,
	}
	mock.lockDoGetHealthClient.Lock()
	mock.calls.DoGetHealthClient = append(mock.calls.DoGetHealthClient, callInfo)
	mock.lockDoGetHealthClient.Unlock()
	return mock.DoGetHealthClientFunc(name, url)
}

// DoGetHealthClientCalls gets all the calls that were made to DoGetHealthClient.
// Check the length with:
//
//	len(mockedInitialiser.DoGetHealthClientCalls())
func (mock *InitialiserMock) DoGetHealthClientCalls() []struct {
	Name string
	URL  string
} {
	var calls []struct {
		Name string
		URL  string
	}
	mock.lockDoGetHealthClient.RLock()
	calls = mock.calls.DoGetHealthClient
	mock.lockDoGetHealthClient.RUnlock()
	return calls
}

// DoGetKafkaConsumer calls DoGetKafkaConsumerFunc.
func (mock *InitialiserMock) DoGetKafkaConsumer(ctx context.Context, kafkaCfg *config.KafkaConfig) (kafka.IConsumerGroup, error) {
	if mock.DoGetKafkaConsumerFunc == nil {
		panic("InitialiserMock.DoGetKafkaConsumerFunc: method is nil but Initialiser.DoGetKafkaConsumer was just called")
	}
	callInfo := struct {
		Ctx      context.Context
		KafkaCfg *config.KafkaConfig
	}{
		Ctx:      ctx,
		KafkaCfg: kafkaCfg,
	}
	mock.lockDoGetKafkaConsumer.Lock()
	mock.calls.DoGetKafkaConsumer = append(mock.calls.DoGetKafkaConsumer, callInfo)
	mock.lockDoGetKafkaConsumer.Unlock()
	return mock.DoGetKafkaConsumerFunc(ctx, kafkaCfg)
}

// DoGetKafkaConsumerCalls gets all the calls that were made to DoGetKafkaConsumer.
// Check the length with:
//
//	len(mockedInitialiser.DoGetKafkaConsumerCalls())
func (mock *InitialiserMock) DoGetKafkaConsumerCalls() []struct {
	Ctx      context.Context
	KafkaCfg *config.KafkaConfig
} {
	var calls []struct {
		Ctx      context.Context
		KafkaCfg *config.KafkaConfig
	}
	mock.lockDoGetKafkaConsumer.RLock()
	calls = mock.calls.DoGetKafkaConsumer
	mock.lockDoGetKafkaConsumer.RUnlock()
	return calls
}

// DoGetKafkaProducer calls DoGetKafkaProducerFunc.
func (mock *InitialiserMock) DoGetKafkaProducer(ctx context.Context, cfg *config.Config) (kafka.IProducer, error) {
	if mock.DoGetKafkaProducerFunc == nil {
		panic("InitialiserMock.DoGetKafkaProducerFunc: method is nil but Initialiser.DoGetKafkaProducer was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Cfg *config.Config
	}{
		Ctx: ctx,
		Cfg: cfg,
	}
	mock.lockDoGetKafkaProducer.Lock()
	mock.calls.DoGetKafkaProducer = append(mock.calls.DoGetKafkaProducer, callInfo)
	mock.lockDoGetKafkaProducer.Unlock()
	return mock.DoGetKafkaProducerFunc(ctx, cfg)
}

// DoGetKafkaProducerCalls gets all the calls that were made to DoGetKafkaProducer.
// Check the length with:
//
//	len(mockedInitialiser.DoGetKafkaProducerCalls())
func (mock *InitialiserMock) DoGetKafkaProducerCalls() []struct {
	Ctx context.Context
	Cfg *config.Config
} {
	var calls []struct {
		Ctx context.Context
		Cfg *config.Config
	}
	mock.lockDoGetKafkaProducer.RLock()
	calls = mock.calls.DoGetKafkaProducer
	mock.lockDoGetKafkaProducer.RUnlock()
	return calls
}

// DoGetKafkaProducerForReindexTaskCounts calls DoGetKafkaProducerForReindexTaskCountsFunc.
func (mock *InitialiserMock) DoGetKafkaProducerForReindexTaskCounts(ctx context.Context, cfg *config.Config) (kafka.IProducer, error) {
	if mock.DoGetKafkaProducerForReindexTaskCountsFunc == nil {
		panic("InitialiserMock.DoGetKafkaProducerForReindexTaskCountsFunc: method is nil but Initialiser.DoGetKafkaProducerForReindexTaskCounts was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Cfg *config.Config
	}{
		Ctx: ctx,
		Cfg: cfg,
	}
	mock.lockDoGetKafkaProducerForReindexTaskCounts.Lock()
	mock.calls.DoGetKafkaProducerForReindexTaskCounts = append(mock.calls.DoGetKafkaProducerForReindexTaskCounts, callInfo)
	mock.lockDoGetKafkaProducerForReindexTaskCounts.Unlock()
	return mock.DoGetKafkaProducerForReindexTaskCountsFunc(ctx, cfg)
}

// DoGetKafkaProducerForReindexTaskCountsCalls gets all the calls that were made to DoGetKafkaProducerForReindexTaskCounts.
// Check the length with:
//
//	len(mockedInitialiser.DoGetKafkaProducerForReindexTaskCountsCalls())
func (mock *InitialiserMock) DoGetKafkaProducerForReindexTaskCountsCalls() []struct {
	Ctx context.Context
	Cfg *config.Config
} {
	var calls []struct {
		Ctx context.Context
		Cfg *config.Config
	}
	mock.lockDoGetKafkaProducerForReindexTaskCounts.RLock()
	calls = mock.calls.DoGetKafkaProducerForReindexTaskCounts
	mock.lockDoGetKafkaProducerForReindexTaskCounts.RUnlock()
	return calls
}

// DoGetZebedeeClient calls DoGetZebedeeClientFunc.
func (mock *InitialiserMock) DoGetZebedeeClient(cfg *config.Config) clients.ZebedeeClient {
	if mock.DoGetZebedeeClientFunc == nil {
		panic("InitialiserMock.DoGetZebedeeClientFunc: method is nil but Initialiser.DoGetZebedeeClient was just called")
	}
	callInfo := struct {
		Cfg *config.Config
	}{
		Cfg: cfg,
	}
	mock.lockDoGetZebedeeClient.Lock()
	mock.calls.DoGetZebedeeClient = append(mock.calls.DoGetZebedeeClient, callInfo)
	mock.lockDoGetZebedeeClient.Unlock()
	return mock.DoGetZebedeeClientFunc(cfg)
}

// DoGetZebedeeClientCalls gets all the calls that were made to DoGetZebedeeClient.
// Check the length with:
//
//	len(mockedInitialiser.DoGetZebedeeClientCalls())
func (mock *InitialiserMock) DoGetZebedeeClientCalls() []struct {
	Cfg *config.Config
} {
	var calls []struct {
		Cfg *config.Config
	}
	mock.lockDoGetZebedeeClient.RLock()
	calls = mock.calls.DoGetZebedeeClient
	mock.lockDoGetZebedeeClient.RUnlock()
	return calls
}
