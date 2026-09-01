package controller

import "github.com/gofiber/fiber/v3"

type AuthController interface {
	Register(c fiber.Ctx) error
	Login(c fiber.Ctx) error
	LogOut(c fiber.Ctx) error
}
