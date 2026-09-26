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
	protectedJobs := jobs.Group("/", middleware.Protected(), middleware.RequireRoles("RECRUITER"))
	protectedJobs.Post("/", jobController.CreateJob)
	protectedJobs.Put("/:id", jobController.UpdateJob)
	protectedJobs.Patch("/:id/status", jobController.UpdateJobStatus)
	protectedJobs.Delete("/:id", jobController.DeleteJob)

	// Recruiter Job Endpoints (Recruiter Only)
	recruiterJobs := router.Group("/recruiter/jobs", middleware.Protected(), middleware.RequireRoles("RECRUITER"))
	recruiterJobs.Get("/", jobController.GetRecruiterJobs)
}
