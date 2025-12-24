-- +goose Up
CREATE TABLE bin
(
    user_id int  NOT NULL,
    name    TEXT NOT NULL UNIQUE,
    bin     TEXT NOT NULL,
    meta    TEXT,
    CONSTRAINT fk_bin_users FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

-- +goose Down
DROP TABLE bin;