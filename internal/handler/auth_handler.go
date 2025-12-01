// internal/handler/auth_handler.go
package handler

import (
	"encoding/json"
	"net/http"

	appErrors "student-portal/internal/commons/errors"
	"student-portal/internal/config"
	"student-portal/internal/models"
	"student-portal/internal/service"
	"student-portal/internal/utils"
)

// AuthHandler handles HTTP requests for authentication.
type AuthHandler struct {
	svc service.UserService
	cfg *config.Config
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(svc service.UserService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{svc: svc, cfg: cfg}
}

// Routes sets up the public routes for authentication.

// Failure 500 {object} errors.AppError
// Router /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, appErrors.ErrBadRequest)
		return
	}
	// Note: Basic validation for required fields is done by the service/repository via error codes.
	// In a real app, external validation library (like validator) should be used here.

	userResp, err := h.svc.RegisterUser(r.Context(), &req)
	if err != nil {
		utils.SendError(w, err)
		return
	}

	utils.SendJSON(w, http.StatusCreated, userResp)
}

// Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, appErrors.ErrBadRequest)
		return
	}

	loginResp, err := h.svc.LoginUser(r.Context(), &req)
	if err != nil {
		utils.SendError(w, err)
		return
	}

	utils.SendJSON(w, http.StatusOK, loginResp)
}
