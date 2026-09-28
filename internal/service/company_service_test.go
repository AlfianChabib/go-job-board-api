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
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestCompanyService_GetCompanies(t *testing.T) {
	ctx := context.Background()

	t.Run("Repository error returns error", func(t *testing.T) {
		repo := new(mockCompanyRepository)
		req := web.GetCompaniesRequest{}
		repo.On("FindAll", ctx, req).Return(nil, int64(0), errors.New("db error"))

		svc := NewCompanyService(repo)
		res, err := svc.GetCompanies(ctx, req)
		assert.Nil(t, res)
		assert.EqualError(t, err, "db error")
	})

	t.Run("Success with default pagination", func(t *testing.T) {
		repo := new(mockCompanyRepository)
		req := web.GetCompaniesRequest{Page: 0, Limit: 0}
		companies := []domain.Company{
			{ID: uuid.New(), Name: "Tech Co", Industry: "IT", Location: "Jakarta", EmployeeSize: domain.EmployeeSize1To50},
		}
		repo.On("FindAll", ctx, req).Return(companies, int64(15), nil)

		svc := NewCompanyService(repo)
		res, err := svc.GetCompanies(ctx, req)
		require.NoError(t, err)
		assert.Len(t, res.Companies, 1)
		assert.Equal(t, int64(15), res.Total)
		assert.Equal(t, 1, res.Page)
		assert.Equal(t, 10, res.Limit)
		assert.Equal(t, 2, res.TotalPages)
	})
}

func TestCompanyService_GetCompanyById(t *testing.T) {
	ctx := context.Background()
	companyID := uuid.New()

	t.Run("Not found returns ErrCompanyNotFound", func(t *testing.T) {
		repo := new(mockCompanyRepository)
		repo.On("FindById", ctx, companyID).Return(nil, gorm.ErrRecordNotFound)

		svc := NewCompanyService(repo)
		res, err := svc.GetCompanyById(ctx, companyID)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrCompanyNotFound, err)
	})

	t.Run("Generic error returns error", func(t *testing.T) {
		repo := new(mockCompanyRepository)
		repo.On("FindById", ctx, companyID).Return(nil, errors.New("db failure"))

		svc := NewCompanyService(repo)
		res, err := svc.GetCompanyById(ctx, companyID)
		assert.Nil(t, res)
		assert.EqualError(t, err, "db failure")
	})

	t.Run("Success returns CompanyDetailResponse", func(t *testing.T) {
		repo := new(mockCompanyRepository)
		company := &domain.Company{
			ID:           companyID,
			Name:         "Acme Corp",
			Industry:     "Finance",
			Location:     "Bandung",
			EmployeeSize: domain.EmployeeSize51To200,
		}
		repo.On("FindById", ctx, companyID).Return(company, nil)

		svc := NewCompanyService(repo)
		res, err := svc.GetCompanyById(ctx, companyID)
		require.NoError(t, err)
		assert.Equal(t, companyID, res.ID)
		assert.Equal(t, "Acme Corp", res.Name)
		assert.NotNil(t, res.Jobs)
		assert.Empty(t, res.Jobs)
	})
}
