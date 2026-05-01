package config

import (
	"time"

	"chronos/envutil"
)

// Config contains Praxis runtime settings loaded from environment variables.
type Config struct {
	ServiceName       string
	HTTPAddr          string
	ReadHeaderTimeout time.Duration
}

// Load reads environment variables with safe defaults.
func Load() Config {
	return Config{
		ServiceName:       envutil.String("PRAXIS_SERVICE_NAME", "praxis"),
		HTTPAddr:          envutil.String("PRAXIS_HTTP_ADDR", ":8082"),
		ReadHeaderTimeout: envutil.Duration("PRAXIS_READ_HEADER_TIMEOUT", 5*time.Second),
	}
}
