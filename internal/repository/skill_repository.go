package repository

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"AlfianChabib/go-job-board-api/internal/model/web"
	"context"
)

type SkillRepository interface {
	FindAll(ctx context.Context, req web.GetDataRequest) ([]domain.Skill, error)
}
