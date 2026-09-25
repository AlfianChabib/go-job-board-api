package main

import (
	"AlfianChabib/go-job-board-api/internal/config"
	"AlfianChabib/go-job-board-api/internal/controller"
	"AlfianChabib/go-job-board-api/internal/middleware"
	"AlfianChabib/go-job-board-api/internal/router"
	"AlfianChabib/go-job-board-api/pkg/validator"
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func NewApp(
	env *config.Env,
	validate validator.StructValidator,
	mw middleware.Middleware,
	authController controller.AuthController,
	candidateController controller.CandidateController,
	recruiterController controller.RecruiterController,
	companyController controller.CompanyController,
	dataController controller.DataController,
) *fiber.App {
	app := fiber.New(fiber.Config{
		StructValidator: validate,
		ErrorHandler:    middleware.NewCustomErrorHandler(validate),
		JSONEncoder:     json.Marshal,
		JSONDecoder:     json.Unmarshal,
	})
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${time}] ${status} - ${latency} ${method} ${path}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowMethods: []string{"GET, POST, HEAD, PUT, DELETE, PATCH, QUERY"},
		AllowOrigins: []string{"*"},
	}))

	router.InitializeRoutes(
		app,
		env,
		mw,
		authController,
		candidateController,
		recruiterController,
		companyController,
		dataController,
	)

	return app
}

func main() {
	env := config.LoadEnv()
	app, err := InitializeApp(env)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	log.Fatal(app.Listen(env.Port, fiber.ListenConfig{
		EnablePrefork: true,
	}))
}
