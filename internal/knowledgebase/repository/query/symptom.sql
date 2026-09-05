-- name: CreateSymptom :one
INSERT INTO symptoms (code, name, description, category_id)
VALUES ($1, $2, $3, $4)
RETURNING id, code, name, description, category_id;

-- name: GetSymptomByID :one
SELECT id, code, name, description, category_id
FROM symptoms
WHERE id = $1;

-- name: GetSymptoms :many 
SELECT id, code, name, description, category_id
FROM symptoms
ORDER BY id;