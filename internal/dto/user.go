package dto

import (
	"time"

	"github.com/euro1061/gohex/internal/domain"
)

// UserRegisterRequest represents the request body for user registration
type UserRegisterRequest struct {
	Name            string `json:"name" validate:"required"`
	Username        string `json:"username" validate:"required,min=3"`
	Password        string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
	Gender          string `json:"gender" validate:"required,oneof=male female"`
	Email           string `json:"email" validate:"required,email"`
}

// UserLoginRequest represents the request body for user login
type UserLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// UserUpdateRequest represents the request body for user update
type UserUpdateRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	Gender   string `json:"gender" validate:"required,oneof=male female"`
	Email    string `json:"email" validate:"required,email"`
}

// UserResponse represents the response body for user-related operations
type UserResponse struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Username    string     `json:"username"`
	Gender      string     `json:"gender"`
	Email       string     `json:"email"`
	LastLoginAt *time.Time `json:"last_login_at"`
}

// ToUser converts UserRegisterRequest to domain.User
func (r *UserRegisterRequest) ToUser() *domain.User {
	return &domain.User{
		Name:     r.Name,
		Username: r.Username,
		Password: r.Password, // Note: Password will be hashed in service layer
		Gender:   r.Gender,
		Email:    r.Email,
	}
}

// FromUser creates UserResponse from domain.User
func UserResponseFromUser(user *domain.User) *UserResponse {
	return &UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Username:    user.Username,
		Gender:      user.Gender,
		Email:       user.Email,
		LastLoginAt: user.LastLoginAt,
	}
}
