package api

type GetUserResponse struct {
	ID    string
	Name  string
	Email string
}

type UpdateUserRequest struct {
	Name  string
	Email string
}

type User interface {
	CreateUser(name, email string) error
	UpdateUser(r *UpdateUserRequest) (*GetUserResponse, error)
	DeleteUser(email string) error
	GetUser(email string) (*GetUserResponse, error)
	ListUsers() ([]GetUserResponse, error)
}
