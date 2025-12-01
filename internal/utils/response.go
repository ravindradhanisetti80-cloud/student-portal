package utils

import (
	"encoding/json"
	"net/http"
	"student-portal/internal/errors"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func WriteSuccess(w http.ResponseWriter, statusCode int, data interface{}) {
	WriteJSON(w, statusCode, Response{
		Success: true,
		Data:    data,
	})
}

func WriteError(w http.ResponseWriter, statusCode int, message string) {
	WriteJSON(w, statusCode, Response{
		Success: false,
		Error:   http.StatusText(statusCode),
		Message: message,
	})
}

func SendError(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*errors.AppError); ok {
		WriteJSON(w, appErr.Code, Response{
			Success: false,
			Error:   http.StatusText(appErr.Code),
			Message: appErr.Message,
		})
	} else {
		WriteError(w, http.StatusInternalServerError, "An unexpected error occurred")
	}
}

func SendJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	WriteSuccess(w, statusCode, data)
}
