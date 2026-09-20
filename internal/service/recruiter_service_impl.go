package service

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"AlfianChabib/go-job-board-api/internal/model/web"
	"AlfianChabib/go-job-board-api/internal/repository"
	"AlfianChabib/go-job-board-api/pkg/errs"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type recruiterService struct {
	repository  repository.RecruiterRepository
	storageRepo repository.StorageRepository
}

func NewRecruiterService(repo repository.RecruiterRepository, storageRepo repository.StorageRepository) RecruiterService {
	return &recruiterService{
		repository:  repo,
		storageRepo: storageRepo,
	}
}

// NewRecruiteService provides backwards compatibility for the initial constructor name
func NewRecruiteService(repo repository.RecruiterRepository, storageRepo ...repository.StorageRepository) RecruiterService {
	var storage repository.StorageRepository
	if len(storageRepo) > 0 {
		storage = storageRepo[0]
	}
	return &recruiterService{
		repository:  repo,
		storageRepo: storage,
	}
}

func (service *recruiterService) CreateCompany(ctx context.Context, recruiterId uuid.UUID, req web.CreateCompanyRequest) (*web.CompanyResponse, error) {
	// Check if recruiter already has a company profile
	existingCompany, err := service.repository.FindByRecruiterId(ctx, recruiterId)
	if err == nil && existingCompany != nil {
		return nil, errs.ErrCompanyAlreadyExists
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	company := domain.Company{
		RecruiterId:  recruiterId,
		Name:         req.Name,
		Location:     req.Location,
		Industry:     req.Industry,
		EmployeeSize: req.EmployeeSize,
		Description:  req.Description,
		Website:      req.Website,
	}

	created, err := service.repository.Create(ctx, company)
	if err != nil {
		return nil, err
	}

	return toCompanyResponse(created), nil
}

func (service *recruiterService) GetCompany(ctx context.Context, recruiterId uuid.UUID) (*web.CompanyResponse, error) {
	company, err := service.repository.FindByRecruiterId(ctx, recruiterId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrCompanyNotFound
		}
		return nil, err
	}

	return toCompanyResponse(company), nil
}

func (service *recruiterService) UpdateCompany(ctx context.Context, recruiterId uuid.UUID, req web.UpdateCompanyRequest) (*web.CompanyResponse, error) {
	company := domain.Company{
		RecruiterId:  recruiterId,
		Name:         req.Name,
		Location:     req.Location,
		Industry:     req.Industry,
		EmployeeSize: req.EmployeeSize,
		Description:  req.Description,
		Website:      req.Website,
	}

	updated, err := service.repository.Update(ctx, company)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrCompanyNotFound
		}
		return nil, err
	}

	return toCompanyResponse(updated), nil
}

func (service *recruiterService) UploadLogo(ctx context.Context, req web.UpdateCompanyLogoRequest) (*web.UploadCompanyLogoResponse, error) {
	// Verify company exists
	_, err := service.repository.FindByRecruiterId(ctx, req.RecruiterId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrCompanyNotFound
		}
		return nil, err
	}

	if service.storageRepo == nil {
		return nil, errs.ErrStorageUnavailable
	}

	fileName := fmt.Sprintf("logo_%s_%d%s", req.RecruiterId, time.Now().Unix(), req.Extension)
	logoUrl, err := service.storageRepo.UploadLogo(ctx, fileName, req.File, req.FileSize, req.ContentType)
	if err != nil {
		return nil, errs.ErrUploadLogoFailed
	}

	err = service.repository.UpdateLogo(ctx, req.RecruiterId, *logoUrl)
	if err != nil {
		return nil, err
	}

	return &web.UploadCompanyLogoResponse{
		LogoUrl: *logoUrl,
	}, nil
}

func (service *recruiterService) UploadBanner(ctx context.Context, req web.UpdateCompanyBannerRequest) (*web.UploadCompanyBannerResponse, error) {
	// Verify company exists
	_, err := service.repository.FindByRecruiterId(ctx, req.RecruiterId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrCompanyNotFound
		}
		return nil, err
	}

	if service.storageRepo == nil {
		return nil, errs.ErrStorageUnavailable
	}

	fileName := fmt.Sprintf("banner_%s_%d%s", req.RecruiterId, time.Now().Unix(), req.Extension)
	bannerUrl, err := service.storageRepo.UploadBanner(ctx, fileName, req.File, req.FileSize, req.ContentType)
	if err != nil {
		return nil, errs.ErrUploadBannerFailed
	}

	err = service.repository.UpdateBanner(ctx, req.RecruiterId, *bannerUrl)
	if err != nil {
		return nil, err
	}

	return &web.UploadCompanyBannerResponse{
		BannerUrl: *bannerUrl,
	}, nil
}

func toCompanyResponse(company *domain.Company) *web.CompanyResponse {
	if company == nil {
		return nil
	}
	return &web.CompanyResponse{
		ID:           company.ID,
		RecruiterId:  company.RecruiterId,
		Name:         company.Name,
		LogoUrl:      company.LogoUrl,
		BannerUrl:    company.BannerUrl,
		Website:      company.Website,
		Industry:     company.Industry,
		EmployeeSize: company.EmployeeSize,
		Description:  company.Description,
		Location:     company.Location,
		CreatedAt:    company.CreatedAt,
		UpdatedAt:    company.UpdatedAt,
	}
}
