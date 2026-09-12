-- name: CreateObservation :one 
INSERT INTO  observations (case_id, kind, ref_id, note) 
VALUES ($1, $2, $3, $4)
RETURNING id, case_id, kind, ref_id, observed_at, note;

-- name: GetObservationByID :many
SELECT id, case_id, kind, ref_id, observed_at, note
FROM observations
WHERE case_id = $1
ORDER BY id;