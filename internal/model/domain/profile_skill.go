package domain

import (
	"time"

	"github.com/google/uuid"
)

type ProfileSkill struct {
	ProfileId uuid.UUID `gorm:"column:profile_id;type:uuid;primaryKey"`
	SkillId   uuid.UUID `gorm:"column:skill_id;type:uuid;primaryKey"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateDate;<-:create"`
	Profile   *Profile  `gorm:"foreignKey:ProfileId;references:ID"`
	Skill     *Skill    `gorm:"foreignKey:SkillId;references:ID"`
}

func (ps *ProfileSkill) TableName() string {
	return "profile_skills"
}
