package repository

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"context"

	"github.com/google/uuid"
)

type AuthRepository interface {
	Register(ctx context.Context, user domain.User) (*domain.User, error)
	IsUserExist(ctx context.Context, email string) bool
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindTokenWithUser(ctx context.Context, userId uuid.UUID, refreshToken string) (*domain.Token, error)
}
