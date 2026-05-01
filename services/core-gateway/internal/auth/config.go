package auth

import (
	"os"
	"strings"
	"time"
)

const (
	envJWTSecret     = "CORE_GATEWAY_JWT_SECRET"
	envJWTIssuer     = "CORE_GATEWAY_JWT_ISSUER"
	envJWTAudience   = "CORE_GATEWAY_JWT_AUDIENCE"
	envJWTClockSkew  = "CORE_GATEWAY_JWT_CLOCK_SKEW"
	defaultJWTSecret = "core-gateway-dev-secret"
)

// Config controls local JWT verification for the gateway.
type Config struct {
	Secret            string
	Issuer            string
	Audiences         []string
	ClockSkew         time.Duration
	AllowedAlgorithms []string
}

// DefaultConfig returns a development-safe verification config with HS256 enabled.
func DefaultConfig() Config {
	return Config{
		Secret:            defaultJWTSecret,
		ClockSkew:         2 * time.Minute,
		AllowedAlgorithms: []string{"HS256"},
	}
}

// LoadConfigFromEnv loads verification settings from process environment.
func LoadConfigFromEnv() Config {
	cfg := DefaultConfig()

	if secret := strings.TrimSpace(os.Getenv(envJWTSecret)); secret != "" {
		cfg.Secret = secret
	}
	if issuer := strings.TrimSpace(os.Getenv(envJWTIssuer)); issuer != "" {
		cfg.Issuer = issuer
	}
	if audiences := splitCSV(os.Getenv(envJWTAudience)); len(audiences) > 0 {
		cfg.Audiences = audiences
	}
	if skew := strings.TrimSpace(os.Getenv(envJWTClockSkew)); skew != "" {
		if parsed, err := time.ParseDuration(skew); err == nil && parsed >= 0 {
			cfg.ClockSkew = parsed
		}
	}

	return cfg
}

func splitCSV(raw string) []string {
	fields := strings.FieldsFunc(raw, func(r rune) bool {
		switch r {
		case ',', ';', ' ':
			return true
		default:
			return false
		}
	})
	values := make([]string, 0, len(fields))
	for _, field := range fields {
		if trimmed := strings.TrimSpace(field); trimmed != "" {
			values = append(values, trimmed)
		}
	}
	return values
}
