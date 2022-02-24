# dp-search-data-finder
Receives events requesting reindex jobs and finds data to populate the new indexes created

### Getting started

* Run `make debug`

The service runs in the background consuming messages from Kafka.
An example event can be created using the helper script, `make produce`.

### Dependencies

* Requires running…
  * [kafka](https://github.com/ONSdigital/dp/blob/main/guides/INSTALLING.md#prerequisites)
* No further dependencies other than those defined in `go.mod`

### Configuration

| Environment variable         | Default                           | Description
| ---------------------------- | --------------------------------- | -----------
| BIND_ADDR                    | localhost:28000                    | The host and port to bind to
| GRACEFUL_SHUTDOWN_TIMEOUT    | 5s                                | The graceful shutdown timeout in seconds (`time.Duration` format)
| HEALTHCHECK_INTERVAL         | 30s                               | Time between self-healthchecks (`time.Duration` format)
| HEALTHCHECK_CRITICAL_TIMEOUT | 90s                               | Time to wait until an unhealthy dependent propagates its state to make this app unhealthy (`time.Duration` format)
| KAFKA_ADDR                   | "localhost:9092"                  | The address of Kafka (accepts list)
| KAFKA_OFFSET_OLDEST          | true                              | Start processing Kafka messages in order from the oldest in the queue
| KAFKA_NUM_WORKERS            | 1                                 | The maximum number of parallel kafka consumers
| KAFKA_SEC_PROTO              | _unset_                           | if set to `TLS`, kafka connections will use TLS ([kafka TLS doc])
| KAFKA_SEC_CA_CERTS           | _unset_                           | CA cert chain for the server cert ([kafka TLS doc])
| KAFKA_SEC_CLIENT_KEY         | _unset_                           | PEM for the client key ([kafka TLS doc])
| KAFKA_SEC_CLIENT_CERT        | _unset_                           | PEM for the client certificate ([kafka TLS doc])
| KAFKA_SEC_SKIP_VERIFY        | false                             | ignores server certificate issues if `true` ([kafka TLS doc])
| HELLO_CALLED_GROUP           | dp-search-data-finder              | The consumer group this application to consume topic messages
| HELLO_CALLED_TOPIC           | hello-called                      | The name of the topic to consume messages from

[kafka TLS doc]: https://github.com/ONSdigital/dp-kafka/tree/main/examples#tls

### Healthcheck

 The `/health` endpoint returns the current status of the service. Dependent services are health checked on an interval defined by the `HEALTHCHECK_INTERVAL` environment variable.

 On a development machine a request to the health check endpoint can be made by:

 `curl localhost:8125/health`

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2022, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.

