-- +goose Up
-- Таблицы блога и права сервисной учётной записи.

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    login_index TEXT NOT NULL UNIQUE,
    login_enc TEXT NOT NULL,
    email_enc TEXT NOT NULL,
    name_enc TEXT NOT NULL,
    bio_enc TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'user'
        CHECK (role IN ('user', 'admin')),
    banned BOOLEAN NOT NULL DEFAULT FALSE,
    token_version INTEGER NOT NULL DEFAULT 0,
    consent_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE posts (
    id BIGSERIAL PRIMARY KEY,
    author_id BIGINT NOT NULL REFERENCES users (id),
    title TEXT NOT NULL
        CHECK (char_length(title) BETWEEN 3 AND 120),
    body TEXT NOT NULL
        CHECK (char_length(body) BETWEEN 1 AND 10000),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE comments (
    id BIGSERIAL PRIMARY KEY,
    post_id BIGINT NOT NULL REFERENCES posts (id) ON DELETE CASCADE,
    author_id BIGINT NOT NULL REFERENCES users (id),
    body TEXT NOT NULL
        CHECK (char_length(body) BETWEEN 1 AND 2000),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX posts_author_idx
    ON posts (author_id, id DESC);

CREATE INDEX comments_post_idx
    ON comments (post_id, id);

GRANT SELECT
    ON users, posts, comments
    TO "user";

GRANT INSERT (
    login_index,
    login_enc,
    email_enc,
    name_enc,
    bio_enc,
    password_hash
)
    ON users
    TO "user";

GRANT UPDATE (
    name_enc,
    bio_enc,
    banned,
    token_version
)
    ON users
    TO "user";

GRANT INSERT, DELETE
    ON posts, comments
    TO "user";

GRANT USAGE, SELECT
    ON ALL SEQUENCES IN SCHEMA public
    TO "user";

-- +goose Down
-- Удаление таблиц при явном откате миграции.

DROP TABLE comments;
DROP TABLE posts;
DROP TABLE users;
