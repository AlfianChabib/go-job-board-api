package controller

import (
	"github.com/gofiber/fiber/v3"
)

type CandidateController interface {
	Get(c fiber.Ctx) error
	Update(c fiber.Ctx) error
	UpdateAvatar(c fiber.Ctx) error
	DeleteAvatar(c fiber.Ctx) error
}
