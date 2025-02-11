-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS  "roles"(
    id UUID DEFAULT uuid_generate_v4() NOT NULL,
    name TEXT NOT NULL,
    level INT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS "roles";
-- +goose StatementEnd
