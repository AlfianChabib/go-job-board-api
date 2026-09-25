package web

import "github.com/google/uuid"

type SkillResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Label string    `json:"label"`
}

type CurrencyResponse struct {
	Entity         string `json:"entity"`
	Currency       string `json:"currency"`
	AlphabeticCode string `json:"alphabetic_code"`
	NumericCode    string `json:"numeric_code,omitempty"`
	MinorUnit      string `json:"minor_unit,omitempty"`
	WithdrawalDate string `json:"withdrawal_date,omitempty"`
}
