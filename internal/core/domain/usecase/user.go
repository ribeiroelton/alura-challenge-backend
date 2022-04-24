package usecase

import (
	"errors"
	"time"

	"github.com/ribeiroelton/alura-challenge-backend/internal/core/domain/model"
	"github.com/ribeiroelton/alura-challenge-backend/internal/core/domain/ports/spi"
	"github.com/ribeiroelton/alura-challenge-backend/pkg/logger"
)

//UserServiceConfig config struct used as param for NewUserService
type UserServiceConfig struct {
	Log logger.Logger
	DB  spi.UserRepository
}

//UserService struct that implements all User api port.
type UserService struct {
	Log logger.Logger
	DB  spi.UserRepository
}

//NewUserService creates a new UserService
func NewUserService(c UserServiceConfig) *UserService {
	return &UserService{
		Log: c.Log,
		DB:  c.DB,
	}
}

//CreateUser creates an user struct and persists it on database
func (s *UserService) CreateUser(name, email string) error {
	ok, err := s.DB.HasUserByEmail(email)
	if err != nil {
		return err
	}
	if ok {
		return errors.New("email already exists")
	}

	m := &model.User{
		Name:      name,
		Email:     email,
		Password:  "123456",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.DB.SaveUser(m); err != nil {
		return err
	}
	return nil
}
