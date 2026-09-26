package controller

import "github.com/gofiber/fiber/v3"

type JobController interface {
	CreateJob(c fiber.Ctx) error
	GetJobs(c fiber.Ctx) error
	GetRecruiterJobs(c fiber.Ctx) error
	GetJobById(c fiber.Ctx) error
	UpdateJob(c fiber.Ctx) error
	UpdateJobStatus(c fiber.Ctx) error
	DeleteJob(c fiber.Ctx) error
}
