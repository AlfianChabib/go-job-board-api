package helper_test

import (
	"AlfianChabib/go-job-board-api/internal/helper"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateFilename(t *testing.T) {
	t.Run("Valid filenames", func(t *testing.T) {
		validNames := []string{
			"resume.pdf",
			"my_avatar-2026.png",
			"CompanyLogo123.JPG",
			"a.b.c.txt",
		}
		for _, name := range validNames {
			err := helper.ValidateFilename(name)
			assert.NoError(t, err, "expected %s to be valid", name)
		}
	})

	t.Run("Empty filename", func(t *testing.T) {
		err := helper.ValidateFilename("")
		assert.EqualError(t, err, "invalid filename: filename cannot be empty")
	})

	t.Run("Filename exceeds 255 characters", func(t *testing.T) {
		longName := strings.Repeat("a", 252) + ".png" // 256 chars
		err := helper.ValidateFilename(longName)
		assert.EqualError(t, err, "invalid filename: exceeds maximum length of 255 characters")
	})

	t.Run("Path traversal filenames", func(t *testing.T) {
		traversalNames := []string{
			"../secret.txt",
			"..hidden",
			"folder/file.png",
			"folder\\file.png",
		}
		for _, name := range traversalNames {
			err := helper.ValidateFilename(name)
			assert.EqualError(t, err, "invalid filename: path traversal is not allowed", "expected path traversal error for %s", name)
		}
	})

	t.Run("Hidden filenames starting with dot", func(t *testing.T) {
		err := helper.ValidateFilename(".env")
		assert.EqualError(t, err, "invalid filename: hidden files not allowed")
	})

	t.Run("Invalid special characters in filename", func(t *testing.T) {
		invalidNames := []string{
			"my file.png",
			"avatar@2x.png",
			"resume$1.pdf",
			"file#1.jpg",
		}
		for _, name := range invalidNames {
			err := helper.ValidateFilename(name)
			assert.EqualError(t, err, "invalid filename: contains invalid characters", "expected invalid character error for %s", name)
		}
	})
}
