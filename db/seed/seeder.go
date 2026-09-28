package main

import (
	"AlfianChabib/go-job-board-api/internal/config"
	"AlfianChabib/go-job-board-api/internal/database"
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"AlfianChabib/go-job-board-api/pkg/utils"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MockCompany struct {
	Name         string              `json:"name"`
	LogoUrl      *string             `json:"logo_url"`
	BannerUrl    *string             `json:"banner_url"`
	Website      *string             `json:"website"`
	Industry     string              `json:"industry"`
	EmployeeSize domain.EmployeeSize `json:"employee_size"`
	Description  *string             `json:"description"`
	Location     string              `json:"location"`
}

type MockRecruiter struct {
	Name     string      `json:"name"`
	Email    string      `json:"email"`
	Password string      `json:"password"`
	Company  MockCompany `json:"company"`
}

type MockJob struct {
	CompanyName        string                `json:"company_name"`
	Title              string                `json:"title"`
	Description        string                `json:"description"`
	Requirements       *string               `json:"requirements"`
	EmploymentType     domain.EmploymentType `json:"employment_type"`
	WorkMode           domain.WorkMode       `json:"work_mode"`
	Location           string                `json:"location"`
	MinSalary          *int64                `json:"min_salary"`
	MaxSalary          *int64                `json:"max_salary"`
	Currency           string                `json:"currency"`
	IsSalaryNegotiable bool                  `json:"is_salary_negotiable"`
	Status             domain.JobStatus      `json:"status"`
}

type MockCandidateProfile struct {
	Headline  *string `json:"headline"`
	Bio       *string `json:"bio"`
	Phone     *string `json:"phone"`
	AvatarUrl *string `json:"avatar_url"`
	ResumeUrl *string `json:"resume_url"`
}

type MockExperience struct {
	CompanyName string  `json:"company_name"`
	Position    string  `json:"position"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date"`
	IsCurrent   bool    `json:"is_current"`
	Description *string `json:"description"`
}

type MockCandidate struct {
	Name        string               `json:"name"`
	Email       string               `json:"email"`
	Password    string               `json:"password"`
	Profile     MockCandidateProfile `json:"profile"`
	Skills      []string             `json:"skills"`
	Experiences []MockExperience     `json:"experiences"`
}

// SeedSkills seeds the master skills dataset from CSV
func SeedSkills(db *gorm.DB) {
	file, err := os.Open("db/seed/skills-dataset.csv")
	if err != nil {
		log.Println("[Skills] Skipping skills-dataset.csv (file not found):", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var skills []domain.Skill
	seen := make(map[string]int)

	_, err = reader.Read()
	if err != nil && err != io.EOF {
		log.Fatal("Error reading CSV header:", err)
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		if len(record) < 3 {
			continue
		}

		name := strings.TrimSpace(record[1])
		if name == "" {
			continue
		}

		label := strings.TrimSpace(record[2])
		abbr := ""
		if len(record) > 3 {
			abbr = strings.TrimSpace(record[3])
		}

		if idx, exists := seen[name]; exists {
			if abbr != "" && skills[idx].Abbreviation == "" {
				skills[idx].Abbreviation = abbr
			}
			if label != "" && skills[idx].Label == "" {
				skills[idx].Label = label
			}
			continue
		}

		seen[name] = len(skills)
		skills = append(skills, domain.Skill{
			Name:         name,
			Label:        label,
			Abbreviation: abbr,
		})
	}

	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"label", "abbreviation"}),
	}).CreateInBatches(&skills, 100)

	if result.Error != nil {
		log.Fatal("Error creating skills:", result.Error)
	}

	fmt.Printf("✅ [Skills] Successfully seeded/updated %d master skills\n", result.RowsAffected)
}

// SeedRecruiters seeds recruiters and their companies from mock_recruiters.json
func SeedRecruiters(db *gorm.DB) {
	bytes, err := os.ReadFile("db/seed/mock_recruiters.json")
	if err != nil {
		log.Fatal("Error reading mock_recruiters.json:", err)
	}

	var recruiters []MockRecruiter
	if err := json.Unmarshal(bytes, &recruiters); err != nil {
		log.Fatal("Error unmarshaling mock_recruiters.json:", err)
	}

	hasher := utils.NewBcryptHasher(bcrypt.DefaultCost)
	createdCount := 0
	for _, r := range recruiters {
		var user domain.User
		err := db.Where("email = ?", r.Email).First(&user).Error
		if err != nil && err == gorm.ErrRecordNotFound {
			user = domain.User{
				Name:  r.Name,
				Email: r.Email,
				Role:  domain.RoleRecruiter,
			}
			if err := db.Create(&user).Error; err != nil {
				log.Printf("Failed to create recruiter %s: %v\n", r.Email, err)
				continue
			}

			hashedPassword, err := hasher.Hash([]byte(r.Password))
			if err != nil {
				log.Printf("Failed to hash password for %s: %v\n", r.Email, err)
				continue
			}
			auth := domain.Auth{
				UserId:   user.ID,
				Email:    user.Email,
				Password: hashedPassword,
			}
			_ = db.Create(&auth).Error
			createdCount++
		}

		// Ensure Company exists for this recruiter
		var company domain.Company
		err = db.Where("recruiter_id = ?", user.ID).First(&company).Error
		if err != nil && err == gorm.ErrRecordNotFound {
			company = domain.Company{
				RecruiterId:  user.ID,
				Name:         r.Company.Name,
				LogoUrl:      r.Company.LogoUrl,
				BannerUrl:    r.Company.BannerUrl,
				Website:      r.Company.Website,
				Industry:     r.Company.Industry,
				EmployeeSize: r.Company.EmployeeSize,
				Description:  r.Company.Description,
				Location:     r.Company.Location,
			}
			if err := db.Create(&company).Error; err != nil {
				log.Printf("Failed to create company %s: %v\n", r.Company.Name, err)
			}
		}
	}

	fmt.Printf("✅ [Recruiters & Companies] Successfully processed %d recruiters & companies (New: %d)\n", len(recruiters), createdCount)
}

// SeedJobs seeds job postings from mock_jobs.json
func SeedJobs(db *gorm.DB) {
	bytes, err := os.ReadFile("db/seed/mock_jobs.json")
	if err != nil {
		log.Fatal("Error reading mock_jobs.json:", err)
	}

	var jobs []MockJob
	if err := json.Unmarshal(bytes, &jobs); err != nil {
		log.Fatal("Error unmarshaling mock_jobs.json:", err)
	}

	createdCount := 0
	for _, j := range jobs {
		var company domain.Company
		if err := db.Where("name = ?", j.CompanyName).First(&company).Error; err != nil {
			log.Printf("Company %s not found for job %s: %v\n", j.CompanyName, j.Title, err)
			continue
		}

		var existingJob domain.Job
		err := db.Where("company_id = ? AND title = ?", company.ID, j.Title).First(&existingJob).Error
		if err != nil && err == gorm.ErrRecordNotFound {
			currency := j.Currency
			if currency == "" {
				currency = "IDR"
			}
			job := domain.Job{
				CompanyId:          company.ID,
				Title:              j.Title,
				Description:        j.Description,
				Requirements:       j.Requirements,
				EmploymentType:     j.EmploymentType,
				WorkMode:           j.WorkMode,
				Location:           j.Location,
				MinSalary:          j.MinSalary,
				MaxSalary:          j.MaxSalary,
				Currency:           currency,
				IsSalaryNegotiable: j.IsSalaryNegotiable,
				Status:             j.Status,
			}
			if err := db.Create(&job).Error; err != nil {
				log.Printf("Failed to create job %s: %v\n", j.Title, err)
				continue
			}
			createdCount++
		}
	}

	fmt.Printf("✅ [Jobs] Successfully processed %d job listings (New: %d)\n", len(jobs), createdCount)
}

// SeedCandidates seeds candidates, profiles, skills, and work experiences from mock_candidates.json
func SeedCandidates(db *gorm.DB) {
	bytes, err := os.ReadFile("db/seed/mock_candidates.json")
	if err != nil {
		log.Fatal("Error reading mock_candidates.json:", err)
	}

	var candidates []MockCandidate
	if err := json.Unmarshal(bytes, &candidates); err != nil {
		log.Fatal("Error unmarshaling mock_candidates.json:", err)
	}

	hasher := utils.NewBcryptHasher(bcrypt.DefaultCost)
	createdCount := 0
	for _, c := range candidates {
		var user domain.User
		err := db.Where("email = ?", c.Email).First(&user).Error
		if err != nil && err == gorm.ErrRecordNotFound {
			user = domain.User{
				Name:  c.Name,
				Email: c.Email,
				Role:  domain.RoleCandidate,
			}
			if err := db.Create(&user).Error; err != nil {
				log.Printf("Failed to create candidate %s: %v\n", c.Email, err)
				continue
			}

			hashedPassword, err := hasher.Hash([]byte(c.Password))
			if err != nil {
				log.Printf("Failed to hash password for %s: %v\n", c.Email, err)
				continue
			}
			auth := domain.Auth{
				UserId:   user.ID,
				Email:    user.Email,
				Password: hashedPassword,
			}
			_ = db.Create(&auth).Error
			createdCount++
		}

		// Ensure Profile exists
		var profile domain.Profile
		err = db.Where("user_id = ?", user.ID).First(&profile).Error
		if err != nil && err == gorm.ErrRecordNotFound {
			profile = domain.Profile{
				UserId:    user.ID,
				Headline:  c.Profile.Headline,
				Bio:       c.Profile.Bio,
				Phone:     c.Profile.Phone,
				AvatarUrl: c.Profile.AvatarUrl,
				ResumeUrl: c.Profile.ResumeUrl,
			}
			if err := db.Create(&profile).Error; err != nil {
				log.Printf("Failed to create profile for %s: %v\n", c.Email, err)
				continue
			}
		}

		// Associate Skills
		for _, skillName := range c.Skills {
			var skill domain.Skill
			if err := db.Where("name ILIKE ?", skillName).First(&skill).Error; err == nil {
				ps := domain.ProfileSkill{
					ProfileId: profile.ID,
					SkillId:   skill.ID,
				}
				_ = db.Clauses(clause.OnConflict{DoNothing: true}).Create(&ps).Error
			}
		}

		// Create Experiences
		for _, exp := range c.Experiences {
			var existingExp domain.Experience
			err := db.Where("profile_id = ? AND company_name = ? AND position = ?", profile.ID, exp.CompanyName, exp.Position).First(&existingExp).Error
			if err != nil && err == gorm.ErrRecordNotFound {
				startDate, _ := time.Parse("2006-01-02", exp.StartDate)
				var endDate *time.Time
				if exp.EndDate != nil && *exp.EndDate != "" {
					parsedEnd, err := time.Parse("2006-01-02", *exp.EndDate)
					if err == nil {
						endDate = &parsedEnd
					}
				}

				experience := domain.Experience{
					ID:          uuid.New(),
					ProfileId:   profile.ID,
					CompanyName: exp.CompanyName,
					Position:    exp.Position,
					StartDate:   startDate,
					EndDate:     endDate,
					IsCurrent:   exp.IsCurrent,
					Description: exp.Description,
				}
				_ = db.Create(&experience).Error
			}
		}
	}

	fmt.Printf("✅ [Candidates] Successfully processed %d candidates & profiles (New: %d)\n", len(candidates), createdCount)
}

// SeedSkils is kept for backwards compatibility
func SeedSkils() {
	env := config.LoadEnv()
	db := database.OpenConnection(env)
	SeedSkills(db)
}

func main() {
	fmt.Println("🌱 [Seeder] Starting database seeding process...")
	env := config.LoadEnv()
	db := database.OpenConnection(env)

	SeedSkills(db)
	SeedRecruiters(db)
	SeedJobs(db)
	SeedCandidates(db)

	fmt.Println("🚀 [Seeder] Seeding process completed successfully!")
}
