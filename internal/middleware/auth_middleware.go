// internal/middleware/auth_middleware.go
package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"student-portal/internal/config"
	"student-portal/internal/constants"
	appErrors "student-portal/internal/errors"
	"student-portal/internal/logger"
	"student-portal/internal/utils"
)

// AuthMiddleware validates the JWT token and sets user claims in the context.
func AuthMiddleware(cfg *config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				handleError(w, appErrors.ErrUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				handleError(w, appErrors.ErrUnauthorized)
				return
			}

			tokenStr := parts[1]
			claims, err := utils.ValidateToken(cfg, tokenStr)
			if err != nil {
				handleError(w, err)
				return
			}

			// Store claims in context
			ctx := context.WithValue(r.Context(), constants.UserClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RoleMiddleware checks if the authenticated user has one of the required roles.
func RoleMiddleware(requiredRoles ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(constants.UserClaimsKey).(*utils.UserClaims)
			if !ok {
				// Should not happen if AuthMiddleware is run first
				handleError(w, appErrors.ErrUnauthorized)
				return
			}

			hasPermission := false
			for _, role := range requiredRoles {
				if claims.Role == string(role) {
					hasPermission = true
					break
				}
			}

			if !hasPermission {
				handleError(w, appErrors.ErrForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GetUserClaims retrieves the UserClaims from the request context.
func GetUserClaims(ctx context.Context) *utils.UserClaims {
	claims, ok := ctx.Value(constants.UserClaimsKey).(*utils.UserClaims)
	if !ok {
		// Log an internal server error if claims are missing when expected
		logger.Logger.Error("Attempted to get claims from context but none were present")
		return nil
	}
	return claims
}

func handleError(w http.ResponseWriter, err error) {
	appErr, ok := err.(*appErrors.AppError)
	if !ok {
		appErr = appErrors.ErrUnauthorized // Default to unauthorized on generic error
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.Code)
	// nolint:errcheck
	json.NewEncoder(w).Encode(appErr)
}
