-- +goose Up
ALTER TABLE products ADD COLUMN about TEXT;

-- +goose Down
ALTER TABLE products DROP COLUMN about;
