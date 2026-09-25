package web

type CreateJobRequest struct {
	Title              string  `json:"title" validate:"required,min=3,max=255"`
	Description        string  `json:"description" validate:"required,min=10"`
	Requirements       *string `json:"requirements,omitempty" validate:"omitempty"`
	EmploymentType     string  `json:"employment_type" validate:"required,oneof=FULL_TIME PART_TIME CONTRACT INTERNSHIP FREELANCE"`
	WorkMode           string  `json:"work_mode" validate:"required,oneof=ON_SITE HYBRID REMOTE"`
	Location           string  `json:"location" validate:"required,max=255"`
	MinSalary          *int64  `json:"min_salary,omitempty" validate:"omitempty,min=0"`
	MaxSalary          *int64  `json:"max_salary,omitempty" validate:"omitempty,gtefield=MinSalary"`
	Currency           string  `json:"currency,omitempty" validate:"omitempty,len=3"`
	IsSalaryNegotiable bool    `json:"is_salary_negotiable"`
}
