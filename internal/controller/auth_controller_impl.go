package controller

import (
	"AlfianChabib/go-job-board-api/internal/model/web"
	"AlfianChabib/go-job-board-api/internal/service"
	"time"

	"github.com/gofiber/fiber/v3"
)

type authControllerImpl struct {
	AuthService service.AuthService
	isProd      bool
}

func NewAuthController(authService service.AuthService, env string) AuthController {
	return &authControllerImpl{
		AuthService: authService,
		isProd:      env == "production",
	}
}

func (controller *authControllerImpl) Register(c fiber.Ctx) error {
	ctx := c.Context()
	var body web.RegisterRequest

	if err := c.Bind().Body(&body); err != nil {
		return err
	}

	user, err := controller.AuthService.Register(ctx, body)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": fiber.StatusCreated,
		"data":   user,
	})
}

func (controller *authControllerImpl) Login(c fiber.Ctx) error {
	ctx := c.Context()
	var body web.LoginRequest

	if err := c.Bind().Body(&body); err != nil {
		return err
	}

	tokens, err := controller.AuthService.Login(ctx, body)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		Path:     "/",
		HTTPOnly: controller.isProd,
		Secure:   controller.isProd,
		SameSite: "Strict",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    tokens,
	})
}

func (controller *authControllerImpl) LogOut(c fiber.Ctx) error {
	ctx := c.Context()

	var data web.LogOutRequest
	if err := c.Bind().Cookie(&data); err != nil {
		return err
	}

	err := controller.AuthService.Logout(ctx, data.RefreshToken)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		MaxAge:   -1,
		Path:     "/",
		HTTPOnly: controller.isProd,
		Secure:   controller.isProd,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    "Logout Success",
	})
}
