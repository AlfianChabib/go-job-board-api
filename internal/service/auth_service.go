package service

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"AlfianChabib/go-job-board-api/internal/model/web"

	"github.com/gofiber/fiber/v3"
)

type AuthService interface {
	Register(ctx fiber.Ctx, data web.RegisterRequest) (*domain.User, error)
	Login(ctx fiber.Ctx, data web.LoginRrequest) error
	Logout(ctx fiber.Ctx) error
}
