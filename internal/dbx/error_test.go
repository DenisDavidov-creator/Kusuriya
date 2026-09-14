package dbx

import (
	"Kusuriya/internal/domain"
	"errors"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

func TestMapPgError(t *testing.T) {
	fk := &pgconn.PgError{Code: "23503", Message: "violates foreign key constraint"}
	other := errors.New("connection refused")

	tests := []struct {
		name string
		err  error
		want error
	}{
		{"no rows -> not found", pgx.ErrNoRows, domain.ErrNotFound},
		{"wrapped no rows", fmt.Errorf("get cause: %w", pgx.ErrNoRows), domain.ErrNotFound},
		{"unique violation -> duplicate code", &pgconn.PgError{Code: "23505"}, domain.ErrDuplicateCode},
		{"check violation -> check pass", &pgconn.PgError{Code: "23514"}, domain.ErrCheckPass},
		{"wrapped pg error", fmt.Errorf("insert: %w", &pgconn.PgError{Code: "23505"}), domain.ErrDuplicateCode},
		{"unmapped pg code is passed through unchanged", fk, fk},
		{"non-pg error is passed through unchanged", other, other},
		{"nil stays nil", nil, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MapPgError(tt.err)

			if tt.want == nil {
				assert.NoError(t, got)
				return
			}
			assert.ErrorIs(t, got, tt.want)
		})
	}
}
