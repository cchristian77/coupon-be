-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE IF NOT EXISTS comments
(
    id         SERIAL PRIMARY KEY,
    created_at TIMESTAMP                    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP                    NOT NULL DEFAULT CURRENT_TIMESTAMP,

    user_id    BIGINT REFERENCES users (id) NOT NULL,
    post_id    BIGINT REFERENCES posts (id) NOT NULL,
    comment    TEXT                         NOT NULL,
    rating     SMALLINT
);


-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

DROP TABLE IF EXISTS comments;