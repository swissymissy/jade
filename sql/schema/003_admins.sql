-- +goose Up
CREATE TABLE admins (
    id              TEXT PRIMARY KEY,
    email           TEXT NOT NULL UNIQUE,
    password_hash   TEXT NOT NULL,
    created_at      TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at      TEXT NOT NULL DEFAULT (datetime('now'))
);

-- +goose StatementBegin
CREATE TRIGGER enforce_single_admin
BEFORE INSERT ON admins
WHEN (SELECT COUNT(*) FROM admins) >= 1
BEGIN
    SELECT RAISE(ABORT, 'Only one admin allowed');
END;
-- +goose StatementEnd

-- +goose Down
DROP TRIGGER enforce_single_admin;
DROP TABLE admins; 