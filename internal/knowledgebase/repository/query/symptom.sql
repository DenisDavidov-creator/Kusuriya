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

-- name: GetSymptomsByCodes :many 
SELECT c.id, sp.likelihood
FROM symptoms s
CROSS JOIN causes c 
LEFT JOIN symptom_profiles AS sp ON s.id = sp.symptom_id AND c.id = sp.cause_id
WHERE s.code = ANY($1::text[])
ORDER BY c.id;

-- name: GetSymptomsByIDs :many 
SELECT c.id, sp.likelihood
FROM symptoms s
CROSS JOIN causes c 
LEFT JOIN symptom_profiles AS sp ON s.id = sp.symptom_id AND c.id = sp.cause_id
WHERE s.id = ANY($1::bigint[])
ORDER BY c.id;