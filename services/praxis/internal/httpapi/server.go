package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"praxis/internal/domain"
)

// Handler serves Praxis health and forecast endpoints.
type Handler struct {
	service string
	store   *domain.StreakStore
	mux     *http.ServeMux
}

// NewHandler builds the Praxis HTTP handler.
func NewHandler(service, dataDir string) (http.Handler, error) {
	store, err := domain.NewStreakStore(dataDir)
	if err != nil {
		return nil, err
	}

	h := &Handler{
		service: service,
		store:   store,
		mux:     http.NewServeMux(),
	}

	h.mux.HandleFunc("/healthz", h.healthz)
	h.mux.HandleFunc("/habit-streaks", h.habitStreaks)
	h.mux.HandleFunc("/relapse-risk", h.relapseRisk)
	return h, nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *Handler) healthz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": h.service,
	})
}

func (h *Handler) habitStreaks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var req habitStreakRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}

		streak, err := req.toDomain()
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		stored, err := h.store.Upsert(streak)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, stored)
	case http.MethodGet:
		writeJSON(w, http.StatusOK, map[string]any{
			"habit_streaks": h.store.List(),
		})
	default:
		methodNotAllowed(w)
	}
}

func (h *Handler) relapseRisk(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}

	habitID := r.URL.Query().Get("habit_id")
	if habitID == "" {
		writeError(w, http.StatusBadRequest, "habit_id is required")
		return
	}

	now := time.Now().UTC()
	if value := r.URL.Query().Get("at"); value != "" {
		parsed, err := time.Parse(time.RFC3339, value)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid at query parameter")
			return
		}
		now = parsed
	}

	forecast, err := h.store.Forecast(habitID, now)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, forecast)
}

type habitStreakRequest struct {
	ID                string `json:"id"`
	HabitID           string `json:"habit_id"`
	CurrentStreakDays int    `json:"current_streak_days"`
	LongestStreakDays int    `json:"longest_streak_days"`
	MissedDaysLast7   int    `json:"missed_days_last_7"`
	LastCompletedAt   string `json:"last_completed_at"`
}

func (r habitStreakRequest) toDomain() (domain.HabitStreak, error) {
	streak := domain.HabitStreak{
		ID:                r.ID,
		HabitID:           r.HabitID,
		CurrentStreakDays: r.CurrentStreakDays,
		LongestStreakDays: r.LongestStreakDays,
		MissedDaysLast7:   r.MissedDaysLast7,
	}
	if r.LastCompletedAt != "" {
		parsed, err := time.Parse(time.RFC3339, r.LastCompletedAt)
		if err != nil {
			return domain.HabitStreak{}, errors.New("last_completed_at must be RFC3339")
		}
		streak.LastCompletedAt = parsed
	}
	return streak, nil
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func methodNotAllowed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}
