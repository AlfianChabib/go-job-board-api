package middleware

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func Protected(jwtManager domain.JwtManager) fiber.Handler {
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
		decodedToken, err := jwtManager.ValidateToken(tokenString)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
		}

		c.Locals("userId", decodedToken.UserID)
		c.Locals("role", decodedToken.Role)

		return c.Next()
	}
}
