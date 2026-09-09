package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Skill struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey"`
	Name      string    `gorm:"column:name;type:varchar;size:100;unique;not null"`
	Label     string    `gorm:"column:label;type:varchar;size:100;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateDate;<-:create"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateDate;autoUpdateDate"`
	Profiles  []Profile `gorm:"many2many:profile_skills;foreignKey:ID;joinForeignKey:SkillId;references:ID;joinReferences:ProfileId"`
}

func (s *Skill) TableName() string {
	return "skills"
}

func (s *Skill) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		uuidV7, _ := uuid.NewV7()
		s.ID = uuidV7
	}
	return nil
}
