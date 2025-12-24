-- +goose Up
CREATE TABLE users
(
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    login            TEXT    NOT NULL UNIQUE,
    password         TEXT    NOT NULL,
    phrase           TEXT    NOT NULL,
    token            TEXT,
    is_authenticated INTEGER NOT NULL DEFAULT 0,
    auth_dt          TEXT    NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose Down
DROP table IF EXISTS users;