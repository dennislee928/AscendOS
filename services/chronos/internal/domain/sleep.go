package domain

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"sync"
	"time"
)

// SleepEvent records a single sleep session for Chronos analysis.
type SleepEvent struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id,omitempty"`
	StartedAt time.Time `json:"started_at"`
	EndedAt   time.Time `json:"ended_at"`
}

// CircadianPhase summarizes the current phase estimate from recent sleep data.
type CircadianPhase struct {
	Phase           string    `json:"phase"`
	Score           float64   `json:"score"`
	SampleCount     int       `json:"sample_count"`
	AverageMidpoint string    `json:"average_midpoint"`
	CalculatedAt    time.Time `json:"calculated_at"`
}

// SleepStore keeps an in-memory history of sleep events.
type SleepStore struct {
	mu     sync.Mutex
	events []SleepEvent
}

// NewSleepStore constructs an empty sleep event store.
func NewSleepStore() *SleepStore {
	return &SleepStore{}
}

// Ingest stores a validated sleep event and returns the stored copy.
func (s *SleepStore) Ingest(event SleepEvent) (SleepEvent, error) {
	if event.ID == "" {
		return SleepEvent{}, errors.New("sleep event id is required")
	}
	if event.EndedAt.IsZero() || event.StartedAt.IsZero() {
		return SleepEvent{}, errors.New("sleep event start and end times are required")
	}
	if !event.EndedAt.After(event.StartedAt) {
		return SleepEvent{}, errors.New("sleep event must end after it starts")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.events {
		if s.events[i].ID == event.ID {
			s.events[i] = event
			return event, nil
		}
	}

	s.events = append(s.events, event)
	return event, nil
}

// Events returns a snapshot copy of the stored sleep events.
func (s *SleepStore) Events() []SleepEvent {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]SleepEvent, len(s.events))
	copy(out, s.events)
	return out
}

// Phase calculates the circadian phase from the stored sleep events.
func (s *SleepStore) Phase(at time.Time) CircadianPhase {
	return CalculateCircadianPhase(s.Events(), at)
}

// CalculateCircadianPhase estimates circadian phase from sleep midpoints.
func CalculateCircadianPhase(events []SleepEvent, at time.Time) CircadianPhase {
	if len(events) == 0 {
		return CircadianPhase{
			Phase:           "insufficient-data",
			Score:           0,
			SampleCount:     0,
			AverageMidpoint: "--:--",
			CalculatedAt:    at.UTC(),
		}
	}

	type sample struct {
		minutes int
	}

	samples := make([]sample, 0, len(events))
	for _, event := range events {
		if event.EndedAt.After(event.StartedAt) {
			midpoint := event.StartedAt.Add(event.EndedAt.Sub(event.StartedAt) / 2)
			samples = append(samples, sample{minutes: minutesOfDay(midpoint)})
		}
	}
	if len(samples) == 0 {
		return CircadianPhase{
			Phase:           "insufficient-data",
			Score:           0,
			SampleCount:     0,
			AverageMidpoint: "--:--",
			CalculatedAt:    at.UTC(),
		}
	}

	sort.Slice(samples, func(i, j int) bool {
		return samples[i].minutes < samples[j].minutes
	})

	total := 0
	for _, sample := range samples {
		total += sample.minutes
	}
	average := float64(total) / float64(len(samples))
	phase := phaseFromMidpoint(average)
	score := phaseScore(average)

	return CircadianPhase{
		Phase:           phase,
		Score:           score,
		SampleCount:     len(samples),
		AverageMidpoint: formatMinutes(average),
		CalculatedAt:    at.UTC(),
	}
}

func phaseFromMidpoint(minutes float64) string {
	switch {
	case minutes < 150:
		return "advanced"
	case minutes <= 240:
		return "aligned"
	default:
		return "delayed"
	}
}

func phaseScore(minutes float64) float64 {
	const target = 180.0
	const maxDistance = 180.0
	score := 1 - math.Abs(minutes-target)/maxDistance
	if score < 0 {
		return 0
	}
	if score > 1 {
		return 1
	}
	return math.Round(score*100) / 100
}

func minutesOfDay(t time.Time) int {
	local := t.In(t.Location())
	return local.Hour()*60 + local.Minute()
}

func formatMinutes(minutes float64) string {
	total := int(math.Round(minutes))
	if total < 0 {
		total = 0
	}
	hours := (total / 60) % 24
	minute := total % 60
	return fmt.Sprintf("%02d:%02d", hours, minute)
}
