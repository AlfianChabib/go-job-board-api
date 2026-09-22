package service

import (
	"AlfianChabib/go-job-board-api/internal/model/web"
	"context"

	"github.com/google/uuid"
)

type CompanyService interface {
	GetCompanies(ctx context.Context, req web.GetCompaniesRequest) (*web.CompanyPaginationResponse, error)
	GetCompanyById(ctx context.Context, companyId uuid.UUID) (*web.CompanyDetailResponse, error)
}
