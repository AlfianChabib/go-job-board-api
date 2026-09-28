package controller

import "github.com/gofiber/fiber/v3"

type ApplicationController interface {
	ApplyJob(c fiber.Ctx) error
	GetJobApplications(c fiber.Ctx) error
	GetCandidateApplications(c fiber.Ctx) error
	GetApplicationById(c fiber.Ctx) error
	UpdateStatus(c fiber.Ctx) error
	WithdrawApplication(c fiber.Ctx) error
}
