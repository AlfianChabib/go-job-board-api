package repository

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"AlfianChabib/go-job-board-api/internal/model/web"
	"context"

	"github.com/google/uuid"
)

type ApplicationRepository interface {
	Create(ctx context.Context, application domain.Application) (*domain.Application, error)
	FindById(ctx context.Context, id uuid.UUID) (*domain.Application, error)
	FindByJobAndCandidate(ctx context.Context, jobId uuid.UUID, candidateId uuid.UUID) (*domain.Application, error)
	FindAllByJobId(ctx context.Context, jobId uuid.UUID, req web.GetJobApplicationsRequest) ([]domain.Application, int64, error)
	FindAllByCandidateId(ctx context.Context, candidateId uuid.UUID, req web.GetCandidateApplicationsRequest) ([]domain.Application, int64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.ApplicationStatus, recruiterNotes *string, rejectionReason *string) (*domain.Application, error)
	Withdraw(ctx context.Context, id uuid.UUID, withdrawnReason *string) (*domain.Application, error)
}
