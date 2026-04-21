-- name: GetAllProducts :many
SELECT * FROM products
WHERE is_available = 1
ORDER BY created_at DESC
LIMIT ?;

-- name: GetProductByID :one
SELECT * FROM products
WHERE id = ?;

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
    OR slug LIKE '%' || ? || '%'
)
ORDER BY created_at DESC
LIMIT ?;

-- name: FilterProductByPrice :many
SELECT * FROM products
WHERE is_available = 1
AND price >= ?
AND price <= ?
ORDER BY created_at DESC
LIMIT ?;

-- name: CreateProduct :one
INSERT INTO products (name, slug, type, price, quantity, description, about, video_url)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateProduct :one
UPDATE products
SET name = ?,
    slug = ?,
    type = ?,
    price = ?,
    quantity = ?,
    description = ?,
    about = ?,
    video_url = ?,
    is_available = ?,
    updated_at = datetime('now')
WHERE id = ?
RETURNING *;

-- name: UpdateProductVideoURL :exec
UPDATE products
SET video_url = ?, updated_at = datetime('now')
WHERE id = ?;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = ?;

-- name: GetAllProductsAdmin :many
SELECT * FROM products
ORDER BY created_at DESC
LIMIT ?;