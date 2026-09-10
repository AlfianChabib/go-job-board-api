package service

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"AlfianChabib/go-job-board-api/internal/model/web"
	"context"

	"github.com/google/uuid"
)

type CandidateService interface {
	Get(ctx context.Context, userId uuid.UUID) (*web.GetCandidateResponse, error)
	Update(ctx context.Context, candidate domain.Profile) (*web.UpdateCandidateResponse, error)
	UploadAvatar(ctx context.Context, req web.UpdateCandidateAvatarRequest) (*string, error)
}
