package web_test

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"AlfianChabib/go-job-board-api/internal/model/web"
	"AlfianChabib/go-job-board-api/pkg/validator"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestWebRequestValidation(t *testing.T) {
	v := validator.NewValidator()

	t.Run("Auth Requests", func(t *testing.T) {
		validReg := web.RegisterRequest{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "password123",
			Role:     domain.RoleCandidate,
		}
		assert.NoError(t, v.Validate(validReg))

		invalidReg := web.RegisterRequest{
			Name:     "Jo",
			Email:    "bad-email",
			Password: "short",
			Role:     "ADMIN",
		}
		assert.Error(t, v.Validate(invalidReg))

		validLogin := web.LoginRequest{
			Email:    "john@example.com",
			Password: "password123",
		}
		assert.NoError(t, v.Validate(validLogin))
		assert.Error(t, v.Validate(web.LoginRequest{}))

		assert.NoError(t, v.Validate(web.LogOutRequest{RefreshToken: "tok"}))
		assert.Error(t, v.Validate(web.LogOutRequest{}))

		assert.NoError(t, v.Validate(web.RefreshTokenRequest{RefreshToken: "rt", AccessToken: "at"}))
		assert.Error(t, v.Validate(web.RefreshTokenRequest{}))
	})

	t.Run("Candidate Requests", func(t *testing.T) {
		assert.NoError(t, v.Validate(web.UpdateCandidateRequest{
			Headline: "Backend Engineer",
			Phone:    "+628123456789",
		}))
		assert.Error(t, v.Validate(web.UpdateCandidateRequest{
			Headline: "",
			Phone:    "not-e164",
		}))

		assert.NoError(t, v.Validate(web.UpdateCandidateSkillsRequest{
			Skills: []web.SkillRequest{{Name: "Go"}},
		}))
		assert.Error(t, v.Validate(web.UpdateCandidateSkillsRequest{}))

		endDate := "2026-01-01T00:00:00Z"
		assert.NoError(t, v.Validate(web.CreateExperienceRequest{
			CompanyName: "Acme",
			Position:    "Engineer",
			StartDate:   "2025-01-01T00:00:00Z",
			EndDate:     &endDate,
			IsCurrent:   false,
		}))
		assert.NoError(t, v.Validate(web.UpdateExperienceRequest{
			ExperienceId: uuid.New(),
			CompanyName:  "Acme",
			Position:     "Engineer",
			StartDate:    "2025-01-01T00:00:00Z",
			EndDate:      &endDate,
			IsCurrent:    false,
		}))
		assert.NoError(t, v.Validate(web.DeleteExperienceRequest{
			ExperienceId: uuid.New(),
		}))
	})

	t.Run("Recruiter & Company Requests", func(t *testing.T) {
		website := "https://acme.com"
		assert.NoError(t, v.Validate(web.CreateCompanyRequest{
			Name:         "Acme Corp",
			Location:     "Jakarta",
			Industry:     "Technology",
			EmployeeSize: domain.EmployeeSize1To50,
			Website:      &website,
		}))
		assert.Error(t, v.Validate(web.CreateCompanyRequest{
			EmployeeSize: "INVALID",
		}))

		assert.NoError(t, v.Validate(web.UpdateCompanyRequest{
			Name:         "Acme Corp",
			Location:     "Jakarta",
			Industry:     "Technology",
			EmployeeSize: domain.EmployeeSize51To200,
			Website:      &website,
		}))

		assert.NoError(t, v.Validate(web.GetCompaniesRequest{Page: 1, Limit: 10}))
		assert.Error(t, v.Validate(web.GetCompaniesRequest{Limit: 200}))
		assert.NoError(t, v.Validate(web.GetCompanyByIdRequest{ID: uuid.New()}))
	})

	t.Run("Job Requests", func(t *testing.T) {
		minSal := int64(5000000)
		maxSal := int64(10000000)
		assert.NoError(t, v.Validate(web.CreateJobRequest{
			Title:          "Golang Developer",
			Description:    "Build scalable APIs in Go",
			EmploymentType: "FULL_TIME",
			WorkMode:       "REMOTE",
			Location:       "Jakarta",
			MinSalary:      &minSal,
			MaxSalary:      &maxSal,
			Currency:       "IDR",
		}))
		assert.Error(t, v.Validate(web.CreateJobRequest{
			Title:          "Go",
			Description:    "short",
			EmploymentType: "INVALID",
			WorkMode:       "INVALID",
		}))

		assert.NoError(t, v.Validate(web.UpdateJobStatusRequest{
			ID:     uuid.New(),
			Status: "OPEN",
		}))
		assert.Error(t, v.Validate(web.UpdateJobStatusRequest{
			ID:     uuid.New(),
			Status: "PENDING",
		}))
	})

	t.Run("Application Requests", func(t *testing.T) {
		assert.NoError(t, v.Validate(web.ApplyJobRequest{
			JobId:     uuid.New(),
			ResumeUrl: "https://storage.example.com/cv.pdf",
		}))
		assert.Error(t, v.Validate(web.ApplyJobRequest{
			ResumeUrl: "not-a-url",
		}))

		assert.NoError(t, v.Validate(web.UpdateApplicationStatusRequest{
			ApplicationId: uuid.New(),
			Status:        "SHORTLISTED",
		}))
		assert.Error(t, v.Validate(web.UpdateApplicationStatusRequest{
			Status: "APPLIED",
		}))

		assert.NoError(t, v.Validate(web.WithdrawApplicationRequest{
			ApplicationId: uuid.New(),
		}))
		assert.NoError(t, v.Validate(web.GetDataRequest{Search: "go", Limit: 20}))
	})
}
