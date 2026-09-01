package main

import (
	"AlfianChabib/go-job-board-api/internal/config"
	"AlfianChabib/go-job-board-api/internal/middleware"
	"AlfianChabib/go-job-board-api/internal/router"
	"AlfianChabib/go-job-board-api/pkg/validator"
	"log"

	"github.com/gofiber/fiber/v3/middleware/recover"

	"github.com/gofiber/fiber/v3"
)

func main() {
	var err error
	var env *config.Env
	env, err = config.LoadEnv()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	app := fiber.New(fiber.Config{
		StructValidator: validator.NewValidator(),
		ErrorHandler:    middleware.ErrorHandler,
	})
	app.Use(recover.New())

	router.Initialize(app)

	log.Fatal(app.Listen(env.Port, fiber.ListenConfig{
		EnablePrefork: true,
	}))
}
