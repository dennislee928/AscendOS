package domain

import (
	"testing"
	"time"
)

func TestApplyReviewUsesSM2StyleIntervals(t *testing.T) {
	card := Card{
		ID:           "card-1",
		IntervalDays: 6,
		EaseFactor:   2.5,
		Repetitions:  2,
		DueAt:        time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC),
	}

	updated, result, err := applyReview(card, 4, time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("applyReview() error = %v", err)
	}

	if updated.Repetitions != 3 {
		t.Fatalf("Repetitions = %d, want %d", updated.Repetitions, 3)
	}
	if result.NextIntervalDays != 15 {
		t.Fatalf("NextIntervalDays = %d, want %d", result.NextIntervalDays, 15)
	}
	if updated.DueAt.Before(time.Date(2026, 5, 16, 12, 0, 0, 0, time.UTC)) {
		t.Fatalf("DueAt = %v, want at or after %v", updated.DueAt, time.Date(2026, 5, 16, 12, 0, 0, 0, time.UTC))
	}
}

func TestCardStoreDueReturnsSortedCards(t *testing.T) {
	store := NewCardStore()
	_, _ = store.Upsert(Card{ID: "b", DueAt: time.Date(2026, 5, 3, 0, 0, 0, 0, time.UTC), EaseFactor: 2.5})
	_, _ = store.Upsert(Card{ID: "a", DueAt: time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC), EaseFactor: 2.5})

	got := store.Due(time.Date(2026, 5, 4, 0, 0, 0, 0, time.UTC))
	if len(got) != 2 {
		t.Fatalf("len(Due()) = %d, want %d", len(got), 2)
	}
	if got[0].ID != "a" || got[1].ID != "b" {
		t.Fatalf("Due order = %q, %q; want a then b", got[0].ID, got[1].ID)
	}
}
