package repository

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"context"
	"errors"

	"gorm.io/gorm"
)

type authRepositoryImpl struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepositoryImpl{db: db}
}

func (repo *authRepositoryImpl) Register(ctx context.Context, user domain.User) (*domain.User, error) {
	err := repo.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		return nil, errors.New("Failed to register new user")
	}

	return &user, nil
}

func (repo *authRepositoryImpl) IsUserExist(ctx context.Context, email string) bool {
	var user domain.User
	err := repo.db.WithContext(ctx).Model(&user).Select("email").Take(&user, "email = ?", email).Error
	if err != nil {
		return false
	}

	return user.Email != ""
}

func (repo *authRepositoryImpl) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := repo.db.WithContext(ctx).Preload("Auth").Take(&user, "email = ?", email).Error
	if err != nil {
		return nil, errors.New("User not found")
	}

	return &user, nil
}
