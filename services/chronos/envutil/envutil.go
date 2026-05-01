package envutil

import (
	"os"
	"time"
)

// String returns the environment value for key, or fallback when unset.
func String(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

// Duration returns the parsed duration for key, or fallback when unset or invalid.
func Duration(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return parsed
}
