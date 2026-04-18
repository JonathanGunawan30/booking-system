package dto

import "github.com/google/uuid"

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}

type UserResponse struct {
	UUID        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Role        string    `json:"role,omitempty"`
	PhoneNumber string    `json:"phone_number"`
}

type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

type RegisterRequest struct {
	Name            string `json:"name" validate:"required,min=3,max=50"`
	Username        string `json:"username" validate:"required,min=3,max=20"`
	Password        string `json:"password" validate:"required,min=6,max=100"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
	Email           string `json:"email" validate:"required,email,max=100"`
	PhoneNumber     string `json:"phone_number" validate:"required,min=10,max=15"`
	RoleID          uint
}

type RegisterResponse struct {
	User UserResponse `json:"user"`
}

type UpdateRequest struct {
	Name            string `json:"name" validate:"required,min=3,max=50"`
	Username        string `json:"username" validate:"required,min=3,max=20"`
	Password        string `json:"password,omitempty"`
	ConfirmPassword string `json:"confirm_password,omitempty"`
	Email           string `json:"email" validate:"required,email,max=100"`
	PhoneNumber     string `json:"phone_number" validate:"required,min=10,max=15"`
	RoleID          uint
}
