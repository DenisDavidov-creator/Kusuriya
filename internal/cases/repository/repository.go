package repository

import (
	sqlcgen "Kusuriya/internal/cases/repository/sqlc"
	"Kusuriya/internal/dbx"
	"Kusuriya/internal/domain"
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CaseRepository struct {
	pool    *pgxpool.Pool
	queries *sqlcgen.Queries
}

func NewCaseRepository(pool *pgxpool.Pool) *CaseRepository {
	return &CaseRepository{
		pool:    pool,
		queries: sqlcgen.New(pool),
	}
}

func (r *CaseRepository) CreateCase(ctx context.Context, d domain.Case) (*domain.Case, error) {
	params := sqlcgen.CreateCaseParams{
		Title: d.Title,
		Description: pgtype.Text{
			String: d.Description,
			Valid:  d.Description != "",
		},
		Status: d.Status,
	}

	result, err := r.queries.CreateCase(ctx, params)
	if err != nil {
		return nil, dbx.MapPgError(err)
	}
	dom := &domain.Case{
		ID:     result.ID,
		Title:  result.Title,
		Status: result.Status,
	}
	if result.Description.Valid {
		dom.Description = result.Description.String
	}
	return dom, nil
}

func (r *CaseRepository) GetCaseByID(ctx context.Context, id int) (*domain.Case, error) {
	result, err := r.queries.GetCaseByID(ctx, int64(id))
	if err != nil {
		return nil, dbx.MapPgError(err)
	}
	dom := &domain.Case{
		ID:     result.ID,
		Title:  result.Title,
		Status: result.Status,
	}
	if result.Description.Valid {
		dom.Description = result.Description.String
	}

	return dom, nil
}

func (r *CaseRepository) CloseCaseByID(ctx context.Context, id int) (*domain.Case, error) {

	params := sqlcgen.CloseCaseByIDParams{
		Status: "closed",
		ID:     int64(id),
	}
	result, err := r.queries.CloseCaseByID(ctx, params)
	if err != nil {
		return nil, dbx.MapPgError(err)
	}
	dom := &domain.Case{
		ID:     result.ID,
		Title:  result.Title,
		Status: result.Status,
	}
	if result.Description.Valid {
		dom.Description = result.Description.String
	}

	return dom, nil
}

func (r *CaseRepository) CreateObservation(ctx context.Context, d domain.Observation) (*domain.Observation, error) {
	params := sqlcgen.CreateObservationParams{
		CaseID: d.CaseID,
		Kind: pgtype.Text{
			String: d.Kind,
			Valid:  d.Kind != "",
		},
		RefID: d.RefID,
		Note: pgtype.Text{
			String: d.Note,
			Valid:  d.Note != "",
		},
	}

	result, err := r.queries.CreateObservation(ctx, params)
	if err != nil {
		return nil, dbx.MapPgError(err)
	}

	dom := &domain.Observation{
		ID:         result.ID,
		CaseID:     result.CaseID,
		RefID:      result.RefID,
		ObservedAt: result.ObservedAt.Time,
	}
	if result.Kind.Valid {
		dom.Kind = result.Kind.String
	}
	if result.Note.Valid {
		dom.Note = result.Note.String
	}
	return dom, nil
}

func (r *CaseRepository) GetObservationByID(ctx context.Context, id int) ([]domain.Observation, error) {
	results, err := r.queries.GetObservationByID(ctx, int64(id))
	if err != nil {
		return nil, dbx.MapPgError(err)
	}

	var doms []domain.Observation
	for _, result := range results {
		dom := domain.Observation{
			ID:         result.ID,
			CaseID:     result.CaseID,
			RefID:      result.RefID,
			ObservedAt: result.ObservedAt.Time,
		}
		if result.Kind.Valid {
			dom.Kind = result.Kind.String
		}
		if result.Note.Valid {
			dom.Note = result.Note.String
		}
		doms = append(doms, dom)
	}
	return doms, nil
}
