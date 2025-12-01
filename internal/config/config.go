package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration values.
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	ServerPort string

	JWTSecret string
	JWTExpiry time.Duration
	AppEnv    string

	// Kafka configuration
	KafkaBrokers string
	KafkaTopic   string
}

// LoadConfig reads configuration from .env file and environment variables.
func LoadConfig() *Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: No .env file found: %v", err)
	}

	//JWT Expiry Parsing
	jwtExpiryStr := getEnv("JWT_EXPIRY", "24h") // Get value or default to "24h"
	expiryDuration, err := time.ParseDuration(jwtExpiryStr)
	if err != nil {
		log.Printf("Warning: Failed to parse JWT_EXPIRY '%s'. Defaulting to 24 hours.", jwtExpiryStr)
		expiryDuration = 24 * time.Hour
	}

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "student_portal"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		ServerPort: getEnv("SERVER_PORT", "8080"), // Changed default back to 8080 for consistency

		JWTSecret: getEnv("JWT_SECRET", "super-secret-key"),
		JWTExpiry: expiryDuration,
		AppEnv:    getEnv("APP_ENV", "development"),

		// Kafka defaults
		KafkaBrokers: getEnv("KAFKA_BROKER", "localhost:9092"),
	}
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// DatabaseURL constructs the PostgreSQL connection URL.
func (c *Config) DatabaseURL() string {
	return "postgresql://" + c.DBUser + ":" + c.DBPassword + "@" + c.DBHost + ":" + c.DBPort + "/" + c.DBName + "?sslmode=" + c.DBSSLMode
}
