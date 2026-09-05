-- name: CreateFactor :one
INSERT INTO factors (code, name, description)
VALUES ($1, $2, $3)
RETURNING id, code, name, description;

-- name: GetFactorByID :one
SELECT id, code, name, description
FROM factors
WHERE id = $1;

-- name: GetFactors :many 
SELECT id, code, name, description
FROM factors
ORDER BY id;