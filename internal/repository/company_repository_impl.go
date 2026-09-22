package repository

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"AlfianChabib/go-job-board-api/internal/model/web"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type companyRepositoryImpl struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) CompanyRepository {
	return &companyRepositoryImpl{
		db: db,
	}
}

func (repo *companyRepositoryImpl) Create(ctx context.Context, company domain.Company) (*domain.Company, error) {
	if err := repo.db.WithContext(ctx).Create(&company).Error; err != nil {
		return nil, err
	}
	return &company, nil
}

func (repo *companyRepositoryImpl) FindByRecruiterId(ctx context.Context, recruiterId uuid.UUID) (*domain.Company, error) {
	var company domain.Company
	if err := repo.db.WithContext(ctx).
		Where("recruiter_id = ?", recruiterId).
		Take(&company).Error; err != nil {
		return nil, err
	}
	return &company, nil
}

func (repo *companyRepositoryImpl) GetCompanyIdByRecruiterId(ctx context.Context, recruiterId uuid.UUID) (uuid.UUID, error) {
	var company domain.Company
	if err := repo.db.WithContext(ctx).
		Select("id").
		Where("recruiter_id = ?", recruiterId).
		Take(&company).Error; err != nil {
		return uuid.Nil, err
	}
	return company.ID, nil
}

func (repo *companyRepositoryImpl) FindById(ctx context.Context, id uuid.UUID) (*domain.Company, error) {
	var company domain.Company
	if err := repo.db.WithContext(ctx).
		Where("id = ?", id).
		Take(&company).Error; err != nil {
		return nil, err
	}
	return &company, nil
}

func (repo *companyRepositoryImpl) FindAll(ctx context.Context, req web.GetCompaniesRequest) ([]domain.Company, int64, error) {
	var companies []domain.Company
	var total int64

	query := repo.db.WithContext(ctx).Model(&domain.Company{})

	if req.Search != "" {
		query = query.Where("name ILIKE ?", "%"+req.Search+"%")
	}

	if req.Industry != "" {
		query = query.Where("industry ILIKE ?", req.Industry)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []domain.Company{}, 0, nil
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&companies).Error; err != nil {
		return nil, 0, err
	}

	return companies, total, nil
}

func (repo *companyRepositoryImpl) Update(ctx context.Context, company domain.Company) (*domain.Company, error) {
	query := repo.db.WithContext(ctx).
		Model(&company).
		Clauses(clause.Returning{})

	if company.ID != uuid.Nil {
		query = query.Where("id = ?", company.ID)
	} else {
		query = query.Where("recruiter_id = ?", company.RecruiterId)
	}

	result := query.Select("name", "location", "industry", "employee_size", "description", "website").
		Updates(&company)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &company, nil
}

func (repo *companyRepositoryImpl) UpdateLogo(ctx context.Context, recruiterId uuid.UUID, logoUrl string) error {
	result := repo.db.WithContext(ctx).
		Model(&domain.Company{}).
		Where("recruiter_id = ?", recruiterId).
		Update("logo_url", logoUrl)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (repo *companyRepositoryImpl) UpdateBanner(ctx context.Context, recruiterId uuid.UUID, bannerUrl string) error {
	result := repo.db.WithContext(ctx).
		Model(&domain.Company{}).
		Where("recruiter_id = ?", recruiterId).
		Update("banner_url", bannerUrl)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
