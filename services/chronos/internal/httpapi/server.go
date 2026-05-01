package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"chronos/internal/domain"
)

// Handler serves Chronos health and domain endpoints.
type Handler struct {
	service string
	store   *domain.SleepStore
	mux     *http.ServeMux
}

// NewHandler builds the Chronos HTTP handler.
func NewHandler(service, dataDir string) (http.Handler, error) {
	store, err := domain.NewSleepStore(dataDir)
	if err != nil {
		return nil, err
	}

	h := &Handler{
		service: service,
		store:   store,
		mux:     http.NewServeMux(),
	}

	h.mux.HandleFunc("/healthz", h.healthz)
	h.mux.HandleFunc("/sleep-events", h.sleepEvents)
	h.mux.HandleFunc("/circadian-phase", h.circadianPhase)
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

func (h *Handler) sleepEvents(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var req sleepEventRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}

		event, err := req.toDomain()
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		stored, err := h.store.Ingest(event)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, stored)
	case http.MethodGet:
		writeJSON(w, http.StatusOK, map[string]any{
			"events": h.store.Events(),
		})
	default:
		methodNotAllowed(w)
	}
}

func (h *Handler) circadianPhase(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}

	calculatedAt := time.Now().UTC()
	if value := r.URL.Query().Get("at"); value != "" {
		parsed, err := time.Parse(time.RFC3339, value)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid at query parameter")
			return
		}
		calculatedAt = parsed
	}

	writeJSON(w, http.StatusOK, h.store.Phase(calculatedAt))
}

type sleepEventRequest struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	StartedAt string `json:"started_at"`
	EndedAt   string `json:"ended_at"`
}

func (r sleepEventRequest) toDomain() (domain.SleepEvent, error) {
	startedAt, err := time.Parse(time.RFC3339, r.StartedAt)
	if err != nil {
		return domain.SleepEvent{}, errors.New("started_at must be RFC3339")
	}
	endedAt, err := time.Parse(time.RFC3339, r.EndedAt)
	if err != nil {
		return domain.SleepEvent{}, errors.New("ended_at must be RFC3339")
	}
	return domain.SleepEvent{
		ID:        r.ID,
		UserID:    r.UserID,
		StartedAt: startedAt,
		EndedAt:   endedAt,
	}, nil
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
