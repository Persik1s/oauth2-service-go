-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS "users"(
    id UUID DEFAULT uuid_generate_v4() NOT NULL,
    username TEXT NOT NULL UNIQUE,
    password BYTEA NOT NULL,
    email TEXT NOT NULL UNIQUE,
    createAt TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS "users";
-- +goose StatementEnd
