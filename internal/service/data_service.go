package service

import (
	"AlfianChabib/go-job-board-api/internal/model/web"
	"context"
)

type DataService interface {
	GetSkills(ctx context.Context) ([]web.SkillResponse, error)
	GetCurrencyCodes(ctx context.Context) ([]web.CurrencyResponse, error)
}
