package config

import (
	"testing"
	"time"
)

func TestLoadUsesDefaultsWhenEnvUnset(t *testing.T) {
	t.Setenv("CHRONOS_SERVICE_NAME", "")
	t.Setenv("CHRONOS_HTTP_ADDR", "")
	t.Setenv("CHRONOS_DATA_DIR", "")

	cfg := Load()

	if cfg.ServiceName != "chronos" {
		t.Fatalf("ServiceName = %q, want %q", cfg.ServiceName, "chronos")
	}
	if cfg.HTTPAddr != ":8080" {
		t.Fatalf("HTTPAddr = %q, want %q", cfg.HTTPAddr, ":8080")
	}
	if cfg.ReadHeaderTimeout != 5*time.Second {
		t.Fatalf("ReadHeaderTimeout = %v, want %v", cfg.ReadHeaderTimeout, 5*time.Second)
	}
	if cfg.DataDir != "./data" {
		t.Fatalf("DataDir = %q, want %q", cfg.DataDir, "./data")
	}
}

func TestLoadUsesEnvironmentOverrides(t *testing.T) {
	t.Setenv("CHRONOS_SERVICE_NAME", "chronos-test")
	t.Setenv("CHRONOS_HTTP_ADDR", "127.0.0.1:9090")
	t.Setenv("CHRONOS_READ_HEADER_TIMEOUT", "2s")
	t.Setenv("CHRONOS_DATA_DIR", "/tmp/chronos-data")

	cfg := Load()

	if cfg.ServiceName != "chronos-test" {
		t.Fatalf("ServiceName = %q, want %q", cfg.ServiceName, "chronos-test")
	}
	if cfg.HTTPAddr != "127.0.0.1:9090" {
		t.Fatalf("HTTPAddr = %q, want %q", cfg.HTTPAddr, "127.0.0.1:9090")
	}
	if cfg.ReadHeaderTimeout != 2*time.Second {
		t.Fatalf("ReadHeaderTimeout = %v, want %v", cfg.ReadHeaderTimeout, 2*time.Second)
	}
	if cfg.DataDir != "/tmp/chronos-data" {
		t.Fatalf("DataDir = %q, want %q", cfg.DataDir, "/tmp/chronos-data")
	}
}
