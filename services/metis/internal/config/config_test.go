package config

import (
	"testing"
	"time"
)

func TestLoadUsesDefaultsWhenEnvUnset(t *testing.T) {
	t.Setenv("METIS_SERVICE_NAME", "")
	t.Setenv("METIS_HTTP_ADDR", "")
	t.Setenv("METIS_DATA_DIR", "")

	cfg := Load()

	if cfg.ServiceName != "metis" {
		t.Fatalf("ServiceName = %q, want %q", cfg.ServiceName, "metis")
	}
	if cfg.HTTPAddr != ":8081" {
		t.Fatalf("HTTPAddr = %q, want %q", cfg.HTTPAddr, ":8081")
	}
	if cfg.ReadHeaderTimeout != 5*time.Second {
		t.Fatalf("ReadHeaderTimeout = %v, want %v", cfg.ReadHeaderTimeout, 5*time.Second)
	}
	if cfg.DataDir != "./data" {
		t.Fatalf("DataDir = %q, want %q", cfg.DataDir, "./data")
	}
}

func TestLoadUsesEnvironmentOverrides(t *testing.T) {
	t.Setenv("METIS_SERVICE_NAME", "metis-test")
	t.Setenv("METIS_HTTP_ADDR", "127.0.0.1:9191")
	t.Setenv("METIS_READ_HEADER_TIMEOUT", "3s")
	t.Setenv("METIS_DATA_DIR", "/tmp/metis-data")

	cfg := Load()

	if cfg.ServiceName != "metis-test" {
		t.Fatalf("ServiceName = %q, want %q", cfg.ServiceName, "metis-test")
	}
	if cfg.HTTPAddr != "127.0.0.1:9191" {
		t.Fatalf("HTTPAddr = %q, want %q", cfg.HTTPAddr, "127.0.0.1:9191")
	}
	if cfg.ReadHeaderTimeout != 3*time.Second {
		t.Fatalf("ReadHeaderTimeout = %v, want %v", cfg.ReadHeaderTimeout, 3*time.Second)
	}
	if cfg.DataDir != "/tmp/metis-data" {
		t.Fatalf("DataDir = %q, want %q", cfg.DataDir, "/tmp/metis-data")
	}
}
