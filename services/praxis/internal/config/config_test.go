package config

import "testing"

func TestLoadUsesDefaultsWhenEnvUnset(t *testing.T) {
	t.Setenv("PRAXIS_SERVICE_NAME", "")
	t.Setenv("PRAXIS_HTTP_ADDR", "")

	cfg := Load()

	if cfg.ServiceName != "praxis" {
		t.Fatalf("ServiceName = %q, want %q", cfg.ServiceName, "praxis")
	}
	if cfg.HTTPAddr != ":8082" {
		t.Fatalf("HTTPAddr = %q, want %q", cfg.HTTPAddr, ":8082")
	}
}

func TestLoadUsesEnvironmentOverrides(t *testing.T) {
	t.Setenv("PRAXIS_SERVICE_NAME", "praxis-test")
	t.Setenv("PRAXIS_HTTP_ADDR", "127.0.0.1:9292")

	cfg := Load()

	if cfg.ServiceName != "praxis-test" {
		t.Fatalf("ServiceName = %q, want %q", cfg.ServiceName, "praxis-test")
	}
	if cfg.HTTPAddr != "127.0.0.1:9292" {
		t.Fatalf("HTTPAddr = %q, want %q", cfg.HTTPAddr, "127.0.0.1:9292")
	}
}
