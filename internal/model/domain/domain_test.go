package domain_test

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm/schema"
)

func TestDomainTableNames(t *testing.T) {
	assert.Equal(t, "applications", (&domain.Application{}).TableName())
	assert.Equal(t, "auths", (&domain.Auth{}).TableName())
	assert.Equal(t, "companies", (&domain.Company{}).TableName())
	assert.Equal(t, "experiences", (&domain.Experience{}).TableName())
	assert.Equal(t, "jobs", (&domain.Job{}).TableName())
	assert.Equal(t, "profiles", (&domain.Profile{}).TableName())
	assert.Equal(t, "profile_skills", (&domain.ProfileSkill{}).TableName())
	assert.Equal(t, "skills", (&domain.Skill{}).TableName())
	assert.Equal(t, "tokens", (&domain.Token{}).TableName())
	assert.Equal(t, "users", (&domain.User{}).TableName())
}

func TestDomainBeforeCreateHooks(t *testing.T) {
	t.Run("Application BeforeCreate", func(t *testing.T) {
		app := &domain.Application{}
		err := app.BeforeCreate(nil)
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, app.ID)
		assert.Equal(t, domain.ApplicationStatusApplied, app.Status)
		assert.False(t, app.AppliedAt.IsZero())
		assert.False(t, app.StatusUpdatedAt.IsZero())

		// Preserves existing values when already set
		presetID := uuid.New()
		presetTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		app2 := &domain.Application{
			ID:              presetID,
			Status:          domain.ApplicationStatusReviewing,
			AppliedAt:       presetTime,
			StatusUpdatedAt: presetTime,
		}
		require.NoError(t, app2.BeforeCreate(nil))
		assert.Equal(t, presetID, app2.ID)
		assert.Equal(t, domain.ApplicationStatusReviewing, app2.Status)
		assert.Equal(t, presetTime, app2.AppliedAt)
		assert.Equal(t, presetTime, app2.StatusUpdatedAt)
	})

	t.Run("Auth BeforeCreate", func(t *testing.T) {
		auth := &domain.Auth{}
		require.NoError(t, auth.BeforeCreate(nil))
		assert.NotEqual(t, uuid.Nil, auth.ID)

		presetID := uuid.New()
		auth2 := &domain.Auth{ID: presetID}
		require.NoError(t, auth2.BeforeCreate(nil))
		assert.Equal(t, presetID, auth2.ID)
	})

	t.Run("Company BeforeCreate", func(t *testing.T) {
		company := &domain.Company{}
		require.NoError(t, company.BeforeCreate(nil))
		assert.NotEqual(t, uuid.Nil, company.ID)

		presetID := uuid.New()
		company2 := &domain.Company{ID: presetID}
		require.NoError(t, company2.BeforeCreate(nil))
		assert.Equal(t, presetID, company2.ID)
	})

	t.Run("Experience BeforeCreate", func(t *testing.T) {
		exp := &domain.Experience{}
		require.NoError(t, exp.BeforeCreate(nil))
		assert.NotEqual(t, uuid.Nil, exp.ID)

		presetID := uuid.New()
		exp2 := &domain.Experience{ID: presetID}
		require.NoError(t, exp2.BeforeCreate(nil))
		assert.Equal(t, presetID, exp2.ID)
	})

	t.Run("Job BeforeCreate", func(t *testing.T) {
		job := &domain.Job{}
		require.NoError(t, job.BeforeCreate(nil))
		assert.NotEqual(t, uuid.Nil, job.ID)
		assert.Equal(t, "IDR", job.Currency)
		assert.Equal(t, domain.JobStatusOpen, job.Status)

		presetID := uuid.New()
		job2 := &domain.Job{
			ID:       presetID,
			Currency: "USD",
			Status:   domain.JobStatusClosed,
		}
		require.NoError(t, job2.BeforeCreate(nil))
		assert.Equal(t, presetID, job2.ID)
		assert.Equal(t, "USD", job2.Currency)
		assert.Equal(t, domain.JobStatusClosed, job2.Status)
	})

	t.Run("Profile BeforeCreate", func(t *testing.T) {
		profile := &domain.Profile{}
		require.NoError(t, profile.BeforeCreate(nil))
		assert.NotEqual(t, uuid.Nil, profile.ID)

		presetID := uuid.New()
		profile2 := &domain.Profile{ID: presetID}
		require.NoError(t, profile2.BeforeCreate(nil))
		assert.Equal(t, presetID, profile2.ID)
	})

	t.Run("Skill BeforeCreate", func(t *testing.T) {
		skill := &domain.Skill{}
		require.NoError(t, skill.BeforeCreate(nil))
		assert.NotEqual(t, uuid.Nil, skill.ID)

		presetID := uuid.New()
		skill2 := &domain.Skill{ID: presetID}
		require.NoError(t, skill2.BeforeCreate(nil))
		assert.Equal(t, presetID, skill2.ID)
	})

	t.Run("Token BeforeCreate", func(t *testing.T) {
		tok := &domain.Token{}
		require.NoError(t, tok.BeforeCreate(nil))
		assert.NotEqual(t, uuid.Nil, tok.ID)

		presetID := uuid.New()
		tok2 := &domain.Token{ID: presetID}
		require.NoError(t, tok2.BeforeCreate(nil))
		assert.Equal(t, presetID, tok2.ID)
	})

	t.Run("User BeforeCreate", func(t *testing.T) {
		user := &domain.User{}
		require.NoError(t, user.BeforeCreate(nil))
		assert.NotEqual(t, uuid.Nil, user.ID)

		presetID := uuid.New()
		user2 := &domain.User{ID: presetID}
		require.NoError(t, user2.BeforeCreate(nil))
		assert.Equal(t, presetID, user2.ID)
	})
}

func TestAllDomainSchemas(t *testing.T) {
	models := []any{
		&domain.Application{},
		&domain.Auth{},
		&domain.Company{},
		&domain.Experience{},
		&domain.Job{},
		&domain.Profile{},
		&domain.ProfileSkill{},
		&domain.Skill{},
		&domain.Token{},
		&domain.User{},
	}

	cache := &sync.Map{}
	for _, m := range models {
		s, err := schema.Parse(m, cache, schema.NamingStrategy{})
		require.NoError(t, err)
		assert.NotNil(t, s)
	}
}
