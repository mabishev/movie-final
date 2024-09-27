-- +goose Up
-- +goose StatementBegin
CREATE TABLE people(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    date_of_birth DATE,
    bio TEXT
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE people;
-- +goose StatementEnd