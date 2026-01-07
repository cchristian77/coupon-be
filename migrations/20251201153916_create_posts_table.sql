-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TYPE post_status AS ENUM ('DRAFT', 'PUBLISHED');

CREATE TABLE IF NOT EXISTS posts
(
    id         SERIAL PRIMARY KEY,
    created_at TIMESTAMP                    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP                    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    user_id    BIGINT REFERENCES users (id) NOT NULL,
    slug       VARCHAR(255) UNIQUE          NOT NULL,
    title      VARCHAR(255)                 NOT NULL,
    body       TEXT                         NOT NULL,
    status     post_status                  NOT NULL DEFAULT 'DRAFT'
);

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

DROP TABLE IF EXISTS posts;

DROP TYPE post_status;