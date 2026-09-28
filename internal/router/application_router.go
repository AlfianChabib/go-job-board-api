package router

import (
	"AlfianChabib/go-job-board-api/internal/controller"
	"AlfianChabib/go-job-board-api/internal/middleware"

	"github.com/gofiber/fiber/v3"
)

func SetupApplicationRoutes(
	router fiber.Router,
	middleware middleware.Middleware,
	applicationController controller.ApplicationController,
) {
	// Candidate Application Routes
	candidate := router.Group("/candidate/applications")
	candidate.Use(middleware.Protected())
	candidate.Use(middleware.RequireRoles("CANDIDATE"))
	candidate.Get("/", applicationController.GetCandidateApplications)
	candidate.Patch("/:id/withdraw", applicationController.WithdrawApplication)

	// Job Applications Routes (Candidate melamar & Recruiter melihat pelamar lowongan)
	jobs := router.Group("/jobs/:id/applications")
	jobs.Use(middleware.Protected())
	jobs.Post("/", middleware.RequireRoles("CANDIDATE"), applicationController.ApplyJob)
	jobs.Get("/", middleware.RequireRoles("RECRUITER"), applicationController.GetJobApplications)

	// General Applications Routes
	applications := router.Group("/applications")
	applications.Use(middleware.Protected())
	applications.Get("/:id", applicationController.GetApplicationById)
	applications.Patch("/:id/status", middleware.RequireRoles("RECRUITER"), applicationController.UpdateStatus)
}
