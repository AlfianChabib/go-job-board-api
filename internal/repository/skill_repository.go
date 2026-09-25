package repository

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"context"
)

type SkillRepository interface {
	FindAll(ctx context.Context) ([]domain.Skill, error)
}
