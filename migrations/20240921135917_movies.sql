-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS movie CASCADE;
CREATE table movies(
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    release_date DATE NOT NULL,
    genre VARCHAR(100),
    description TEXT
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS movies CASCADE;
-- +goose StatementEnd