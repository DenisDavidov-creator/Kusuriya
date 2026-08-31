package repository

import (
	"Kusuriya/internal/domain"
	sqlcgen "Kusuriya/internal/knowledgebase/repository/sqlc"
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type KnowledgebaseRepository struct {
	pool    *pgxpool.Pool
	queries *sqlcgen.Queries
}

func NewKnowledgebaseRepository(pool *pgxpool.Pool) *KnowledgebaseRepository {
	return &KnowledgebaseRepository{
		pool:    pool,
		queries: sqlcgen.New(pool),
	}
}

func (r *KnowledgebaseRepository) CreateCause(ctx context.Context, d domain.Cause) (*domain.Cause, error) {
	var remedyID pgtype.Int8
	if d.RemedyID != nil {
		remedyID = pgtype.Int8{Int64: *d.RemedyID, Valid: true}
	}
	params := sqlcgen.CreateCauseParams{
		Code:  d.Code,
		Name:  d.Name,
		Type:  d.Type,
		Prior: d.Prior,
		Description: pgtype.Text{
			String: d.Description,
			Valid:  d.Description != "",
		},
		RemedyID: remedyID,
	}

	result, err := r.queries.CreateCause(ctx, params)

	if err != nil {
		return nil, MapPgError(err)
	}

	dom := &domain.Cause{
		ID:    result.ID,
		Code:  result.Code,
		Name:  result.Name,
		Type:  result.Type,
		Prior: result.Prior,
	}
	if result.Description.Valid {
		dom.Description = result.Description.String
	}
	if result.RemedyID.Valid {
		id := result.RemedyID.Int64
		dom.RemedyID = &id
	}

	return dom, nil
}
