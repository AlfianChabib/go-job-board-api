package router_test

import (
	"AlfianChabib/go-job-board-api/internal/config"
	"AlfianChabib/go-job-board-api/internal/middleware"
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"AlfianChabib/go-job-board-api/internal/router"
	"AlfianChabib/go-job-board-api/pkg/utils"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubAuthController struct{}

func (s *stubAuthController) Register(c fiber.Ctx) error     { return c.SendStatus(fiber.StatusOK) }
func (s *stubAuthController) Login(c fiber.Ctx) error        { return c.SendStatus(fiber.StatusOK) }
func (s *stubAuthController) LogOut(c fiber.Ctx) error       { return c.SendStatus(fiber.StatusOK) }
func (s *stubAuthController) RefreshToken(c fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) }

type stubCandidateController struct{}

func (s *stubCandidateController) Get(c fiber.Ctx) error              { return c.SendStatus(fiber.StatusOK) }
func (s *stubCandidateController) Update(c fiber.Ctx) error           { return c.SendStatus(fiber.StatusOK) }
func (s *stubCandidateController) UpdateAvatar(c fiber.Ctx) error     { return c.SendStatus(fiber.StatusOK) }
func (s *stubCandidateController) DeleteAvatar(c fiber.Ctx) error     { return c.SendStatus(fiber.StatusOK) }
func (s *stubCandidateController) UpdateSkills(c fiber.Ctx) error     { return c.SendStatus(fiber.StatusOK) }
func (s *stubCandidateController) GetExperiences(c fiber.Ctx) error   { return c.SendStatus(fiber.StatusOK) }
func (s *stubCandidateController) CreateExperience(c fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) }
func (s *stubCandidateController) UpdateExperience(c fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) }
func (s *stubCandidateController) DeleteExperience(c fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) }

type stubRecruiterController struct{}

func (s *stubRecruiterController) CreateCompany(c fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) }
func (s *stubRecruiterController) GetCompany(c fiber.Ctx) error    { return c.SendStatus(fiber.StatusOK) }
func (s *stubRecruiterController) UpdateCompany(c fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) }
func (s *stubRecruiterController) UpdateLogo(c fiber.Ctx) error    { return c.SendStatus(fiber.StatusOK) }
func (s *stubRecruiterController) UpdateBanner(c fiber.Ctx) error  { return c.SendStatus(fiber.StatusOK) }

type stubCompanyController struct{}

func (s *stubCompanyController) GetCompanies(c fiber.Ctx) error   { return c.SendStatus(fiber.StatusOK) }
func (s *stubCompanyController) GetCompanyById(c fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) }

type stubDataController struct{}

func (s *stubDataController) GetSkills(c fiber.Ctx) error        { return c.SendStatus(fiber.StatusOK) }
func (s *stubDataController) GetCurrencyCodes(c fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) }

type stubJobController struct{}

func (s *stubJobController) CreateJob(c fiber.Ctx) error        { return c.SendStatus(fiber.StatusOK) }
func (s *stubJobController) GetJobs(c fiber.Ctx) error          { return c.SendStatus(fiber.StatusOK) }
func (s *stubJobController) GetRecruiterJobs(c fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) }
func (s *stubJobController) GetJobById(c fiber.Ctx) error       { return c.SendStatus(fiber.StatusOK) }
func (s *stubJobController) UpdateJob(c fiber.Ctx) error        { return c.SendStatus(fiber.StatusOK) }
func (s *stubJobController) UpdateJobStatus(c fiber.Ctx) error  { return c.SendStatus(fiber.StatusOK) }
func (s *stubJobController) DeleteJob(c fiber.Ctx) error        { return c.SendStatus(fiber.StatusOK) }

type stubApplicationController struct{}

func (s *stubApplicationController) ApplyJob(c fiber.Ctx) error                 { return c.SendStatus(fiber.StatusOK) }
func (s *stubApplicationController) GetJobApplications(c fiber.Ctx) error       { return c.SendStatus(fiber.StatusOK) }
func (s *stubApplicationController) GetCandidateApplications(c fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) }
func (s *stubApplicationController) GetApplicationById(c fiber.Ctx) error       { return c.SendStatus(fiber.StatusOK) }
func (s *stubApplicationController) UpdateStatus(c fiber.Ctx) error             { return c.SendStatus(fiber.StatusOK) }
func (s *stubApplicationController) WithdrawApplication(c fiber.Ctx) error      { return c.SendStatus(fiber.StatusOK) }

func TestInitializeRoutes(t *testing.T) {
	app := fiber.New()
	jwtMgr := utils.NewJwtManager("secret", 3600, "refresh", 7200)
	mw := middleware.NewMiddleware(jwtMgr)

	router.InitializeRoutes(
		app,
		&config.Env{},
		mw,
		&stubAuthController{},
		&stubCandidateController{},
		&stubRecruiterController{},
		&stubCompanyController{},
		&stubDataController{},
		&stubJobController{},
		&stubApplicationController{},
	)

	candidatePair, err := jwtMgr.GenerateTokenPair(uuid.New(), domain.RoleCandidate)
	require.NoError(t, err)
	recruiterPair, err := jwtMgr.GenerateTokenPair(uuid.New(), domain.RoleRecruiter)
	require.NoError(t, err)

	t.Run("GET / returns Hello world!", func(t *testing.T) {
		resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("GET /api/user/ returns hello", func(t *testing.T) {
		resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/api/user/", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(t, "hello", string(body))
	})

	t.Run("Public routes are reachable", func(t *testing.T) {
		publicRoutes := []struct {
			method string
			path   string
		}{
			{http.MethodGet, "/api/currencies"},
			{http.MethodGet, "/api/skills"},
			{http.MethodPost, "/api/auth/register"},
			{http.MethodPost, "/api/auth/login"},
			{http.MethodPost, "/api/auth/logout"},
			{http.MethodPost, "/api/auth/refresh"},
			{http.MethodGet, "/api/companies/"},
			{http.MethodGet, "/api/companies/123"},
			{http.MethodGet, "/api/jobs/"},
			{http.MethodGet, "/api/jobs/123"},
		}

		for _, rt := range publicRoutes {
			req := httptest.NewRequest(rt.method, rt.path, nil)
			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode, "failed route %s %s", rt.method, rt.path)
		}
	})

	t.Run("Candidate protected routes are reachable with candidate token", func(t *testing.T) {
		candRoutes := []struct {
			method string
			path   string
		}{
			{http.MethodGet, "/api/candidate/"},
			{http.MethodGet, "/api/candidate/profile"},
			{http.MethodPut, "/api/candidate/"},
			{http.MethodPut, "/api/candidate/profile"},
			{http.MethodPut, "/api/candidate/avatar"},
			{http.MethodPatch, "/api/candidate/avatar"},
			{http.MethodDelete, "/api/candidate/avatar"},
			{http.MethodPut, "/api/candidate/skills"},
			{http.MethodGet, "/api/candidate/experiences"},
			{http.MethodPost, "/api/candidate/experiences"},
			{http.MethodPut, "/api/candidate/experiences/123"},
			{http.MethodDelete, "/api/candidate/experiences/123"},
			{http.MethodGet, "/api/candidate/applications/"},
			{http.MethodPatch, "/api/candidate/applications/123/withdraw"},
			{http.MethodPost, "/api/jobs/123/applications/"},
			{http.MethodGet, "/api/applications/123"},
		}

		for _, rt := range candRoutes {
			req := httptest.NewRequest(rt.method, rt.path, nil)
			req.Header.Set("Authorization", "Bearer "+candidatePair.AccessToken)
			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode, "failed route %s %s", rt.method, rt.path)
		}
	})

	t.Run("Recruiter protected routes are reachable with recruiter token", func(t *testing.T) {
		recRoutes := []struct {
			method string
			path   string
		}{
			{http.MethodPost, "/api/recruiter/company"},
			{http.MethodGet, "/api/recruiter/company"},
			{http.MethodPut, "/api/recruiter/company"},
			{http.MethodPut, "/api/recruiter/company/logo"},
			{http.MethodPut, "/api/recruiter/company/banner"},
			{http.MethodPost, "/api/jobs/"},
			{http.MethodPut, "/api/jobs/123"},
			{http.MethodPatch, "/api/jobs/123/status"},
			{http.MethodDelete, "/api/jobs/123"},
			{http.MethodGet, "/api/recruiter/jobs/"},
			{http.MethodGet, "/api/jobs/123/applications/"},
			{http.MethodPatch, "/api/applications/123/status"},
		}

		for _, rt := range recRoutes {
			req := httptest.NewRequest(rt.method, rt.path, nil)
			req.Header.Set("Authorization", "Bearer "+recruiterPair.AccessToken)
			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode, "failed route %s %s", rt.method, rt.path)
		}
	})
}
