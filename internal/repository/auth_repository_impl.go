package repository

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type AuthRepositoryImpl struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepositoryImpl {
	return &AuthRepositoryImpl{
		db: db,
	}
}

func (repository *AuthRepositoryImpl) Register(ctx fiber.Ctx, user domain.User) (*domain.User, error) {
	repository.db.Model(&user).Create(&user)
	return &user, nil
}

func (repository *AuthRepositoryImpl) IsUserAlreadyExist(ctx fiber.Ctx, email string) bool {
	var user domain.User
	repository.db.Model(&user).Select("email").First(&user, "email = ?", email)
	return user.Email != ""
}
