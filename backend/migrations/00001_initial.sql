-- +goose Up

-- Пользователи, заявки и сохраняемые в PostgreSQL сессии.

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    -- Поля содержат только шифротекст AES-256-GCM; длина ввода проверяется до шифрования.
    surname TEXT NOT NULL CHECK (surname LIKE 'v1:%'),
    name TEXT NOT NULL CHECK (name LIKE 'v1:%'),
    patronymic TEXT NOT NULL CHECK (patronymic LIKE 'v1:%'),
    login TEXT NOT NULL CHECK (login LIKE 'v1:%'),
    login_index TEXT NOT NULL CHECK (login_index ~ '^[0-9a-f]{64}$'),
    password_hash TEXT NOT NULL,
    phone TEXT NOT NULL CHECK (phone LIKE 'v1:%'),
    email TEXT NOT NULL CHECK (email LIKE 'v1:%'),
    consent_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    consent_version TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX users_login_idx
    ON users (login_index);

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
    address TEXT NOT NULL CHECK (address LIKE 'v1:%'),
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
