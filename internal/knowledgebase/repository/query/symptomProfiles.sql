-- name: CreateSymptomProfiles :one
INSERT INTO symptom_profiles (cause_id, symptom_id, likelihood)
VALUES ($1, $2, $3)
RETURNING id, cause_id, symptom_id, likelihood;

-- name: DeleteSymptomProfiles :exec
DELETE FROM symptom_profiles
WHERE cause_id = $1;