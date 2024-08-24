-- +goose Up
-- +goose StatementBegin
CREATE table ratings(
    userID INT references users(id),
    movieID INT references movie(id),
    rating INT check (
        rating >= 1
        and rating <= 10
    ),
    PRIMARY KEY (userID, movieID)
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE ratings;
-- +goose StatementEnd