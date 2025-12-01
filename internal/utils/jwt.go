// internal/utils/jwt.go
package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"student-portal/internal/config"
	appErrors "student-portal/internal/errors"
)

// UserClaims defines the claims structure for the JWT.
type UserClaims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken creates a new JWT for the given user ID, email, and role.
func GenerateToken(cfg *config.Config, userID int64, email, role string) (string, error) {
	expirationTime := time.Now().Add(cfg.JWTExpiry)

	claims := &UserClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "student-portal-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", appErrors.ErrInternalServerError
	}

	return tokenString, nil
}

// ValidateToken parses and validates a JWT string.
func ValidateToken(cfg *config.Config, tokenStr string) (*UserClaims, error) {
	claims := &UserClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, appErrors.ErrUnauthorized
	}

	if !token.Valid {
		return nil, appErrors.ErrUnauthorized
	}

	return claims, nil
}
