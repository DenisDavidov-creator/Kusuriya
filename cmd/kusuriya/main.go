package main

import (
	"Kusuriya/internal/knowledgebase/handler"
	"Kusuriya/internal/knowledgebase/repository"
	"Kusuriya/internal/knowledgebase/service"
	"Kusuriya/internal/logger"
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	log := logger.New()
	slog.SetDefault(log)

	r := chi.NewRouter()
	pool := ConntectToDB()
	defer pool.Close()

	rk := repository.NewKnowledgebaseRepository(pool)

	sk := service.NewKnowledgebaseService(rk)

	hk := handler.NewKnowledgebaseHandler(sk)

	r.Post("/something", hk.CreateCause)

	log.Info("starting server", "port", 8080)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Error("server failed", "err", err)
		os.Exit(1)
	}
}

func ConntectToDB() *pgxpool.Pool {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("We can't open the pool connection %v", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("db doesn't answer %v", err)
	}
	log.Println("connect to db")
	return pool
}
