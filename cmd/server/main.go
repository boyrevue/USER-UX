package main

import (
	"log"
	"net/http"

	"client-ux/internal/api/routes"
	"client-ux/internal/config"
)

func main() {
	cfg := config.Load()

	router := routes.SetupRoutes()

	log.Printf("🚀 CLIENT-UX server starting on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
