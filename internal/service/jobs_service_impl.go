package service

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"AlfianChabib/go-job-board-api/internal/model/web"
	"AlfianChabib/go-job-board-api/internal/repository"
	"AlfianChabib/go-job-board-api/pkg/errs"
	"context"
	"errors"
	"math"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type jobServiceImpl struct {
	jobRepo     repository.JobRepository
	companyRepo repository.CompanyRepository
}

func NewJobService(jobRepo repository.JobRepository, companyRepo repository.CompanyRepository) JobService {
	return &jobServiceImpl{
		jobRepo:     jobRepo,
		companyRepo: companyRepo,
	}
}

func (s *jobServiceImpl) CreateJob(ctx context.Context, recruiterUserId uuid.UUID, req web.CreateJobRequest) (*web.JobResponse, error) {
	company, err := s.companyRepo.FindByRecruiterId(ctx, recruiterUserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrCompanyNotFound
		}
		return nil, errs.ErrInternalServer
	}

	if req.MinSalary != nil && req.MaxSalary != nil && *req.MaxSalary < *req.MinSalary {
		return nil, errs.ErrInvalidSalaryRange
	}

	currency := req.Currency
	if currency == "" {
		currency = "IDR"
	}

	job := domain.Job{
		CompanyId:          company.ID,
		Title:              req.Title,
		Description:        req.Description,
		Requirements:       req.Requirements,
		EmploymentType:     domain.EmploymentType(req.EmploymentType),
		WorkMode:           domain.WorkMode(req.WorkMode),
		Location:           req.Location,
		MinSalary:          req.MinSalary,
		MaxSalary:          req.MaxSalary,
		Currency:           currency,
		IsSalaryNegotiable: req.IsSalaryNegotiable,
		Status:             domain.JobStatusOpen,
	}

	created, err := s.jobRepo.Create(ctx, job)
	if err != nil {
		return nil, errs.ErrInternalServer
	}

	created.Company = company
	return toJobResponse(created), nil
}

func (s *jobServiceImpl) GetJobs(ctx context.Context, req web.GetJobsRequest) (*web.JobPaginationResponse, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	jobs, total, err := s.jobRepo.FindAll(ctx, req)
	if err != nil {
		return nil, errs.ErrInternalServer
	}

	jobResponses := make([]web.JobResponse, 0, len(jobs))
	for i := range jobs {
		jobResponses = append(jobResponses, *toJobResponse(&jobs[i]))
	}

	totalPages := 0
	if limit > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(limit)))
	}

	return &web.JobPaginationResponse{
		Jobs:       jobResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *jobServiceImpl) GetRecruiterJobs(ctx context.Context, recruiterUserId uuid.UUID, req web.GetRecruiterJobsRequest) (*web.JobPaginationResponse, error) {
	company, err := s.companyRepo.FindByRecruiterId(ctx, recruiterUserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrCompanyNotFound
		}
		return nil, errs.ErrInternalServer
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	jobs, total, err := s.jobRepo.FindAllByCompanyId(ctx, company.ID, req)
	if err != nil {
		return nil, errs.ErrInternalServer
	}

	jobResponses := make([]web.JobResponse, 0, len(jobs))
	for i := range jobs {
		jobResponses = append(jobResponses, *toJobResponse(&jobs[i]))
	}

	totalPages := 0
	if limit > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(limit)))
	}

	return &web.JobPaginationResponse{
		Jobs:       jobResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *jobServiceImpl) GetJobById(ctx context.Context, jobId uuid.UUID) (*web.JobResponse, error) {
	job, err := s.jobRepo.FindById(ctx, jobId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrJobNotFound
		}
		return nil, errs.ErrInternalServer
	}

	return toJobResponse(job), nil
}

func (s *jobServiceImpl) UpdateJob(ctx context.Context, recruiterUserId uuid.UUID, jobId uuid.UUID, req web.UpdateJobRequest) (*web.JobResponse, error) {
	company, err := s.companyRepo.FindByRecruiterId(ctx, recruiterUserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrCompanyNotFound
		}
		return nil, errs.ErrInternalServer
	}

	existingJob, err := s.jobRepo.FindById(ctx, jobId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrJobNotFound
		}
		return nil, errs.ErrInternalServer
	}

	// IDOR check: Verify job belongs to recruiter's company
	if existingJob.CompanyId != company.ID {
		return nil, errs.ErrJobForbidden
	}

	if req.MinSalary != nil && req.MaxSalary != nil && *req.MaxSalary < *req.MinSalary {
		return nil, errs.ErrInvalidSalaryRange
	}

	currency := req.Currency
	if currency == "" {
		currency = "IDR"
	}

	existingJob.Title = req.Title
	existingJob.Description = req.Description
	existingJob.Requirements = req.Requirements
	existingJob.EmploymentType = domain.EmploymentType(req.EmploymentType)
	existingJob.WorkMode = domain.WorkMode(req.WorkMode)
	existingJob.Location = req.Location
	existingJob.MinSalary = req.MinSalary
	existingJob.MaxSalary = req.MaxSalary
	existingJob.Currency = currency
	existingJob.IsSalaryNegotiable = req.IsSalaryNegotiable

	updated, err := s.jobRepo.Update(ctx, *existingJob)
	if err != nil {
		return nil, errs.ErrInternalServer
	}

	updated.Company = company
	return toJobResponse(updated), nil
}

func (s *jobServiceImpl) UpdateJobStatus(ctx context.Context, recruiterUserId uuid.UUID, jobId uuid.UUID, status string) (*web.JobResponse, error) {
	if status != string(domain.JobStatusOpen) && status != string(domain.JobStatusClosed) {
		return nil, errs.ErrInvalidJobStatus
	}

	company, err := s.companyRepo.FindByRecruiterId(ctx, recruiterUserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrCompanyNotFound
		}
		return nil, errs.ErrInternalServer
	}

	existingJob, err := s.jobRepo.FindById(ctx, jobId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrJobNotFound
		}
		return nil, errs.ErrInternalServer
	}

	// IDOR check: Verify job belongs to recruiter's company
	if existingJob.CompanyId != company.ID {
		return nil, errs.ErrJobForbidden
	}

	jobStatus := domain.JobStatus(status)
	if err := s.jobRepo.UpdateStatus(ctx, jobId, jobStatus); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrJobNotFound
		}
		return nil, errs.ErrInternalServer
	}

	existingJob.Status = jobStatus
	existingJob.Company = company
	return toJobResponse(existingJob), nil
}

func (s *jobServiceImpl) DeleteJob(ctx context.Context, recruiterUserId uuid.UUID, jobId uuid.UUID) error {
	company, err := s.companyRepo.FindByRecruiterId(ctx, recruiterUserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.ErrCompanyNotFound
		}
		return errs.ErrInternalServer
	}

	existingJob, err := s.jobRepo.FindById(ctx, jobId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.ErrJobNotFound
		}
		return errs.ErrInternalServer
	}

	// IDOR check: Verify job belongs to recruiter's company
	if existingJob.CompanyId != company.ID {
		return errs.ErrJobForbidden
	}

	if err := s.jobRepo.Delete(ctx, jobId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.ErrJobNotFound
		}
		return errs.ErrInternalServer
	}

	return nil
}

func toJobResponse(job *domain.Job) *web.JobResponse {
	if job == nil {
		return nil
	}

	var companySummary *web.JobCompanySummary
	if job.Company != nil {
		companySummary = &web.JobCompanySummary{
			ID:       job.Company.ID,
			Name:     job.Company.Name,
			LogoUrl:  job.Company.LogoUrl,
			Location: job.Company.Location,
		}
	}

	return &web.JobResponse{
		ID:                 job.ID,
		CompanyId:          job.CompanyId,
		Title:              job.Title,
		Description:        job.Description,
		Requirements:       job.Requirements,
		EmploymentType:     job.EmploymentType,
		WorkMode:           job.WorkMode,
		Location:           job.Location,
		MinSalary:          job.MinSalary,
		MaxSalary:          job.MaxSalary,
		Currency:           job.Currency,
		IsSalaryNegotiable: job.IsSalaryNegotiable,
		Status:             job.Status,
		CreatedAt:          job.CreatedAt,
		UpdatedAt:          job.UpdatedAt,
		Company:            companySummary,
	}
}
