package assets_test

import (
	assets "AlfianChabib/go-job-board-api"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCurrencyCodesCSVEmbedded(t *testing.T) {
	assert.NotEmpty(t, assets.CurrencyCodesCSV)
	assert.Contains(t, assets.CurrencyCodesCSV, "IDR")
	assert.Contains(t, assets.CurrencyCodesCSV, "USD")
}
