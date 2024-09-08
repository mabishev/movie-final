-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN name VARCHAR(50),
    ADD COLUMN surname VARCHAR(50);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN name,
    DROP COLUMN surname;
-- +goose StatementEnd