-- +goose Up
CREATE table bin
(
    id   TEXT NOT NULL UNIQUE,
    bin  TEXT NOT NULL,
    meta TEXT
);

-- +goose Down
DROP table IF EXISTS bin;