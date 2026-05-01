package httpapi

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthz(t *testing.T) {
	r := NewRouter()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestHealthzAddsRequestID(t *testing.T) {
	r := NewRouter()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	req.Header.Set(requestIDHeader, "req-123")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	if got := w.Header().Get(requestIDHeader); got != "req-123" {
		t.Fatalf("expected response request id %q, got %q", "req-123", got)
	}
}

func TestModulesListsCatalog(t *testing.T) {
	r := NewRouter()
	req := httptest.NewRequest(http.MethodGet, "/modules", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var body struct {
		Modules []struct {
			Name     string `json:"name"`
			Key      string `json:"key"`
			Language string `json:"language"`
			Domain   string `json:"domain"`
		} `json:"modules"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if got, want := len(body.Modules), 8; got != want {
		t.Fatalf("expected %d modules, got %d", want, got)
	}

	first := body.Modules[0]
	if first.Name != "chronos" || first.Key != "CHRONOS" || first.Language != "Go" || first.Domain != "Sleep science" {
		t.Fatalf("unexpected first module: %+v", first)
	}
}

func TestModulesGetByName(t *testing.T) {
	r := NewRouter()
	req := httptest.NewRequest(http.MethodGet, "/modules/argentum", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var body struct {
		Name     string `json:"name"`
		Key      string `json:"key"`
		Language string `json:"language"`
		Domain   string `json:"domain"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if body.Name != "argentum" || body.Key != "ARGENTUM" || body.Language != "Python" || body.Domain != "Finance" {
		t.Fatalf("unexpected module payload: %+v", body)
	}
}

func TestModulesGetByNameReturns404(t *testing.T) {
	r := NewRouter()
	req := httptest.NewRequest(http.MethodGet, "/modules/unknown", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if got := body["error"]; got != "module not found" {
		t.Fatalf("expected module not found error, got %q", got)
	}
}

func TestMeRequiresBearer(t *testing.T) {
	r := NewRouter()
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestMeRejectsBlankBearerToken(t *testing.T) {
	r := NewRouter()
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer   ")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestMeAcceptsStubJWT(t *testing.T) {
	r := NewRouter()

	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{"sub":"user-123"}`))
	token := strings.Join([]string{header, payload, "sig"}, ".")

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
