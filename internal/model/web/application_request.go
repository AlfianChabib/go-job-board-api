package web

import "github.com/google/uuid"

type ApplyJobRequest struct {
	JobId          uuid.UUID `uri:"id"`
	ResumeUrl      string    `json:"resume_url" validate:"required,url"`
	ResumeFilename *string   `json:"resume_filename,omitempty" validate:"omitempty,max=255"`
	CoverLetter    *string   `json:"cover_letter,omitempty"`
	ExpectedSalary *int64    `json:"expected_salary,omitempty" validate:"omitempty,min=0"`
}

type UpdateApplicationStatusRequest struct {
	ApplicationId   uuid.UUID `uri:"id"`
	Status          string    `json:"status" validate:"required,oneof=REVIEWING SHORTLISTED INTERVIEWING OFFERED HIRED REJECTED"`
	RecruiterNotes  *string   `json:"recruiter_notes,omitempty"`
	RejectionReason *string   `json:"rejection_reason,omitempty" validate:"omitempty,max=255"`
}

type WithdrawApplicationRequest struct {
	ApplicationId   uuid.UUID `uri:"id"`
	WithdrawnReason *string   `json:"withdrawn_reason,omitempty" validate:"omitempty,max=500"`
}

type GetJobApplicationsRequest struct {
	JobId  uuid.UUID `uri:"id"`
	Page   int       `query:"page" validate:"omitempty,min=1"`
	Limit  int       `query:"limit" validate:"omitempty,min=1,max=100"`
	Status string    `query:"status" validate:"omitempty,oneof=APPLIED REVIEWING SHORTLISTED INTERVIEWING OFFERED HIRED REJECTED WITHDRAWN"`
	Search string    `query:"search" validate:"omitempty"`
}

type GetCandidateApplicationsRequest struct {
	Page   int    `query:"page" validate:"omitempty,min=1"`
	Limit  int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Status string `query:"status" validate:"omitempty,oneof=APPLIED REVIEWING SHORTLISTED INTERVIEWING OFFERED HIRED REJECTED WITHDRAWN"`
}

type GetApplicationByIdRequest struct {
	ApplicationId uuid.UUID `uri:"id"`
}
