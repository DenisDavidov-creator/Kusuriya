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

-- name: GetCausesWithRemediesByIDs :many
SELECT c.id, c.code, c.name, c.type, c.prior, c.description, r.id AS remedy_id, r.name AS remedy_name, r.description AS remedy_description, r.notes AS remedy_notes
FROM causes c
LEFT JOIN remedies r ON c.remedy_id = r.id
WHERE c.id = ANY($1::bigint[])
ORDER BY c.id;