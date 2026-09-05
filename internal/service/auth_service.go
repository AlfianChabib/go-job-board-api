package service

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"AlfianChabib/go-job-board-api/internal/model/web"
	"context"
)

type AuthService interface {
	Register(ctx context.Context, data web.RegisterRequest) (*web.RegisterResponse, error)
	Login(ctx context.Context, data web.LoginRequest) (*domain.TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
	RefreshToken(ctx context.Context, data web.RefreshTokenRequest) (*web.RefreshTokenResponse, error)
}
