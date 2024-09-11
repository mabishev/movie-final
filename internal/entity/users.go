package entity

import (
	"errors"
	"time"
)

type CreateUser struct {
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

type User struct {
	ID          int64
	Name        string
	Surname     string
	Sex         string
	DateOfBirth time.Time
	Country     string
	City        string
}

var ErrUserNotFound error = errors.New("user not found")

//  FOR TEST:
//  "email": "test@gmail.com",
//  "password": "test"
