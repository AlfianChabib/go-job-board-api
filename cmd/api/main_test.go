package main

import (
	"AlfianChabib/go-job-board-api/internal/config"
	"AlfianChabib/go-job-board-api/internal/middleware"
	"AlfianChabib/go-job-board-api/pkg/validator"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProvideFunctionsAndNewApp(t *testing.T) {
	env := &config.Env{
		AppEnv:           "test",
		AccessSecretKey:  "acc-secret",
		AccessDuration:   3600,
		RefreshSecretKey: "ref-secret",
		RefreshDuration:  7200,
		MinioPublicUrl:   "http://localhost:9000",
		MinioAvatarBucket: "avatars",
		MinioCvBucket:     "cvs",
	}

	hasher := ProvideBcryptHasher()
	require.NotNil(t, hasher)

	jwtMgr := ProvideJwtManager(env)
	require.NotNil(t, jwtMgr)

	authCtrl := ProvideAuthController(nil, env)
	require.NotNil(t, authCtrl)

	candCtrl := ProvideCandidateController(nil)
	require.NotNil(t, candCtrl)

	recCtrl := ProvideRecruiterController(nil)
	require.NotNil(t, recCtrl)

	compCtrl := ProvideCompanyController(nil)
	require.NotNil(t, compCtrl)

	dataCtrl := ProvideDataController(nil)
	require.NotNil(t, dataCtrl)

	jobCtrl := ProvideJobController(nil)
	require.NotNil(t, jobCtrl)

	appCtrl := ProvideApplicationController(nil)
	require.NotNil(t, appCtrl)

	storageRepo := ProvideStoragerepository(nil, env)
	require.NotNil(t, storageRepo)

	v := validator.NewValidator()
	mw := middleware.NewMiddleware(jwtMgr)

	app := NewApp(
		env,
		v,
		mw,
		authCtrl,
		candCtrl,
		recCtrl,
		compCtrl,
		dataCtrl,
		jobCtrl,
		appCtrl,
	)
	require.NotNil(t, app)

	resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/", nil))
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}
