package validation

import (
	"fmt"
	"strings"
)

// IsEmpty checks if a string is empty or contains only whitespace.
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// IsValidEmail performs basic email validation.
func IsValidEmail(email string) bool {
	if IsEmpty(email) {
		return false
	}
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// IsValidUUID checks if a string is a valid UUID format.
func IsValidUUID(uuid string) bool {
	if len(uuid) != 36 {
		return false
	}
	parts := strings.Split(uuid, "-")
	if len(parts) != 5 {
		return false
	}
	// Basic format check: 8-4-4-4-12
	expectedLengths := []int{8, 4, 4, 4, 12}
	for i, part := range parts {
		if len(part) != expectedLengths[i] {
			return false
		}
	}
	return true
}

// ValidateRange checks if a number is within a range.
func ValidateRange(value, min, max int) error {
	if value < min || value > max {
		return fmt.Errorf("value %d is out of range [%d, %d]", value, min, max)
	}
	return nil
}

// ValidatePositive checks if a number is positive.
func ValidatePositive(value int) error {
	if value <= 0 {
		return fmt.Errorf("value must be positive, got %d", value)
	}
	return nil
}

// ValidateNonEmpty checks if a slice is not empty.
func ValidateNonEmpty[T any](slice []T, fieldName string) error {
	if len(slice) == 0 {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}
	return nil
}
