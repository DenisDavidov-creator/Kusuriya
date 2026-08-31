package repository

import (
	"Kusuriya/internal/domain"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func MapPgError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return domain.ErrDuplicateCode
		}
	}
	return err
}
