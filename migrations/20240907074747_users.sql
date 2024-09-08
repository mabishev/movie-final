-- +goose Up
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN login;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN login VARCHAR(50);
-- +goose StatementEnd