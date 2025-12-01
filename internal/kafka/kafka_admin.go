// internal/utils/kafka_admin.go
package utils

import (
	"context"
	"time"

	"student-portal/internal/commons/constants"
	"student-portal/internal/commons/logger"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// CreateTopics connects to the Kafka broker and ensures all required topics exist.
func CreateTopics(brokers []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Create a connection to the first broker (admin client)
	conn, err := kafka.DialContext(ctx, "tcp", brokers[0])
	if err != nil {
		logger.Logger.Error("Failed to dial Kafka broker for admin operations", zap.Error(err))
		return err
	}
	defer conn.Close()

	// 2. Define topics to be created
	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             constants.TopicUserEvents,
			NumPartitions:     1, // Start with 1, scale up later
			ReplicationFactor: 1, // 1 for local docker setup
		},
		// Add other topics here (e.g., constants.TopicCourseEnrollments)
	}

	// 3. Create the topics
	err = conn.CreateTopics(topicConfigs...)
	if err != nil {
		logger.Logger.Error("Failed to create Kafka topics", zap.Error(err))
		return err
	}

	logger.Logger.Info("Successfully ensured all required Kafka topics exist")
	return nil
}
