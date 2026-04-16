-- name: GetAllProducts :many
SELECT * FROM products
WHERE is_available = 1
ORDER BY created_at DESC
LIMIT ?;

-- name: GetProductByID :one
SELECT * FROM products
WHERE id = ? AND is_available = 1;

-- name: GetProductBySlug :one
SELECT * FROM products
WHERE slug = ? AND is_available = 1;

-- name: SearchProduct :many
SELECT * FROM products
WHERE is_available = 1
AND (
    name LIKE '%'|| ? || '%'
    OR description LIKE '%' || ? || '%'
    OR type LIKE '%' || ? || '%'
)
ORDER BY created_at DESC;

-- name: FilterProductByPrice :many
SELECT * FROM products
WHERE is_available = 1
AND price >= ?
AND price <= ?
ORDER BY created_at DESC;

-- name: CreateProduct :one
INSERT INTO products (name, slug, type, price, quantity, description, video_url)
VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateProduct :one
UPDATE products
SET name = ?,
    slug = ?,
    type = ?,
    price = ?,
    quantity = ?,
    description = ?,
    video_url = ?,
    is_available = ?,
    updated_at = datetime('now')
WHERE id = ?
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = ?;