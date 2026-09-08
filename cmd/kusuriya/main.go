package main

import (
	hdh "Kusuriya/internal/diagnose/handler"
	sds "Kusuriya/internal/diagnose/service"
	hkh "Kusuriya/internal/knowledgebase/handler"
	"Kusuriya/internal/knowledgebase/repository"
	sks "Kusuriya/internal/knowledgebase/service"
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

	sk := sks.NewKnowledgebaseService(rk)
	sd := sds.NewDiagnoseService(rk)

	hk := hkh.NewKnowledgebaseHandler(sk)
	hd := hdh.NewDiagnoseHandler(sd)

	r.Post("/cause", hk.CreateCause)
	r.Get("/cause/{id}", hk.GetCauseByID)
	r.Get("/cause", hk.GetCauses)
	r.Put("/cause/{id}", hk.UpdateCause)
	r.Delete("/cause/{id}", hk.DeleteCause)

	r.Put("/cause/{id}/symptom-profiles", hk.CreateSymptomProfiles)
	r.Put("/cause/{id}/factor-profiles", hk.CreateFactorProfiles)

	r.Post("/remedy", hk.CreateRemedy)
	r.Get("/remedy/{id}", hk.GetRemedyByID)
	r.Get("/remedy", hk.GetRemedies)

	r.Post("/factor", hk.CreateFactor)
	r.Get("/factor/{id}", hk.GetFactorByID)
	r.Get("/factor", hk.GetFactors)

	r.Post("/category", hk.CreateCategory)
	r.Get("/category/{id}", hk.GetCategoryByID)
	r.Get("/category", hk.GetCategories)

	r.Post("/symptom", hk.CreateSymptom)
	r.Get("/symptom/{id}", hk.GetSymptomByID)
	r.Get("/symptom", hk.GetSymptoms)

	r.Post("/diagnose", hd.DeterminDiagnose)
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
