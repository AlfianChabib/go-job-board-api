package repository

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"AlfianChabib/go-job-board-api/internal/model/web"
	"context"

	"github.com/google/uuid"
)

type CompanyRepository interface {
	Create(ctx context.Context, company domain.Company) (*domain.Company, error)
	FindByRecruiterId(ctx context.Context, recruiterId uuid.UUID) (*domain.Company, error)
	GetCompanyIdByRecruiterId(ctx context.Context, recruiterId uuid.UUID) (uuid.UUID, error)
	FindById(ctx context.Context, id uuid.UUID) (*domain.Company, error)
	FindAll(ctx context.Context, req web.GetCompaniesRequest) ([]domain.Company, int64, error)
	Update(ctx context.Context, company domain.Company) (*domain.Company, error)
	UpdateLogo(ctx context.Context, recruiterId uuid.UUID, logoUrl string) error
	UpdateBanner(ctx context.Context, recruiterId uuid.UUID, bannerUrl string) error
}
