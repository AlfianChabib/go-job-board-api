package router

import (
	"AlfianChabib/go-job-board-api/internal/controller"
	"AlfianChabib/go-job-board-api/internal/middleware"

	"github.com/gofiber/fiber/v3"
)

func SetupRecruiterRoutes(
	router fiber.Router,
	middleware middleware.Middleware,
	recruiterController controller.RecruiterController,
) {
	recruiter := router.Group("/recruiter")
	recruiter.Use(middleware.Protected())
	recruiter.Use(middleware.RequireRoles("RECRUITER"))

	recruiter.Post("/company", recruiterController.CreateCompany)
	recruiter.Get("/company", recruiterController.GetCompany)
	recruiter.Put("/company", recruiterController.UpdateCompany)
	recruiter.Put("/company/logo", recruiterController.UpdateLogo)
	recruiter.Put("/company/banner", recruiterController.UpdateBanner)
}
