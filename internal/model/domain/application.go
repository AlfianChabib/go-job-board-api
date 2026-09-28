package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ApplicationStatus string

const (
	ApplicationStatusApplied      ApplicationStatus = "APPLIED"
	ApplicationStatusReviewing    ApplicationStatus = "REVIEWING"
	ApplicationStatusShortlisted  ApplicationStatus = "SHORTLISTED"
	ApplicationStatusInterviewing ApplicationStatus = "INTERVIEWING"
	ApplicationStatusOffered      ApplicationStatus = "OFFERED"
	ApplicationStatusHired        ApplicationStatus = "HIRED"
	ApplicationStatusRejected     ApplicationStatus = "REJECTED"
	ApplicationStatusWithdrawn    ApplicationStatus = "WITHDRAWN"
)

type Application struct {
	ID              uuid.UUID         `gorm:"column:id;type:uuid;primaryKey"`
	JobId           uuid.UUID         `gorm:"column:job_id;type:uuid;not null"`
	CandidateId     uuid.UUID         `gorm:"column:candidate_id;type:uuid;not null"`
	ResumeUrl       string            `gorm:"column:resume_url;type:varchar(255);not null"`
	ResumeFilename  *string           `gorm:"column:resume_filename;type:varchar(255)"`
	CoverLetter     *string           `gorm:"column:cover_letter;type:text"`
	ExpectedSalary  *int64            `gorm:"column:expected_salary;type:bigint"`
	Status          ApplicationStatus `gorm:"column:status;type:varchar(50);default:'APPLIED';not null"`
	RecruiterNotes  *string           `gorm:"column:recruiter_notes;type:text"`
	RejectionReason *string           `gorm:"column:rejection_reason;type:varchar(255)"`
	WithdrawnReason *string           `gorm:"column:withdrawn_reason;type:varchar(255)"`
	AppliedAt       time.Time         `gorm:"column:applied_at;autoCreateDate;<-:create"`
	StatusUpdatedAt time.Time         `gorm:"column:status_updated_at;autoCreateDate;autoUpdateDate"`
	CreatedAt       time.Time         `gorm:"column:created_at;autoCreateDate;<-:create"`
	UpdatedAt       time.Time         `gorm:"column:updated_at;autoCreateDate;autoUpdateDate"`
	DeletedAt       gorm.DeletedAt    `gorm:"column:deleted_at;index"`

	// Relasi
	Job       *Job  `gorm:"foreignKey:JobId;references:ID"`
	Candidate *User `gorm:"foreignKey:CandidateId;references:ID"`
}

func (a *Application) TableName() string {
	return "applications"
}

func (a *Application) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		uuidV7, _ := uuid.NewV7()
		a.ID = uuidV7
	}
	if a.Status == "" {
		a.Status = ApplicationStatusApplied
	}
	now := time.Now()
	if a.AppliedAt.IsZero() {
		a.AppliedAt = now
	}
	if a.StatusUpdatedAt.IsZero() {
		a.StatusUpdatedAt = now
	}
	return nil
}
