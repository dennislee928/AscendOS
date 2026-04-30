package config

import "os"

// Config contains Praxis runtime settings loaded from environment variables.
type Config struct {
	ServiceName string
	HTTPAddr    string
}

// Load reads environment variables with safe defaults.
func Load() Config {
	return Config{
		ServiceName: getenv("PRAXIS_SERVICE_NAME", "praxis"),
		HTTPAddr:    getenv("PRAXIS_HTTP_ADDR", ":8082"),
	}
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
