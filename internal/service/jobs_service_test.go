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

func TestJobService_CreateJob(t *testing.T) {
	ctx := context.Background()
	recruiterID := uuid.New()
	companyID := uuid.New()
	company := &domain.Company{ID: companyID, Name: "Acme", Location: "Jakarta"}

	minSal := int64(10000000)
	maxSal := int64(5000000)
	validMaxSal := int64(20000000)

	t.Run("Company not found and generic error", func(t *testing.T) {
		jobRepo := new(mockJobRepository)
		compRepo := new(mockCompanyRepository)
		svc := NewJobService(jobRepo, compRepo)

		compRepo.On("FindByRecruiterId", ctx, recruiterID).Return(nil, gorm.ErrRecordNotFound).Once()
		res, err := svc.CreateJob(ctx, recruiterID, web.CreateJobRequest{})
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrCompanyNotFound, err)

		compRepo.On("FindByRecruiterId", ctx, recruiterID).Return(nil, errors.New("db err")).Once()
		res, err = svc.CreateJob(ctx, recruiterID, web.CreateJobRequest{})
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrInternalServer, err)
	})

	t.Run("Invalid salary range (MaxSalary < MinSalary)", func(t *testing.T) {
		jobRepo := new(mockJobRepository)
		compRepo := new(mockCompanyRepository)
		svc := NewJobService(jobRepo, compRepo)

		compRepo.On("FindByRecruiterId", ctx, recruiterID).Return(company, nil)
		res, err := svc.CreateJob(ctx, recruiterID, web.CreateJobRequest{
			MinSalary: &minSal,
			MaxSalary: &maxSal,
		})
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrInvalidSalaryRange, err)
	})

	t.Run("Create repo error and success with default currency IDR", func(t *testing.T) {
		jobRepo := new(mockJobRepository)
		compRepo := new(mockCompanyRepository)
		svc := NewJobService(jobRepo, compRepo)

		req := web.CreateJobRequest{
			Title:          "Backend Engineer",
			Description:    "Build Go APIs",
			EmploymentType: "FULL_TIME",
			WorkMode:       "REMOTE",
			Location:       "Jakarta",
			MinSalary:      &minSal,
			MaxSalary:      &validMaxSal,
		}

		compRepo.On("FindByRecruiterId", ctx, recruiterID).Return(company, nil)
		jobRepo.On("Create", ctx, mock.AnythingOfType("domain.Job")).Return(nil, errors.New("db err")).Once()
		res, err := svc.CreateJob(ctx, recruiterID, req)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrInternalServer, err)

		jobID := uuid.New()
		jobRepo.On("Create", ctx, mock.AnythingOfType("domain.Job")).Return(&domain.Job{
			ID:        jobID,
			CompanyId: companyID,
			Title:     req.Title,
			Currency:  "IDR",
			Status:    domain.JobStatusOpen,
		}, nil).Once()

		res, err = svc.CreateJob(ctx, recruiterID, req)
		require.NoError(t, err)
		assert.Equal(t, jobID, res.ID)
		assert.Equal(t, "IDR", res.Currency)
		require.NotNil(t, res.Company)
		assert.Equal(t, companyID, res.Company.ID)
	})
}

func TestJobService_GetJobsAndRecruiterJobs(t *testing.T) {
	ctx := context.Background()
	recruiterID := uuid.New()
	companyID := uuid.New()

	t.Run("GetJobs error and success", func(t *testing.T) {
		jobRepo := new(mockJobRepository)
		compRepo := new(mockCompanyRepository)
		svc := NewJobService(jobRepo, compRepo)
		req := web.GetJobsRequest{}

		jobRepo.On("FindAll", ctx, req).Return(nil, int64(0), errors.New("db err")).Once()
		res, err := svc.GetJobs(ctx, req)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrInternalServer, err)

		jobRepo.On("FindAll", ctx, req).Return([]domain.Job{{ID: uuid.New(), Title: "Go Dev"}}, int64(1), nil).Once()
		res, err = svc.GetJobs(ctx, req)
		require.NoError(t, err)
		assert.Len(t, res.Jobs, 1)
		assert.Equal(t, 1, res.Page)
		assert.Equal(t, 10, res.Limit)
		assert.Equal(t, 1, res.TotalPages)
	})

	t.Run("GetRecruiterJobs company not found, company error, repo error, and success", func(t *testing.T) {
		jobRepo := new(mockJobRepository)
		compRepo := new(mockCompanyRepository)
		svc := NewJobService(jobRepo, compRepo)
		req := web.GetRecruiterJobsRequest{}

		compRepo.On("FindByRecruiterId", ctx, recruiterID).Return(nil, gorm.ErrRecordNotFound).Once()
		res, err := svc.GetRecruiterJobs(ctx, recruiterID, req)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrCompanyNotFound, err)

		compRepo.On("FindByRecruiterId", ctx, recruiterID).Return(nil, errors.New("db err")).Once()
		res, err = svc.GetRecruiterJobs(ctx, recruiterID, req)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrInternalServer, err)

		compRepo.On("FindByRecruiterId", ctx, recruiterID).Return(&domain.Company{ID: companyID}, nil)
		jobRepo.On("FindAllByCompanyId", ctx, companyID, req).Return(nil, int64(0), errors.New("db err")).Once()
		res, err = svc.GetRecruiterJobs(ctx, recruiterID, req)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrInternalServer, err)

		jobRepo.On("FindAllByCompanyId", ctx, companyID, req).Return([]domain.Job{{ID: uuid.New(), Title: "Go Dev"}}, int64(1), nil).Once()
		res, err = svc.GetRecruiterJobs(ctx, recruiterID, req)
		require.NoError(t, err)
		assert.Len(t, res.Jobs, 1)
	})

	t.Run("GetJobById not found, error, and success", func(t *testing.T) {
		jobRepo := new(mockJobRepository)
		compRepo := new(mockCompanyRepository)
		svc := NewJobService(jobRepo, compRepo)
		jobID := uuid.New()

		jobRepo.On("FindById", ctx, jobID).Return(nil, gorm.ErrRecordNotFound).Once()
		res, err := svc.GetJobById(ctx, jobID)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrJobNotFound, err)

		jobRepo.On("FindById", ctx, jobID).Return(nil, errors.New("db err")).Once()
		res, err = svc.GetJobById(ctx, jobID)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrInternalServer, err)

		jobRepo.On("FindById", ctx, jobID).Return(&domain.Job{ID: jobID, Title: "Go Dev"}, nil).Once()
		res, err = svc.GetJobById(ctx, jobID)
		require.NoError(t, err)
		assert.Equal(t, jobID, res.ID)
	})
}

func TestJobService_UpdateStatusAndDelete(t *testing.T) {
	ctx := context.Background()
	recruiterID := uuid.New()
	companyID := uuid.New()
	otherCompanyID := uuid.New()
	jobID := uuid.New()
	company := &domain.Company{ID: companyID, Name: "Acme"}

	t.Run("UpdateJob scenarios", func(t *testing.T) {
		jobRepo := new(mockJobRepository)
		compRepo := new(mockCompanyRepository)
		svc := NewJobService(jobRepo, compRepo)
		req := web.UpdateJobRequest{Title: "Updated Go Dev", EmploymentType: "FULL_TIME", WorkMode: "REMOTE"}

		// 1. Company not found
		compRepo.On("FindByRecruiterId", ctx, recruiterID).Return(nil, gorm.ErrRecordNotFound).Once()
		_, err := svc.UpdateJob(ctx, recruiterID, jobID, req)
		assert.Equal(t, errs.ErrCompanyNotFound, err)

		// 2. Company generic error
		compRepo.On("FindByRecruiterId", ctx, recruiterID).Return(nil, errors.New("db err")).Once()
		_, err = svc.UpdateJob(ctx, recruiterID, jobID, req)
		assert.Equal(t, errs.ErrInternalServer, err)

		// 3. Job not found
		compRepo.On("FindByRecruiterId", ctx, recruiterID).Return(company, nil)
		jobRepo.On("FindById", ctx, jobID).Return(nil, gorm.ErrRecordNotFound).Once()
		_, err = svc.UpdateJob(ctx, recruiterID, jobID, req)
		assert.Equal(t, errs.ErrJobNotFound, err)

		// 4. Job generic error
		jobRepo.On("FindById", ctx, jobID).Return(nil, errors.New("db err")).Once()
		_, err = svc.UpdateJob(ctx, recruiterID, jobID, req)
		assert.Equal(t, errs.ErrInternalServer, err)

		// 5. IDOR forbidden
		jobRepo.On("FindById", ctx, jobID).Return(&domain.Job{ID: jobID, CompanyId: otherCompanyID}, nil).Once()
		_, err = svc.UpdateJob(ctx, recruiterID, jobID, req)
		assert.Equal(t, errs.ErrJobForbidden, err)

		// 6. Invalid salary range
		minS := int64(20)
		maxS := int64(10)
		jobRepo.On("FindById", ctx, jobID).Return(&domain.Job{ID: jobID, CompanyId: companyID}, nil).Once()
		_, err = svc.UpdateJob(ctx, recruiterID, jobID, web.UpdateJobRequest{MinSalary: &minS, MaxSalary: &maxS})
		assert.Equal(t, errs.ErrInvalidSalaryRange, err)

		// 7. Update repo error
		jobRepo.On("FindById", ctx, jobID).Return(&domain.Job{ID: jobID, CompanyId: companyID}, nil).Once()
		jobRepo.On("Update", ctx, mock.AnythingOfType("domain.Job")).Return(nil, errors.New("db err")).Once()
		_, err = svc.UpdateJob(ctx, recruiterID, jobID, req)
		assert.Equal(t, errs.ErrInternalServer, err)

		// 8. Success
		jobRepo.On("FindById", ctx, jobID).Return(&domain.Job{ID: jobID, CompanyId: companyID}, nil).Once()
		jobRepo.On("Update", ctx, mock.AnythingOfType("domain.Job")).Return(&domain.Job{ID: jobID, CompanyId: companyID, Title: req.Title}, nil).Once()
		res, err := svc.UpdateJob(ctx, recruiterID, jobID, req)
		require.NoError(t, err)
		assert.Equal(t, req.Title, res.Title)
	})

	t.Run("UpdateJobStatus scenarios", func(t *testing.T) {
		jobRepo := new(mockJobRepository)
		compRepo := new(mockCompanyRepository)
		svc := NewJobService(jobRepo, compRepo)

		// 1. Invalid status
		_, err := svc.UpdateJobStatus(ctx, recruiterID, jobID, "INVALID")
		assert.Equal(t, errs.ErrInvalidJobStatus, err)

		// 2. Company not found
		compRepo.On("FindByRecruiterId", ctx, recruiterID).Return(nil, gorm.ErrRecordNotFound).Once()
		_, err = svc.UpdateJobStatus(ctx, recruiterID, jobID, "CLOSED")
		assert.Equal(t, errs.ErrCompanyNotFound, err)

		// 3. Company generic error
		compRepo.On("FindByRecruiterId", ctx, recruiterID).Return(nil, errors.New("db err")).Once()
		_, err = svc.UpdateJobStatus(ctx, recruiterID, jobID, "CLOSED")
		assert.Equal(t, errs.ErrInternalServer, err)

		// 4. Job not found
		compRepo.On("FindByRecruiterId", ctx, recruiterID).Return(company, nil)
		jobRepo.On("FindById", ctx, jobID).Return(nil, gorm.ErrRecordNotFound).Once()
		_, err = svc.UpdateJobStatus(ctx, recruiterID, jobID, "CLOSED")
		assert.Equal(t, errs.ErrJobNotFound, err)

		// 5. Job generic error
		jobRepo.On("FindById", ctx, jobID).Return(nil, errors.New("db err")).Once()
		_, err = svc.UpdateJobStatus(ctx, recruiterID, jobID, "CLOSED")
		assert.Equal(t, errs.ErrInternalServer, err)

		// 6. IDOR forbidden
		jobRepo.On("FindById", ctx, jobID).Return(&domain.Job{ID: jobID, CompanyId: otherCompanyID}, nil).Once()
		_, err = svc.UpdateJobStatus(ctx, recruiterID, jobID, "CLOSED")
		assert.Equal(t, errs.ErrJobForbidden, err)

		// 7. UpdateStatus not found
		jobRepo.On("FindById", ctx, jobID).Return(&domain.Job{ID: jobID, CompanyId: companyID}, nil).Once()
		jobRepo.On("UpdateStatus", ctx, jobID, domain.JobStatusClosed).Return(gorm.ErrRecordNotFound).Once()
		_, err = svc.UpdateJobStatus(ctx, recruiterID, jobID, "CLOSED")
		assert.Equal(t, errs.ErrJobNotFound, err)

		// 8. UpdateStatus generic error
		jobRepo.On("FindById", ctx, jobID).Return(&domain.Job{ID: jobID, CompanyId: companyID}, nil).Once()
		jobRepo.On("UpdateStatus", ctx, jobID, domain.JobStatusClosed).Return(errors.New("db err")).Once()
		_, err = svc.UpdateJobStatus(ctx, recruiterID, jobID, "CLOSED")
		assert.Equal(t, errs.ErrInternalServer, err)

		// 9. Success
		jobRepo.On("FindById", ctx, jobID).Return(&domain.Job{ID: jobID, CompanyId: companyID}, nil).Once()
		jobRepo.On("UpdateStatus", ctx, jobID, domain.JobStatusClosed).Return(nil).Once()
		res, err := svc.UpdateJobStatus(ctx, recruiterID, jobID, "CLOSED")
		require.NoError(t, err)
		assert.Equal(t, domain.JobStatusClosed, res.Status)
	})

	t.Run("DeleteJob scenarios", func(t *testing.T) {
		jobRepo := new(mockJobRepository)
		compRepo := new(mockCompanyRepository)
		svc := NewJobService(jobRepo, compRepo)

		// 1. Company not found
		compRepo.On("FindByRecruiterId", ctx, recruiterID).Return(nil, gorm.ErrRecordNotFound).Once()
		assert.Equal(t, errs.ErrCompanyNotFound, svc.DeleteJob(ctx, recruiterID, jobID))

		// 2. Company generic error
		compRepo.On("FindByRecruiterId", ctx, recruiterID).Return(nil, errors.New("db err")).Once()
		assert.Equal(t, errs.ErrInternalServer, svc.DeleteJob(ctx, recruiterID, jobID))

		// 3. Job not found
		compRepo.On("FindByRecruiterId", ctx, recruiterID).Return(company, nil)
		jobRepo.On("FindById", ctx, jobID).Return(nil, gorm.ErrRecordNotFound).Once()
		assert.Equal(t, errs.ErrJobNotFound, svc.DeleteJob(ctx, recruiterID, jobID))

		// 4. Job generic error
		jobRepo.On("FindById", ctx, jobID).Return(nil, errors.New("db err")).Once()
		assert.Equal(t, errs.ErrInternalServer, svc.DeleteJob(ctx, recruiterID, jobID))

		// 5. IDOR forbidden
		jobRepo.On("FindById", ctx, jobID).Return(&domain.Job{ID: jobID, CompanyId: otherCompanyID}, nil).Once()
		assert.Equal(t, errs.ErrJobForbidden, svc.DeleteJob(ctx, recruiterID, jobID))

		// 6. Delete not found
		jobRepo.On("FindById", ctx, jobID).Return(&domain.Job{ID: jobID, CompanyId: companyID}, nil).Once()
		jobRepo.On("Delete", ctx, jobID).Return(gorm.ErrRecordNotFound).Once()
		assert.Equal(t, errs.ErrJobNotFound, svc.DeleteJob(ctx, recruiterID, jobID))

		// 7. Delete generic error
		jobRepo.On("FindById", ctx, jobID).Return(&domain.Job{ID: jobID, CompanyId: companyID}, nil).Once()
		jobRepo.On("Delete", ctx, jobID).Return(errors.New("db err")).Once()
		assert.Equal(t, errs.ErrInternalServer, svc.DeleteJob(ctx, recruiterID, jobID))

		// 8. Success
		jobRepo.On("FindById", ctx, jobID).Return(&domain.Job{ID: jobID, CompanyId: companyID}, nil).Once()
		jobRepo.On("Delete", ctx, jobID).Return(nil).Once()
		assert.NoError(t, svc.DeleteJob(ctx, recruiterID, jobID))
	})

	t.Run("toJobResponse nil check", func(t *testing.T) {
		assert.Nil(t, toJobResponse(nil))
	})
}
