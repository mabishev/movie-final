-- +goose Up
-- +goose StatementBegin
CREATE TABLE movie_people_roles(
    movie_id INT REFERENCES movies(id),
    person_id INT REFERENCES people(id),
    role_id INT REFERENCES roles(id)
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE movie_people_roles;
-- +goose StatementEnd