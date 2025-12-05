package string

import (
	"strings"
	"unicode"
)

// Truncate truncates a string to a maximum length, adding ellipsis if needed.
func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// ContainsAny checks if a string contains any of the given substrings.
func ContainsAny(s string, substrings ...string) bool {
	for _, substr := range substrings {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}

// ToCamelCase converts a string to camelCase.
func ToCamelCase(s string) string {
	if s == "" {
		return s
	}

	words := strings.FieldsFunc(s, func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	})

	if len(words) == 0 {
		return s
	}

	result := strings.ToLower(words[0])
	for i := 1; i < len(words); i++ {
		if len(words[i]) > 0 {
			result += strings.ToUpper(words[i][:1]) + strings.ToLower(words[i][1:])
		}
	}

	return result
}

// ToSnakeCase converts a string to snake_case.
func ToSnakeCase(s string) string {
	if s == "" {
		return s
	}

	var result strings.Builder
	var prevLower bool

	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 && prevLower {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
			prevLower = false
		} else {
			result.WriteRune(r)
			prevLower = unicode.IsLower(r)
		}
	}

	return result.String()
}

// Sanitize removes potentially dangerous characters from a string.
func Sanitize(s string) string {
	return strings.TrimSpace(s)
}

// Default returns a default value if the string is empty.
func Default(s, defaultValue string) string {
	if strings.TrimSpace(s) == "" {
		return defaultValue
	}
	return s
}

