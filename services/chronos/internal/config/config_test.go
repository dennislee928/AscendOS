package config

import "testing"

func TestLoadUsesDefaultsWhenEnvUnset(t *testing.T) {
	t.Setenv("CHRONOS_SERVICE_NAME", "")
	t.Setenv("CHRONOS_HTTP_ADDR", "")

	cfg := Load()

	if cfg.ServiceName != "chronos" {
		t.Fatalf("ServiceName = %q, want %q", cfg.ServiceName, "chronos")
	}
	if cfg.HTTPAddr != ":8080" {
		t.Fatalf("HTTPAddr = %q, want %q", cfg.HTTPAddr, ":8080")
	}
}

func TestLoadUsesEnvironmentOverrides(t *testing.T) {
	t.Setenv("CHRONOS_SERVICE_NAME", "chronos-test")
	t.Setenv("CHRONOS_HTTP_ADDR", "127.0.0.1:9090")

	cfg := Load()

	if cfg.ServiceName != "chronos-test" {
		t.Fatalf("ServiceName = %q, want %q", cfg.ServiceName, "chronos-test")
	}
	if cfg.HTTPAddr != "127.0.0.1:9090" {
		t.Fatalf("HTTPAddr = %q, want %q", cfg.HTTPAddr, "127.0.0.1:9090")
	}
}
