package domain

import (
	"bytes"
	"encoding/json"
	"errors"
	"math"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

// Card represents a spaced-repetition item managed by Metis.
type Card struct {
	ID             string    `json:"id"`
	Front          string    `json:"front"`
	Back           string    `json:"back"`
	DueAt          time.Time `json:"due_at"`
	IntervalDays   int       `json:"interval_days"`
	Repetitions    int       `json:"repetitions"`
	EaseFactor     float64   `json:"ease_factor"`
	LastReviewedAt time.Time `json:"last_reviewed_at,omitempty"`
}

// ReviewResult reports the next review schedule for a card.
type ReviewResult struct {
	Card             Card      `json:"card"`
	Quality          int       `json:"quality"`
	PreviousInterval int       `json:"previous_interval_days"`
	NextIntervalDays int       `json:"next_interval_days"`
	NextDueAt        time.Time `json:"next_due_at"`
}

// CardStore keeps a file-backed card collection for the scheduler endpoints.
type CardStore struct {
	mu       sync.Mutex
	dataPath string
	cards    map[string]Card
}

// NewCardStore loads cards from the configured data directory.
func NewCardStore(dataDir string) (*CardStore, error) {
	if dataDir == "" {
		return nil, errors.New("data dir is required")
	}

	store := &CardStore{
		dataPath: filepath.Join(dataDir, "cards.json"),
		cards:    make(map[string]Card),
	}
	if err := store.load(); err != nil {
		return nil, err
	}
	return store, nil
}

// Upsert stores or replaces a card definition.
func (s *CardStore) Upsert(card Card) (Card, error) {
	if card.ID == "" {
		return Card{}, errors.New("card id is required")
	}
	if card.EaseFactor == 0 {
		card.EaseFactor = 2.5
	}
	if card.EaseFactor < 1.3 {
		card.EaseFactor = 1.3
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	next := cloneCards(s.cards)
	next[card.ID] = card
	if err := s.persistLocked(next); err != nil {
		return Card{}, err
	}
	s.cards = next
	return card, nil
}

// Review applies an SM-2 style review to a card and stores the new state.
func (s *CardStore) Review(id string, quality int, at time.Time) (ReviewResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	card, ok := s.cards[id]
	if !ok {
		return ReviewResult{}, errors.New("card not found")
	}

	updated, result, err := applyReview(card, quality, at)
	if err != nil {
		return ReviewResult{}, err
	}

	next := cloneCards(s.cards)
	next[id] = updated
	if err := s.persistLocked(next); err != nil {
		return ReviewResult{}, err
	}
	s.cards = next
	result.Card = updated
	return result, nil
}

// Due returns the cards that are due at or before the provided time.
func (s *CardStore) Due(at time.Time) []Card {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]Card, 0, len(s.cards))
	for _, card := range s.cards {
		if !card.DueAt.After(at) {
			out = append(out, card)
		}
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].DueAt.Equal(out[j].DueAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].DueAt.Before(out[j].DueAt)
	})
	return out
}

func applyReview(card Card, quality int, at time.Time) (Card, ReviewResult, error) {
	if quality < 0 || quality > 5 {
		return Card{}, ReviewResult{}, errors.New("quality must be between 0 and 5")
	}
	if card.EaseFactor == 0 {
		card.EaseFactor = 2.5
	}

	previous := card.IntervalDays
	if previous <= 0 {
		previous = 1
	}

	ease := card.EaseFactor
	if quality < 3 {
		card.Repetitions = 0
		card.IntervalDays = 1
		ease = math.Max(1.3, ease-0.2)
	} else {
		card.Repetitions++
		switch card.Repetitions {
		case 1:
			card.IntervalDays = 1
		case 2:
			card.IntervalDays = 6
		default:
			card.IntervalDays = int(math.Round(float64(previous) * ease))
			if card.IntervalDays < 1 {
				card.IntervalDays = 1
			}
		}
		delta := 0.1 - float64(5-quality)*(0.08+float64(5-quality)*0.02)
		ease = math.Max(1.3, ease+delta)
	}

	card.EaseFactor = math.Round(ease*100) / 100
	card.LastReviewedAt = at.UTC()
	card.DueAt = card.LastReviewedAt.Add(time.Duration(card.IntervalDays) * 24 * time.Hour)

	result := ReviewResult{
		Quality:          quality,
		PreviousInterval: previous,
		NextIntervalDays: card.IntervalDays,
		NextDueAt:        card.DueAt,
	}
	return card, result, nil
}

func cloneCards(cards map[string]Card) map[string]Card {
	next := make(map[string]Card, len(cards))
	for id, card := range cards {
		next[id] = card
	}
	return next
}

func (s *CardStore) load() error {
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
		Cards map[string]Card `json:"cards"`
	}
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return err
	}
	if snapshot.Cards == nil {
		snapshot.Cards = make(map[string]Card)
	}
	s.cards = cloneCards(snapshot.Cards)
	return nil
}

func (s *CardStore) persistLocked(cards map[string]Card) error {
	if err := os.MkdirAll(filepath.Dir(s.dataPath), 0o755); err != nil {
		return err
	}

	payload, err := json.Marshal(struct {
		Cards map[string]Card `json:"cards"`
	}{
		Cards: cards,
	})
	if err != nil {
		return err
	}

	return os.WriteFile(s.dataPath, payload, 0o644)
}
