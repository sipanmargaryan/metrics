package helpers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/twmb/franz-go/pkg/kgo"
)

// Message is the common JSON envelope for every published event.
type Message struct {
	MessageID string `json:"message_id"`
	Data      any    `json:"data"`
}

// EventProducer is the publishing dependency passed to application services.
type EventProducer interface {
	Produce(ctx context.Context, topic, key string, event any) error
}

// GenerateKafkaKey returns a stable message key for an entity ID.
func GenerateKafkaKey(id uint) string {
	return strconv.FormatUint(uint64(id), 10)
}

// KafkaProducer publishes JSON events. Reuse one producer across requests and
// close it after all publishers have stopped during application shutdown.
type KafkaProducer struct {
	client *kgo.Client
}

// NewKafkaProducer configures a producer for a plaintext Kafka cluster.
// Idempotent writes are enabled by franz-go's default configuration.
// Construction does not verify cluster connectivity.
// Options override the default batch byte limit and 10 ms linger.
func NewKafkaProducer(brokers []string, options ...kgo.ProducerOpt) (*KafkaProducer, error) {
	addresses := make([]string, 0, len(brokers))
	for _, broker := range brokers {
		if address := strings.TrimSpace(broker); address != "" {
			addresses = append(addresses, address)
		}
	}
	if len(addresses) == 0 {
		return nil, errors.New("kafka: at least one broker is required")
	}

	clientOptions := []kgo.Opt{
		kgo.SeedBrokers(addresses...),
		kgo.RequiredAcks(kgo.AllISRAcks()),
		kgo.ProducerBatchMaxBytes(1_000_012),
		kgo.ProducerLinger(10 * time.Millisecond),
	}
	for _, option := range options {
		clientOptions = append(clientOptions, option)
	}
	client, err := kgo.NewClient(clientOptions...)
	if err != nil {
		return nil, fmt.Errorf("kafka: initialize producer: %w", err)
	}
	return &KafkaProducer{client: client}, nil
}

// Produce encodes event as JSON and waits for broker acknowledgement.
// Use the same key to route related events to the same partition.
// Internal retries are idempotent. A new Produce call is a new message, even
// with the same key, so callers must account for duplicates on application retries.
// Every call generates a UUIDv7 message ID, retained through internal retries.
func (p *KafkaProducer) Produce(ctx context.Context, topic, key string, event any) error {
	record, err := newKafkaRecord(topic, key, event)
	if err != nil {
		return err
	}
	// ProduceSync bypasses linger. Queue with Produce and wait for its callback
	// so concurrent requests can batch while callers still receive delivery errors.
	done := make(chan error, 1)
	p.client.Produce(ctx, record, func(_ *kgo.Record, err error) {
		done <- err
	})
	if err := <-done; err != nil {
		return fmt.Errorf("kafka: produce to %q: %w", topic, err)
	}
	return nil
}

func newKafkaRecord(topic, key string, event any) (*kgo.Record, error) {
	if strings.TrimSpace(topic) == "" {
		return nil, errors.New("kafka: topic is required")
	}
	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("kafka: generate message ID: %w", err)
	}
	value, err := json.Marshal(Message{MessageID: id.String(), Data: event})
	if err != nil {
		return nil, fmt.Errorf("kafka: encode event: %w", err)
	}
	return &kgo.Record{
		Topic: topic,
		Key:   []byte(key),
		Value: value,
	}, nil
}

// Close releases the producer's connections after publishers have stopped.
func (p *KafkaProducer) Close() error {
	p.client.Close()
	return nil
}
