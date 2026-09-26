package repository

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"AlfianChabib/go-job-board-api/internal/model/web"
	"context"

	"github.com/google/uuid"
)

type JobRepository interface {
	Create(ctx context.Context, job domain.Job) (*domain.Job, error)
	FindAll(ctx context.Context, req web.GetJobsRequest) ([]domain.Job, int64, error)
	FindAllByCompanyId(ctx context.Context, companyId uuid.UUID, req web.GetRecruiterJobsRequest) ([]domain.Job, int64, error)
	FindById(ctx context.Context, id uuid.UUID) (*domain.Job, error)
	Update(ctx context.Context, job domain.Job) (*domain.Job, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.JobStatus) error
	Delete(ctx context.Context, id uuid.UUID) error
}
