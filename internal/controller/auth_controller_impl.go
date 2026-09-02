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

func (controller *AuthControllerImpl) Register(ctx fiber.Ctx) error {
	var body web.RegisterRequest

	if err := ctx.Bind().Body(&body); err != nil {
		return err
	}

	user, err := controller.AuthService.Register(ctx, body)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": fiber.StatusCreated,
		"data":   user,
	})
}

func (controller *AuthControllerImpl) Login(c fiber.Ctx) error {
	panic("TODO: Implement")
}

func (authController *AuthControllerImpl) LogOut(c fiber.Ctx) error {
	panic("TODO: Implement")
}
