package router

import "github.com/gofiber/fiber/v3"

func SetupUserRoutes(router fiber.Router) {
	user := router.Group("/user")
	user.Get("/", func(c fiber.Ctx) error {
		return c.SendString("hello")
	})
}
