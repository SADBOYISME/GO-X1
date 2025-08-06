package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique;not null" validate:"required,min=3,max=50"`
	Email     string    `json:"email" gorm:"unique;not null" validate:"required,email"`
	Password  string    `json:"password" gorm:"not null" validate:"required,min=6"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserResponse is used for API responses (excludes sensitive data)
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateUserRequest is used for creating new users
type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// UpdateUserRequest is used for updating users
type UpdateUserRequest struct {
	Username string `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
	Password string `json:"password,omitempty" validate:"omitempty,min=6"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}