package entity

import "errors"

type Movie struct {
	ID     int64
	Name   string
	Year   int
	Rating int64
}

var ErrNotFound error = errors.New("movie not found")
