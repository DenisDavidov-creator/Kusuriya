package main

import (
	"Kusuriya/internal/logger"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	log := logger.New()

	r := chi.NewRouter()

	r.Get("/something", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("something \n"))
	})

	log.Info("starting server", "port", 8080)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Error("server failed", "err", err)
		os.Exit(1)
	}

}
