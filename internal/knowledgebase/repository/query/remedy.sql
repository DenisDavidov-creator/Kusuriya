-- name: CreateRemedy :one
INSERT INTO remedies (name, description, notes)
VALUES ($1, $2, $3)
RETURNING id, name, description, notes;

-- name: GetRemedyByID :one
SELECT id, name, description, notes
FROM remedies
WHERE id = $1;

-- name: GetRemedies :many 
SELECT id, name, description, notes
FROM remedies
ORDER BY id;