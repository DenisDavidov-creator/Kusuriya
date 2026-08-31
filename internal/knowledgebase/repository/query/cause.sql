-- name: CreateCause :one
INSERT INTO causes (code, name, type, prior, description, remedy_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, code, name, type, prior, description, remedy_id;

