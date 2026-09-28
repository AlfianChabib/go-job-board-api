package router

import (
	"AlfianChabib/go-job-board-api/internal/controller"
	"AlfianChabib/go-job-board-api/internal/middleware"

	"github.com/gofiber/fiber/v3"
)

func SetupJobRoutes(
	router fiber.Router,
	middleware middleware.Middleware,
	jobController controller.JobController,
) {
	// Public Job Endpoints
	jobs := router.Group("/jobs")
	jobs.Get("/", jobController.GetJobs)
	jobs.Get("/:id", jobController.GetJobById)

	// Protected Job Endpoints (Recruiter Only)
	jobs.Post("/", middleware.Protected(), middleware.RequireRoles("RECRUITER"), jobController.CreateJob)
	jobs.Put("/:id", middleware.Protected(), middleware.RequireRoles("RECRUITER"), jobController.UpdateJob)
	jobs.Patch("/:id/status", middleware.Protected(), middleware.RequireRoles("RECRUITER"), jobController.UpdateJobStatus)
	jobs.Delete("/:id", middleware.Protected(), middleware.RequireRoles("RECRUITER"), jobController.DeleteJob)

	// Recruiter Job Endpoints (Recruiter Only)
	recruiterJobs := router.Group("/recruiter/jobs", middleware.Protected(), middleware.RequireRoles("RECRUITER"))
	recruiterJobs.Get("/", jobController.GetRecruiterJobs)
}
