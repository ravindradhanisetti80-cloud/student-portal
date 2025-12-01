package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

// KafkaProducer handles Kafka message publishing
type KafkaProducer struct {
	writer *kafka.Writer
}

// NewKafkaProducer creates a new Kafka producer
func NewKafkaProducer(brokers []string) *KafkaProducer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Balancer:     &kafka.LeastBytes{},
		Async:        true, // Enable async publishing for better performance
		RequiredAcks: kafka.RequireOne,
	}

	return &KafkaProducer{writer: writer}
}

// Close closes the Kafka writer
func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}

// PublishMessage publishes a message to a specific topic
func (p *KafkaProducer) PublishMessage(ctx context.Context, topic string, key string, value interface{}) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	message := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: jsonValue,
		Time:  time.Now(),
	}

	return p.writer.WriteMessages(ctx, message)
}

// AuthEvent represents authentication events
type AuthEvent struct {
	EventType string    `json:"event_type"`
	UserID    int64     `json:"user_id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	Timestamp time.Time `json:"timestamp"`
	IPAddress string    `json:"ip_address,omitempty"`
	UserAgent string    `json:"user_agent,omitempty"`
}

// PublishLoginEvent publishes a login event to Kafka
func (p *KafkaProducer) PublishLoginEvent(ctx context.Context, userID int64, email, name, role string) error {
	event := AuthEvent{
		EventType: "user_login",
		UserID:    userID,
		Email:     email,
		Name:      name,
		Role:      role,
		Timestamp: time.Now(),
	}

	return p.PublishMessage(ctx, "user-auth-events", email, event)
}

// PublishRegisterEvent publishes a register event to Kafka
func (p *KafkaProducer) PublishRegisterEvent(ctx context.Context, userID int64, email, name, role string) error {
	event := AuthEvent{
		EventType: "user_register",
		UserID:    userID,
		Email:     email,
		Name:      name,
		Role:      role,
		Timestamp: time.Now(),
	}

	return p.PublishMessage(ctx, "user-auth-events", email, event)
}

// PublishUpdateEvent publishes an update event to Kafka
func (p *KafkaProducer) PublishUpdateEvent(ctx context.Context, userID int64, email, name, role string) error {
	event := AuthEvent{
		EventType: "user_update",
		UserID:    userID,
		Email:     email,
		Name:      name,
		Role:      role,
		Timestamp: time.Now(),
	}

	return p.PublishMessage(ctx, "user-auth-events", email, event)
}
