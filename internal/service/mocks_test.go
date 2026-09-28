package service

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"AlfianChabib/go-job-board-api/internal/model/web"
	"context"
	"io"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockAuthRepository
type mockAuthRepository struct {
	mock.Mock
}

func (m *mockAuthRepository) Register(ctx context.Context, user domain.User) (*domain.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockAuthRepository) IsUserExist(ctx context.Context, email string) bool {
	args := m.Called(ctx, email)
	return args.Bool(0)
}

func (m *mockAuthRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

// MockTokenRepository
type mockTokenRepository struct {
	mock.Mock
}

func (m *mockTokenRepository) Save(ctx context.Context, token domain.Token) (*domain.Token, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Token), args.Error(1)
}

func (m *mockTokenRepository) RevokeToken(ctx context.Context, refreshToken string) error {
	args := m.Called(ctx, refreshToken)
	return args.Error(0)
}

func (m *mockTokenRepository) FindTokenWithUser(ctx context.Context, userId uuid.UUID, refreshToken string) (*domain.Token, error) {
	args := m.Called(ctx, userId, refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Token), args.Error(1)
}

// MockPasswordHasher
type mockPasswordHasher struct {
	mock.Mock
}

func (m *mockPasswordHasher) Hash(password []byte) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *mockPasswordHasher) Compare(hashedPassword, password []byte) bool {
	args := m.Called(hashedPassword, password)
	return args.Bool(0)
}

// MockJwtManager
type mockJwtManager struct {
	mock.Mock
}

func (m *mockJwtManager) GenerateTokenPair(userID uuid.UUID, role domain.UserRole) (*domain.TokenPair, error) {
	args := m.Called(userID, role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TokenPair), args.Error(1)
}

func (m *mockJwtManager) ValidateAccessToken(tokenString string) (*domain.JwtCustomClaims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.JwtCustomClaims), args.Error(1)
}

func (m *mockJwtManager) ValidateRefreshToken(tokenString string) (*domain.JwtCustomClaims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.JwtCustomClaims), args.Error(1)
}

// MockCandidateRepository
type mockCandidateRepository struct {
	mock.Mock
}

func (m *mockCandidateRepository) Get(ctx context.Context, userId uuid.UUID) (*domain.Profile, error) {
	args := m.Called(ctx, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Profile), args.Error(1)
}

func (m *mockCandidateRepository) GetProfileIdByUserId(ctx context.Context, userId uuid.UUID) (uuid.UUID, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *mockCandidateRepository) Update(ctx context.Context, candidate domain.Profile) (*domain.Profile, error) {
	args := m.Called(ctx, candidate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Profile), args.Error(1)
}

func (m *mockCandidateRepository) UploadAvatar(ctx context.Context, userId uuid.UUID, avatarUrl string) error {
	args := m.Called(ctx, userId, avatarUrl)
	return args.Error(0)
}

func (m *mockCandidateRepository) DeleteAvatar(ctx context.Context, userId uuid.UUID) error {
	args := m.Called(ctx, userId)
	return args.Error(0)
}

func (m *mockCandidateRepository) UpdateSkills(ctx context.Context, userId uuid.UUID, skills web.UpdateCandidateSkillsRequest) ([]domain.Skill, error) {
	args := m.Called(ctx, userId, skills)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Skill), args.Error(1)
}

func (m *mockCandidateRepository) GetExperiences(ctx context.Context, userId uuid.UUID) ([]domain.Experience, error) {
	args := m.Called(ctx, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Experience), args.Error(1)
}

func (m *mockCandidateRepository) CreateExperience(ctx context.Context, userId uuid.UUID, experience domain.Experience) error {
	args := m.Called(ctx, userId, experience)
	return args.Error(0)
}

func (m *mockCandidateRepository) UpdateExperience(ctx context.Context, experienceId uuid.UUID, experience domain.Experience) error {
	args := m.Called(ctx, experienceId, experience)
	return args.Error(0)
}

func (m *mockCandidateRepository) DeleteExperience(ctx context.Context, profileId uuid.UUID, experienceId uuid.UUID) error {
	args := m.Called(ctx, profileId, experienceId)
	return args.Error(0)
}

// MockStorageRepository
type mockStorageRepository struct {
	mock.Mock
}

func (m *mockStorageRepository) UploadAvatar(ctx context.Context, fileName string, file io.Reader, size int64, contentType string) (*string, error) {
	args := m.Called(ctx, fileName, file, size, contentType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*string), args.Error(1)
}

func (m *mockStorageRepository) DeleteAvatar(ctx context.Context, fileName string) error {
	args := m.Called(ctx, fileName)
	return args.Error(0)
}

func (m *mockStorageRepository) UploadLogo(ctx context.Context, fileName string, file io.Reader, size int64, contentType string) (*string, error) {
	args := m.Called(ctx, fileName, file, size, contentType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*string), args.Error(1)
}

func (m *mockStorageRepository) UploadBanner(ctx context.Context, fileName string, file io.Reader, size int64, contentType string) (*string, error) {
	args := m.Called(ctx, fileName, file, size, contentType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*string), args.Error(1)
}

// MockCompanyRepository
type mockCompanyRepository struct {
	mock.Mock
}

func (m *mockCompanyRepository) Create(ctx context.Context, company domain.Company) (*domain.Company, error) {
	args := m.Called(ctx, company)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Company), args.Error(1)
}

func (m *mockCompanyRepository) FindByRecruiterId(ctx context.Context, recruiterId uuid.UUID) (*domain.Company, error) {
	args := m.Called(ctx, recruiterId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Company), args.Error(1)
}

func (m *mockCompanyRepository) GetCompanyIdByRecruiterId(ctx context.Context, recruiterId uuid.UUID) (uuid.UUID, error) {
	args := m.Called(ctx, recruiterId)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *mockCompanyRepository) FindById(ctx context.Context, id uuid.UUID) (*domain.Company, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Company), args.Error(1)
}

func (m *mockCompanyRepository) FindAll(ctx context.Context, req web.GetCompaniesRequest) ([]domain.Company, int64, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]domain.Company), args.Get(1).(int64), args.Error(2)
}

func (m *mockCompanyRepository) Update(ctx context.Context, company domain.Company) (*domain.Company, error) {
	args := m.Called(ctx, company)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Company), args.Error(1)
}

func (m *mockCompanyRepository) UpdateLogo(ctx context.Context, recruiterId uuid.UUID, logoUrl string) error {
	args := m.Called(ctx, recruiterId, logoUrl)
	return args.Error(0)
}

func (m *mockCompanyRepository) UpdateBanner(ctx context.Context, recruiterId uuid.UUID, bannerUrl string) error {
	args := m.Called(ctx, recruiterId, bannerUrl)
	return args.Error(0)
}

// MockSkillRepository
type mockSkillRepository struct {
	mock.Mock
}

func (m *mockSkillRepository) FindAll(ctx context.Context, req web.GetDataRequest) ([]domain.Skill, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Skill), args.Error(1)
}

// MockJobRepository
type mockJobRepository struct {
	mock.Mock
}

func (m *mockJobRepository) Create(ctx context.Context, job domain.Job) (*domain.Job, error) {
	args := m.Called(ctx, job)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Job), args.Error(1)
}

func (m *mockJobRepository) FindAll(ctx context.Context, req web.GetJobsRequest) ([]domain.Job, int64, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]domain.Job), args.Get(1).(int64), args.Error(2)
}

func (m *mockJobRepository) FindAllByCompanyId(ctx context.Context, companyId uuid.UUID, req web.GetRecruiterJobsRequest) ([]domain.Job, int64, error) {
	args := m.Called(ctx, companyId, req)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]domain.Job), args.Get(1).(int64), args.Error(2)
}

func (m *mockJobRepository) FindById(ctx context.Context, id uuid.UUID) (*domain.Job, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Job), args.Error(1)
}

func (m *mockJobRepository) Update(ctx context.Context, job domain.Job) (*domain.Job, error) {
	args := m.Called(ctx, job)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Job), args.Error(1)
}

func (m *mockJobRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.JobStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *mockJobRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockApplicationRepository
type mockApplicationRepository struct {
	mock.Mock
}

func (m *mockApplicationRepository) Create(ctx context.Context, application domain.Application) (*domain.Application, error) {
	args := m.Called(ctx, application)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Application), args.Error(1)
}

func (m *mockApplicationRepository) FindById(ctx context.Context, id uuid.UUID) (*domain.Application, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Application), args.Error(1)
}

func (m *mockApplicationRepository) FindByJobAndCandidate(ctx context.Context, jobId uuid.UUID, candidateId uuid.UUID) (*domain.Application, error) {
	args := m.Called(ctx, jobId, candidateId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Application), args.Error(1)
}

func (m *mockApplicationRepository) FindAllByJobId(ctx context.Context, jobId uuid.UUID, req web.GetJobApplicationsRequest) ([]domain.Application, int64, error) {
	args := m.Called(ctx, jobId, req)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]domain.Application), args.Get(1).(int64), args.Error(2)
}

func (m *mockApplicationRepository) FindAllByCandidateId(ctx context.Context, candidateId uuid.UUID, req web.GetCandidateApplicationsRequest) ([]domain.Application, int64, error) {
	args := m.Called(ctx, candidateId, req)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]domain.Application), args.Get(1).(int64), args.Error(2)
}

func (m *mockApplicationRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.ApplicationStatus, recruiterNotes *string, rejectionReason *string) (*domain.Application, error) {
	args := m.Called(ctx, id, status, recruiterNotes, rejectionReason)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Application), args.Error(1)
}

func (m *mockApplicationRepository) Withdraw(ctx context.Context, id uuid.UUID, withdrawnReason *string) (*domain.Application, error) {
	args := m.Called(ctx, id, withdrawnReason)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Application), args.Error(1)
}
