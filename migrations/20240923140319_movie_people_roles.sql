-- +goose Up
-- +goose StatementBegin
CREATE TABLE movie_crew(
    movie_id INT REFERENCES movies(id),
    person_id INT REFERENCES people(id),
    role_id INT REFERENCES role(id)
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE movie_crew;
-- +goose StatementEnd