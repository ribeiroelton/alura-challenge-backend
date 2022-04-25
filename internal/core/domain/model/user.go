package model

import "time"

type User struct {
	ID        string
	Name      string `validate:"required"`
	Email     string `validate:"required,email"`
	Password  string
	CreatedAt time.Time `validate:"required"`
	UpdatedAt time.Time `validate:"required"`
}
