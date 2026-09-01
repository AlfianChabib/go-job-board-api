package controller

import (
	"AlfianChabib/go-job-board-api/internal/model/web"
	"AlfianChabib/go-job-board-api/internal/service"

	"github.com/gofiber/fiber/v3"
)

type AuthControllerImpl struct {
	AuthService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthControllerImpl {
	return &AuthControllerImpl{
		AuthService: authService,
	}
}

func (authController *AuthControllerImpl) Register(c fiber.Ctx) error {
	var body web.RegisterRequest

	if err := c.Bind().Body(&body); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": fiber.StatusCreated,
		"data":   &body,
	})
}

func (authController *AuthControllerImpl) Login(c fiber.Ctx) error {
	panic("TODO: Implement")
}

func (authController *AuthControllerImpl) LogOut(c fiber.Ctx) error {
	panic("TODO: Implement")
}
