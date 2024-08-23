package users

import "errors"

type User struct {
	ID          int32  `json:id`
	Email       string `json:email`
	Password    string `json:password`
	Sex         string `json:sex`
	DateOfBirth string `json:dateofbirth`
	Country     string `json:country`
	City        string `json:country`
}

var ErrUserNotFound error = errors.New("user not found")
