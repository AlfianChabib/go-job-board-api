package service

import (
	"AlfianChabib/go-job-board-api/internal/model/web"
	"context"
)

type DataService interface {
	GetSkills(ctx context.Context, req web.GetDataRequest) ([]web.SkillResponse, error)
	GetCurrencyCodes(ctx context.Context, req web.GetDataRequest) ([]web.CurrencyResponse, error)
}
