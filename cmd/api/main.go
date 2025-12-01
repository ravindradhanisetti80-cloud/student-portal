// cmd/api/main.go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"student-portal/internal/config"
	"student-portal/internal/handler"
	"student-portal/internal/logger"
	"student-portal/internal/repository"
	"student-portal/internal/routes"
	"student-portal/internal/service"
	"student-portal/internal/utils" // Contains Kafka utilities

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap" // For structured error logging
)

func main() {
	// 1. Load Configuration
	cfg := config.LoadConfig()

	// 2. Initialize Logger
	logger.InitLogger(cfg.AppEnv)
	defer logger.SyncLogger()

	logger.Logger.Info(
		fmt.Sprintf("Starting server in %s environment", cfg.AppEnv),
	)

	// --- SETUP CONTEXT FOR GRACEFUL SHUTDOWN ---
	// This context is used to signal the server, Kafka consumer, and topic creation to stop/timeout.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 3. Initialize Database Connection
	dbPool, err := initDB(cfg)
	if err != nil {
		logger.Logger.Fatal(fmt.Sprintf("Failed to connect to database: %v", err))
	}
	defer dbPool.Close()

	// 4. Kafka Initialization
	brokerList := strings.Split(cfg.KafkaBrokers, ",")

	// 4a. Check Broker Reachability & Create Topics (Health Check)
	if err := utils.CreateTopics(brokerList); err != nil {
		logger.Logger.Fatal(fmt.Sprintf("Failed to connect to Kafka or create topics: %v", err))
	}

	// 4b. Initialize Kafka Producer
	kafkaProducer := utils.NewKafkaProducer(brokerList)
	defer kafkaProducer.Close()

	// 4c. Start Kafka Consumer (Pass the shutdown context)
	kafkaConsumerReader := utils.StartConsumer(ctx, cfg)
	defer func() {
		// Ensure the consumer reader is closed during shutdown
		if err := kafkaConsumerReader.Close(); err != nil {
			logger.Logger.Error("Failed to gracefully close Kafka consumer reader", zap.Error(err))
		} else {
			logger.Logger.Info("Kafka consumer reader closed successfully")
		}
	}()

	// 5. Dependency Injection
	userRepo := repository.NewUserRepository(dbPool)
	userService := service.NewUserService(userRepo, cfg, kafkaProducer)
	authHandler := handler.NewAuthHandler(userService, cfg)
	userHandler := handler.NewUserHandler(userService, cfg)

	// 6. Setup Router
	r := routes.SetupRouter(cfg, authHandler, userHandler)

	// 7. Start Server
	server := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r,
	}

	go func() {
		logger.Logger.Info(fmt.Sprintf("Server starting on :%s", cfg.ServerPort))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger.Fatal(fmt.Sprintf("Could not listen on %s: %v", cfg.ServerPort, err))
		}
	}()

	// Wait for interrupt signal (comes from the context setup at the start of main)
	<-ctx.Done()

	logger.Logger.Info("Server shutting down...")

	// 8. Shutdown gracefully
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(timeoutCtx); err != nil {
		logger.Logger.Error(fmt.Sprintf("Server shutdown failed: %v", err))
	} else {
		logger.Logger.Info("Server gracefully stopped")
	}
}

// initDB initializes the PostgreSQL connection pool. (Unchanged)
func initDB(cfg *config.Config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbPool, err := pgxpool.New(ctx, cfg.DatabaseURL())
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := dbPool.Ping(ctx); err != nil {
		return nil, err
	}

	logger.Logger.Info("Successfully connected to PostgreSQL!")
	return dbPool, nil
}
