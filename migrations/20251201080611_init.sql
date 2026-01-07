-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TYPE user_role AS ENUM ('ADMIN', 'USER');

CREATE TABLE IF NOT EXISTS users
(
    id         SERIAL PRIMARY KEY,
    created_at TIMESTAMP           NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP           NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    username   VARCHAR(100) UNIQUE NOT NULL,
    email      VARCHAR(255) UNIQUE NOT NULL,
    password   TEXT                NOT NULL,
    full_name  VARCHAR(255),
    role       user_role           NOT NULL DEFAULT 'USER'
);

CREATE TABLE IF NOT EXISTS sessions
(
    id                      SERIAL PRIMARY KEY,
    session_id              UUID                                  DEFAULT gen_random_uuid() UNIQUE,
    user_id                 BIGINT REFERENCES users (id) NOT NULL,
    access_token            TEXT                         NOT NULL,
    access_token_expires_at TIMESTAMPTZ                  NOT NULL,
    access_token_created_at TIMESTAMPTZ                  NOT NULL DEFAULT (now()),
    user_agent              VARCHAR(255)                 NOT NULL,
    client_ip               VARCHAR(255)                 NOT NULL
);


-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

DROP TABLE IF EXISTS sessions;

DROP TABLE IF EXISTS users;

DROP TYPE user_role;
