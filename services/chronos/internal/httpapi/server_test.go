package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlerIngestsSleepEventAndReturnsPhase(t *testing.T) {
	handler, err := NewHandler("chronos", t.TempDir())
	if err != nil {
		t.Fatalf("NewHandler() error = %v", err)
	}

	postReq := httptest.NewRequest(http.MethodPost, "/sleep-events", strings.NewReader(`{
		"id":"sleep-1",
		"user_id":"user-1",
		"started_at":"2026-05-01T23:00:00Z",
		"ended_at":"2026-05-02T07:00:00Z"
	}`))
	postRec := httptest.NewRecorder()
	handler.ServeHTTP(postRec, postReq)

	if postRec.Code != http.StatusCreated {
		t.Fatalf("POST /sleep-events status = %d, want %d", postRec.Code, http.StatusCreated)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/circadian-phase", nil)
	getRec := httptest.NewRecorder()
	handler.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Fatalf("GET /circadian-phase status = %d, want %d", getRec.Code, http.StatusOK)
	}
	if body := getRec.Body.String(); !strings.Contains(body, `"phase":"aligned"`) {
		t.Fatalf("GET /circadian-phase body = %s, want aligned phase", body)
	}
}

func TestHandlerRejectsInvalidSleepEvent(t *testing.T) {
	handler, err := NewHandler("chronos", t.TempDir())
	if err != nil {
		t.Fatalf("NewHandler() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/sleep-events", strings.NewReader(`{"id":"bad","started_at":"oops"}`))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestHandlerReloadsPersistedSleepEvents(t *testing.T) {
	dataDir := t.TempDir()

	handler, err := NewHandler("chronos", dataDir)
	if err != nil {
		t.Fatalf("NewHandler() error = %v", err)
	}
	postReq := httptest.NewRequest(http.MethodPost, "/sleep-events", strings.NewReader(`{
		"id":"sleep-1",
		"user_id":"user-1",
		"started_at":"2026-05-01T23:00:00Z",
		"ended_at":"2026-05-02T07:00:00Z"
	}`))
	postRec := httptest.NewRecorder()
	handler.ServeHTTP(postRec, postReq)
	if postRec.Code != http.StatusCreated {
		t.Fatalf("POST /sleep-events status = %d, want %d", postRec.Code, http.StatusCreated)
	}

	reloaded, err := NewHandler("chronos", dataDir)
	if err != nil {
		t.Fatalf("NewHandler() reload error = %v", err)
	}
	getReq := httptest.NewRequest(http.MethodGet, "/sleep-events", nil)
	getRec := httptest.NewRecorder()
	reloaded.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Fatalf("GET /sleep-events status = %d, want %d", getRec.Code, http.StatusOK)
	}
	if body := getRec.Body.String(); !strings.Contains(body, `"sleep-1"`) {
		t.Fatalf("GET /sleep-events body = %s, want persisted event", body)
	}
}
