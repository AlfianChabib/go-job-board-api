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

type applicationServiceImpl struct {
	applicationRepo repository.ApplicationRepository
	jobRepo         repository.JobRepository
	companyRepo     repository.CompanyRepository
}

func NewApplicationService(
	applicationRepo repository.ApplicationRepository,
	jobRepo repository.JobRepository,
	companyRepo repository.CompanyRepository,
) ApplicationService {
	return &applicationServiceImpl{
		applicationRepo: applicationRepo,
		jobRepo:         jobRepo,
		companyRepo:     companyRepo,
	}
}

func (s *applicationServiceImpl) ApplyJob(ctx context.Context, candidateUserId uuid.UUID, req web.ApplyJobRequest) (*web.CandidateApplicationResponse, error) {
	job, err := s.jobRepo.FindById(ctx, req.JobId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrJobNotFound
		}
		return nil, errs.ErrInternalServer
	}

	if job.Status != domain.JobStatusOpen {
		return nil, errs.ErrJobClosed
	}

	existing, err := s.applicationRepo.FindByJobAndCandidate(ctx, req.JobId, candidateUserId)
	if err == nil && existing != nil {
		return nil, errs.ErrAlreadyApplied
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.ErrInternalServer
	}

	app := domain.Application{
		JobId:          req.JobId,
		CandidateId:    candidateUserId,
		ResumeUrl:      req.ResumeUrl,
		ResumeFilename: req.ResumeFilename,
		CoverLetter:    req.CoverLetter,
		ExpectedSalary: req.ExpectedSalary,
		Status:         domain.ApplicationStatusApplied,
	}

	created, err := s.applicationRepo.Create(ctx, app)
	if err != nil {
		return nil, errs.ErrInternalServer
	}

	created.Job = job
	return toCandidateApplicationResponse(created), nil
}

func (s *applicationServiceImpl) GetJobApplications(ctx context.Context, recruiterUserId uuid.UUID, req web.GetJobApplicationsRequest) (*web.RecruiterApplicationPaginationResponse, error) {
	company, err := s.companyRepo.FindByRecruiterId(ctx, recruiterUserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrCompanyNotFound
		}
		return nil, errs.ErrInternalServer
	}

	job, err := s.jobRepo.FindById(ctx, req.JobId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrJobNotFound
		}
		return nil, errs.ErrInternalServer
	}

	if job.CompanyId != company.ID {
		return nil, errs.ErrJobForbidden
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	apps, total, err := s.applicationRepo.FindAllByJobId(ctx, req.JobId, req)
	if err != nil {
		return nil, errs.ErrInternalServer
	}

	appResponses := make([]web.RecruiterApplicationResponse, 0, len(apps))
	for i := range apps {
		appResponses = append(appResponses, *toRecruiterApplicationResponse(&apps[i]))
	}

	totalPages := 0
	if limit > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(limit)))
	}

	return &web.RecruiterApplicationPaginationResponse{
		Applications: appResponses,
		Total:        total,
		Page:         page,
		Limit:        limit,
		TotalPages:   totalPages,
	}, nil
}

func (s *applicationServiceImpl) GetCandidateApplications(ctx context.Context, candidateUserId uuid.UUID, req web.GetCandidateApplicationsRequest) (*web.CandidateApplicationPaginationResponse, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	apps, total, err := s.applicationRepo.FindAllByCandidateId(ctx, candidateUserId, req)
	if err != nil {
		return nil, errs.ErrInternalServer
	}

	appResponses := make([]web.CandidateApplicationResponse, 0, len(apps))
	for i := range apps {
		appResponses = append(appResponses, *toCandidateApplicationResponse(&apps[i]))
	}

	totalPages := 0
	if limit > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(limit)))
	}

	return &web.CandidateApplicationPaginationResponse{
		Applications: appResponses,
		Total:        total,
		Page:         page,
		Limit:        limit,
		TotalPages:   totalPages,
	}, nil
}

func (s *applicationServiceImpl) GetApplicationById(ctx context.Context, userId uuid.UUID, userRole string, applicationId uuid.UUID) (any, error) {
	app, err := s.applicationRepo.FindById(ctx, applicationId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrApplicationNotFound
		}
		return nil, errs.ErrInternalServer
	}

	if userRole == string(domain.RoleCandidate) {
		if app.CandidateId != userId {
			return nil, errs.ErrApplicationForbidden
		}
		return toCandidateApplicationResponse(app), nil
	}

	if userRole == string(domain.RoleRecruiter) {
		company, err := s.companyRepo.FindByRecruiterId(ctx, userId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errs.ErrCompanyNotFound
			}
			return nil, errs.ErrInternalServer
		}
		if app.Job == nil || app.Job.CompanyId != company.ID {
			return nil, errs.ErrApplicationForbidden
		}
		return toRecruiterApplicationResponse(app), nil
	}

	return nil, errs.ErrForbidden
}

func (s *applicationServiceImpl) UpdateStatus(ctx context.Context, recruiterUserId uuid.UUID, req web.UpdateApplicationStatusRequest) (*web.RecruiterApplicationResponse, error) {
	app, err := s.applicationRepo.FindById(ctx, req.ApplicationId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrApplicationNotFound
		}
		return nil, errs.ErrInternalServer
	}

	company, err := s.companyRepo.FindByRecruiterId(ctx, recruiterUserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrCompanyNotFound
		}
		return nil, errs.ErrInternalServer
	}

	if app.Job == nil || app.Job.CompanyId != company.ID {
		return nil, errs.ErrApplicationForbidden
	}

	if app.Status == domain.ApplicationStatusWithdrawn {
		return nil, errs.ErrApplicationAlreadyWithdrawn
	}

	switch domain.ApplicationStatus(req.Status) {
	case domain.ApplicationStatusReviewing,
		domain.ApplicationStatusShortlisted,
		domain.ApplicationStatusInterviewing,
		domain.ApplicationStatusOffered,
		domain.ApplicationStatusHired,
		domain.ApplicationStatusRejected:
		// valid
	default:
		return nil, errs.ErrInvalidApplicationStatus
	}

	updated, err := s.applicationRepo.UpdateStatus(ctx, req.ApplicationId, domain.ApplicationStatus(req.Status), req.RecruiterNotes, req.RejectionReason)
	if err != nil {
		return nil, errs.ErrInternalServer
	}

	return toRecruiterApplicationResponse(updated), nil
}

func (s *applicationServiceImpl) WithdrawApplication(ctx context.Context, candidateUserId uuid.UUID, req web.WithdrawApplicationRequest) (*web.CandidateApplicationResponse, error) {
	app, err := s.applicationRepo.FindById(ctx, req.ApplicationId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrApplicationNotFound
		}
		return nil, errs.ErrInternalServer
	}

	if app.CandidateId != candidateUserId {
		return nil, errs.ErrApplicationForbidden
	}

	if app.Status == domain.ApplicationStatusWithdrawn {
		return nil, errs.ErrApplicationAlreadyWithdrawn
	}

	if app.Status == domain.ApplicationStatusHired {
		return nil, errs.ErrInvalidAction
	}

	updated, err := s.applicationRepo.Withdraw(ctx, req.ApplicationId, req.WithdrawnReason)
	if err != nil {
		return nil, errs.ErrInternalServer
	}

	return toCandidateApplicationResponse(updated), nil
}

func toCandidateApplicationResponse(app *domain.Application) *web.CandidateApplicationResponse {
	if app == nil {
		return nil
	}

	var jobResponse *web.JobResponse
	if app.Job != nil {
		jobResponse = toJobResponse(app.Job)
	}

	return &web.CandidateApplicationResponse{
		ID:              app.ID,
		JobId:           app.JobId,
		Job:             jobResponse,
		ResumeUrl:       app.ResumeUrl,
		ResumeFilename:  app.ResumeFilename,
		CoverLetter:     app.CoverLetter,
		ExpectedSalary:  app.ExpectedSalary,
		Status:          string(app.Status),
		RejectionReason: app.RejectionReason,
		WithdrawnReason: app.WithdrawnReason,
		AppliedAt:       app.AppliedAt,
		StatusUpdatedAt: app.StatusUpdatedAt,
	}
}

func toRecruiterApplicationResponse(app *domain.Application) *web.RecruiterApplicationResponse {
	if app == nil {
		return nil
	}

	var candidateSummary *web.CandidateSummary
	if app.Candidate != nil {
		candidateSummary = &web.CandidateSummary{
			ID:    app.Candidate.ID,
			Name:  app.Candidate.Name,
			Email: app.Candidate.Email,
		}
		if app.Candidate.Profile != nil {
			candidateSummary.Headline = app.Candidate.Profile.Headline
			candidateSummary.Phone = app.Candidate.Profile.Phone
			candidateSummary.AvatarUrl = app.Candidate.Profile.AvatarUrl
		}
	}

	return &web.RecruiterApplicationResponse{
		ID:              app.ID,
		JobId:           app.JobId,
		CandidateId:     app.CandidateId,
		Candidate:       candidateSummary,
		ResumeUrl:       app.ResumeUrl,
		ResumeFilename:  app.ResumeFilename,
		CoverLetter:     app.CoverLetter,
		ExpectedSalary:  app.ExpectedSalary,
		Status:          string(app.Status),
		RecruiterNotes:  app.RecruiterNotes,
		RejectionReason: app.RejectionReason,
		WithdrawnReason: app.WithdrawnReason,
		AppliedAt:       app.AppliedAt,
		StatusUpdatedAt: app.StatusUpdatedAt,
	}
}
