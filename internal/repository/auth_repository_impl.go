package repository

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"errors"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type authRepositoryImpl struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepositoryImpl{db: db}
}

func (repo *authRepositoryImpl) Register(ctx fiber.Ctx, user domain.User) (*domain.User, error) {
	err := repo.db.Model(&user).Create(&user).Error
	if err != nil {
		return nil, errors.New("Failed to register new user")
	}

	return &user, nil
}

func (repo *authRepositoryImpl) IsUserAlreadyExist(ctx fiber.Ctx, email string) bool {
	var user domain.User
	err := repo.db.Model(&user).Select("email").Take(&user, "email = ?", email).Error
	if err != nil {
		return false
	}

	return user.Email != ""
}
