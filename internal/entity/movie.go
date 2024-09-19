package entity

import "errors"

type Movie struct {
	ID   int64
	Name string
	Year int
}

type MovieWithRating struct {
	Movies Movie
	Rating int64
}

var ErrMovieNotFound error = errors.New("movie not found")
