package repository

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
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

func (r *skillRepositoryImpl) FindAll(ctx context.Context) ([]domain.Skill, error) {
	var skills []domain.Skill
	if err := r.db.WithContext(ctx).Order("name ASC").Find(&skills).Error; err != nil {
		return nil, err
	}
	return skills, nil
}
