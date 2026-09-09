package controller

import (
	"AlfianChabib/go-job-board-api/internal/helper/response"
	"AlfianChabib/go-job-board-api/internal/model/web"
	"AlfianChabib/go-job-board-api/internal/service"
	"time"

	"github.com/gofiber/fiber/v3"
)

type authController struct {
	AuthService service.AuthService
	isProd      bool
}

func NewAuthController(authService service.AuthService, env string) AuthController {
	return &authController{
		AuthService: authService,
		isProd:      env == "production",
	}
}

func (controller *authController) Register(c fiber.Ctx) error {
	ctx := c.Context()
	var req web.RegisterRequest

	if err := c.Bind().Body(&req); err != nil {
		return err
	}

	user, err := controller.AuthService.Register(ctx, req)
	if err != nil {
		return err
	}

	return response.OK(c, "Register Success", user)
}

func (controller *authController) Login(c fiber.Ctx) error {
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

	return response.Success(c, fiber.StatusOK, "Login success", tokens)
}

func (controller *authController) LogOut(c fiber.Ctx) error {
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

	return response.Message(c, fiber.StatusOK, "Logout Success")
}

func (controller *authController) RefreshToken(c fiber.Ctx) error {
	ctx := c.Context()
	var req web.RefreshTokenRequest
	if err := c.Bind().All(&req); err != nil {
		return err
	}

	tokens, err := controller.AuthService.RefreshToken(ctx, req)
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

	return response.Success(c, fiber.StatusOK, "Success refresh token", tokens)
}
