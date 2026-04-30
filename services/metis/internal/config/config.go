package config

import "os"

// Config contains Metis runtime settings loaded from environment variables.
type Config struct {
	ServiceName string
	HTTPAddr    string
}

// Load reads environment variables with safe defaults.
func Load() Config {
	return Config{
		ServiceName: getenv("METIS_SERVICE_NAME", "metis"),
		HTTPAddr:    getenv("METIS_HTTP_ADDR", ":8081"),
	}
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
