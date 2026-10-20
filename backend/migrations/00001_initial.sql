-- +goose Up

-- Пользователи, заявки и сохраняемые в PostgreSQL сессии.

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    surname TEXT NOT NULL
        CHECK (char_length(surname) BETWEEN 1 AND 20),
    name TEXT NOT NULL
        CHECK (char_length(name) BETWEEN 1 AND 20),
    patronymic TEXT NOT NULL
        CHECK (char_length(patronymic) BETWEEN 1 AND 20),
    login TEXT NOT NULL
        CHECK (char_length(login) BETWEEN 1 AND 20),
    password_hash TEXT NOT NULL,
    phone TEXT NOT NULL
        CHECK (char_length(phone) BETWEEN 1 AND 20),
    email TEXT NOT NULL
        CHECK (char_length(email) BETWEEN 1 AND 20),
    consent_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    consent_version TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX users_login_idx
    ON users (lower(login));

CREATE TABLE requests (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL
        REFERENCES users (id),
    service TEXT NOT NULL
        CHECK (
            service IN (
                'maintenance',
                'deep',
                'renovation'
            )
        ),
    address TEXT NOT NULL
        CHECK (char_length(address) BETWEEN 1 AND 100),
    cleaning_date DATE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX requests_owner_idx
    ON requests (user_id, id DESC);

CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    user_id BIGINT NOT NULL
        REFERENCES users (id)
        ON DELETE CASCADE,
    expires_at TIMESTAMPTZ NOT NULL DEFAULT now() + INTERVAL '30 minutes'
);

CREATE INDEX sessions_expiry_idx
    ON sessions (expires_at);

GRANT SELECT, INSERT
    ON users, requests
    TO "user";

GRANT SELECT, INSERT, UPDATE, DELETE
    ON sessions
    TO "user";

GRANT USAGE, SELECT
    ON ALL SEQUENCES IN SCHEMA public
    TO "user";

-- +goose Down

DROP TABLE sessions;

DROP TABLE requests;

DROP TABLE users;
