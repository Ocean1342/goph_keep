-- +goose Up
CREATE table creeds
(
    id       TEXT NOT NULL UNIQUE,
    login    TEXT NOT NULL,
    password TEXT NOT NULL,
    meta     TEXT
);

-- +goose Down
DROP table IF EXISTS creeds;