package entity

import "errors"

type Movie struct {
	ID          int64
	Name        string
	Year        int
	Description string
}

type MovieWithRating struct {
	Movies Movie
	Rating int64
}

var ErrMovieNotFound error = errors.New("movie not found")
