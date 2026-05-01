package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlerIngestsSleepEventAndReturnsPhase(t *testing.T) {
	handler := NewHandler("chronos")

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
	handler := NewHandler("chronos")

	req := httptest.NewRequest(http.MethodPost, "/sleep-events", strings.NewReader(`{"id":"bad","started_at":"oops"}`))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}
