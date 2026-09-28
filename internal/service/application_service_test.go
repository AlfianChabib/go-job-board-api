package service

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"AlfianChabib/go-job-board-api/internal/model/web"
	"AlfianChabib/go-job-board-api/pkg/errs"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestApplicationService_ApplyJob(t *testing.T) {
	ctx := context.Background()
	candidateID := uuid.New()
	jobID := uuid.New()
	req := web.ApplyJobRequest{
		JobId:     jobID,
		ResumeUrl: "https://example.com/cv.pdf",
	}

	t.Run("Job not found, generic error, and job closed", func(t *testing.T) {
		appRepo := new(mockApplicationRepository)
		jobRepo := new(mockJobRepository)
		compRepo := new(mockCompanyRepository)
		svc := NewApplicationService(appRepo, jobRepo, compRepo)

		jobRepo.On("FindById", ctx, jobID).Return(nil, gorm.ErrRecordNotFound).Once()
		_, err := svc.ApplyJob(ctx, candidateID, req)
		assert.Equal(t, errs.ErrJobNotFound, err)

		jobRepo.On("FindById", ctx, jobID).Return(nil, errors.New("db err")).Once()
		_, err = svc.ApplyJob(ctx, candidateID, req)
		assert.Equal(t, errs.ErrInternalServer, err)

		jobRepo.On("FindById", ctx, jobID).Return(&domain.Job{ID: jobID, Status: domain.JobStatusClosed}, nil).Once()
		_, err = svc.ApplyJob(ctx, candidateID, req)
		assert.Equal(t, errs.ErrJobClosed, err)
	})

	t.Run("Already applied, find error, create error, and success", func(t *testing.T) {
		appRepo := new(mockApplicationRepository)
		jobRepo := new(mockJobRepository)
		compRepo := new(mockCompanyRepository)
		svc := NewApplicationService(appRepo, jobRepo, compRepo)

		openJob := &domain.Job{ID: jobID, Status: domain.JobStatusOpen, Title: "Go Engineer"}
		jobRepo.On("FindById", ctx, jobID).Return(openJob, nil)

		// 1. Already applied
		appRepo.On("FindByJobAndCandidate", ctx, jobID, candidateID).Return(&domain.Application{ID: uuid.New()}, nil).Once()
		_, err := svc.ApplyJob(ctx, candidateID, req)
		assert.Equal(t, errs.ErrAlreadyApplied, err)

		// 2. FindByJobAndCandidate unexpected error
		appRepo.On("FindByJobAndCandidate", ctx, jobID, candidateID).Return(nil, errors.New("db err")).Once()
		_, err = svc.ApplyJob(ctx, candidateID, req)
		assert.Equal(t, errs.ErrInternalServer, err)

		// 3. Create error
		appRepo.On("FindByJobAndCandidate", ctx, jobID, candidateID).Return(nil, gorm.ErrRecordNotFound).Once()
		appRepo.On("Create", ctx, mock.AnythingOfType("domain.Application")).Return(nil, errors.New("create err")).Once()
		_, err = svc.ApplyJob(ctx, candidateID, req)
		assert.Equal(t, errs.ErrInternalServer, err)

		// 4. Success
		appID := uuid.New()
		appRepo.On("FindByJobAndCandidate", ctx, jobID, candidateID).Return(nil, gorm.ErrRecordNotFound).Once()
		appRepo.On("Create", ctx, mock.AnythingOfType("domain.Application")).Return(&domain.Application{
			ID:          appID,
			JobId:       jobID,
			CandidateId: candidateID,
			ResumeUrl:   req.ResumeUrl,
			Status:      domain.ApplicationStatusApplied,
		}, nil).Once()

		res, err := svc.ApplyJob(ctx, candidateID, req)
		require.NoError(t, err)
		assert.Equal(t, appID, res.ID)
		assert.Equal(t, "APPLIED", res.Status)
		require.NotNil(t, res.Job)
		assert.Equal(t, jobID, res.Job.ID)
	})
}

func TestApplicationService_GetApplicationsList(t *testing.T) {
	ctx := context.Background()
	recruiterID := uuid.New()
	candidateID := uuid.New()
	companyID := uuid.New()
	jobID := uuid.New()

	t.Run("GetJobApplications scenarios", func(t *testing.T) {
		appRepo := new(mockApplicationRepository)
		jobRepo := new(mockJobRepository)
		compRepo := new(mockCompanyRepository)
		svc := NewApplicationService(appRepo, jobRepo, compRepo)
		req := web.GetJobApplicationsRequest{JobId: jobID}

		// 1. Company not found
		compRepo.On("FindByRecruiterId", ctx, recruiterID).Return(nil, gorm.ErrRecordNotFound).Once()
		_, err := svc.GetJobApplications(ctx, recruiterID, req)
		assert.Equal(t, errs.ErrCompanyNotFound, err)

		// 2. Company generic error
		compRepo.On("FindByRecruiterId", ctx, recruiterID).Return(nil, errors.New("db err")).Once()
		_, err = svc.GetJobApplications(ctx, recruiterID, req)
		assert.Equal(t, errs.ErrInternalServer, err)

		// 3. Job not found
		compRepo.On("FindByRecruiterId", ctx, recruiterID).Return(&domain.Company{ID: companyID}, nil)
		jobRepo.On("FindById", ctx, jobID).Return(nil, gorm.ErrRecordNotFound).Once()
		_, err = svc.GetJobApplications(ctx, recruiterID, req)
		assert.Equal(t, errs.ErrJobNotFound, err)

		// 4. Job generic error
		jobRepo.On("FindById", ctx, jobID).Return(nil, errors.New("db err")).Once()
		_, err = svc.GetJobApplications(ctx, recruiterID, req)
		assert.Equal(t, errs.ErrInternalServer, err)

		// 5. Job belongs to another company
		jobRepo.On("FindById", ctx, jobID).Return(&domain.Job{ID: jobID, CompanyId: uuid.New()}, nil).Once()
		_, err = svc.GetJobApplications(ctx, recruiterID, req)
		assert.Equal(t, errs.ErrJobForbidden, err)

		// 6. Repo FindAllByJobId error
		jobRepo.On("FindById", ctx, jobID).Return(&domain.Job{ID: jobID, CompanyId: companyID}, nil)
		appRepo.On("FindAllByJobId", ctx, jobID, req).Return(nil, int64(0), errors.New("db err")).Once()
		_, err = svc.GetJobApplications(ctx, recruiterID, req)
		assert.Equal(t, errs.ErrInternalServer, err)

		// 7. Success with Candidate Profile populated
		headline := "Go Dev"
		apps := []domain.Application{
			{
				ID:          uuid.New(),
				JobId:       jobID,
				CandidateId: candidateID,
				Status:      domain.ApplicationStatusApplied,
				Candidate: &domain.User{
					ID:    candidateID,
					Name:  "Alice",
					Email: "alice@example.com",
					Profile: &domain.Profile{
						Headline: &headline,
					},
				},
			},
		}
		appRepo.On("FindAllByJobId", ctx, jobID, req).Return(apps, int64(1), nil).Once()
		res, err := svc.GetJobApplications(ctx, recruiterID, req)
		require.NoError(t, err)
		assert.Len(t, res.Applications, 1)
		require.NotNil(t, res.Applications[0].Candidate)
		assert.Equal(t, &headline, res.Applications[0].Candidate.Headline)
	})

	t.Run("GetCandidateApplications error and success", func(t *testing.T) {
		appRepo := new(mockApplicationRepository)
		svc := NewApplicationService(appRepo, new(mockJobRepository), new(mockCompanyRepository))
		req := web.GetCandidateApplicationsRequest{}

		appRepo.On("FindAllByCandidateId", ctx, candidateID, req).Return(nil, int64(0), errors.New("db err")).Once()
		_, err := svc.GetCandidateApplications(ctx, candidateID, req)
		assert.Equal(t, errs.ErrInternalServer, err)

		appRepo.On("FindAllByCandidateId", ctx, candidateID, req).Return([]domain.Application{
			{ID: uuid.New(), CandidateId: candidateID, Status: domain.ApplicationStatusApplied},
		}, int64(1), nil).Once()
		res, err := svc.GetCandidateApplications(ctx, candidateID, req)
		require.NoError(t, err)
		assert.Len(t, res.Applications, 1)
	})
}

func TestApplicationService_GetByIdUpdateAndWithdraw(t *testing.T) {
	ctx := context.Background()
	candidateID := uuid.New()
	recruiterID := uuid.New()
	companyID := uuid.New()
	appID := uuid.New()
	jobID := uuid.New()

	t.Run("GetApplicationById for Candidate, Recruiter, and Invalid Role", func(t *testing.T) {
		appRepo := new(mockApplicationRepository)
		compRepo := new(mockCompanyRepository)
		svc := NewApplicationService(appRepo, new(mockJobRepository), compRepo)

		// 1. Not found & generic error
		appRepo.On("FindById", ctx, appID).Return(nil, gorm.ErrRecordNotFound).Once()
		_, err := svc.GetApplicationById(ctx, candidateID, string(domain.RoleCandidate), appID)
		assert.Equal(t, errs.ErrApplicationNotFound, err)

		appRepo.On("FindById", ctx, appID).Return(nil, errors.New("db err")).Once()
		_, err = svc.GetApplicationById(ctx, candidateID, string(domain.RoleCandidate), appID)
		assert.Equal(t, errs.ErrInternalServer, err)

		// 2. Candidate forbidden & success
		app := &domain.Application{
			ID:          appID,
			CandidateId: candidateID,
			JobId:       jobID,
			Job:         &domain.Job{ID: jobID, CompanyId: companyID},
		}
		appRepo.On("FindById", ctx, appID).Return(app, nil)

		_, err = svc.GetApplicationById(ctx, uuid.New(), string(domain.RoleCandidate), appID)
		assert.Equal(t, errs.ErrApplicationForbidden, err)

		res, err := svc.GetApplicationById(ctx, candidateID, string(domain.RoleCandidate), appID)
		require.NoError(t, err)
		assert.IsType(t, &web.CandidateApplicationResponse{}, res)

		// 3. Recruiter company not found, error, forbidden, and success
		compRepo.On("FindByRecruiterId", ctx, recruiterID).Return(nil, gorm.ErrRecordNotFound).Once()
		_, err = svc.GetApplicationById(ctx, recruiterID, string(domain.RoleRecruiter), appID)
		assert.Equal(t, errs.ErrCompanyNotFound, err)

		compRepo.On("FindByRecruiterId", ctx, recruiterID).Return(nil, errors.New("db err")).Once()
		_, err = svc.GetApplicationById(ctx, recruiterID, string(domain.RoleRecruiter), appID)
		assert.Equal(t, errs.ErrInternalServer, err)

		compRepo.On("FindByRecruiterId", ctx, recruiterID).Return(&domain.Company{ID: uuid.New()}, nil).Once()
		_, err = svc.GetApplicationById(ctx, recruiterID, string(domain.RoleRecruiter), appID)
		assert.Equal(t, errs.ErrApplicationForbidden, err)

		compRepo.On("FindByRecruiterId", ctx, recruiterID).Return(&domain.Company{ID: companyID}, nil).Once()
		resRec, err := svc.GetApplicationById(ctx, recruiterID, string(domain.RoleRecruiter), appID)
		require.NoError(t, err)
		assert.IsType(t, &web.RecruiterApplicationResponse{}, resRec)

		// 4. Unknown role returns ErrForbidden
		_, err = svc.GetApplicationById(ctx, recruiterID, "UNKNOWN_ROLE", appID)
		assert.Equal(t, errs.ErrForbidden, err)
	})

	t.Run("UpdateStatus scenarios", func(t *testing.T) {
		appRepo := new(mockApplicationRepository)
		compRepo := new(mockCompanyRepository)
		svc := NewApplicationService(appRepo, new(mockJobRepository), compRepo)

		// 1. Not found & generic error
		appRepo.On("FindById", ctx, appID).Return(nil, gorm.ErrRecordNotFound).Once()
		_, err := svc.UpdateStatus(ctx, recruiterID, web.UpdateApplicationStatusRequest{ApplicationId: appID})
		assert.Equal(t, errs.ErrApplicationNotFound, err)

		appRepo.On("FindById", ctx, appID).Return(nil, errors.New("db err")).Once()
		_, err = svc.UpdateStatus(ctx, recruiterID, web.UpdateApplicationStatusRequest{ApplicationId: appID})
		assert.Equal(t, errs.ErrInternalServer, err)

		// 2. Company not found & generic error
		activeApp := &domain.Application{
			ID:     appID,
			Status: domain.ApplicationStatusApplied,
			Job:    &domain.Job{ID: jobID, CompanyId: companyID},
		}
		appRepo.On("FindById", ctx, appID).Return(activeApp, nil).Once()
		compRepo.On("FindByRecruiterId", ctx, recruiterID).Return(nil, gorm.ErrRecordNotFound).Once()
		_, err = svc.UpdateStatus(ctx, recruiterID, web.UpdateApplicationStatusRequest{ApplicationId: appID})
		assert.Equal(t, errs.ErrCompanyNotFound, err)

		appRepo.On("FindById", ctx, appID).Return(activeApp, nil).Once()
		compRepo.On("FindByRecruiterId", ctx, recruiterID).Return(nil, errors.New("db err")).Once()
		_, err = svc.UpdateStatus(ctx, recruiterID, web.UpdateApplicationStatusRequest{ApplicationId: appID})
		assert.Equal(t, errs.ErrInternalServer, err)

		// 3. Forbidden (different company)
		appRepo.On("FindById", ctx, appID).Return(activeApp, nil).Once()
		compRepo.On("FindByRecruiterId", ctx, recruiterID).Return(&domain.Company{ID: uuid.New()}, nil).Once()
		_, err = svc.UpdateStatus(ctx, recruiterID, web.UpdateApplicationStatusRequest{ApplicationId: appID})
		assert.Equal(t, errs.ErrApplicationForbidden, err)

		// 4. Already withdrawn
		compRepo.On("FindByRecruiterId", ctx, recruiterID).Return(&domain.Company{ID: companyID}, nil)
		appRepo.On("FindById", ctx, appID).Return(&domain.Application{
			ID:     appID,
			Status: domain.ApplicationStatusWithdrawn,
			Job:    &domain.Job{ID: jobID, CompanyId: companyID},
		}, nil).Once()
		_, err = svc.UpdateStatus(ctx, recruiterID, web.UpdateApplicationStatusRequest{ApplicationId: appID, Status: "REVIEWING"})
		assert.Equal(t, errs.ErrApplicationAlreadyWithdrawn, err)

		// 5. Invalid status
		appRepo.On("FindById", ctx, appID).Return(activeApp, nil)
		_, err = svc.UpdateStatus(ctx, recruiterID, web.UpdateApplicationStatusRequest{ApplicationId: appID, Status: "INVALID"})
		assert.Equal(t, errs.ErrInvalidApplicationStatus, err)

		// 6. Repo UpdateStatus error
		appRepo.On("UpdateStatus", ctx, appID, domain.ApplicationStatusShortlisted, (*string)(nil), (*string)(nil)).
			Return(nil, errors.New("db err")).Once()
		_, err = svc.UpdateStatus(ctx, recruiterID, web.UpdateApplicationStatusRequest{ApplicationId: appID, Status: "SHORTLISTED"})
		assert.Equal(t, errs.ErrInternalServer, err)

		// 7. Success
		appRepo.On("UpdateStatus", ctx, appID, domain.ApplicationStatusShortlisted, (*string)(nil), (*string)(nil)).
			Return(&domain.Application{ID: appID, Status: domain.ApplicationStatusShortlisted}, nil).Once()
		res, err := svc.UpdateStatus(ctx, recruiterID, web.UpdateApplicationStatusRequest{ApplicationId: appID, Status: "SHORTLISTED"})
		require.NoError(t, err)
		assert.Equal(t, "SHORTLISTED", res.Status)
	})

	t.Run("WithdrawApplication scenarios", func(t *testing.T) {
		appRepo := new(mockApplicationRepository)
		svc := NewApplicationService(appRepo, new(mockJobRepository), new(mockCompanyRepository))
		req := web.WithdrawApplicationRequest{ApplicationId: appID}

		// 1. Not found & generic error
		appRepo.On("FindById", ctx, appID).Return(nil, gorm.ErrRecordNotFound).Once()
		_, err := svc.WithdrawApplication(ctx, candidateID, req)
		assert.Equal(t, errs.ErrApplicationNotFound, err)

		appRepo.On("FindById", ctx, appID).Return(nil, errors.New("db err")).Once()
		_, err = svc.WithdrawApplication(ctx, candidateID, req)
		assert.Equal(t, errs.ErrInternalServer, err)

		// 2. Forbidden
		appRepo.On("FindById", ctx, appID).Return(&domain.Application{ID: appID, CandidateId: uuid.New()}, nil).Once()
		_, err = svc.WithdrawApplication(ctx, candidateID, req)
		assert.Equal(t, errs.ErrApplicationForbidden, err)

		// 3. Already withdrawn
		appRepo.On("FindById", ctx, appID).Return(&domain.Application{
			ID:          appID,
			CandidateId: candidateID,
			Status:      domain.ApplicationStatusWithdrawn,
		}, nil).Once()
		_, err = svc.WithdrawApplication(ctx, candidateID, req)
		assert.Equal(t, errs.ErrApplicationAlreadyWithdrawn, err)

		// 4. Already hired -> ErrInvalidAction
		appRepo.On("FindById", ctx, appID).Return(&domain.Application{
			ID:          appID,
			CandidateId: candidateID,
			Status:      domain.ApplicationStatusHired,
		}, nil).Once()
		_, err = svc.WithdrawApplication(ctx, candidateID, req)
		assert.Equal(t, errs.ErrInvalidAction, err)

		// 5. Withdraw repo error
		appRepo.On("FindById", ctx, appID).Return(&domain.Application{
			ID:          appID,
			CandidateId: candidateID,
			Status:      domain.ApplicationStatusApplied,
		}, nil)
		appRepo.On("Withdraw", ctx, appID, (*string)(nil)).Return(nil, errors.New("db err")).Once()
		_, err = svc.WithdrawApplication(ctx, candidateID, req)
		assert.Equal(t, errs.ErrInternalServer, err)

		// 6. Success
		appRepo.On("Withdraw", ctx, appID, (*string)(nil)).Return(&domain.Application{
			ID:          appID,
			CandidateId: candidateID,
			Status:      domain.ApplicationStatusWithdrawn,
		}, nil).Once()
		res, err := svc.WithdrawApplication(ctx, candidateID, req)
		require.NoError(t, err)
		assert.Equal(t, "WITHDRAWN", res.Status)
	})

	t.Run("Nil mapper checks", func(t *testing.T) {
		assert.Nil(t, toCandidateApplicationResponse(nil))
		assert.Nil(t, toRecruiterApplicationResponse(nil))
	})
}
