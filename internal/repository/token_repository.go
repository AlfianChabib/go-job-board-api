package repository

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"context"

	"github.com/google/uuid"
)

type TokenRepository interface {
	Save(ctx context.Context, token domain.Token) (*domain.Token, error)
	RevokeToken(ctx context.Context, refreshToken string) error
	FindTokenWithUser(ctx context.Context, userId uuid.UUID, refreshToken string) (*domain.Token, error)
}
