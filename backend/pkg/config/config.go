package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Load reads .env file and loads variables into environment.
// If .env doesn't exist, it silently continues (useful for production where env vars are set externally).
func Load() error {
	if err := godotenv.Load(); err != nil {
		// In production, .env might not exist and env vars are set via system/k8s
		// So we don't fail here, just log if needed
		return nil
	}
	return nil
}

// MustLoad panics if .env file cannot be loaded (useful for local dev).
func MustLoad() {
	if err := godotenv.Load(); err != nil {
		panic(fmt.Sprintf("failed to load .env: %v", err))
	}
}

// Get returns environment variable value, with optional default.
func Get(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

// MustGet returns environment variable or panics if missing.
func MustGet(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("required env var %s is not set", key))
	}
	return v
}

// GetInt returns environment variable as int, with optional default.
func GetInt(key string, defaultValue int) int {
	if v := os.Getenv(key); v != "" {
		var result int
		if _, err := fmt.Sscanf(v, "%d", &result); err == nil {
			return result
		}
	}
	return defaultValue
}

// GetBool returns environment variable as bool, with optional default.
func GetBool(key string, defaultValue bool) bool {
	if v := os.Getenv(key); v != "" {
		switch v {
		case "true", "1", "yes", "on":
			return true
		case "false", "0", "no", "off":
			return false
		}
	}
	return defaultValue
}
