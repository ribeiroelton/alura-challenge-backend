package spi

import (
	"github.com/ribeiroelton/alura-challenge-backend/internal/core/domain/model"
)

type UserRepository interface {
	SaveUser(*model.User) error
	UpdateUserByEmail(*model.User) error
	DeleteUserByEmail(email string) error
	GetUserByEmail(email string) (*model.User, error)
	HasUserByEmail(email string) (bool, error)
	ListUsers() ([]model.User, error)
}
