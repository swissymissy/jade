-- +goose Up
ALTER TABLE admins 
ADD COLUMN recovery_hash TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE admins 
DROP COLUMN recovery_hash;