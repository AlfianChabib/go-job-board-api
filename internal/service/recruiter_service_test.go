package service

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"AlfianChabib/go-job-board-api/internal/model/web"
	"AlfianChabib/go-job-board-api/pkg/errs"
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestRecruiterService_CreateCompany(t *testing.T) {
	ctx := context.Background()
	recruiterID := uuid.New()
	req := web.CreateCompanyRequest{
		Name:         "Acme Corp",
		Location:     "Jakarta",
		Industry:     "Tech",
		EmployeeSize: domain.EmployeeSize1To50,
	}

	t.Run("Company already exists returns ErrCompanyAlreadyExists", func(t *testing.T) {
		repo := new(mockCompanyRepository)
		repo.On("FindByRecruiterId", ctx, recruiterID).Return(&domain.Company{ID: uuid.New()}, nil)

		svc := NewRecruiterService(repo, new(mockStorageRepository))
		res, err := svc.CreateCompany(ctx, recruiterID, req)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrCompanyAlreadyExists, err)
	})

	t.Run("FindByRecruiterId unexpected error returns error", func(t *testing.T) {
		repo := new(mockCompanyRepository)
		repo.On("FindByRecruiterId", ctx, recruiterID).Return(nil, errors.New("db error"))

		svc := NewRecruiterService(repo, new(mockStorageRepository))
		res, err := svc.CreateCompany(ctx, recruiterID, req)
		assert.Nil(t, res)
		assert.EqualError(t, err, "db error")
	})

	t.Run("Create error returns error", func(t *testing.T) {
		repo := new(mockCompanyRepository)
		repo.On("FindByRecruiterId", ctx, recruiterID).Return(nil, gorm.ErrRecordNotFound)
		repo.On("Create", ctx, mock.AnythingOfType("domain.Company")).Return(nil, errors.New("insert error"))

		svc := NewRecruiterService(repo, new(mockStorageRepository))
		res, err := svc.CreateCompany(ctx, recruiterID, req)
		assert.Nil(t, res)
		assert.EqualError(t, err, "insert error")
	})

	t.Run("Success creates and returns CompanyResponse", func(t *testing.T) {
		repo := new(mockCompanyRepository)
		companyID := uuid.New()
		repo.On("FindByRecruiterId", ctx, recruiterID).Return(nil, gorm.ErrRecordNotFound)
		repo.On("Create", ctx, mock.AnythingOfType("domain.Company")).Return(&domain.Company{
			ID:           companyID,
			RecruiterId:  recruiterID,
			Name:         req.Name,
			Location:     req.Location,
			Industry:     req.Industry,
			EmployeeSize: req.EmployeeSize,
		}, nil)

		svc := NewRecruiteService(repo, new(mockStorageRepository))
		res, err := svc.CreateCompany(ctx, recruiterID, req)
		require.NoError(t, err)
		assert.Equal(t, companyID, res.ID)
		assert.Equal(t, req.Name, res.Name)
	})
}

func TestRecruiterService_GetAndUpdateCompany(t *testing.T) {
	ctx := context.Background()
	recruiterID := uuid.New()

	t.Run("GetCompany not found, generic error, and success", func(t *testing.T) {
		repo := new(mockCompanyRepository)
		svc := NewRecruiterService(repo, nil)

		repo.On("FindByRecruiterId", ctx, recruiterID).Return(nil, gorm.ErrRecordNotFound).Once()
		res, err := svc.GetCompany(ctx, recruiterID)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrCompanyNotFound, err)

		repo.On("FindByRecruiterId", ctx, recruiterID).Return(nil, errors.New("db err")).Once()
		res, err = svc.GetCompany(ctx, recruiterID)
		assert.Nil(t, res)
		assert.EqualError(t, err, "db err")

		repo.On("FindByRecruiterId", ctx, recruiterID).Return(&domain.Company{ID: uuid.New(), Name: "Acme"}, nil).Once()
		res, err = svc.GetCompany(ctx, recruiterID)
		require.NoError(t, err)
		assert.Equal(t, "Acme", res.Name)
	})

	t.Run("UpdateCompany not found, generic error, and success", func(t *testing.T) {
		repo := new(mockCompanyRepository)
		svc := NewRecruiterService(repo, nil)
		req := web.UpdateCompanyRequest{
			Name:         "Acme Updated",
			Location:     "Bandung",
			Industry:     "IT",
			EmployeeSize: domain.EmployeeSize51To200,
		}

		repo.On("Update", ctx, mock.AnythingOfType("domain.Company")).Return(nil, gorm.ErrRecordNotFound).Once()
		res, err := svc.UpdateCompany(ctx, recruiterID, req)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrCompanyNotFound, err)

		repo.On("Update", ctx, mock.AnythingOfType("domain.Company")).Return(nil, errors.New("db err")).Once()
		res, err = svc.UpdateCompany(ctx, recruiterID, req)
		assert.Nil(t, res)
		assert.EqualError(t, err, "db err")

		repo.On("Update", ctx, mock.AnythingOfType("domain.Company")).Return(&domain.Company{
			ID:   uuid.New(),
			Name: req.Name,
		}, nil).Once()
		res, err = svc.UpdateCompany(ctx, recruiterID, req)
		require.NoError(t, err)
		assert.Equal(t, "Acme Updated", res.Name)
	})
}

func TestRecruiterService_UploadLogoAndBanner(t *testing.T) {
	ctx := context.Background()
	recruiterID := uuid.New()
	logoReq := web.UpdateCompanyLogoRequest{
		RecruiterId: recruiterID,
		File:        bytes.NewReader([]byte("logo")),
		FileSize:    4,
		ContentType: "image/png",
		Extension:   ".png",
	}
	bannerReq := web.UpdateCompanyBannerRequest{
		RecruiterId: recruiterID,
		File:        bytes.NewReader([]byte("banner")),
		FileSize:    6,
		ContentType: "image/jpeg",
		Extension:   ".jpg",
	}

	t.Run("UploadLogo company not found, generic error, storage unavailable, upload fail, repo fail, and success", func(t *testing.T) {
		repo := new(mockCompanyRepository)

		// 1. Not found
		repo.On("FindByRecruiterId", ctx, recruiterID).Return(nil, gorm.ErrRecordNotFound).Once()
		svcNoStorage := NewRecruiteService(repo)
		res, err := svcNoStorage.UploadLogo(ctx, logoReq)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrCompanyNotFound, err)

		// 2. Generic error
		repo.On("FindByRecruiterId", ctx, recruiterID).Return(nil, errors.New("db err")).Once()
		res, err = svcNoStorage.UploadLogo(ctx, logoReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "db err")

		// 3. Storage unavailable
		repo.On("FindByRecruiterId", ctx, recruiterID).Return(&domain.Company{}, nil).Once()
		res, err = svcNoStorage.UploadLogo(ctx, logoReq)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrStorageUnavailable, err)

		// 4. Storage upload error
		storage := new(mockStorageRepository)
		svc := NewRecruiterService(repo, storage)
		repo.On("FindByRecruiterId", ctx, recruiterID).Return(&domain.Company{}, nil).Once()
		storage.On("UploadLogo", ctx, mock.AnythingOfType("string"), logoReq.File, logoReq.FileSize, logoReq.ContentType).
			Return(nil, errors.New("s3 err")).Once()
		res, err = svc.UploadLogo(ctx, logoReq)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrUploadLogoFailed, err)

		// 5. Repo UpdateLogo error
		logoUrl := "http://localhost:9000/avatar/logo.png"
		repo.On("FindByRecruiterId", ctx, recruiterID).Return(&domain.Company{}, nil).Once()
		storage.On("UploadLogo", ctx, mock.AnythingOfType("string"), logoReq.File, logoReq.FileSize, logoReq.ContentType).
			Return(&logoUrl, nil).Once()
		repo.On("UpdateLogo", ctx, recruiterID, logoUrl).Return(errors.New("db update err")).Once()
		res, err = svc.UploadLogo(ctx, logoReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "db update err")

		// 6. Success
		repo.On("FindByRecruiterId", ctx, recruiterID).Return(&domain.Company{}, nil).Once()
		storage.On("UploadLogo", ctx, mock.AnythingOfType("string"), logoReq.File, logoReq.FileSize, logoReq.ContentType).
			Return(&logoUrl, nil).Once()
		repo.On("UpdateLogo", ctx, recruiterID, logoUrl).Return(nil).Once()
		res, err = svc.UploadLogo(ctx, logoReq)
		require.NoError(t, err)
		assert.Equal(t, logoUrl, res.LogoUrl)
	})

	t.Run("UploadBanner company not found, generic error, storage unavailable, upload fail, repo fail, and success", func(t *testing.T) {
		repo := new(mockCompanyRepository)
		svcNoStorage := NewRecruiteService(repo)

		// 1. Not found
		repo.On("FindByRecruiterId", ctx, recruiterID).Return(nil, gorm.ErrRecordNotFound).Once()
		res, err := svcNoStorage.UploadBanner(ctx, bannerReq)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrCompanyNotFound, err)

		// 2. Generic error
		repo.On("FindByRecruiterId", ctx, recruiterID).Return(nil, errors.New("db err")).Once()
		res, err = svcNoStorage.UploadBanner(ctx, bannerReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "db err")

		// 3. Storage unavailable
		repo.On("FindByRecruiterId", ctx, recruiterID).Return(&domain.Company{}, nil).Once()
		res, err = svcNoStorage.UploadBanner(ctx, bannerReq)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrStorageUnavailable, err)

		// 4. Storage upload error
		storage := new(mockStorageRepository)
		svc := NewRecruiterService(repo, storage)
		repo.On("FindByRecruiterId", ctx, recruiterID).Return(&domain.Company{}, nil).Once()
		storage.On("UploadBanner", ctx, mock.AnythingOfType("string"), bannerReq.File, bannerReq.FileSize, bannerReq.ContentType).
			Return(nil, errors.New("s3 err")).Once()
		res, err = svc.UploadBanner(ctx, bannerReq)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrUploadBannerFailed, err)

		// 5. Repo UpdateBanner error
		bannerUrl := "http://localhost:9000/avatar/banner.jpg"
		repo.On("FindByRecruiterId", ctx, recruiterID).Return(&domain.Company{}, nil).Once()
		storage.On("UploadBanner", ctx, mock.AnythingOfType("string"), bannerReq.File, bannerReq.FileSize, bannerReq.ContentType).
			Return(&bannerUrl, nil).Once()
		repo.On("UpdateBanner", ctx, recruiterID, bannerUrl).Return(errors.New("db update err")).Once()
		res, err = svc.UploadBanner(ctx, bannerReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "db update err")

		// 6. Success
		repo.On("FindByRecruiterId", ctx, recruiterID).Return(&domain.Company{}, nil).Once()
		storage.On("UploadBanner", ctx, mock.AnythingOfType("string"), bannerReq.File, bannerReq.FileSize, bannerReq.ContentType).
			Return(&bannerUrl, nil).Once()
		repo.On("UpdateBanner", ctx, recruiterID, bannerUrl).Return(nil).Once()
		res, err = svc.UploadBanner(ctx, bannerReq)
		require.NoError(t, err)
		assert.Equal(t, bannerUrl, res.BannerUrl)
	})

	t.Run("toCompanyResponse nil check", func(t *testing.T) {
		assert.Nil(t, toCompanyResponse(nil))
	})
}
