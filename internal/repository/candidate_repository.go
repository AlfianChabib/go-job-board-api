package repository

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"AlfianChabib/go-job-board-api/internal/model/web"
	"context"

	"github.com/google/uuid"
)

type CandidateRepository interface {
	Get(ctx context.Context, userId uuid.UUID) (*domain.Profile, error)
	Update(ctx context.Context, candidate domain.Profile) (*domain.Profile, error)
	UploadAvatar(ctx context.Context, userId uuid.UUID, avatarUrl string) error
	DeleteAvatar(ctx context.Context, userId uuid.UUID) error
	UpdateSkills(ctx context.Context, userId uuid.UUID, skills web.UpdateCandidateSkillsRequest) (*[]domain.Skill, error)
	GetExperiences(ctx context.Context, userId uuid.UUID) (*[]domain.Experience, error)
	CreateExperience(ctx context.Context, userId uuid.UUID, experience domain.Experience) error
	UpdateExperience(ctx context.Context, experienceId uuid.UUID, experience domain.Experience) error
}
