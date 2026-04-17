-- name: CreateProductImage :one
INSERT INTO product_images (product_id, s3_key, cover)
VALUES (?, ? , ?)
RETURNING *;

-- name: GetImagesByProductID :many
SELECT * FROM product_images 
WHERE product_id = ?
ORDER BY created_at ASC;

-- name: CountImagesByProductID :one
SELECT COUNT(*) FROM product_images
WHERE product_id = ?;

-- name: DeleteImage :exec
DELETE FROM product_images WHERE id = ?;

-- name: GetImageByID :one
SELECT * FROM product_images WHERE id = ?;