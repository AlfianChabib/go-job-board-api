package service

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"AlfianChabib/go-job-board-api/internal/model/web"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDataService_GetSkills(t *testing.T) {
	ctx := context.Background()

	t.Run("Repository error returns error", func(t *testing.T) {
		skillRepo := new(mockSkillRepository)
		req := web.GetDataRequest{Search: "go", Limit: 5}
		skillRepo.On("FindAll", ctx, req).Return(nil, errors.New("db error"))

		svc := NewDataService(skillRepo)
		res, err := svc.GetSkills(ctx, req)
		assert.Nil(t, res)
		assert.EqualError(t, err, "db error")
	})

	t.Run("Success returns mapped SkillResponse slice", func(t *testing.T) {
		skillRepo := new(mockSkillRepository)
		req := web.GetDataRequest{Search: "go", Limit: 5}
		skillID := uuid.New()
		skillRepo.On("FindAll", ctx, req).Return([]domain.Skill{
			{
				ID:           skillID,
				Name:         "golang",
				Label:        "Go (Programming Language)",
				Abbreviation: "Go",
			},
		}, nil)

		svc := NewDataService(skillRepo)
		res, err := svc.GetSkills(ctx, req)
		require.NoError(t, err)
		require.Len(t, res, 1)
		assert.Equal(t, skillID, res[0].ID)
		assert.Equal(t, "golang", res[0].Name)
		assert.Equal(t, "Go (Programming Language)", res[0].Label)
		assert.Equal(t, "Go", res[0].Abbreviation)
	})
}

func TestDataService_GetCurrencyCodes(t *testing.T) {
	ctx := context.Background()
	svc := NewDataService(new(mockSkillRepository))

	t.Run("Default limit returns 10 currencies when Limit <= 0", func(t *testing.T) {
		res, err := svc.GetCurrencyCodes(ctx, web.GetDataRequest{Limit: 0})
		require.NoError(t, err)
		assert.Len(t, res, 10)
	})

	t.Run("Search by alphabetic code", func(t *testing.T) {
		res, err := svc.GetCurrencyCodes(ctx, web.GetDataRequest{Search: "IDR", Limit: 5})
		require.NoError(t, err)
		require.NotEmpty(t, res)
		found := false
		for _, c := range res {
			if c.AlphabeticCode == "IDR" {
				found = true
				break
			}
		}
		assert.True(t, found)
	})

	t.Run("Search with no matches returns non-nil empty slice", func(t *testing.T) {
		res, err := svc.GetCurrencyCodes(ctx, web.GetDataRequest{Search: "ZZZZ_NON_EXISTENT_CURRENCY", Limit: 5})
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.Empty(t, res)
	})
}
