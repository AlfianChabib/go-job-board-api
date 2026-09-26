package router

import (
	"AlfianChabib/go-job-board-api/internal/controller"
	"AlfianChabib/go-job-board-api/internal/middleware"

	"github.com/gofiber/fiber/v3"
)

func SetupCandidateRoutes(
	router fiber.Router,
	middleware middleware.Middleware,
	candidateController controller.CandidateController,
) {
	candidate := router.Group("/candidate")
	candidate.Use(middleware.Protected())
	candidate.Use(middleware.RequireRoles("CANDIDATE"))

	candidate.Get("/", candidateController.Get)
	candidate.Get("/profile", candidateController.Get)
	candidate.Put("/", candidateController.Update)
	candidate.Put("/profile", candidateController.Update)

	candidate.Put("/avatar", candidateController.UpdateAvatar)
	candidate.Patch("/avatar", candidateController.UpdateAvatar)
	candidate.Delete("/avatar", candidateController.DeleteAvatar)

	candidate.Put("/skills", candidateController.UpdateSkills)

	candidate.Get("/experiences", candidateController.GetExperiences)
	candidate.Post("/experiences", candidateController.CreateExperience)
	candidate.Put("/experiences/:experienceId", candidateController.UpdateExperience)
	candidate.Delete("/experiences/:experienceId", candidateController.DeleteExperience)
}
