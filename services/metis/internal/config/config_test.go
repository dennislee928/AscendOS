package config

import "testing"

func TestLoadUsesDefaultsWhenEnvUnset(t *testing.T) {
	t.Setenv("METIS_SERVICE_NAME", "")
	t.Setenv("METIS_HTTP_ADDR", "")

	cfg := Load()

	if cfg.ServiceName != "metis" {
		t.Fatalf("ServiceName = %q, want %q", cfg.ServiceName, "metis")
	}
	if cfg.HTTPAddr != ":8081" {
		t.Fatalf("HTTPAddr = %q, want %q", cfg.HTTPAddr, ":8081")
	}
}

func TestLoadUsesEnvironmentOverrides(t *testing.T) {
	t.Setenv("METIS_SERVICE_NAME", "metis-test")
	t.Setenv("METIS_HTTP_ADDR", "127.0.0.1:9191")

	cfg := Load()

	if cfg.ServiceName != "metis-test" {
		t.Fatalf("ServiceName = %q, want %q", cfg.ServiceName, "metis-test")
	}
	if cfg.HTTPAddr != "127.0.0.1:9191" {
		t.Fatalf("HTTPAddr = %q, want %q", cfg.HTTPAddr, "127.0.0.1:9191")
	}
}
