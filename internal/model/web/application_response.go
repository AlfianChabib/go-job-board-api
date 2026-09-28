package web

import (
	"time"

	"github.com/google/uuid"
)

// CandidateSummary provides candidate info for recruiter view.
type CandidateSummary struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Headline  *string   `json:"headline,omitempty"`
	Phone     *string   `json:"phone,omitempty"`
	AvatarUrl *string   `json:"avatar_url,omitempty"`
}

// CandidateApplicationResponse is sent to candidate.
// NOTE: RecruiterNotes is strictly omitted for privacy and security.
type CandidateApplicationResponse struct {
	ID              uuid.UUID    `json:"id"`
	JobId           uuid.UUID    `json:"job_id"`
	Job             *JobResponse `json:"job,omitempty"`
	ResumeUrl       string       `json:"resume_url"`
	ResumeFilename  *string      `json:"resume_filename,omitempty"`
	CoverLetter     *string      `json:"cover_letter,omitempty"`
	ExpectedSalary  *int64       `json:"expected_salary,omitempty"`
	Status          string       `json:"status"`
	RejectionReason *string      `json:"rejection_reason,omitempty"`
	WithdrawnReason *string      `json:"withdrawn_reason,omitempty"`
	AppliedAt       time.Time    `json:"applied_at"`
	StatusUpdatedAt time.Time    `json:"status_updated_at"`
}

// RecruiterApplicationResponse is sent to recruiter.
// Contains candidate details and internal recruiter notes.
type RecruiterApplicationResponse struct {
	ID              uuid.UUID         `json:"id"`
	JobId           uuid.UUID         `json:"job_id"`
	CandidateId     uuid.UUID         `json:"candidate_id"`
	Candidate       *CandidateSummary `json:"candidate,omitempty"`
	ResumeUrl       string            `json:"resume_url"`
	ResumeFilename  *string           `json:"resume_filename,omitempty"`
	CoverLetter     *string           `json:"cover_letter,omitempty"`
	ExpectedSalary  *int64            `json:"expected_salary,omitempty"`
	Status          string            `json:"status"`
	RecruiterNotes  *string           `json:"recruiter_notes,omitempty"`
	RejectionReason *string           `json:"rejection_reason,omitempty"`
	WithdrawnReason *string           `json:"withdrawn_reason,omitempty"`
	AppliedAt       time.Time         `json:"applied_at"`
	StatusUpdatedAt time.Time         `json:"status_updated_at"`
}

type CandidateApplicationPaginationResponse struct {
	Applications []CandidateApplicationResponse `json:"applications"`
	Total        int64                          `json:"total"`
	Page         int                            `json:"page"`
	Limit        int                            `json:"limit"`
	TotalPages   int                            `json:"total_pages"`
}

type RecruiterApplicationPaginationResponse struct {
	Applications []RecruiterApplicationResponse `json:"applications"`
	Total        int64                          `json:"total"`
	Page         int                            `json:"page"`
	Limit        int                            `json:"limit"`
	TotalPages   int                            `json:"total_pages"`
}
