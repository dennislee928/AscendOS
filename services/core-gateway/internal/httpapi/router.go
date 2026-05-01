package httpapi

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"core-gateway/internal/auth"
	natspub "core-gateway/internal/nats"
	"github.com/gin-gonic/gin"
)

const requestIDHeader = "X-Request-Id"

func NewRouter(pub *natspub.Publisher) *gin.Engine {
	catalog := newModuleCatalog()
	authConfig := auth.LoadConfigFromEnv()

	r := gin.New()
	r.Use(requestIDMiddleware(), gin.Logger(), gin.Recovery())

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/readyz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ready": true})
	})

	r.GET("/metrics", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/plain; version=0.0.4", []byte("# core-gateway metrics placeholder\n"))
	})

	r.GET("/modules", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"modules": catalog.all()})
	})

	r.GET("/modules/:name", func(c *gin.Context) {
		name := c.Param("name")
		module, ok := catalog.get(name)
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "module not found"})
			return
		}

		c.JSON(http.StatusOK, module)
		if pub != nil {
			pub.Publish("ascendos.gateway.events", map[string]any{"type": "module_access", "module": name, "ts": time.Now().UTC()})
		}
	})

	r.GET("/me", func(c *gin.Context) {
		authz := c.GetHeader("Authorization")
		if !strings.HasPrefix(authz, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}

		token := strings.TrimPrefix(authz, "Bearer ")
		if token == "" || strings.TrimSpace(token) != token {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid bearer token"})
			return
		}

		claims, err := auth.ParseJWT(token, authConfig)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		sub, ok := claims["sub"].(string)
		if !ok || strings.TrimSpace(sub) == "" {
			claims["sub"] = "unknown"
		}

		c.JSON(http.StatusOK, gin.H{"claims": claims})
		if pub != nil {
			pub.Publish("ascendos.gateway.events", map[string]any{"type": "module_access", "module": "me", "ts": time.Now().UTC()})
		}
	})

	return r
}

func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := strings.TrimSpace(c.GetHeader(requestIDHeader))
		if requestID == "" {
			requestID = generateRequestID()
		}

		c.Set(requestIDHeader, requestID)
		c.Writer.Header().Set(requestIDHeader, requestID)
		c.Next()
	}
}

func generateRequestID() string {
	var entropy [8]byte
	if _, err := rand.Read(entropy[:]); err == nil {
		return hex.EncodeToString(entropy[:]) + "-" + time.Now().UTC().Format("20060102T150405.000000000Z")
	}

	return time.Now().UTC().Format("20060102T150405.000000000Z")
}
