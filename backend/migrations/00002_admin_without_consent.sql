-- +goose Up
-- CLI создаёт служебный аккаунт без фиктивной даты пользовательского согласия.

ALTER TABLE users
    ALTER COLUMN consent_at DROP NOT NULL;

-- +goose Down
-- Откат возможен только при отсутствии аккаунтов с consent_at IS NULL.
-- Отсутствующее согласие не заменяется вымышленной датой.

ALTER TABLE users
    ALTER COLUMN consent_at SET NOT NULL;
