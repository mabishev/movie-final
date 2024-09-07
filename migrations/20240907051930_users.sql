-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN sex varchar(10),
    ADD COLUMN dateofbirth DATE,
    ADD COLUMN country varchar(50),
    ADD COLUMN city varchar(50);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN sex,
    DROP COLUMN dateofbirth,
    DROP COLUMN country,
    DROP COLUMN city;
-- +goose StatementEnd