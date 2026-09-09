package middleware

import (
	"AlfianChabib/go-job-board-api/internal/helper/response"
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"slices"
	"strings"

	"github.com/gofiber/fiber/v3"
)

type Middleware interface {
	Protected() fiber.Handler
	RequireRoles(allowedRoles ...domain.UserRole) fiber.Handler
}

type middleware struct {
	jwtManager domain.JwtManager
}

func NewMiddleware(jwtManager domain.JwtManager) Middleware {
	return &middleware{
		jwtManager: jwtManager,
	}
}

func (middleware *middleware) Protected() fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
		}

		tokenString := parts[1]
		decodedToken, err := middleware.jwtManager.ValidateAccessToken(tokenString)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}

		c.Locals("userId", decodedToken.UserID)
		c.Locals("role", decodedToken.Role)

		return c.Next()
	}
}

func (middleware *middleware) RequireRoles(allowedRoles ...domain.UserRole) fiber.Handler {
	return func(c fiber.Ctx) error {
		userRole, ok := c.Locals("role").(domain.UserRole)
		if !ok || userRole == "" {
			return response.Error(c, fiber.StatusUnauthorized, "Access denied: role not found in session.")
		}

		isAllowed := slices.Contains(allowedRoles, userRole)

		if !isAllowed {
			return response.Error(c, fiber.StatusForbidden, "You do not have access rights (Forbidden) for this action.")
		}

		return c.Next()
	}
}
