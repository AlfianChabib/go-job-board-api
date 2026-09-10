package repository

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type candidateRepository struct {
	db *gorm.DB
}

func NewCandidateRepository(db *gorm.DB) CandidateRepository {
	return &candidateRepository{
		db: db,
	}
}

func (repo *candidateRepository) Get(ctx context.Context, userId uuid.UUID) (*domain.Profile, error) {
	var candidate domain.Profile
	err := repo.db.WithContext(ctx).Take(&candidate, "user_id = ?", userId).Error
	if err != nil {
		return nil, err
	}

	return &candidate, nil
}

func (repo *candidateRepository) Update(ctx context.Context, candidate domain.Profile) (*domain.Profile, error) {
	err := repo.db.WithContext(ctx).
		Model(&candidate).
		Clauses(clause.Returning{}). // <-- Isi otomatis ID & kolom lainnya dari DB
		Where("user_id = ?", candidate.UserId).
		Select("headline", "phone").
		Updates(&candidate).Error
	if err != nil {
		return nil, err
	}
	return &candidate, nil
}

func (repo *candidateRepository) UploadAvatar(ctx context.Context, userId uuid.UUID, avatarUrl string) error {
	err := repo.db.WithContext(ctx).
		Model(&domain.Profile{}).
		Where("user_id = ?", userId).
		Update("avatar_url", avatarUrl).Error
	if err != nil {
		return err
	}

	return nil
}

func (repo *candidateRepository) DeleteAvatar(ctx context.Context, userId uuid.UUID) error {
	err := repo.db.WithContext(ctx).
		Model(&domain.Profile{}).
		Where("user_id = ?", userId).
		Update("avatar_url", nil).Error
	if err != nil {
		return err
	}

	return nil
}
