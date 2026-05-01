package config

import (
	"time"

	"chronos/envutil"
)

// Config contains Metis runtime settings loaded from environment variables.
type Config struct {
	ServiceName       string
	HTTPAddr          string
	ReadHeaderTimeout time.Duration
	DataDir           string
}

// Load reads environment variables with safe defaults.
func Load() Config {
	return Config{
		ServiceName:       envutil.String("METIS_SERVICE_NAME", "metis"),
		HTTPAddr:          envutil.String("METIS_HTTP_ADDR", ":8081"),
		ReadHeaderTimeout: envutil.Duration("METIS_READ_HEADER_TIMEOUT", 5*time.Second),
		DataDir:           envutil.String("METIS_DATA_DIR", "./data"),
	}
}
