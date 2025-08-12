package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"stock/backend/pkg/handlers"

	cp_middleware "stock/backend/pkg/middleware"

	chi_middleware "github.com/go-chi/chi/v5/middleware"
)

func main() {

	godotenv.Overload()

	r := chi.NewRouter()
	r.Use(chi_middleware.Logger)
	r.Use(cp_middleware.ApplyCorsHandler())

	r.Route("/v1", func(r chi.Router) {
		r.Route("/api", func(r chi.Router) {
			r.Get("/stocks/list", handlers.GetStocksHandler)
			r.Get("/stocks/recommendations", handlers.GetBasicRecommendationsHandler)
		})
	})

	startSever(r)

}

func startSever(r *chi.Mux) {
	listenBy := os.Getenv("LISTENER")
	port := os.Getenv("PORT")
	if listenBy == "SOCKET" {
		if port == "" {
			port = "/tmp/stock-backend.sock"
		}
		setupSignalHandler(port)
		listener, err := net.Listen("unix", port)
		if err != nil {
			panic(err)
		}
		defer listener.Close()
		log.Printf("[Escuchando en archvivo %s]", port)
		http.Serve(listener, r)
	} else {
		if port == "" {
			port = "3000"
		}
		host := os.Getenv("HOST")
		log.Printf("[Escuchando en host:port (%s:%s)]", host, port)
		http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), r)
	}
}

func setupSignalHandler(socketPath string) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		// Eliminar el archivo de socket cuando se recibe una señal de interrupción o SIGTERM
		if err := os.Remove(socketPath); err != nil {
			fmt.Println("Error al eliminar el archivo de socket:", err)
		}
		os.Exit(1)
	}()
}
