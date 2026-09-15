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
	DeleteAvatar(ctx context.Context, userId uuid.UUID) error
	UpdateSkills(ctx context.Context, userId uuid.UUID, skills web.UpdateCandidateSkillsRequest) (*[]domain.Skill, error)
	GetExperiences(ctx context.Context, userId uuid.UUID) (*[]domain.Experience, error)
}
