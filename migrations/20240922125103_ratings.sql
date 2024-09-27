-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS ratings;
CREATE table ratings(
    userID INT references users(id),
    movieID INT references movies(id),
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