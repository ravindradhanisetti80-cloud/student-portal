// internal/utils/password.go
package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword generates a bcrypt hash of the plain-text password.
func HashPassword(password string) (string, error) {
	// bcrypt cost 14 is a reasonable default. Higher cost means more security but slower hashing.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost+4)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPasswordHash compares a plain-text password with its bcrypt hash.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
