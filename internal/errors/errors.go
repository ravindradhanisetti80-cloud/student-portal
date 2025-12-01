// internal/errors/errors.go
package errors

import (
	"fmt"
	"net/http"
)

// AppError is a custom error type for centralized error handling.
type AppError struct {
	Code    int    `json:"-"` // HTTP status code
	Message string `json:"error"`
}

func (e *AppError) Error() string {
	return e.Message
}

// New creates a new AppError with a specific code and message.
func New(code int, format string, args ...interface{}) *AppError {
	return &AppError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

// Predefined common errors
var (
	ErrBadRequest          = New(http.StatusBadRequest, "Invalid request payload or parameters")
	ErrUnauthorized        = New(http.StatusUnauthorized, "Authentication required")
	ErrForbidden           = New(http.StatusForbidden, "Insufficient permissions")
	ErrNotFound            = New(http.StatusNotFound, "Resource not found")
	ErrConflict            = New(http.StatusConflict, "Resource already exists")
	ErrInternalServerError = New(http.StatusInternalServerError, "An unexpected error occurred")
	ErrInvalidCredentials  = New(http.StatusUnauthorized, "Invalid email or password")
	ErrEmailExists         = New(http.StatusConflict, "Email already exists")
)
