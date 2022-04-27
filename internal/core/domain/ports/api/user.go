package api

type GetUserResponse struct {
	Name  string
	Email string
}

type UpdateUserRequest struct {
	Name  string
	Email string
}

type User interface {
	CreateUser(name, email string) error
	DeleteUser(email string) error
	GetUser(email string) (*GetUserResponse, error)
	ListUsers() ([]GetUserResponse, error)
}
