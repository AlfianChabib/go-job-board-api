package service

import (
	"AlfianChabib/go-job-board-api/internal/model/web"

	"github.com/gofiber/fiber/v3"
)

type AuthService interface {
	Register(ctx fiber.Ctx, data web.RegisterRequest) (*web.RegisterResponse, error)
	Login(ctx fiber.Ctx, data web.LoginRrequest) error
	Logout(ctx fiber.Ctx) error
}
