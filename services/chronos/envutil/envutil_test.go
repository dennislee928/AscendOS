package envutil

import (
	"testing"
	"time"
)

func TestStringReturnsFallbackWhenUnset(t *testing.T) {
	t.Setenv("CHRONOS_TEST_ENVUTIL", "")

	got := String("CHRONOS_TEST_ENVUTIL", "fallback")

	if got != "fallback" {
		t.Fatalf("String() = %q, want %q", got, "fallback")
	}
}

func TestStringReturnsEnvironmentValue(t *testing.T) {
	t.Setenv("CHRONOS_TEST_ENVUTIL", "present")

	got := String("CHRONOS_TEST_ENVUTIL", "fallback")

	if got != "present" {
		t.Fatalf("String() = %q, want %q", got, "present")
	}
}

func TestDurationReturnsFallbackWhenUnset(t *testing.T) {
	t.Setenv("CHRONOS_TEST_DURATION", "")

	got := Duration("CHRONOS_TEST_DURATION", 5*time.Second)

	if got != 5*time.Second {
		t.Fatalf("Duration() = %v, want %v", got, 5*time.Second)
	}
}

func TestDurationReturnsParsedValue(t *testing.T) {
	t.Setenv("CHRONOS_TEST_DURATION", "750ms")

	got := Duration("CHRONOS_TEST_DURATION", 5*time.Second)

	if got != 750*time.Millisecond {
		t.Fatalf("Duration() = %v, want %v", got, 750*time.Millisecond)
	}
}

func TestDurationFallsBackOnInvalidValue(t *testing.T) {
	t.Setenv("CHRONOS_TEST_DURATION", "not-a-duration")

	got := Duration("CHRONOS_TEST_DURATION", 5*time.Second)

	if got != 5*time.Second {
		t.Fatalf("Duration() = %v, want %v", got, 5*time.Second)
	}
}
