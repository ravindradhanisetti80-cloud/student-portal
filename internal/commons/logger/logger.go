// internal/logger/logger.go
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is the application-wide structured logger.
var Logger *zap.Logger

// InitLogger initializes the Zap logger based on the application environment.
func InitLogger(env string) {
	var err error
	var cfg zap.Config

	if env == "production" {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		// Development configuration for human readability
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	Logger, err = cfg.Build()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	// Replace the global standard logger with Zap's sugar wrapper for compatibility
	zap.ReplaceGlobals(Logger)
}

// SyncLogger flushes the buffer, should be deferred in main.
func SyncLogger() {
	if Logger != nil {
		// nolint:errcheck
		Logger.Sync()
	}
}
