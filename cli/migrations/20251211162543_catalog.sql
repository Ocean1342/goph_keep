-- +goose Up
CREATE TABLE catalog
(
    user_id   INTEGER NOT NULL,
    data_name TEXT    NOT NULL,
    data_id   TEXT    NOT NULL UNIQUE,
    data_type TEXT    NOT NULL,
    is_synced INTEGER NOT NULL DEFAULT 0
);
-- +goose Down
DROP table IF EXISTS catalog;
