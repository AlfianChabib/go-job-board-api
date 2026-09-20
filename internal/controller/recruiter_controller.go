package controller

import "github.com/gofiber/fiber/v3"

type RecruiterController interface {
	CreateCompany(c fiber.Ctx) error
	GetCompany(c fiber.Ctx) error
	UpdateCompany(c fiber.Ctx) error
	UpdateLogo(c fiber.Ctx) error
	UpdateBanner(c fiber.Ctx) error
}
