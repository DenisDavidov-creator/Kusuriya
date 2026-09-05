-- name: CreateCategory :one
INSERT INTO categories (name)
VALUES ($1)
RETURNING id, name;

-- name: GetCategoryByID :one
SELECT id, name
FROM categories
WHERE id = $1;

-- name: GetCategories :many 
SELECT id, name
FROM categories
ORDER BY id;