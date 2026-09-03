-- name: CreateCause :one
INSERT INTO causes (code, name, type, prior, description, remedy_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, code, name, type, prior, description, remedy_id;

-- name: GetCauseByID :one 
SELECT id, code, name, type, prior, description, remedy_id FROM causes
WHERE id = $1;

-- name: GetCauses :many
SELECT id, code, name, type, prior, description, remedy_id FROM causes
ORDER BY id;

-- name: UpdateCause :one
UPDATE causes 
SET code = $1,
    name = $2,
    type = $3,
    prior = $4,
    description = $5,
    remedy_id = $6
WHERE id = $7
RETURNING id, code, name, type, prior, description, remedy_id;

-- name: DeleteCause :exec 
DELETE FROM causes
WHERE id = $1;