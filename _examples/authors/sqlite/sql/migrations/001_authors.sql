-- +goose Up
CREATE TABLE IF NOT EXISTS authors (
    id   integer    PRIMARY KEY AUTOINCREMENT,
    name text   NOT NULL,
    bio  text
);

-- +goose Down
DROP TABLE IF EXISTS authors;