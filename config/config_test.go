package config

import (
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConfig(t *testing.T) {
	Convey("Given an environment with no environment variables set", t, func() {
		os.Clearenv()
		cfg, err := Get()
		Convey("When the config values are retrieved", func() {
			Convey("Then there should be no error returned", func() {
				So(err, ShouldBeNil)
			})
			Convey("Then the values should be set to the expected defaults", func() {
				So(cfg.APIRouterURL, ShouldEqual, "http://localhost:23200/v1")
				So(cfg.BindAddr, ShouldEqual, "localhost:28000")
				So(cfg.EnablePublishContentUpdatedTopic, ShouldEqual, false)
				So(cfg.GracefulShutdownTimeout, ShouldEqual, 5*time.Second)
				So(cfg.HealthCheckInterval, ShouldEqual, 30*time.Second)
				So(cfg.HealthCheckCriticalTimeout, ShouldEqual, 90*time.Second)
				So(cfg.KafkaConfig.Brokers, ShouldHaveLength, 3)
				So(cfg.KafkaConfig.Brokers[0], ShouldEqual, "localhost:9092")
				So(cfg.KafkaConfig.ContentUpdatedTopic, ShouldEqual, "content-updated")
				So(cfg.KafkaConfig.ConsumerGroup, ShouldEqual, "dp-search-data-finder")
				So(cfg.KafkaConfig.NumWorkers, ShouldEqual, 1)
				So(cfg.KafkaConfig.ReindexRequestedTopic, ShouldEqual, "reindex-requested")
				So(cfg.KafkaConfig.SecProtocol, ShouldEqual, "")
				So(cfg.KafkaConfig.Version, ShouldEqual, "1.0.2")
				So(cfg.ServiceAuthToken, ShouldEqual, "")
				So(cfg.ZebedeeRequestLimit, ShouldEqual, 10)
			})
			Convey("Then a second call to config should return the same config", func() {
				newCfg, newErr := Get()
				So(newErr, ShouldBeNil)
				So(newCfg, ShouldResemble, cfg)
			})
		})
	})
}
