package usecase

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/ribeiroelton/alura-challenge-backend/internal/core/domain/model"
	"github.com/ribeiroelton/alura-challenge-backend/internal/core/domain/ports/api"
	"github.com/ribeiroelton/alura-challenge-backend/internal/core/domain/ports/spi"
	"golang.org/x/crypto/bcrypt"
)

//UserServiceConfig config struct used as param for NewUserService
type UserServiceConfig struct {
	Log    spi.Logger
	DB     spi.UserRepository
	Mailer spi.Mailer
}

//UserService struct that implements all User api port.
type UserService struct {
	Log    spi.Logger
	DB     spi.UserRepository
	mailer spi.Mailer
}

//NewUserService creates a new UserService
func NewUserService(c UserServiceConfig) *UserService {
	return &UserService{
		Log:    c.Log,
		DB:     c.DB,
		mailer: c.Mailer,
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
	pass := s.newRandonPassword()
	hashed, err := s.hashPassword(pass)
	if err != nil {
		return err
	}

	if err := s.mailer.Send(pass); err != nil {
		return err
	}

	m := &model.User{
		Name:      name,
		Email:     email,
		Password:  hashed,
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
			Name:  u.Name,
			Email: u.Email,
		}
		res = append(res, r)
	}
	return res, nil
}

func (s *UserService) newRandonPassword() string {
	rand.Seed(time.Now().UnixNano())
	pass := []string{}
	for i := 0; i < 6; i++ {
		pass = append(pass, strconv.Itoa(rand.Intn(9)))
	}
	return strings.Join(pass, "")

}

func (s *UserService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s *UserService) passwordCheck(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
