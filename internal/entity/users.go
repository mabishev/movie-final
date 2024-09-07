package entity

import (
	"errors"
	"time"
)

type User struct {
	ID          int64
	Email       string
	Password    string
	Name        string
	Surname     string
	Sex         string
	DateOfBirth time.Time
	Country     string
	City        string
}

var ErrUserNotFound error = errors.New("user not found")
