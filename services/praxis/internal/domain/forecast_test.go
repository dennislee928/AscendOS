package domain

import (
	"testing"
	"time"
)

func TestForecastRelapseBuildsRiskScore(t *testing.T) {
	streak := HabitStreak{
		ID:                "streak-1",
		HabitID:           "habit-1",
		CurrentStreakDays: 3,
		LongestStreakDays: 12,
		MissedDaysLast7:   2,
		LastCompletedAt:   time.Date(2026, 4, 28, 9, 0, 0, 0, time.UTC),
	}

	got := ForecastRelapse(streak, time.Date(2026, 5, 1, 9, 0, 0, 0, time.UTC))

	if got.RiskLevel != "moderate" && got.RiskLevel != "high" {
		t.Fatalf("RiskLevel = %q, want moderate or high", got.RiskLevel)
	}
	if got.RiskScore <= 0 {
		t.Fatalf("RiskScore = %v, want positive", got.RiskScore)
	}
	if got.HabitID != "habit-1" {
		t.Fatalf("HabitID = %q, want %q", got.HabitID, "habit-1")
	}
}

func TestStreakStoreForecastRequiresKnownHabit(t *testing.T) {
	store := NewStreakStore()
	if _, err := store.Forecast("missing", time.Now().UTC()); err == nil {
		t.Fatal("Forecast() error = nil, want error")
	}
}
