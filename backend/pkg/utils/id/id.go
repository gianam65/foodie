package id

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// GenerateUUID generates a new UUID string.
func GenerateUUID() string {
	return uuid.New().String()
}

// IsValidUUID checks if a string is a valid UUID.
func IsValidUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}

// GenerateShortID generates a short random ID (16 characters).
func GenerateShortID() (string, error) {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate ID: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// MustGenerateShortID generates a short random ID or panics.
func MustGenerateShortID() string {
	id, err := GenerateShortID()
	if err != nil {
		panic(err)
	}
	return id
}

// SanitizeID removes unwanted characters from an ID string.
func SanitizeID(id string) string {
	// Remove whitespace and common problematic characters
	id = strings.TrimSpace(id)
	id = strings.ReplaceAll(id, " ", "")
	id = strings.ReplaceAll(id, "\n", "")
	id = strings.ReplaceAll(id, "\r", "")
	return id
}
