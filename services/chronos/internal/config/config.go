package config

import (
	"time"

	"chronos/envutil"
)

// Config contains Chronos runtime settings loaded from environment variables.
type Config struct {
	ServiceName       string
	HTTPAddr          string
	ReadHeaderTimeout time.Duration
	DataDir           string
}

// Load reads environment variables with safe defaults.
func Load() Config {
	return Config{
		ServiceName:       envutil.String("CHRONOS_SERVICE_NAME", "chronos"),
		HTTPAddr:          envutil.String("CHRONOS_HTTP_ADDR", ":8080"),
		ReadHeaderTimeout: envutil.Duration("CHRONOS_READ_HEADER_TIMEOUT", 5*time.Second),
		DataDir:           envutil.String("CHRONOS_DATA_DIR", "./data"),
	}
}
