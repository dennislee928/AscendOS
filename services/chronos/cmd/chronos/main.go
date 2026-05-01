package main

import (
	"log"
	"net/http"

	"chronos/internal/config"
	"chronos/internal/httpapi"
)

func main() {
	cfg := config.Load()

	handler, err := httpapi.NewHandler(cfg.ServiceName, cfg.DataDir)
	if err != nil {
		log.Fatalf("handler error: %v", err)
	}

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           handler,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
	}

	log.Printf("%s listening on %s", cfg.ServiceName, cfg.HTTPAddr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
