package repository

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"context"

	"github.com/gofiber/fiber/v3/log"
	"gorm.io/gorm"
)

type tokenRepositoryImpl struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepositoryImpl{db: db}
}

func (repo *tokenRepositoryImpl) Save(ctx context.Context, token domain.Token) (*domain.Token, error) {
	err := repo.db.Model(&token).Save(&token).Error
	if err != nil {
		log.Info(err)

		return nil, err
	}
	log.Info(token)

	return &token, nil
}
