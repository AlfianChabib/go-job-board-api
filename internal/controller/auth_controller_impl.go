package controller

import (
	"AlfianChabib/go-job-board-api/internal/model/web"
	"AlfianChabib/go-job-board-api/internal/service"

	"github.com/gofiber/fiber/v3"
)

type authControllerImpl struct {
	AuthService service.AuthService
}

func NewAuthController(authService service.AuthService) AuthController {
	return &authControllerImpl{
		AuthService: authService,
	}
}

func (controller *authControllerImpl) Register(ctx fiber.Ctx) error {
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

func (controller *authControllerImpl) Login(c fiber.Ctx) error {
	panic("TODO: Implement")
}

func (a *authControllerImpl) LogOut(c fiber.Ctx) error {
	panic("TODO: Implement")
}
