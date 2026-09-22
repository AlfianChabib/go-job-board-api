package controller

import "github.com/gofiber/fiber/v3"

type CompanyController interface {
	GetCompanies(c fiber.Ctx) error
	GetCompanyById(c fiber.Ctx) error
}
