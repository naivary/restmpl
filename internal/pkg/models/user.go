package models

import "time"

type User struct {
	ID        string    `jsonapi:"primary,user"`
	CreatedAt time.Time `jsonapi:"attr,createdAt"`
	Username  string    `jsonapi:"attr,username"`
	Email     string    `jsonapi:"attr,email"`
	Password  string    `jsonapi:"attr,password"`
}

func NewUser() User {
	return User{
		CreatedAt: time.Now(),
	}
}
