package controller

import "github.com/gofiber/fiber/v3"

type DataController interface {
	GetSkills(c fiber.Ctx) error
	GetCurrencyCodes(c fiber.Ctx) error
}
