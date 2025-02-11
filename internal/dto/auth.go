package dto

import "github.com/google/uuid"

type SignUpRequestDto struct {
	Username string `json:"username" validate:"required,min=3`
	Password string `json:"password" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
}

type SignUpResponeDto struct {
	ClientId int       `json:"client_id"`
	UserId   uuid.UUID `json:"user_id"`
}

type SignInRequestDto struct {
	Username string `json:"username" validate:"required,min=3`
	Password string `json:"password" validate:"required,min=3"`
}
