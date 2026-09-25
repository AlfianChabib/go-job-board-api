package repository

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"AlfianChabib/go-job-board-api/internal/model/web"
	"context"

	"gorm.io/gorm"
)

type skillRepositoryImpl struct {
	db *gorm.DB
}

func NewSkillRepository(db *gorm.DB) SkillRepository {
	return &skillRepositoryImpl{
		db: db,
	}
}

func (r *skillRepositoryImpl) FindAll(ctx context.Context, req web.GetDataRequest) ([]domain.Skill, error) {
	var skills []domain.Skill
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	query := r.db.WithContext(ctx).Model(&domain.Skill{})

	if req.Search != "" {
		searchTerm := "%" + req.Search + "%"
		query = query.Where("name ILIKE ? OR label ILIKE ? OR abbreviation ILIKE ?", searchTerm, searchTerm, searchTerm)
	}

	if err := query.Order("name ASC").Limit(limit).Find(&skills).Error; err != nil {
		return nil, err
	}
	return skills, nil
}
