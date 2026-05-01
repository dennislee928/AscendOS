package config

import (
	"testing"
	"time"
)

func TestLoadUsesDefaultsWhenEnvUnset(t *testing.T) {
	t.Setenv("PRAXIS_SERVICE_NAME", "")
	t.Setenv("PRAXIS_HTTP_ADDR", "")
	t.Setenv("PRAXIS_DATA_DIR", "")

	cfg := Load()

	if cfg.ServiceName != "praxis" {
		t.Fatalf("ServiceName = %q, want %q", cfg.ServiceName, "praxis")
	}
	if cfg.HTTPAddr != ":8082" {
		t.Fatalf("HTTPAddr = %q, want %q", cfg.HTTPAddr, ":8082")
	}
	if cfg.ReadHeaderTimeout != 5*time.Second {
		t.Fatalf("ReadHeaderTimeout = %v, want %v", cfg.ReadHeaderTimeout, 5*time.Second)
	}
	if cfg.DataDir != "./data" {
		t.Fatalf("DataDir = %q, want %q", cfg.DataDir, "./data")
	}
}

func TestLoadUsesEnvironmentOverrides(t *testing.T) {
	t.Setenv("PRAXIS_SERVICE_NAME", "praxis-test")
	t.Setenv("PRAXIS_HTTP_ADDR", "127.0.0.1:9292")
	t.Setenv("PRAXIS_READ_HEADER_TIMEOUT", "4s")
	t.Setenv("PRAXIS_DATA_DIR", "/tmp/praxis-data")

	cfg := Load()

	if cfg.ServiceName != "praxis-test" {
		t.Fatalf("ServiceName = %q, want %q", cfg.ServiceName, "praxis-test")
	}
	if cfg.HTTPAddr != "127.0.0.1:9292" {
		t.Fatalf("HTTPAddr = %q, want %q", cfg.HTTPAddr, "127.0.0.1:9292")
	}
	if cfg.ReadHeaderTimeout != 4*time.Second {
		t.Fatalf("ReadHeaderTimeout = %v, want %v", cfg.ReadHeaderTimeout, 4*time.Second)
	}
	if cfg.DataDir != "/tmp/praxis-data" {
		t.Fatalf("DataDir = %q, want %q", cfg.DataDir, "/tmp/praxis-data")
	}
}
