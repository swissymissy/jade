-- +goose Up
CREATE TABLE products (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    name            TEXT NOT NULL,
    slug            TEXT NOT NULL UNIQUE,
    type            TEXT NOT NULL,
    price           REAL NOT NULL CHECK(price>0),
    quantity        INTEGER NOT NULL DEFAULT 0 CHECK(quantity>=0),
    description     TEXT,
    is_available    INTEGER NOT NULL DEFAULT 1,
    video_url       TEXT,
    created_at      TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at      TEXT NOT NULL DEFAULT (datetime('now'))
);

-- +goose Down
DROP TABLE products;