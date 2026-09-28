package exeption_test

import (
	"AlfianChabib/go-job-board-api/internal/exeption"
	"errors"
	"testing"
)

func TestResonseError(t *testing.T) {
	exeption.ResonseError(nil)
	exeption.ResonseError(errors.New("test error"))
}
