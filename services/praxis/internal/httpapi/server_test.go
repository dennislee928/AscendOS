package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlerStoresStreakAndReturnsForecast(t *testing.T) {
	handler := NewHandler("praxis")

	createReq := httptest.NewRequest(http.MethodPost, "/habit-streaks", strings.NewReader(`{
		"id":"streak-1",
		"habit_id":"habit-1",
		"current_streak_days":5,
		"longest_streak_days":9,
		"missed_days_last_7":1,
		"last_completed_at":"2026-04-30T09:00:00Z"
	}`))
	createRec := httptest.NewRecorder()
	handler.ServeHTTP(createRec, createReq)
	if createRec.Code != http.StatusCreated {
		t.Fatalf("POST /habit-streaks status = %d, want %d", createRec.Code, http.StatusCreated)
	}

	forecastReq := httptest.NewRequest(http.MethodGet, "/relapse-risk?habit_id=habit-1&at=2026-05-01T09:00:00Z", nil)
	forecastRec := httptest.NewRecorder()
	handler.ServeHTTP(forecastRec, forecastReq)
	if forecastRec.Code != http.StatusOK {
		t.Fatalf("GET /relapse-risk status = %d, want %d", forecastRec.Code, http.StatusOK)
	}
	if body := forecastRec.Body.String(); !strings.Contains(body, `"habit_id":"habit-1"`) {
		t.Fatalf("GET /relapse-risk body = %s, want habit_id payload", body)
	}
}

func TestHandlerRejectsMissingHabitID(t *testing.T) {
	handler := NewHandler("praxis")

	req := httptest.NewRequest(http.MethodGet, "/relapse-risk", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestHandlerListsHabitStreaks(t *testing.T) {
	handler := NewHandler("praxis")

	req := httptest.NewRequest(http.MethodPost, "/habit-streaks", strings.NewReader(`{
		"id":"streak-2",
		"habit_id":"habit-2",
		"current_streak_days":8,
		"longest_streak_days":11,
		"missed_days_last_7":0
	}`))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("POST /habit-streaks status = %d, want %d", rec.Code, http.StatusCreated)
	}

	listReq := httptest.NewRequest(http.MethodGet, "/habit-streaks", nil)
	listRec := httptest.NewRecorder()
	handler.ServeHTTP(listRec, listReq)

	if listRec.Code != http.StatusOK {
		t.Fatalf("GET /habit-streaks status = %d, want %d", listRec.Code, http.StatusOK)
	}
	if body := listRec.Body.String(); !strings.Contains(body, `"habit_streaks"`) {
		t.Fatalf("GET /habit-streaks body = %s, want streak list", body)
	}
}
