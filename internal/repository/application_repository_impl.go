package repository

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"AlfianChabib/go-job-board-api/internal/model/web"
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type applicationRepositoryImpl struct {
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) ApplicationRepository {
	return &applicationRepositoryImpl{db: db}
}

func (r *applicationRepositoryImpl) Create(ctx context.Context, application domain.Application) (*domain.Application, error) {
	if err := r.db.WithContext(ctx).Create(&application).Error; err != nil {
		return nil, err
	}
	return &application, nil
}

func (r *applicationRepositoryImpl) FindById(ctx context.Context, id uuid.UUID) (*domain.Application, error) {
	var application domain.Application
	if err := r.db.WithContext(ctx).
		Preload("Job.Company").
		Preload("Candidate.Profile").
		Where("id = ?", id).
		Take(&application).Error; err != nil {
		return nil, err
	}
	return &application, nil
}

func (r *applicationRepositoryImpl) FindByJobAndCandidate(ctx context.Context, jobId uuid.UUID, candidateId uuid.UUID) (*domain.Application, error) {
	var application domain.Application
	if err := r.db.WithContext(ctx).
		Where("job_id = ? AND candidate_id = ?", jobId, candidateId).
		Take(&application).Error; err != nil {
		return nil, err
	}
	return &application, nil
}

func (r *applicationRepositoryImpl) FindAllByJobId(ctx context.Context, jobId uuid.UUID, req web.GetJobApplicationsRequest) ([]domain.Application, int64, error) {
	var applications []domain.Application
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Application{}).Where("applications.job_id = ?", jobId)

	if req.Status != "" {
		query = query.Where("applications.status = ?", req.Status)
	}

	if req.Search != "" {
		searchTerm := "%" + req.Search + "%"
		query = query.Joins("JOIN users ON users.id = applications.candidate_id").
			Where("users.name ILIKE ? OR users.email ILIKE ?", searchTerm, searchTerm)
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

	if err := query.Select("applications.*").
		Preload("Candidate.Profile").
		Preload("Job.Company").
		Order("applications.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&applications).Error; err != nil {
		return nil, 0, err
	}

	return applications, total, nil
}

func (r *applicationRepositoryImpl) FindAllByCandidateId(ctx context.Context, candidateId uuid.UUID, req web.GetCandidateApplicationsRequest) ([]domain.Application, int64, error) {
	var applications []domain.Application
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Application{}).Where("candidate_id = ?", candidateId)

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

	if err := query.Preload("Job.Company").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&applications).Error; err != nil {
		return nil, 0, err
	}

	return applications, total, nil
}

func (r *applicationRepositoryImpl) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.ApplicationStatus, recruiterNotes *string, rejectionReason *string) (*domain.Application, error) {
	var application domain.Application
	if err := r.db.WithContext(ctx).Preload("Job.Company").Preload("Candidate.Profile").Where("id = ?", id).Take(&application).Error; err != nil {
		return nil, err
	}

	application.Status = status
	application.StatusUpdatedAt = time.Now()
	if recruiterNotes != nil {
		application.RecruiterNotes = recruiterNotes
	}
	if rejectionReason != nil {
		application.RejectionReason = rejectionReason
	}

	if err := r.db.WithContext(ctx).Save(&application).Error; err != nil {
		return nil, err
	}

	return &application, nil
}

func (r *applicationRepositoryImpl) Withdraw(ctx context.Context, id uuid.UUID, withdrawnReason *string) (*domain.Application, error) {
	var application domain.Application
	if err := r.db.WithContext(ctx).Preload("Job.Company").Preload("Candidate.Profile").Where("id = ?", id).Take(&application).Error; err != nil {
		return nil, err
	}

	application.Status = domain.ApplicationStatusWithdrawn
	application.StatusUpdatedAt = time.Now()
	if withdrawnReason != nil {
		application.WithdrawnReason = withdrawnReason
	}

	if err := r.db.WithContext(ctx).Save(&application).Error; err != nil {
		return nil, err
	}

	return &application, nil
}
