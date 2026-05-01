package domain

import (
	"errors"
	"math"
	"sort"
	"sync"
	"time"
)

// HabitStreak describes the streak state used for relapse forecasting.
type HabitStreak struct {
	ID                string    `json:"id"`
	HabitID           string    `json:"habit_id"`
	CurrentStreakDays int       `json:"current_streak_days"`
	LongestStreakDays int       `json:"longest_streak_days"`
	MissedDaysLast7   int       `json:"missed_days_last_7"`
	LastCompletedAt   time.Time `json:"last_completed_at"`
}

// RelapseForecast summarizes the predicted relapse risk for a habit.
type RelapseForecast struct {
	HabitID       string    `json:"habit_id"`
	RiskScore     float64   `json:"risk_score"`
	RiskLevel     string    `json:"risk_level"`
	Drivers       []string  `json:"drivers"`
	CalculatedAt  time.Time `json:"calculated_at"`
	CurrentStreak int       `json:"current_streak_days"`
	LongestStreak int       `json:"longest_streak_days"`
}

// StreakStore keeps habit streaks in memory for forecast endpoints.
type StreakStore struct {
	mu      sync.Mutex
	streaks map[string]HabitStreak
}

// NewStreakStore constructs an empty streak store.
func NewStreakStore() *StreakStore {
	return &StreakStore{streaks: make(map[string]HabitStreak)}
}

// Upsert stores or replaces a habit streak definition.
func (s *StreakStore) Upsert(streak HabitStreak) (HabitStreak, error) {
	if streak.ID == "" {
		return HabitStreak{}, errors.New("streak id is required")
	}
	if streak.HabitID == "" {
		return HabitStreak{}, errors.New("habit id is required")
	}
	if streak.LongestStreakDays < streak.CurrentStreakDays {
		streak.LongestStreakDays = streak.CurrentStreakDays
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.streaks[streak.HabitID] = streak
	return streak, nil
}

// Forecast calculates relapse risk from the stored streak for a habit.
func (s *StreakStore) Forecast(habitID string, at time.Time) (RelapseForecast, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	streak, ok := s.streaks[habitID]
	if !ok {
		return RelapseForecast{}, errors.New("habit streak not found")
	}

	return ForecastRelapse(streak, at), nil
}

// ForecastRelapse converts streak state into a normalized risk forecast.
func ForecastRelapse(streak HabitStreak, at time.Time) RelapseForecast {
	score := 0.55
	if streak.CurrentStreakDays > 0 {
		score -= math.Min(float64(streak.CurrentStreakDays)/21.0, 1) * 0.25
	}
	if streak.LongestStreakDays > 0 {
		score -= math.Min(float64(streak.LongestStreakDays)/60.0, 1) * 0.1
	}
	score += math.Min(float64(streak.MissedDaysLast7), 7) * 0.08

	if !streak.LastCompletedAt.IsZero() {
		daysSince := int(at.Sub(streak.LastCompletedAt).Hours() / 24)
		if daysSince > 1 {
			score += math.Min(float64(daysSince-1)*0.06, 0.24)
		}
	}

	score = clamp(score, 0, 1)

	drivers := make([]string, 0, 3)
	if streak.CurrentStreakDays >= 7 {
		drivers = append(drivers, "current streak is stable")
	} else {
		drivers = append(drivers, "current streak is short")
	}
	if streak.MissedDaysLast7 > 0 {
		drivers = append(drivers, "recent misses increase risk")
	}
	if !streak.LastCompletedAt.IsZero() && at.Sub(streak.LastCompletedAt) > 48*time.Hour {
		drivers = append(drivers, "recent completion is stale")
	}

	riskLevel := "high"
	switch {
	case score < 0.33:
		riskLevel = "low"
	case score < 0.66:
		riskLevel = "moderate"
	}

	return RelapseForecast{
		HabitID:       streak.HabitID,
		RiskScore:     math.Round(score*100) / 100,
		RiskLevel:     riskLevel,
		Drivers:       uniqueDrivers(drivers),
		CalculatedAt:  at.UTC(),
		CurrentStreak: streak.CurrentStreakDays,
		LongestStreak: streak.LongestStreakDays,
	}
}

func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func uniqueDrivers(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	sort.Strings(values)
	out := values[:1]
	for _, value := range values[1:] {
		if value != out[len(out)-1] {
			out = append(out, value)
		}
	}
	return out
}
