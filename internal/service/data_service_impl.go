package service

import (
	assets "AlfianChabib/go-job-board-api"
	"AlfianChabib/go-job-board-api/internal/model/web"
	"AlfianChabib/go-job-board-api/internal/repository"
	"context"
	"encoding/csv"
	"io"
	"strings"
	"sync"
)

type dataService struct {
	skillRepository repository.SkillRepository
	currencyCache   []web.CurrencyResponse
	once            sync.Once
}

func NewDataService(skillRepository repository.SkillRepository) DataService {
	s := &dataService{
		skillRepository: skillRepository,
	}
	s.loadCurrencies()
	return s
}

func (service *dataService) loadCurrencies() {
	service.once.Do(func() {
		reader := csv.NewReader(strings.NewReader(assets.CurrencyCodesCSV))
		// skip header
		_, err := reader.Read()
		if err != nil {
			return
		}

		var currencies []web.CurrencyResponse
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				continue
			}

			// Format: Entity,Currency,AlphabeticCode,NumericCode,MinorUnit,WithdrawalDate
			if len(record) < 3 {
				continue
			}

			code := strings.TrimSpace(record[2])
			if code == "" {
				continue
			}

			c := web.CurrencyResponse{
				Entity:         strings.TrimSpace(record[0]),
				Currency:       strings.TrimSpace(record[1]),
				AlphabeticCode: code,
			}

			if len(record) > 3 {
				c.NumericCode = strings.TrimSpace(record[3])
			}
			if len(record) > 4 {
				c.MinorUnit = strings.TrimSpace(record[4])
			}
			if len(record) > 5 {
				c.WithdrawalDate = strings.TrimSpace(record[5])
			}

			currencies = append(currencies, c)
		}
		service.currencyCache = currencies
	})
}

func (service *dataService) GetSkills(ctx context.Context) ([]web.SkillResponse, error) {
	skills, err := service.skillRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	skillResponses := make([]web.SkillResponse, 0, len(skills))
	for _, s := range skills {
		skillResponses = append(skillResponses, web.SkillResponse{
			ID:    s.ID,
			Name:  s.Name,
			Label: s.Label,
		})
	}

	return skillResponses, nil
}

func (service *dataService) GetCurrencyCodes(ctx context.Context) ([]web.CurrencyResponse, error) {
	return service.currencyCache, nil
}
