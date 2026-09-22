package service

import (
	"AlfianChabib/go-job-board-api/internal/model/web"
	"AlfianChabib/go-job-board-api/internal/repository"
	"AlfianChabib/go-job-board-api/pkg/errs"
	"context"
	"errors"
	"math"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type companyServiceImpl struct {
	companyRepo repository.CompanyRepository
}

func NewCompanyService(companyRepo repository.CompanyRepository) CompanyService {
	return &companyServiceImpl{
		companyRepo: companyRepo,
	}
}

func (service *companyServiceImpl) GetCompanies(ctx context.Context, req web.GetCompaniesRequest) (*web.CompanyPaginationResponse, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	companies, total, err := service.companyRepo.FindAll(ctx, req)
	if err != nil {
		return nil, err
	}

	companyResponses := make([]web.CompanyResponse, 0, len(companies))
	for _, c := range companies {
		companyResponses = append(companyResponses, *toCompanyResponse(&c))
	}

	totalPages := 0
	if limit > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(limit)))
	}

	return &web.CompanyPaginationResponse{
		Companies:  companyResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (service *companyServiceImpl) GetCompanyById(ctx context.Context, companyId uuid.UUID) (*web.CompanyDetailResponse, error) {
	company, err := service.companyRepo.FindById(ctx, companyId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrCompanyNotFound
		}
		return nil, err
	}

	companyRes := toCompanyResponse(company)

	return &web.CompanyDetailResponse{
		CompanyResponse: *companyRes,
		Jobs:            make([]web.CompanyJobResponse, 0),
	}, nil
}
