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

-- name: GetFactorsByCodes :many 
SELECT c.id, fp.modifier
FROM factors f
CROSS JOIN causes c 
LEFT JOIN factor_profiles AS fp ON f.id = fp.factor_id AND c.id = fp.cause_id
WHERE f.code = ANY($1::text[])
ORDER BY c.id;

-- name: GetFactorsByIDs :many 
SELECT c.id, fp.modifier
FROM factors f
CROSS JOIN causes c 
LEFT JOIN factor_profiles AS fp ON f.id = fp.factor_id AND c.id = fp.cause_id
WHERE f.id = ANY($1::bigint[])
ORDER BY c.id;