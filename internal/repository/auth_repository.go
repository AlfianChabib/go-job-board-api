package repository

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"

	"github.com/gofiber/fiber/v3"
)

type AuthRepository interface {
	Register(ctx fiber.Ctx, user domain.User) (*domain.User, error)
	IsUserAlreadyExist(ctx fiber.Ctx, email string) bool
}
