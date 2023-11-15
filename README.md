# dp-search-data-finder
Receives events requesting reindex jobs and finds data to populate the new indexes created

### Getting started

* Run `make debug`

The service runs in the background consuming messages from Kafka.
An example event can be created using the helper script, `make produce`.

### Dependencies

* Requires running…
  * [kafka](https://github.com/ONSdigital/dp/blob/main/guides/INSTALLING.md#prerequisites)
  * [zebedee](https://github.com/ONSdigital/zebedee)
* No further dependencies other than those defined in `go.mod`

### Configuration

| Environment variable          | Default                           | Description
| ----------------------------- | --------------------------------- | -----------
| API_ROUTER_URL                | "http://localhost:23200/v1"       | API router URL
| BIND_ADDR                     | "localhost:28000"                 | The host and port to bind to
| ENABLE_PUBLISH_CONTENT_UPDATED_TOPIC | false                      | Enables content update topic to be published to
| CONTENT_UPDATED_TOPIC_FLAG    | false                             | produce events only if set to `true`
| GRACEFUL_SHUTDOWN_TIMEOUT     | 5s                                | The graceful shutdown timeout in seconds (`time.Duration` format)
| HEALTHCHECK_CRITICAL_TIMEOUT  | 90s                               | Time to wait until an unhealthy dependent propagates its state to make this app unhealthy (`time.Duration` format)
| HEALTHCHECK_INTERVAL          | 30s                               | Time between self-healthchecks (`time.Duration` format)
| KAFKA_ADDR                    | "localhost:9092"                  | The address of Kafka (accepts list)
| KAFKA_CONTENT_UPDATED_TOPIC   | "content-updated"                 | The name of the topic to produce messages for
| KAFKA_NUM_WORKERS             | 1                                 | The maximum number of parallel kafka consumers
| KAFKA_OFFSET_OLDEST           | true                              | Start processing Kafka messages in order from the oldest in the queue
| KAFKA_REINDEX_REQUESTED_GROUP | "dp-search-data-finder"           | The consumer group for this application to consume reindex-requested messages
| KAFKA_REINDEX_REQUESTED_TOPIC | "reindex-requested"               | The name of the topic to consume messages from
| KAFKA_SEC_PROTO               | _unset_                           | if set to `TLS`, kafka connections will use TLS ([kafka TLS doc])
| KAFKA_SEC_CA_CERTS            | _unset_                           | CA cert chain for the server cert ([kafka TLS doc])
| KAFKA_SEC_CLIENT_CERT         | _unset_                           | PEM for the client certificate ([kafka TLS doc])
| KAFKA_SEC_CLIENT_KEY          | _unset_                           | PEM for the client key ([kafka TLS doc])
| KAFKA_SEC_SKIP_VERIFY         | false                             | ignores server certificate issues if `true` ([kafka TLS doc])
| KAFKA_VERSION                 | "1.0.2"                           | The kafka version that this service expects to connect to
| ZEBEDEE_CLIENT_TIMEOUT        | 30s                               | Time to wait for the zebedee client to respond to requests e.g. the published index request
| ZEBEDEE_URL                   | http://localhost:8082             | The URL to zebedee
| ZEBEDEE_REQUEST_LIMIT         | 10                                | Optional limit for development purposes on zebedee requests. Setting this to 0 removes the limiter.

[kafka TLS doc]: https://github.com/ONSdigital/dp-kafka/tree/main/examples#tls

### Healthcheck

The `/health` endpoint returns the current status of the service. Dependent services are health checked on an interval defined by the `HEALTHCHECK_INTERVAL` environment variable.

On a development machine a request to the health check endpoint can be made by:

`curl localhost:28000/health`

### Content Updated Topic

a. To check if service produced an event for content-updated Topic

make test

b. set EnablePublishContentUpdatedTopic in config as true/false

c. build and run the service
go build
./dp-search-data-finder

d. Send the event
make produce
Type uri (any text)

e. Check the service logs if either of the following appears and there is no error in service logs   
EnablePublishContentUpdatedTopic Flag is enabled/disabled

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2022, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
