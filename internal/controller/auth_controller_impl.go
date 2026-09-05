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
	var req web.RegisterRequest

	if err := c.Bind().Body(&req); err != nil {
		return err
	}

	user, err := controller.AuthService.Register(ctx, req)
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
	var req web.LoginRequest

	if err := c.Bind().Body(&req); err != nil {
		return err
	}

	tokens, err := controller.AuthService.Login(ctx, req)
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

	var req web.LogOutRequest
	if err := c.Bind().Cookie(&req); err != nil {
		return err
	}

	err := controller.AuthService.Logout(ctx, req.RefreshToken)
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
		"message": "Logout Success",
	})
}

func (controller *authControllerImpl) RefreshToken(c fiber.Ctx) error {
	ctx := c.Context()
	var req web.RefreshTokenRequest
	if err := c.Bind().All(&req); err != nil {
		return err
	}

	response, err := controller.AuthService.RefreshToken(ctx, req)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    response.RefreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		Path:     "/",
		HTTPOnly: controller.isProd,
		Secure:   controller.isProd,
		SameSite: "Strict",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}
