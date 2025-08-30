package main

import (
	"log"
	"net/http"

	"github.com/client-ux/internal/api/routes"
	"github.com/client-ux/internal/config"
)

func main() {
	cfg := config.Load()
	
	router := routes.SetupRoutes()
	
	log.Printf("🚀 CLIENT-UX server starting on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
