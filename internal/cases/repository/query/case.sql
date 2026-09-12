-- name: CreateCase :one 
INSERT INTO cases (title, description, status)
VALUES ($1, $2, $3)
RETURNING id, title, description, status;

-- name: GetCaseByID :one
SELECT id, title, description, status
FROM cases
WHERE id = $1;

-- name: CloseCaseByID :one 
UPDATE cases
SET status = $1
WHERE id = $2
RETURNING id, title, description, status;
