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
	store, err := NewStreakStore(t.TempDir())
	if err != nil {
		t.Fatalf("NewStreakStore() error = %v", err)
	}
	if _, err := store.Forecast("missing", time.Now().UTC()); err == nil {
		t.Fatal("Forecast() error = nil, want error")
	}
}

func TestStreakStorePersistsAcrossInstances(t *testing.T) {
	dataDir := t.TempDir()

	store, err := NewStreakStore(dataDir)
	if err != nil {
		t.Fatalf("NewStreakStore() error = %v", err)
	}
	streak := HabitStreak{
		ID:                "streak-1",
		HabitID:           "habit-1",
		CurrentStreakDays: 5,
		LongestStreakDays: 9,
		MissedDaysLast7:   1,
		LastCompletedAt:   time.Date(2026, 4, 30, 9, 0, 0, 0, time.UTC),
	}
	if _, err := store.Upsert(streak); err != nil {
		t.Fatalf("Upsert() error = %v", err)
	}

	reloaded, err := NewStreakStore(dataDir)
	if err != nil {
		t.Fatalf("NewStreakStore() reload error = %v", err)
	}

	got, err := reloaded.Forecast("habit-1", time.Date(2026, 5, 1, 9, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("Forecast() error = %v", err)
	}
	if got.HabitID != "habit-1" {
		t.Fatalf("HabitID = %q, want %q", got.HabitID, "habit-1")
	}
}
