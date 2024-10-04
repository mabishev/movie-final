-- +goose Up
-- +goose StatementBegin
CREATE TABLE role(
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE role;
-- +goose StatementEnd