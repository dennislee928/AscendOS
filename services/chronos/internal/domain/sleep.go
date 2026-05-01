package domain

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
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

// SleepStore keeps a file-backed history of sleep events.
type SleepStore struct {
	mu       sync.Mutex
	dataPath string
	events   []SleepEvent
}

// NewSleepStore loads sleep events from the configured data directory.
func NewSleepStore(dataDir string) (*SleepStore, error) {
	if dataDir == "" {
		return nil, errors.New("data dir is required")
	}

	store := &SleepStore{
		dataPath: filepath.Join(dataDir, "sleep-events.json"),
	}
	if err := store.load(); err != nil {
		return nil, err
	}
	return store, nil
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

	next := append([]SleepEvent(nil), s.events...)
	for i := range next {
		if next[i].ID == event.ID {
			next[i] = event
			if err := s.persistLocked(next); err != nil {
				return SleepEvent{}, err
			}
			s.events = next
			return event, nil
		}
	}

	next = append(next, event)
	if err := s.persistLocked(next); err != nil {
		return SleepEvent{}, err
	}
	s.events = next
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

func (s *SleepStore) load() error {
	data, err := os.ReadFile(s.dataPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		return nil
	}

	var snapshot struct {
		Events []SleepEvent `json:"events"`
	}
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return err
	}
	s.events = append([]SleepEvent(nil), snapshot.Events...)
	return nil
}

func (s *SleepStore) persistLocked(events []SleepEvent) error {
	if err := os.MkdirAll(filepath.Dir(s.dataPath), 0o755); err != nil {
		return err
	}

	payload, err := json.Marshal(struct {
		Events []SleepEvent `json:"events"`
	}{
		Events: events,
	})
	if err != nil {
		return err
	}

	return os.WriteFile(s.dataPath, payload, 0o644)
}
