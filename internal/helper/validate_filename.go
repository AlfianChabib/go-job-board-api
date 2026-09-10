package helper

import (
	"fmt"
	"strings"
)

// validateFilename checks if the filename contains only alphanumeric characters, underscores, dashes, and dots.
func ValidateFilename(filename string) error {
	// Check if the filename is empty
	if filename == "" {
		return fmt.Errorf("invalid filename: filename cannot be empty")
	}

	// Check filename length
	if len(filename) > 255 {
		return fmt.Errorf("invalid filename: exceeds maximum length of %d characters", 255)
	}

	// Prevent path traversal
	if strings.Contains(filename, "/") || strings.Contains(filename, "\\") || strings.HasPrefix(filename, "..") {
		return fmt.Errorf("invalid filename: path traversal is not allowed")
	}

	// Prevent hidden files (files starting with a dot)
	if strings.HasPrefix(filename, ".") {
		return fmt.Errorf("invalid filename: hidden files not allowed")
	}

	// Check each character in the filename for invalid characters
	for _, char := range filename {
		if !isAlphaNumericOrSpecial(char) {
			return fmt.Errorf("invalid filename: contains invalid characters")
		}
	}
	return nil
}

// isAlphaNumericOrSpecial checks if a character is alphanumeric or a valid special character.
func isAlphaNumericOrSpecial(char rune) bool {
	// Validate the character (alphanumeric, dash, underscore, or dot)
	switch {
	case 'A' <= char && char <= 'Z':
		return true
	case 'a' <= char && char <= 'z':
		return true
	case '0' <= char && char <= '9':
		return true
	case char == '-' || char == '_':
		return true
	case char == '.':
		return true
	default:
		return false
	}
}
