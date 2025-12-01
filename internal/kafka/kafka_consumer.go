// internal/utils/kafka_consumer.go
package utils

import (
	"context"
	"strings"
	"time"

	"student-portal/internal/commons/constants" // Assuming you created this in Step 1C
	"student-portal/internal/commons/logger"
	"student-portal/internal/config"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// StartConsumer initializes and starts listening on the user events topic.
// Returns the reader for graceful shutdown.
// NOTE: In a real application, you would manage multiple topics and groups.
func StartConsumer(ctx context.Context, cfg *config.Config) *kafka.Reader {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     strings.Split(cfg.KafkaBrokers, ","),
		Topic:       constants.TopicUserEvents,
		GroupID:     "student-portal-group", // Unique ID for this service instance
		MinBytes:    10e3,                   // 10KB
		MaxBytes:    10e6,                   // 10MB
		MaxAttempts: 3,
		Dialer: &kafka.Dialer{
			Timeout: 10 * time.Second,
		},
	})

	logger.Logger.Info("Kafka consumer started",
		zap.String("topic", constants.TopicUserEvents),
		zap.String("group", "student-portal-group"),
	)

	// Start consuming messages in a goroutine
	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.Logger.Info("Kafka consumer shutting down due to context cancellation")
				return
			default:
				m, err := r.ReadMessage(context.Background())

				if err != nil {
					// Handle recoverable errors (e.g., EOF, connection loss)
					logger.Logger.Error("Error reading Kafka message", zap.Error(err))
					time.Sleep(5 * time.Second) // Wait before retrying
					continue
				}

				// --- Message Processing Logic ---
				// In a real application, you would deserialize m.Value and call a service layer function
				valueStr := string(m.Value)
				if len(valueStr) > 20 {
					valueStr = valueStr[:20]
				}
				logger.Logger.Info("Message received",
					zap.String("topic", m.Topic),
					zap.Int64("offset", m.Offset),
					zap.String("key", string(m.Key)),
					zap.String("value_preview", valueStr), // Log only a snippet
				)
			}
		}
	}()

	// Add a function to close the reader for graceful shutdown (to be called in main.go)
	logger.Logger.Info("[Kafka] Consumer initialized. Closing function available.")
	return r
}

// CloseConsumer attempts to gracefully close the Kafka reader connection.
func CloseConsumer(cfg *config.Config) {
	// NOTE: This requires storing the reader instance, which is advanced.
	// For simplicity here, we'll assume manual management in main.go if needed,
	// or rely on process termination if the reader is scoped globally.
	logger.Logger.Info("Kafka consumer shutdown initiated.")
}
