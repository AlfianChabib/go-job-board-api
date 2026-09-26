package service

import (
	"AlfianChabib/go-job-board-api/internal/model/web"
	"context"

	"github.com/google/uuid"
)

type JobService interface {
	CreateJob(ctx context.Context, recruiterUserId uuid.UUID, req web.CreateJobRequest) (*web.JobResponse, error)
	GetJobs(ctx context.Context, req web.GetJobsRequest) (*web.JobPaginationResponse, error)
	GetRecruiterJobs(ctx context.Context, recruiterUserId uuid.UUID, req web.GetRecruiterJobsRequest) (*web.JobPaginationResponse, error)
	GetJobById(ctx context.Context, jobId uuid.UUID) (*web.JobResponse, error)
	UpdateJob(ctx context.Context, recruiterUserId uuid.UUID, jobId uuid.UUID, req web.UpdateJobRequest) (*web.JobResponse, error)
	UpdateJobStatus(ctx context.Context, recruiterUserId uuid.UUID, jobId uuid.UUID, status string) (*web.JobResponse, error)
	DeleteJob(ctx context.Context, recruiterUserId uuid.UUID, jobId uuid.UUID) error
}
