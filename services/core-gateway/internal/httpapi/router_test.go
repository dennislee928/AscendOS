package httpapi

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHealthz(t *testing.T) {
	r := NewRouter(nil)
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestHealthzAddsRequestID(t *testing.T) {
	r := NewRouter(nil)
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
	r := NewRouter(nil)
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
	r := NewRouter(nil)
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
	r := NewRouter(nil)
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
	r := NewRouter(nil)
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestMeRejectsBlankBearerToken(t *testing.T) {
	r := NewRouter(nil)
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer   ")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestMeAcceptsValidatedJWT(t *testing.T) {
	secret := "test-secret"
	setJWTEnv(t, secret, "https://issuer.example", "gateway", 30*time.Second)

	r := NewRouter(nil)

	token := signJWT(t, secret, map[string]any{
		"sub": "user-123",
		"iss": "https://issuer.example",
		"aud": "gateway",
		"iat": time.Now().UTC().Add(-time.Second).Unix(),
		"nbf": time.Now().UTC().Add(-time.Second).Unix(),
		"exp": time.Now().UTC().Add(2 * time.Minute).Unix(),
	})

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var body struct {
		Claims map[string]any `json:"claims"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if got := body.Claims["sub"]; got != "user-123" {
		t.Fatalf("expected subject to be preserved, got %#v", got)
	}
}

func TestMeRejectsSignatureMismatch(t *testing.T) {
	setJWTEnv(t, "correct-secret", "", "", 0)

	r := NewRouter(nil)
	token := signJWT(t, "wrong-secret", map[string]any{
		"sub": "user-123",
		"exp": time.Now().UTC().Add(2 * time.Minute).Unix(),
	})

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestMeRejectsIssuerMismatch(t *testing.T) {
	setJWTEnv(t, "test-secret", "https://issuer.example", "", 0)

	r := NewRouter(nil)
	token := signJWT(t, "test-secret", map[string]any{
		"sub": "user-123",
		"iss": "https://issuer.example/other",
		"exp": time.Now().UTC().Add(2 * time.Minute).Unix(),
	})

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestMeRejectsAudienceMismatch(t *testing.T) {
	setJWTEnv(t, "test-secret", "", "gateway", 0)

	r := NewRouter(nil)
	token := signJWT(t, "test-secret", map[string]any{
		"sub": "user-123",
		"aud": []string{"metrics"},
		"exp": time.Now().UTC().Add(2 * time.Minute).Unix(),
	})

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestMeHonorsClockSkew(t *testing.T) {
	setJWTEnv(t, "test-secret", "", "", 45*time.Second)

	r := NewRouter(nil)
	token := signJWT(t, "test-secret", map[string]any{
		"sub": "user-123",
		"exp": time.Now().UTC().Add(-30 * time.Second).Unix(),
	})

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestMeDefaultsMissingSubject(t *testing.T) {
	setJWTEnv(t, "test-secret", "", "", 0)

	r := NewRouter(nil)
	token := signJWT(t, "test-secret", map[string]any{
		"exp": time.Now().UTC().Add(2 * time.Minute).Unix(),
	})

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var body struct {
		Claims map[string]any `json:"claims"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if got := body.Claims["sub"]; got != "unknown" {
		t.Fatalf("expected synthesized subject, got %#v", got)
	}
}

func signJWT(t *testing.T, secret string, claims map[string]any) string {
	t.Helper()

	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payloadBytes, err := json.Marshal(claims)
	if err != nil {
		t.Fatalf("failed to marshal claims: %v", err)
	}
	payload := base64.RawURLEncoding.EncodeToString(payloadBytes)
	signingInput := header + "." + payload

	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(signingInput))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return strings.Join([]string{header, payload, signature}, ".")
}

func setJWTEnv(t *testing.T, secret, issuer, audience string, skew time.Duration) {
	t.Helper()
	t.Setenv("CORE_GATEWAY_JWT_SECRET", secret)
	t.Setenv("CORE_GATEWAY_JWT_ISSUER", issuer)
	t.Setenv("CORE_GATEWAY_JWT_AUDIENCE", audience)
	t.Setenv("CORE_GATEWAY_JWT_CLOCK_SKEW", skew.String())
}
