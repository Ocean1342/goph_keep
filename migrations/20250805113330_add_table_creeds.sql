-- +goose Up
CREATE TABLE creeds
(
    user_id  int  NOT NULL,
    name     TEXT NOT NULL UNIQUE,
    login    TEXT NOT NULL,
    password TEXT NOT NULL,
    meta     TEXT,
    CONSTRAINT fk_creeds_users FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

-- +goose Down
DROP TABLE creeds;