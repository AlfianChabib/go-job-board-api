package repository

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"AlfianChabib/go-job-board-api/internal/model/web"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type jobRepositoryImpl struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) JobRepository {
	return &jobRepositoryImpl{
		db: db,
	}
}

func (repo *jobRepositoryImpl) Create(ctx context.Context, job domain.Job) (*domain.Job, error) {
	if err := repo.db.WithContext(ctx).Create(&job).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

func (repo *jobRepositoryImpl) FindAll(ctx context.Context, req web.GetJobsRequest) ([]domain.Job, int64, error) {
	var jobs []domain.Job
	var total int64

	query := repo.db.WithContext(ctx).Model(&domain.Job{}).Where("status = ?", domain.JobStatusOpen)

	if req.Search != "" {
		searchTerm := "%" + req.Search + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ?", searchTerm, searchTerm)
	}

	if req.EmploymentType != "" {
		query = query.Where("employment_type = ?", req.EmploymentType)
	}

	if req.WorkMode != "" {
		query = query.Where("work_mode = ?", req.WorkMode)
	}

	if req.Location != "" {
		query = query.Where("location ILIKE ?", "%"+req.Location+"%")
	}

	if req.MinSalary != nil {
		query = query.Where("min_salary >= ?", *req.MinSalary)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
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

	if err := query.Preload("Company").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&jobs).Error; err != nil {
		return nil, 0, err
	}

	return jobs, total, nil
}

func (repo *jobRepositoryImpl) FindAllByCompanyId(ctx context.Context, companyId uuid.UUID, req web.GetRecruiterJobsRequest) ([]domain.Job, int64, error) {
	var jobs []domain.Job
	var total int64

	query := repo.db.WithContext(ctx).Model(&domain.Job{}).Where("company_id = ?", companyId)

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
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

	if err := query.Preload("Company").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&jobs).Error; err != nil {
		return nil, 0, err
	}

	return jobs, total, nil
}

func (repo *jobRepositoryImpl) FindById(ctx context.Context, id uuid.UUID) (*domain.Job, error) {
	var job domain.Job
	if err := repo.db.WithContext(ctx).
		Preload("Company").
		Where("id = ?", id).
		Take(&job).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

func (repo *jobRepositoryImpl) Update(ctx context.Context, job domain.Job) (*domain.Job, error) {
	if err := repo.db.WithContext(ctx).
		Model(&domain.Job{}).
		Where("id = ?", job.ID).
		Select("Title", "Description", "Requirements", "EmploymentType", "WorkMode", "Location", "MinSalary", "MaxSalary", "Currency", "IsSalaryNegotiable").
		Updates(&job).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

func (repo *jobRepositoryImpl) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.JobStatus) error {
	res := repo.db.WithContext(ctx).
		Model(&domain.Job{}).
		Where("id = ?", id).
		Update("status", status)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (repo *jobRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	res := repo.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&domain.Job{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
