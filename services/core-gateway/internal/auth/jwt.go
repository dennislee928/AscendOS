package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

// Claims contains a validated JWT payload.
type Claims map[string]any

type jwtHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ,omitempty"`
	Kid string `json:"kid,omitempty"`
}

// ParseJWT verifies and validates a JWT using the provided local configuration.
func ParseJWT(token string, cfg Config) (Claims, error) {
	cfg = normalizeConfig(cfg)

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	header, err := decodeSegment[sysHeader](parts[0])
	if err != nil {
		return nil, errors.New("invalid token header")
	}

	if err := validateAlgorithm(header.Alg, cfg.AllowedAlgorithms); err != nil {
		return nil, err
	}
	if header.Typ != "" && !strings.EqualFold(header.Typ, "JWT") {
		return nil, errors.New("invalid token type")
	}

	signingInput := parts[0] + "." + parts[1]
	if err := verifySignature(header.Alg, signingInput, parts[2], cfg.Secret); err != nil {
		return nil, err
	}

	claims, err := decodeClaims(parts[1])
	if err != nil {
		return nil, err
	}

	if err := validateRegisteredClaims(claims, cfg, time.Now().UTC()); err != nil {
		return nil, err
	}

	return claims, nil
}

type sysHeader = jwtHeader

func normalizeConfig(cfg Config) Config {
	if cfg.ClockSkew < 0 {
		cfg.ClockSkew = 0
	}
	if strings.TrimSpace(cfg.Secret) == "" {
		cfg.Secret = defaultJWTSecret
	}
	if len(cfg.AllowedAlgorithms) == 0 {
		cfg.AllowedAlgorithms = []string{"HS256"}
	}
	return cfg
}

func decodeSegment[T any](segment string) (T, error) {
	var value T
	raw, err := base64.RawURLEncoding.DecodeString(segment)
	if err != nil {
		return value, err
	}
	dec := json.NewDecoder(strings.NewReader(string(raw)))
	dec.UseNumber()
	if err := dec.Decode(&value); err != nil {
		return value, err
	}
	return value, nil
}

func decodeClaims(segment string) (Claims, error) {
	raw, err := base64.RawURLEncoding.DecodeString(segment)
	if err != nil {
		return nil, errors.New("invalid token payload")
	}

	dec := json.NewDecoder(strings.NewReader(string(raw)))
	dec.UseNumber()
	claims := Claims{}
	if err := dec.Decode(&claims); err != nil {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}

func validateAlgorithm(alg string, allowed []string) error {
	normalized := strings.TrimSpace(strings.ToUpper(alg))
	if normalized == "" {
		return errors.New("missing token algorithm")
	}
	for _, candidate := range allowed {
		if normalized == strings.ToUpper(strings.TrimSpace(candidate)) {
			if normalized == "NONE" {
				return errors.New("unsigned tokens are not accepted")
			}
			if normalized == "HS256" {
				return nil
			}
			return fmt.Errorf("unsupported token algorithm: %s", alg)
		}
	}
	return fmt.Errorf("unsupported token algorithm: %s", alg)
}

func verifySignature(alg, signingInput, signature, secret string) error {
	switch strings.ToUpper(strings.TrimSpace(alg)) {
	case "HS256":
		mac := hmac.New(sha256.New, []byte(secret))
		_, _ = mac.Write([]byte(signingInput))
		expected := mac.Sum(nil)

		got, err := base64.RawURLEncoding.DecodeString(signature)
		if err != nil {
			return errors.New("invalid token signature")
		}

		if !hmac.Equal(expected, got) {
			return errors.New("invalid token signature")
		}
		return nil
	default:
		return fmt.Errorf("unsupported token algorithm: %s", alg)
	}
}

func validateRegisteredClaims(claims Claims, cfg Config, now time.Time) error {
	if cfg.Issuer != "" {
		issuer, _ := claims["iss"].(string)
		if issuer == "" || issuer != cfg.Issuer {
			return errors.New("invalid token issuer")
		}
	}

	if len(cfg.Audiences) > 0 {
		if !claimAudienceMatches(claims["aud"], cfg.Audiences) {
			return errors.New("invalid token audience")
		}
	}

	if raw, ok := claims["exp"]; ok {
		exp, err := parseNumericDate(raw)
		if err != nil {
			return errors.New("invalid token expiration")
		}
		if now.After(exp.Add(cfg.ClockSkew)) {
			return errors.New("token expired")
		}
	}

	if raw, ok := claims["nbf"]; ok {
		nbf, err := parseNumericDate(raw)
		if err != nil {
			return errors.New("invalid token not-before")
		}
		if now.Before(nbf.Add(-cfg.ClockSkew)) {
			return errors.New("token not yet valid")
		}
	}

	if raw, ok := claims["iat"]; ok {
		iat, err := parseNumericDate(raw)
		if err != nil {
			return errors.New("invalid token issued-at")
		}
		if now.Before(iat.Add(-cfg.ClockSkew)) {
			return errors.New("token issued in the future")
		}
	}

	return nil
}

func claimAudienceMatches(raw any, audiences []string) bool {
	if len(audiences) == 0 {
		return true
	}

	switch value := raw.(type) {
	case string:
		return containsString(audiences, value)
	case []any:
		for _, candidate := range value {
			if audience, ok := candidate.(string); ok && containsString(audiences, audience) {
				return true
			}
		}
	case []string:
		for _, audience := range value {
			if containsString(audiences, audience) {
				return true
			}
		}
	}

	return false
}

func containsString(values []string, candidate string) bool {
	for _, value := range values {
		if value == candidate {
			return true
		}
	}
	return false
}

func parseNumericDate(raw any) (time.Time, error) {
	switch value := raw.(type) {
	case json.Number:
		return unixFromNumber(value)
	case float64:
		return unixFromFloat(value), nil
	case int64:
		return time.Unix(value, 0).UTC(), nil
	case int:
		return time.Unix(int64(value), 0).UTC(), nil
	case string:
		if strings.TrimSpace(value) == "" {
			return time.Time{}, errors.New("empty numeric date")
		}
		number := json.Number(value)
		return unixFromNumber(number)
	default:
		return time.Time{}, errors.New("unsupported numeric date type")
	}
}

func unixFromNumber(number json.Number) (time.Time, error) {
	if strings.Contains(number.String(), ".") {
		f, err := number.Float64()
		if err != nil {
			return time.Time{}, err
		}
		return unixFromFloat(f), nil
	}
	i, err := number.Int64()
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(i, 0).UTC(), nil
}

func unixFromFloat(value float64) time.Time {
	seconds, fraction := math.Modf(value)
	return time.Unix(int64(seconds), int64(fraction*float64(time.Second))).UTC()
}
