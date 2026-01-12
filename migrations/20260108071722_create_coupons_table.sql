-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE IF NOT EXISTS coupons
(
    id               SERIAL PRIMARY KEY,
    created_at       TIMESTAMP           NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP           NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at       TIMESTAMP,

    name             VARCHAR(255) NOT NULL,
    amount           INT                 NOT NULL,
    remaining_amount INT                 NOT NULL,

    CONSTRAINT amount_must_be_positive CHECK(amount >= 0),
    CONSTRAINT remaining_amount_must_be_positive CHECK(remaining_amount >= 0)
);

CREATE UNIQUE INDEX coupon_name_unique_idx ON coupons(name) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS user_claims
(
    id         SERIAL PRIMARY KEY,
    created_at TIMESTAMP                      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP                      NOT NULL DEFAULT CURRENT_TIMESTAMP,

    user_id    BIGINT REFERENCES users (id)   NOT NULL,
    coupon_id  BIGINT REFERENCES coupons (id) NOT NULL,
    UNIQUE (user_id, coupon_id)
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

DROP TABLE IF EXISTS user_claims;

DROP TABLE IF EXISTS coupons
