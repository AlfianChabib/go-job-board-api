package service

import (
	"AlfianChabib/go-job-board-api/internal/model/web"
	"context"

	"github.com/google/uuid"
)

type RecruiterService interface {
	CreateCompany(ctx context.Context, recruiterId uuid.UUID, req web.CreateCompanyRequest) (*web.CompanyResponse, error)
	GetCompany(ctx context.Context, recruiterId uuid.UUID) (*web.CompanyResponse, error)
	UpdateCompany(ctx context.Context, recruiterId uuid.UUID, req web.UpdateCompanyRequest) (*web.CompanyResponse, error)
	UploadLogo(ctx context.Context, req web.UpdateCompanyLogoRequest) (*web.UploadCompanyLogoResponse, error)
	UploadBanner(ctx context.Context, req web.UpdateCompanyBannerRequest) (*web.UploadCompanyBannerResponse, error)
}
