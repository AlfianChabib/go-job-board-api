package errs_test

import (
	"AlfianChabib/go-job-board-api/pkg/errs"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppError_Error(t *testing.T) {
	err := errs.New(http.StatusBadRequest, "custom error message")
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusBadRequest, err.Code)
	assert.Equal(t, "custom error message", err.Message)
	assert.Equal(t, "custom error message", err.Error())
}

func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		name         string
		err          *errs.AppError
		expectedCode int
	}{
		{"ErrUnauthorized", errs.ErrUnauthorized, http.StatusUnauthorized},
		{"ErrInvalidCredentials", errs.ErrInvalidCredentials, http.StatusUnauthorized},
		{"ErrTokenExpired", errs.ErrTokenExpired, http.StatusUnauthorized},
		{"ErrTokenRevoked", errs.ErrTokenRevoked, http.StatusUnauthorized},
		{"ErrForbidden", errs.ErrForbidden, http.StatusForbidden},
		{"ErrUserAlreadyExists", errs.ErrUserAlreadyExists, http.StatusConflict},
		{"ErrInvalidToken", errs.ErrInvalidToken, http.StatusUnauthorized},
		{"ErrInvalidInput", errs.ErrInvalidInput, http.StatusBadRequest},
		{"ErrMissingHeader", errs.ErrMissingHeader, http.StatusBadRequest},
		{"ErrUnprocessable", errs.ErrUnprocessable, http.StatusUnprocessableEntity},
		{"ErrNotFound", errs.ErrNotFound, http.StatusNotFound},
		{"ErrDuplicateData", errs.ErrDuplicateData, http.StatusConflict},
		{"ErrDataInUse", errs.ErrDataInUse, http.StatusConflict},
		{"ErrFileTooLarge", errs.ErrFileTooLarge, http.StatusRequestEntityTooLarge},
		{"ErrInvalidFileType", errs.ErrInvalidFileType, http.StatusUnsupportedMediaType},
		{"ErrAlreadyApplied", errs.ErrAlreadyApplied, http.StatusConflict},
		{"ErrProfileIncomplete", errs.ErrProfileIncomplete, http.StatusBadRequest},
		{"ErrInvalidAction", errs.ErrInvalidAction, http.StatusBadRequest},
		{"ErrExperienceNotFound", errs.ErrExperienceNotFound, http.StatusNotFound},
		{"ErrInternalServer", errs.ErrInternalServer, http.StatusInternalServerError},
		{"ErrStorageUnavailable", errs.ErrStorageUnavailable, http.StatusInternalServerError},
		{"ErrCompanyNotFound", errs.ErrCompanyNotFound, http.StatusNotFound},
		{"ErrCompanyAlreadyExists", errs.ErrCompanyAlreadyExists, http.StatusConflict},
		{"ErrUploadLogoFailed", errs.ErrUploadLogoFailed, http.StatusInternalServerError},
		{"ErrUploadBannerFailed", errs.ErrUploadBannerFailed, http.StatusInternalServerError},
		{"ErrJobNotFound", errs.ErrJobNotFound, http.StatusNotFound},
		{"ErrJobForbidden", errs.ErrJobForbidden, http.StatusForbidden},
		{"ErrJobClosed", errs.ErrJobClosed, http.StatusBadRequest},
		{"ErrInvalidJobStatus", errs.ErrInvalidJobStatus, http.StatusBadRequest},
		{"ErrInvalidSalaryRange", errs.ErrInvalidSalaryRange, http.StatusBadRequest},
		{"ErrApplicationNotFound", errs.ErrApplicationNotFound, http.StatusNotFound},
		{"ErrApplicationForbidden", errs.ErrApplicationForbidden, http.StatusForbidden},
		{"ErrInvalidApplicationStatus", errs.ErrInvalidApplicationStatus, http.StatusBadRequest},
		{"ErrApplicationAlreadyWithdrawn", errs.ErrApplicationAlreadyWithdrawn, http.StatusBadRequest},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.NotNil(t, tc.err)
			assert.Equal(t, tc.expectedCode, tc.err.Code)
			assert.NotEmpty(t, tc.err.Message)
			assert.Equal(t, tc.err.Message, tc.err.Error())
		})
	}
}
