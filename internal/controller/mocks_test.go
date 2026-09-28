package controller

import (
	"AlfianChabib/go-job-board-api/internal/middleware"
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"AlfianChabib/go-job-board-api/internal/model/web"
	"AlfianChabib/go-job-board-api/pkg/validator"
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/textproto"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

func setupTestApp() *fiber.App {
	v := validator.NewValidator()
	return fiber.New(fiber.Config{
		StructValidator: v,
		ErrorHandler:    middleware.NewCustomErrorHandler(v),
	})
}

func withSession(userID uuid.UUID, role domain.UserRole) fiber.Handler {
	return func(c fiber.Ctx) error {
		c.Locals("userId", userID)
		c.Locals("role", role)
		return c.Next()
	}
}

func createMultipartPayload(fieldName, fileName, contentType string, content []byte) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldName, fileName))
	h.Set("Content-Type", contentType)

	part, err := writer.CreatePart(h)
	if err != nil {
		return nil, "", err
	}
	if _, err := part.Write(content); err != nil {
		return nil, "", err
	}
	if err := writer.Close(); err != nil {
		return nil, "", err
	}
	return body, writer.FormDataContentType(), nil
}

// MockAuthService
type mockAuthService struct {
	mock.Mock
}

func (m *mockAuthService) Register(ctx context.Context, data web.RegisterRequest) (*web.RegisterResponse, error) {
	args := m.Called(ctx, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*web.RegisterResponse), args.Error(1)
}

func (m *mockAuthService) Login(ctx context.Context, data web.LoginRequest) (*domain.TokenPair, error) {
	args := m.Called(ctx, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TokenPair), args.Error(1)
}

func (m *mockAuthService) Logout(ctx context.Context, refreshToken string) error {
	args := m.Called(ctx, refreshToken)
	return args.Error(0)
}

func (m *mockAuthService) RefreshToken(ctx context.Context, data web.RefreshTokenRequest) (*web.RefreshTokenResponse, error) {
	args := m.Called(ctx, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*web.RefreshTokenResponse), args.Error(1)
}

// MockCandidateService
type mockCandidateService struct {
	mock.Mock
}

func (m *mockCandidateService) Get(ctx context.Context, userId uuid.UUID) (*web.GetCandidateResponse, error) {
	args := m.Called(ctx, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*web.GetCandidateResponse), args.Error(1)
}

func (m *mockCandidateService) Update(ctx context.Context, candidate domain.Profile) (*web.UpdateCandidateResponse, error) {
	args := m.Called(ctx, candidate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*web.UpdateCandidateResponse), args.Error(1)
}

func (m *mockCandidateService) UploadAvatar(ctx context.Context, req web.UpdateCandidateAvatarRequest) (*string, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*string), args.Error(1)
}

func (m *mockCandidateService) DeleteAvatar(ctx context.Context, userId uuid.UUID) error {
	args := m.Called(ctx, userId)
	return args.Error(0)
}

func (m *mockCandidateService) UpdateSkills(ctx context.Context, userId uuid.UUID, skills web.UpdateCandidateSkillsRequest) ([]domain.Skill, error) {
	args := m.Called(ctx, userId, skills)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Skill), args.Error(1)
}

func (m *mockCandidateService) GetExperiences(ctx context.Context, userId uuid.UUID) ([]domain.Experience, error) {
	args := m.Called(ctx, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Experience), args.Error(1)
}

func (m *mockCandidateService) CreateExperience(ctx context.Context, userId uuid.UUID, req web.CreateExperienceRequest) error {
	args := m.Called(ctx, userId, req)
	return args.Error(0)
}

func (m *mockCandidateService) UpdateExperience(ctx context.Context, userId uuid.UUID, req web.UpdateExperienceRequest) error {
	args := m.Called(ctx, userId, req)
	return args.Error(0)
}

func (m *mockCandidateService) DeleteExperience(ctx context.Context, userId uuid.UUID, experienceId uuid.UUID) error {
	args := m.Called(ctx, userId, experienceId)
	return args.Error(0)
}

// MockCompanyService
type mockCompanyService struct {
	mock.Mock
}

func (m *mockCompanyService) GetCompanies(ctx context.Context, req web.GetCompaniesRequest) (*web.CompanyPaginationResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*web.CompanyPaginationResponse), args.Error(1)
}

func (m *mockCompanyService) GetCompanyById(ctx context.Context, companyId uuid.UUID) (*web.CompanyDetailResponse, error) {
	args := m.Called(ctx, companyId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*web.CompanyDetailResponse), args.Error(1)
}

// MockDataService
type mockDataService struct {
	mock.Mock
}

func (m *mockDataService) GetSkills(ctx context.Context, req web.GetDataRequest) ([]web.SkillResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]web.SkillResponse), args.Error(1)
}

func (m *mockDataService) GetCurrencyCodes(ctx context.Context, req web.GetDataRequest) ([]web.CurrencyResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]web.CurrencyResponse), args.Error(1)
}

// MockJobService
type mockJobService struct {
	mock.Mock
}

func (m *mockJobService) CreateJob(ctx context.Context, recruiterUserId uuid.UUID, req web.CreateJobRequest) (*web.JobResponse, error) {
	args := m.Called(ctx, recruiterUserId, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*web.JobResponse), args.Error(1)
}

func (m *mockJobService) GetJobs(ctx context.Context, req web.GetJobsRequest) (*web.JobPaginationResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*web.JobPaginationResponse), args.Error(1)
}

func (m *mockJobService) GetRecruiterJobs(ctx context.Context, recruiterUserId uuid.UUID, req web.GetRecruiterJobsRequest) (*web.JobPaginationResponse, error) {
	args := m.Called(ctx, recruiterUserId, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*web.JobPaginationResponse), args.Error(1)
}

func (m *mockJobService) GetJobById(ctx context.Context, jobId uuid.UUID) (*web.JobResponse, error) {
	args := m.Called(ctx, jobId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*web.JobResponse), args.Error(1)
}

func (m *mockJobService) UpdateJob(ctx context.Context, recruiterUserId uuid.UUID, jobId uuid.UUID, req web.UpdateJobRequest) (*web.JobResponse, error) {
	args := m.Called(ctx, recruiterUserId, jobId, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*web.JobResponse), args.Error(1)
}

func (m *mockJobService) UpdateJobStatus(ctx context.Context, recruiterUserId uuid.UUID, jobId uuid.UUID, status string) (*web.JobResponse, error) {
	args := m.Called(ctx, recruiterUserId, jobId, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*web.JobResponse), args.Error(1)
}

func (m *mockJobService) DeleteJob(ctx context.Context, recruiterUserId uuid.UUID, jobId uuid.UUID) error {
	args := m.Called(ctx, recruiterUserId, jobId)
	return args.Error(0)
}

// MockRecruiterService
type mockRecruiterService struct {
	mock.Mock
}

func (m *mockRecruiterService) CreateCompany(ctx context.Context, recruiterId uuid.UUID, req web.CreateCompanyRequest) (*web.CompanyResponse, error) {
	args := m.Called(ctx, recruiterId, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*web.CompanyResponse), args.Error(1)
}

func (m *mockRecruiterService) GetCompany(ctx context.Context, recruiterId uuid.UUID) (*web.CompanyResponse, error) {
	args := m.Called(ctx, recruiterId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*web.CompanyResponse), args.Error(1)
}

func (m *mockRecruiterService) UpdateCompany(ctx context.Context, recruiterId uuid.UUID, req web.UpdateCompanyRequest) (*web.CompanyResponse, error) {
	args := m.Called(ctx, recruiterId, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*web.CompanyResponse), args.Error(1)
}

func (m *mockRecruiterService) UploadLogo(ctx context.Context, req web.UpdateCompanyLogoRequest) (*web.UploadCompanyLogoResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*web.UploadCompanyLogoResponse), args.Error(1)
}

func (m *mockRecruiterService) UploadBanner(ctx context.Context, req web.UpdateCompanyBannerRequest) (*web.UploadCompanyBannerResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*web.UploadCompanyBannerResponse), args.Error(1)
}

// MockApplicationService
type mockApplicationService struct {
	mock.Mock
}

func (m *mockApplicationService) ApplyJob(ctx context.Context, candidateUserId uuid.UUID, req web.ApplyJobRequest) (*web.CandidateApplicationResponse, error) {
	args := m.Called(ctx, candidateUserId, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*web.CandidateApplicationResponse), args.Error(1)
}

func (m *mockApplicationService) GetJobApplications(ctx context.Context, recruiterUserId uuid.UUID, req web.GetJobApplicationsRequest) (*web.RecruiterApplicationPaginationResponse, error) {
	args := m.Called(ctx, recruiterUserId, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*web.RecruiterApplicationPaginationResponse), args.Error(1)
}

func (m *mockApplicationService) GetCandidateApplications(ctx context.Context, candidateUserId uuid.UUID, req web.GetCandidateApplicationsRequest) (*web.CandidateApplicationPaginationResponse, error) {
	args := m.Called(ctx, candidateUserId, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*web.CandidateApplicationPaginationResponse), args.Error(1)
}

func (m *mockApplicationService) GetApplicationById(ctx context.Context, userId uuid.UUID, userRole string, applicationId uuid.UUID) (any, error) {
	args := m.Called(ctx, userId, userRole, applicationId)
	return args.Get(0), args.Error(1)
}

func (m *mockApplicationService) UpdateStatus(ctx context.Context, recruiterUserId uuid.UUID, req web.UpdateApplicationStatusRequest) (*web.RecruiterApplicationResponse, error) {
	args := m.Called(ctx, recruiterUserId, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*web.RecruiterApplicationResponse), args.Error(1)
}

func (m *mockApplicationService) WithdrawApplication(ctx context.Context, candidateUserId uuid.UUID, req web.WithdrawApplicationRequest) (*web.CandidateApplicationResponse, error) {
	args := m.Called(ctx, candidateUserId, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*web.CandidateApplicationResponse), args.Error(1)
}
