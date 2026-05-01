package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlerSchedulesAndReturnsDueCards(t *testing.T) {
	handler, err := NewHandler("metis", t.TempDir())
	if err != nil {
		t.Fatalf("NewHandler() error = %v", err)
	}

	createReq := httptest.NewRequest(http.MethodPost, "/cards", strings.NewReader(`{
		"id":"card-1",
		"front":"front",
		"back":"back",
		"due_at":"2026-05-01T00:00:00Z",
		"interval_days":1,
		"repetitions":0,
		"ease_factor":2.5
	}`))
	createRec := httptest.NewRecorder()
	handler.ServeHTTP(createRec, createReq)
	if createRec.Code != http.StatusCreated {
		t.Fatalf("POST /cards status = %d, want %d", createRec.Code, http.StatusCreated)
	}

	reviewReq := httptest.NewRequest(http.MethodPost, "/cards/review", strings.NewReader(`{
		"card_id":"card-1",
		"quality":4,
		"at":"2026-05-02T12:00:00Z"
	}`))
	reviewRec := httptest.NewRecorder()
	handler.ServeHTTP(reviewRec, reviewReq)
	if reviewRec.Code != http.StatusOK {
		t.Fatalf("POST /cards/review status = %d, want %d", reviewRec.Code, http.StatusOK)
	}

	dueReq := httptest.NewRequest(http.MethodGet, "/cards/due?at=2026-05-03T12:00:00Z", nil)
	dueRec := httptest.NewRecorder()
	handler.ServeHTTP(dueRec, dueReq)
	if dueRec.Code != http.StatusOK {
		t.Fatalf("GET /cards/due status = %d, want %d", dueRec.Code, http.StatusOK)
	}
	if body := dueRec.Body.String(); !strings.Contains(body, `"cards"`) {
		t.Fatalf("GET /cards/due body = %s, want cards payload", body)
	}
}

func TestHandlerRejectsInvalidReviewQuality(t *testing.T) {
	handler, err := NewHandler("metis", t.TempDir())
	if err != nil {
		t.Fatalf("NewHandler() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/cards/review", strings.NewReader(`{
		"card_id":"missing",
		"quality":9
	}`))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestHandlerReloadsPersistedCardState(t *testing.T) {
	dataDir := t.TempDir()

	handler, err := NewHandler("metis", dataDir)
	if err != nil {
		t.Fatalf("NewHandler() error = %v", err)
	}
	createReq := httptest.NewRequest(http.MethodPost, "/cards", strings.NewReader(`{
		"id":"card-1",
		"front":"front",
		"back":"back",
		"due_at":"2026-05-01T00:00:00Z",
		"interval_days":1,
		"repetitions":0,
		"ease_factor":2.5
	}`))
	createRec := httptest.NewRecorder()
	handler.ServeHTTP(createRec, createReq)
	if createRec.Code != http.StatusCreated {
		t.Fatalf("POST /cards status = %d, want %d", createRec.Code, http.StatusCreated)
	}

	reloaded, err := NewHandler("metis", dataDir)
	if err != nil {
		t.Fatalf("NewHandler() reload error = %v", err)
	}
	dueReq := httptest.NewRequest(http.MethodGet, "/cards/due?at=2026-05-02T00:00:00Z", nil)
	dueRec := httptest.NewRecorder()
	reloaded.ServeHTTP(dueRec, dueReq)
	if dueRec.Code != http.StatusOK {
		t.Fatalf("GET /cards/due status = %d, want %d", dueRec.Code, http.StatusOK)
	}
	if body := dueRec.Body.String(); !strings.Contains(body, `"card-1"`) {
		t.Fatalf("GET /cards/due body = %s, want persisted card", body)
	}
}
