// internal/handler/user_handler.go
package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"student-portal/internal/config"
	appErrors "student-portal/internal/errors"
	"student-portal/internal/middleware"
	"student-portal/internal/models"
	"student-portal/internal/service"
	"student-portal/internal/utils"

	"github.com/go-chi/chi/v5"
)

// UserHandler handles HTTP requests for user management.
type UserHandler struct {
	svc service.UserService
	cfg *config.Config
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(svc service.UserService, cfg *config.Config) *UserHandler {
	return &UserHandler{svc: svc, cfg: cfg}
}

// GetOwnProfile retrieves the authenticated user's profile.
func (h *UserHandler) GetOwnProfile(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	if claims == nil {
		utils.SendError(w, appErrors.ErrUnauthorized)
		return
	}

	userResp, err := h.svc.GetUserByID(r.Context(), claims.UserID)
	if err != nil {
		utils.SendError(w, err)
		return
	}

	utils.SendJSON(w, http.StatusOK, userResp)
}

// UpdateOwnProfile updates the authenticated user's profile.
func (h *UserHandler) UpdateOwnProfile(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	if claims == nil {
		utils.SendError(w, appErrors.ErrUnauthorized)
		return
	}

	var req models.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, appErrors.ErrBadRequest)
		return
	}

	userResp, err := h.svc.UpdateProfile(r.Context(), claims.UserID, &req)
	if err != nil {
		utils.SendError(w, err)
		return
	}

	utils.SendJSON(w, http.StatusOK, userResp)
}

// ListUsers retrieves all users (Admin Only).
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	query := utils.NewPaginationQuery(r)

	users, totalCount, err := h.svc.ListUsers(r.Context(), query.Limit, query.Offset)
	if err != nil {
		utils.SendError(w, err)
		return
	}

	resp := utils.NewPaginationResponse(users, query, totalCount)
	utils.SendJSON(w, http.StatusOK, resp)
}

// GetUserByID retrieves a user by ID (Admin Only).
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.SendError(w, appErrors.ErrBadRequest)
		return
	}

	userResp, err := h.svc.GetUserByID(r.Context(), id)
	if err != nil {
		utils.SendError(w, err)
		return
	}

	utils.SendJSON(w, http.StatusOK, userResp)
}

// UpdateUser updates a user by ID (Admin Only).
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.SendError(w, appErrors.ErrBadRequest)
		return
	}

	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, appErrors.ErrBadRequest)
		return
	}

	userResp, err := h.svc.UpdateUser(r.Context(), id, &req)
	if err != nil {
		utils.SendError(w, err)
		return
	}

	utils.SendJSON(w, http.StatusOK, userResp)
}

// DeleteUser deletes a user by ID (Admin Only).
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.SendError(w, appErrors.ErrBadRequest)
		return
	}

	if err := h.svc.DeleteUser(r.Context(), id); err != nil {
		utils.SendError(w, err)
		return
	}

	utils.SendJSON(w, http.StatusNoContent, nil)
}
