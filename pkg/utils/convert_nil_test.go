package utils_test

import (
	"AlfianChabib/go-job-board-api/pkg/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNilIfEmpty(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		emptyStr := ""
		valStr := "hello"

		assert.Nil(t, utils.NilIfEmpty(&emptyStr))
		assert.Equal(t, &valStr, utils.NilIfEmpty(&valStr))
		assert.Nil(t, utils.NilIfEmpty[string](nil))
	})

	t.Run("int / int64", func(t *testing.T) {
		zeroInt := 0
		valInt := 100
		var zeroInt64 int64 = 0
		var valInt64 int64 = 5000000

		assert.Nil(t, utils.NilIfEmpty(&zeroInt))
		assert.Equal(t, &valInt, utils.NilIfEmpty(&valInt))

		assert.Nil(t, utils.NilIfEmpty(&zeroInt64))
		assert.Equal(t, &valInt64, utils.NilIfEmpty(&valInt64))
	})

	t.Run("uuid.UUID", func(t *testing.T) {
		zeroUUID := uuid.Nil
		validUUID := uuid.New()

		assert.Nil(t, utils.NilIfEmpty(&zeroUUID))
		assert.Equal(t, &validUUID, utils.NilIfEmpty(&validUUID))
	})
}

func TestCleanseEmpty(t *testing.T) {
	emptyStr := ""
	pStr := &emptyStr
	utils.CleanseEmpty(&pStr)
	assert.Nil(t, pStr)

	valStr := "active"
	pValStr := &valStr
	utils.CleanseEmpty(&pValStr)
	assert.NotNil(t, pValStr)
	assert.Equal(t, "active", *pValStr)
}
