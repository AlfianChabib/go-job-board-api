package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Profile struct {
	ID          uuid.UUID    `gorm:"column:id;type:uuid;primaryKey"`
	UserId      uuid.UUID    `gorm:"column:user_id;type:uuid;unique;not null"`
	AvatarUrl   *string      `gorm:"column:avatar_url;type:varchar;size:255"`
	Headline    *string      `gorm:"column:headline;type:varchar;size:255"`
	Bio         *string      `gorm:"column:bio;type:text"`
	Phone       *string      `gorm:"column:phone;type:varchar;size:100"`
	ResumeUrl   *string      `gorm:"column:resume_url;type:varchar;size:255"`
	CreatedAt   time.Time    `gorm:"column:created_at;autoCreateDate;<-:create"`
	UpdatedAt   time.Time    `gorm:"column:updated_at;autoCreateDate;autoUpdateDate"`
	User        *User        `gorm:"foreignKey:UserId;references:ID"`
	Skills      []Skill      `gorm:"many2many:profile_skills;foreignKey:ID;joinForeignKey:ProfileId;references:ID;joinReferences:SkillId"`
	Experiences []Experience `gorm:"foreignKey:ProfileId;references:ID"`
}

func (p *Profile) TableName() string {
	return "profiles"
}

func (p *Profile) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		uuidV7, _ := uuid.NewV7()
		p.ID = uuidV7
	}
	return nil
}
