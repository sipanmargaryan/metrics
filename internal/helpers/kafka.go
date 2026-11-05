package helpers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

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
func NewKafkaProducer(brokers []string) (*KafkaProducer, error) {
	addresses := make([]string, 0, len(brokers))
	for _, broker := range brokers {
		if address := strings.TrimSpace(broker); address != "" {
			addresses = append(addresses, address)
		}
	}
	if len(addresses) == 0 {
		return nil, errors.New("kafka: at least one broker is required")
	}

	client, err := kgo.NewClient(
		kgo.SeedBrokers(addresses...),
		kgo.RequiredAcks(kgo.AllISRAcks()), //hoq many isr should acknowledge
		kgo.ProducerLinger(10*time.Millisecond),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka: initialize producer: %w", err)
	}
	return &KafkaProducer{client: client}, nil
}

// Produce encodes event as JSON and waits for broker acknowledgement.
// Use the same key to route related events to the same partition.
// Internal retries are idempotent. A new Produce call is a new message, even
// with the same key, so callers must account for duplicates on application retries.
func (p *KafkaProducer) Produce(ctx context.Context, topic, key string, event any) error {
	if strings.TrimSpace(topic) == "" {
		return errors.New("kafka: topic is required")
	}
	value, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("kafka: encode event: %w", err)
	}
	if err := p.client.ProduceSync(ctx, &kgo.Record{
		Topic: topic,
		Key:   []byte(key),
		Value: value,
	}).FirstErr(); err != nil {
		return fmt.Errorf("kafka: produce to %q: %w", topic, err)
	}
	return nil
}

// Close releases the producer's connections after publishers have stopped.
func (p *KafkaProducer) Close() error {
	p.client.Close()
	return nil
}
