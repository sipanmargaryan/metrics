# Metrics

A Go API for managing warehouses, products, and orders with PostgreSQL.
Includes Prometheus metrics, Loki, and Grafana monitoring. Requires Docker
with the Compose plugin and GNU Make.

```sh
make migrate-up    # Start PostgreSQL and apply migrations
make up            # Start the API and monitoring services
make down          # Stop services and preserve data
make help          # List available commands
make env           # Create .env from env_example if missing
make build         # Build the application image
make restart       # Restart the API
make logs          # Follow all service logs
make logs SERVICE=app # Follow API logs
make ps            # Show service status
make config        # Validate Compose configuration
make db-up         # Start PostgreSQL only
make db-shell      # Open the PostgreSQL shell
make test          # Run Go tests in Docker
make load-test     # Send 100 API requests (requires Python 3)
make swagger       # Generate Swagger documentation (requires Go)
make migrate-down # Revert the latest migration (stop the API first)
make dev           # Alias for make up
make dev-down      # Alias for make down
```

Every Kafka publish automatically wraps the event in a common JSON envelope:

```json
{
  "message_id": "019a581d-3600-7000-8000-000000000001",
  "data": { "order_id": 42 }
}
```

`message_id` is a UUIDv7 generated inside the producer; callers keep passing
their event payload as before. Consumers now read the payload from `data`.
The Kafka key still controls partition routing independently of the message ID.

Kafka's internal retries preserve the ID. Calling `Produce` again generates a
new ID, so this supports deduplicating redelivery of a published message, but
does not deduplicate separate application publish calls. Consumers should record
processed IDs and apply business changes in the same database transaction, then
commit Kafka offsets after the transaction succeeds.

Configure batching when creating the shared producer:

```go
producer, err := helpers.NewKafkaProducer(cfg.Kafka.Brokers,
    kgo.ProducerBatchMaxBytes(64*1024), // 64 KiB per partition batch
    kgo.ProducerLinger(10*time.Millisecond),
)
```

Import `github.com/twmb/franz-go/pkg/kgo` and `time` for these options. Without
overrides, the batch limit is 1,000,012 bytes and linger is 10 ms. The batch
limit also caps individual uncompressed records, including envelope overhead;
choose a value that fits your largest message and the broker's limits. Linger
allows more messages to accumulate per partition, trading latency for batching.
Concurrent publish calls can share a batch; sequential calls that each wait for
acknowledgement cannot accumulate together. Publishing still waits for broker
acknowledgement and returns delivery errors.
