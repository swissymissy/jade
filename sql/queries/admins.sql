-- name: GetAdminByEmail :one
SELECT * FROM admins WHERE email = ?;

-- name: GetAdminByID :one
SELECT * FROM admins WHERE ID = ?;

-- name: CreateAdmin :one
INSERT INTO admins (id, email, password_hash, recovery_hash)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: UpdateAdminEmail :exec
UPDATE admins
SET email = ? , updated_at = datetime('now')
WHERE id = ?;

-- name: UpdateAdminPassword :exec
UPDATE admins 
SET password_hash = ?, updated_at = datetime('now')
WHERE id = ?;

-- name: GetAdmin :one
SELECT * FROM admins LIMIT 1;

-- name: UpdateAdminRecoveryHash :exec
UPDATE admins
SET recovery_hash = ?, updated_at = datetime('now')
WHERE id = ?;