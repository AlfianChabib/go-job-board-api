package router

import (
	"AlfianChabib/go-job-board-api/internal/controller"

	"github.com/gofiber/fiber/v3"
)

func SetupCandidateRoutes(router fiber.Router, protected fiber.Handler, candidateController controller.CandidateController) {
	candidate := router.Group("/candidate")
	candidate.Use(protected)

	candidate.Get("/", candidateController.Get)
	candidate.Put("/", candidateController.Update)
}
