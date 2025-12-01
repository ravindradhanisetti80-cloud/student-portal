// internal/models/user.go
package models

import (
	"time"
)

// User represents the structure of the users table in the database.
type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Omit from JSON response
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RegisterRequest is the structure for the registration request body.
type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"required,oneof=student admin"`
}

// LoginRequest is the structure for the login request body.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UpdateProfileRequest is the structure for the update profile request body.
type UpdateProfileRequest struct {
	Name  *string `json:"name"`
	Email *string `json:"email" validate:"omitempty,email"`
}

// UpdateUserRequest is the structure for the admin update user request body.
type UpdateUserRequest struct {
	Name  *string `json:"name"`
	Email *string `json:"email" validate:"omitempty,email"`
	Role  *string `json:"role" validate:"omitempty,oneof=student admin"`
}

// UserResponse is the standard response structure for a User, omitting the password.
type UserResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LoginResponse contains the JWT token and user info.
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// ToResponse converts a User model to a UserResponse DTO.
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// ToResponsePtr converts a User model to a pointer to UserResponse DTO.
// This is used primarily in the Service Layer to return models to Handlers.
func (u *User) ToResponsePtr() *UserResponse {
	res := u.ToResponse()
	return &res
}
