package repository

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"context"

	"github.com/google/uuid"
)

type CandidateRepository interface {
	Get(ctx context.Context, candidateId uuid.UUID) (*domain.Profile, error)
	Update(ctx context.Context, candidate domain.Profile) (*domain.Profile, error)
	UploadAvatar(ctx context.Context, userId uuid.UUID, avatarUrl string) error
}
