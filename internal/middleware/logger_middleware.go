// internal/middleware/logger_middleware.go
package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
	"student-portal/internal/logger"
)

// Logger is a simple structure to hold the response status.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

// RequestLogger is a Chi middleware that logs every request with structured logging.
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &statusRecorder{
			ResponseWriter: w,
			status:         http.StatusOK, // Default status
		}

		next.ServeHTTP(recorder, r)

		logger.Logger.Info("Request Handled",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Int("status", recorder.status),
			zap.Duration("duration", time.Since(start)),
			zap.String("remote_addr", r.RemoteAddr),
		)
	})
}
