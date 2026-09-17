package errs

import (
	"net/http"
)

type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func New(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

var (
	// ==========================================
	// 1. AUTHENTICATION & AUTHORIZATION (401, 403)
	// ==========================================
	ErrUnauthorized       = &AppError{Code: http.StatusUnauthorized, Message: "Your session is invalid or has expired. Please log in again."}
	ErrInvalidCredentials = &AppError{Code: http.StatusUnauthorized, Message: "Invalid email or password."}
	ErrTokenExpired       = &AppError{Code: http.StatusUnauthorized, Message: "Access token has expired."}
	ErrTokenRevoked       = &AppError{Code: http.StatusUnauthorized, Message: "Session has been revoked."}
	ErrForbidden          = &AppError{Code: http.StatusForbidden, Message: "You do not have permission to perform this action."}

	// ==========================================
	// 2. CLIENT REQUEST & VALIDATION (400, 422)
	// ==========================================
	ErrInvalidInput  = &AppError{Code: http.StatusBadRequest, Message: "Invalid request input format."}
	ErrMissingHeader = &AppError{Code: http.StatusBadRequest, Message: "A required header is missing from the request."}
	ErrUnprocessable = &AppError{Code: http.StatusUnprocessableEntity, Message: "The request is well-formed but cannot be processed due to business logic errors."}

	// ==========================================
	// 3. RESOURCE & DATABASE (404, 409)
	// ==========================================
	ErrNotFound      = &AppError{Code: http.StatusNotFound, Message: "Resource not found."}
	ErrDuplicateData = &AppError{Code: http.StatusConflict, Message: "Data already exists or is duplicated."}
	ErrDataInUse     = &AppError{Code: http.StatusConflict, Message: "Resource cannot be deleted because it is currently in use by another entity."}

	// ==========================================
	// 4. FILE UPLOADS & MEDIA (413, 415)
	// ==========================================
	ErrFileTooLarge    = &AppError{Code: http.StatusRequestEntityTooLarge, Message: "File size exceeds the maximum allowed limit."}
	ErrInvalidFileType = &AppError{Code: http.StatusUnsupportedMediaType, Message: "Unsupported file format."}

	// ==========================================
	// 5. BUSINESS LOGIC (Specific to Job Board) (400, 409)
	// ==========================================
	ErrAlreadyApplied    = &AppError{Code: http.StatusConflict, Message: "You have already applied for this job posting."}
	ErrProfileIncomplete = &AppError{Code: http.StatusBadRequest, Message: "Your profile is incomplete. Please complete your profile before applying."}
	ErrInvalidAction     = &AppError{Code: http.StatusBadRequest, Message: "This action is not allowed in the current state."}

	// ==========================================
	// 6. Candidate
	// ==========================================
	ErrExperienceNotFound = &AppError{Code: http.StatusNotFound, Message: "Experience not found"}
)
