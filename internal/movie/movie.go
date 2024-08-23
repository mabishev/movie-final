package movie

import "errors"

type Movie struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Year int    `json:"year"`
}

var ErrNotFound error = errors.New("movie not found")
