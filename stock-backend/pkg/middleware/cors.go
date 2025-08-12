package middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/cors"
)

func ApplyCorsHandler() func(http.Handler) http.Handler {

	allowedOrigin := os.Getenv("ALLOW_ORIGIN")
	log.Printf("[ALLOW_ORIGIN: (%s)]", allowedOrigin)

	var CORS_HANDLER func(http.Handler) http.Handler = cors.Handler(cors.Options{
		AllowedOrigins:   []string{allowedOrigin},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	})

	return CORS_HANDLER
}
