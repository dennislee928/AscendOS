package config

import "os"

// Config contains Chronos runtime settings loaded from environment variables.
type Config struct {
	ServiceName string
	HTTPAddr    string
}

// Load reads environment variables with safe defaults.
func Load() Config {
	return Config{
		ServiceName: getenv("CHRONOS_SERVICE_NAME", "chronos"),
		HTTPAddr:    getenv("CHRONOS_HTTP_ADDR", ":8080"),
	}
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
