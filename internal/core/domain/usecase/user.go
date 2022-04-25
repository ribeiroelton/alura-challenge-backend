package usecase

import (
	"errors"
	"time"

	"github.com/ribeiroelton/alura-challenge-backend/internal/core/domain/model"
	"github.com/ribeiroelton/alura-challenge-backend/internal/core/domain/ports/api"
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
func NewUserService(c UserServiceConfig) api.User {
	return &UserService{
		Log: c.Log,
		DB:  c.DB,
	}
}

//CreateUser creates an user struct and persists it on database
//TODO Implement bcrpyt random password
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

//DeleteUser deletes an user by its email
func (s *UserService) DeleteUser(email string) error {
	err := s.DB.DeleteUserByEmail(email)
	if err != nil {
		return err
	}
	return nil
}

//GetUser returns an user by its email
func (s *UserService) GetUser(email string) (*api.GetUserResponse, error) {
	u, err := s.DB.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return &api.GetUserResponse{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}, nil
}

//GetUser returns a slice of users. Empty slice if no users found
func (s *UserService) ListUsers() ([]api.GetUserResponse, error) {
	us, err := s.DB.ListUsers()
	if err != nil {
		return nil, err
	}

	res := []api.GetUserResponse{}

	for _, u := range us {
		r := api.GetUserResponse{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
		}
		res = append(res, r)
	}
	return res, nil
}

//UpdateUser updates an user by its email.
//TODO Implement UpdateUser
func (s *UserService) UpdateUser(r *api.UpdateUserRequest) (*api.GetUserResponse, error) {
	return &api.GetUserResponse{}, nil
}

//hashPassword creates an new random password
//TODO Implement hashPassword
func (s *UserService) hashPassword(password string) (string, error) {
	return "", nil
}

//comparePassword validates if password is the same as the hash
//TODO Implement hashPassword
func (s *UserService) passwordCheck(password, hash string) (bool, error) {
	return true, nil
}
