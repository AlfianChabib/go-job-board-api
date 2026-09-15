package repository

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v3/log"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type tokenRepositoryImpl struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepositoryImpl{db: db}
}

func (repo *tokenRepositoryImpl) Save(ctx context.Context, token domain.Token) (*domain.Token, error) {
	err := repo.db.WithContext(ctx).Model(&token).Save(&token).Error
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (repo *tokenRepositoryImpl) RevokeToken(ctx context.Context, refreshToken string) error {
	now := time.Now()
	result := repo.db.WithContext(ctx).Model(&domain.Token{}).
		Where("refresh_token = ? AND revoked_at IS NULL", refreshToken).
		Update("revoked_at", now)

	if result.Error != nil {
		return errors.New("Internal server error")
	}
	if result.RowsAffected == 0 {
		return errors.New("Sesi tidak valid atau sudah berakhir")
	}
	return nil
}

func (repo *tokenRepositoryImpl) FindTokenWithUser(ctx context.Context, userId uuid.UUID, refreshToken string) (*domain.Token, error) {
	var token domain.Token
	err := repo.db.WithContext(ctx).
		Where("refresh_token = ? AND user_id = ?", refreshToken, userId).
		Preload("User").
		First(&token).Error
	if err != nil {
		log.Info(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("token tidak valid atau sudah dicabut")
		}

		return nil, err
	}

	return &token, nil
}
