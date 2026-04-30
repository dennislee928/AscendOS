package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
)

// Claims is a minimal JWT claims map for stubbed auth.
type Claims map[string]any

// ParseJWTStub extracts claims from the JWT payload without signature verification.
func ParseJWTStub(token string) (Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("invalid token payload")
	}

	claims := Claims{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
