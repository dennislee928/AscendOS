package main

import (
	"log"
	"net/http"

	"metis/internal/config"
	"metis/internal/httpapi"
)

func main() {
	cfg := config.Load()

	handler := httpapi.NewHandler(cfg.ServiceName)

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
