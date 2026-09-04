package repository

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"context"
)

type TokenRepository interface {
	Save(ctx context.Context, token domain.Token) (*domain.Token, error)
}
