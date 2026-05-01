package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"metis/internal/domain"
)

// Handler serves Metis health and scheduling endpoints.
type Handler struct {
	service string
	store   *domain.CardStore
	mux     *http.ServeMux
}

// NewHandler builds the Metis HTTP handler.
func NewHandler(service, dataDir string) (http.Handler, error) {
	store, err := domain.NewCardStore(dataDir)
	if err != nil {
		return nil, err
	}

	h := &Handler{
		service: service,
		store:   store,
		mux:     http.NewServeMux(),
	}

	h.mux.HandleFunc("/healthz", h.healthz)
	h.mux.HandleFunc("/cards", h.cards)
	h.mux.HandleFunc("/cards/review", h.reviewCard)
	h.mux.HandleFunc("/cards/due", h.dueCards)
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

func (h *Handler) cards(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var req cardRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}

		card, err := req.toDomain()
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		stored, err := h.store.Upsert(card)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, stored)
	case http.MethodGet:
		writeJSON(w, http.StatusOK, map[string]any{
			"cards": h.store.Due(time.Now().UTC()),
		})
	default:
		methodNotAllowed(w)
	}
}

func (h *Handler) reviewCard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w)
		return
	}

	var req reviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.CardID == "" {
		writeError(w, http.StatusBadRequest, "card_id is required")
		return
	}

	now := time.Now().UTC()
	if req.At != "" {
		parsed, err := time.Parse(time.RFC3339, req.At)
		if err != nil {
			writeError(w, http.StatusBadRequest, "at must be RFC3339")
			return
		}
		now = parsed
	}

	result, err := h.store.Review(req.CardID, req.Quality, now)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) dueCards(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
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

	writeJSON(w, http.StatusOK, map[string]any{
		"cards": h.store.Due(now),
	})
}

type cardRequest struct {
	ID           string  `json:"id"`
	Front        string  `json:"front"`
	Back         string  `json:"back"`
	DueAt        string  `json:"due_at"`
	IntervalDays int     `json:"interval_days"`
	Repetitions  int     `json:"repetitions"`
	EaseFactor   float64 `json:"ease_factor"`
}

func (r cardRequest) toDomain() (domain.Card, error) {
	card := domain.Card{
		ID:           r.ID,
		Front:        r.Front,
		Back:         r.Back,
		IntervalDays: r.IntervalDays,
		Repetitions:  r.Repetitions,
		EaseFactor:   r.EaseFactor,
	}
	if r.DueAt != "" {
		parsed, err := time.Parse(time.RFC3339, r.DueAt)
		if err != nil {
			return domain.Card{}, errors.New("due_at must be RFC3339")
		}
		card.DueAt = parsed
	}
	if card.DueAt.IsZero() {
		card.DueAt = time.Now().UTC()
	}
	return card, nil
}

type reviewRequest struct {
	CardID  string `json:"card_id"`
	Quality int    `json:"quality"`
	At      string `json:"at"`
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
