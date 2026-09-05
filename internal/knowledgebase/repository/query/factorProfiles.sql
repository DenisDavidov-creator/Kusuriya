-- name: CreateFactorProfiles :one
INSERT INTO factor_profiles (cause_id, factor_id, modifier)
VALUES ($1, $2, $3)
RETURNING id, cause_id, factor_id, modifier;

-- name: DeleteFactorProfiles :exec
DELETE FROM factor_profiles
WHERE cause_id = $1;