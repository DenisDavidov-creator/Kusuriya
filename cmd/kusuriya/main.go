package main

import (
	hch "Kusuriya/internal/cases/handler"
	rcr "Kusuriya/internal/cases/repository"
	scs "Kusuriya/internal/cases/service"
	hdh "Kusuriya/internal/diagnose/handler"
	sds "Kusuriya/internal/diagnose/service"
	"Kusuriya/internal/httpx"
	hkh "Kusuriya/internal/knowledgebase/handler"
	rkr "Kusuriya/internal/knowledgebase/repository"
	sks "Kusuriya/internal/knowledgebase/service"
	"Kusuriya/internal/logger"
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)

// @title           Pet Project API
// @version         1.0
// @description     A simple API for my pet project
// @host            localhost:8080
// @BasePath        /api/v1
func main() {
	// _ = godotenv.Load()

	log := logger.New()
	slog.SetDefault(log)

	pool := ConntectToDB()
	defer pool.Close()

	r := chi.NewRouter()
	rk := rkr.NewKnowledgebaseRepository(pool)
	rc := rcr.NewCaseRepository(pool)

	sk := sks.NewKnowledgebaseService(rk)
	sd := sds.NewDiagnoseService(rk)
	sc := scs.NewCaseService(rc, rk)

	hk := hkh.NewKnowledgebaseHandler(sk)
	hd := hdh.NewDiagnoseHandler(sd)
	hc := hch.NewCaseHandler(sc)

	r.Use(middleware.RequestID)
	r.Use(httpx.StructuredLogger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:8090"}, // the Swagger UI
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}))

	r.Route("/api/v1", func(r chi.Router) {
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

		r.Post("/diagnose", hd.DetermineDiagnose)

		r.Post("/case", hc.CreateCase)
		r.Get("/case/{id}", hc.GetCaseByID)
		r.Patch("/case/{id}/close", hc.CloseCaseByID)
		r.Post("/case/{id}/observation", hc.CreateObservation)
		r.Get("/case/{id}/observation", hc.GetObservationByID)
	})

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      8 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server failed", "err", err)
			stop()
		}
	}()
	log.Info("starting server", "port", 8080)

	<-ctx.Done()
	log.Info("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("forced shutdown", "err", err)
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
		slog.Error("We can't open the pool connection", "error", err.Error())
		os.Exit(1)
	}

	if err := pool.Ping(context.Background()); err != nil {
		slog.Error("db desn't answer", "error", err.Error())
		os.Exit(1)
	}
	log.Println("connect to db")
	return pool
}
