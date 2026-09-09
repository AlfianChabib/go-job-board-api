package web

import "AlfianChabib/go-job-board-api/internal/model/domain"

type RegisterRequest struct {
	Name     string          `json:"name" validate:"required,min=3,max=255"`
	Email    string          `json:"email" validate:"required,email"`
	Password string          `json:"password" validate:"required,min=8"`
	Role     domain.UserRole `json:"role" validate:"required,oneof=RECRUITER CANDIDATE"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LogOutRequest struct {
	RefreshToken string `cookie:"refresh_token" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `cookie:"refresh_token" json:"refresh_token" validate:"required"`
	AccessToken  string `header:"authorization" json:"access_token" validate:"required"`
}
