package request

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func GetAuthLocals(c fiber.Ctx) (uuid.UUID, domain.UserRole, error) {
	userId, ok := c.Locals("userId").(uuid.UUID)
	if !ok || userId == uuid.Nil {
		return uuid.Nil, "", fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	role, ok := c.Locals("role").(domain.UserRole)
	if !ok || role == "" {
		return uuid.Nil, "", fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	return userId, role, nil
}

type LocalSession struct {
	UserId uuid.UUID       `json:"user_id"`
	Role   domain.UserRole `json:"role"`
}

func GetLocalSession(c fiber.Ctx) (LocalSession, error) {
	userId, ok := c.Locals("userId").(uuid.UUID)
	if !ok || userId == uuid.Nil {
		return LocalSession{}, fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	role, ok := c.Locals("role").(domain.UserRole)
	if !ok || role == "" {
		return LocalSession{}, fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	return LocalSession{
		UserId: userId,
		Role:   role,
	}, nil
}
