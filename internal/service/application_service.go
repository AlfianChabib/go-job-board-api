package service

import (
	"AlfianChabib/go-job-board-api/internal/model/web"
	"context"

	"github.com/google/uuid"
)

type ApplicationService interface {
	ApplyJob(ctx context.Context, candidateUserId uuid.UUID, req web.ApplyJobRequest) (*web.CandidateApplicationResponse, error)
	GetJobApplications(ctx context.Context, recruiterUserId uuid.UUID, req web.GetJobApplicationsRequest) (*web.RecruiterApplicationPaginationResponse, error)
	GetCandidateApplications(ctx context.Context, candidateUserId uuid.UUID, req web.GetCandidateApplicationsRequest) (*web.CandidateApplicationPaginationResponse, error)
	GetApplicationById(ctx context.Context, userId uuid.UUID, userRole string, applicationId uuid.UUID) (any, error)
	UpdateStatus(ctx context.Context, recruiterUserId uuid.UUID, req web.UpdateApplicationStatusRequest) (*web.RecruiterApplicationResponse, error)
	WithdrawApplication(ctx context.Context, candidateUserId uuid.UUID, req web.WithdrawApplicationRequest) (*web.CandidateApplicationResponse, error)
}
