package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
	"github.com/sigame/game/internal/domain"
)

// Producer handles publishing events to Kafka
type Producer struct {
	writer *kafka.Writer
	topic  string
}

// NewProducer creates a new Kafka producer
func NewProducer(brokers []string, topic string) *Producer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &Producer{
		writer: writer,
		topic:  topic,
	}
}

// Close closes the producer
func (p *Producer) Close() error {
	return p.writer.Close()
}

// PublishGameEvent publishes a game event to Kafka
func (p *Producer) PublishGameEvent(ctx context.Context, event *domain.GameEvent) error {
	// Serialize event to JSON
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Create Kafka message
	msg := kafka.Message{
		Key:   []byte(event.GameID.String()),
		Value: data,
	}

	// Send message
	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	log.Printf("Published event %s for game %s to Kafka", event.EventType, event.GameID)
	return nil
}

// PublishGameEvents publishes multiple events in batch
func (p *Producer) PublishGameEvents(ctx context.Context, events []*domain.GameEvent) error {
	if len(events) == 0 {
		return nil
	}

	messages := make([]kafka.Message, 0, len(events))

	for _, event := range events {
		data, err := json.Marshal(event)
		if err != nil {
			log.Printf("Failed to marshal event: %v", err)
			continue
		}

		messages = append(messages, kafka.Message{
			Key:   []byte(event.GameID.String()),
			Value: data,
		})
	}

	if err := p.writer.WriteMessages(ctx, messages...); err != nil {
		return fmt.Errorf("failed to write messages: %w", err)
	}

	log.Printf("Published %d events to Kafka", len(messages))
	return nil
}

