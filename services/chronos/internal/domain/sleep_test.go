package domain

import (
	"testing"
	"time"
)

func TestCalculateCircadianPhase(t *testing.T) {
	events := []SleepEvent{
		{
			ID:        "one",
			StartedAt: time.Date(2026, 5, 1, 23, 0, 0, 0, time.UTC),
			EndedAt:   time.Date(2026, 5, 2, 7, 0, 0, 0, time.UTC),
		},
		{
			ID:        "two",
			StartedAt: time.Date(2026, 5, 2, 22, 30, 0, 0, time.UTC),
			EndedAt:   time.Date(2026, 5, 3, 6, 30, 0, 0, time.UTC),
		},
	}

	got := CalculateCircadianPhase(events, time.Date(2026, 5, 3, 12, 0, 0, 0, time.UTC))

	if got.Phase != "aligned" {
		t.Fatalf("Phase = %q, want %q", got.Phase, "aligned")
	}
	if got.SampleCount != 2 {
		t.Fatalf("SampleCount = %d, want %d", got.SampleCount, 2)
	}
	if got.AverageMidpoint != "02:45" {
		t.Fatalf("AverageMidpoint = %q, want %q", got.AverageMidpoint, "02:45")
	}
	if got.Score <= 0 {
		t.Fatalf("Score = %v, want positive", got.Score)
	}
}

func TestSleepStoreIngestValidatesAndReplaces(t *testing.T) {
	store, err := NewSleepStore(t.TempDir())
	if err != nil {
		t.Fatalf("NewSleepStore() error = %v", err)
	}
	event := SleepEvent{
		ID:        "evt-1",
		StartedAt: time.Date(2026, 5, 1, 22, 0, 0, 0, time.UTC),
		EndedAt:   time.Date(2026, 5, 2, 6, 0, 0, 0, time.UTC),
	}

	if _, err := store.Ingest(event); err != nil {
		t.Fatalf("Ingest() error = %v", err)
	}

	replaced := event
	replaced.EndedAt = replaced.EndedAt.Add(time.Hour)
	if _, err := store.Ingest(replaced); err != nil {
		t.Fatalf("Ingest() replacement error = %v", err)
	}

	got := store.Events()
	if len(got) != 1 {
		t.Fatalf("len(Events()) = %d, want %d", len(got), 1)
	}
	if got[0].EndedAt != replaced.EndedAt {
		t.Fatalf("stored event end = %v, want %v", got[0].EndedAt, replaced.EndedAt)
	}
}

func TestSleepStorePersistsAcrossInstances(t *testing.T) {
	dataDir := t.TempDir()

	store, err := NewSleepStore(dataDir)
	if err != nil {
		t.Fatalf("NewSleepStore() error = %v", err)
	}
	event := SleepEvent{
		ID:        "evt-1",
		StartedAt: time.Date(2026, 5, 1, 22, 0, 0, 0, time.UTC),
		EndedAt:   time.Date(2026, 5, 2, 6, 0, 0, 0, time.UTC),
	}
	if _, err := store.Ingest(event); err != nil {
		t.Fatalf("Ingest() error = %v", err)
	}

	reloaded, err := NewSleepStore(dataDir)
	if err != nil {
		t.Fatalf("NewSleepStore() reload error = %v", err)
	}

	got := reloaded.Events()
	if len(got) != 1 {
		t.Fatalf("len(Events()) = %d, want %d", len(got), 1)
	}
	if got[0].ID != event.ID {
		t.Fatalf("stored event ID = %q, want %q", got[0].ID, event.ID)
	}
}
