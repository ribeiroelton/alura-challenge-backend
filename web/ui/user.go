package ui

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ribeiroelton/alura-challenge-backend/internal/core/domain/ports/api"
	"github.com/ribeiroelton/alura-challenge-backend/internal/core/domain/ports/spi"
)

//UserHandlerConfig config struct used as param on NewUserHandler
type UserHandlerConfig struct {
	Service api.User
	Log     spi.Logger
}

//UserHandler handles and requests related to User api interface
type UserHandler struct {
	service api.User
	log     spi.Logger
}

//NewUserHandler creates and UserHandler
func NewUserHandler(c *UserHandlerConfig) *UserHandler {
	return &UserHandler{
		service: c.Service,
		log:     c.Log,
	}
}

func (h *UserHandler) GetUsers(c echo.Context) error {
	data := make(map[string]interface{})

	lu, err := h.service.ListUsers()
	if err != nil {
		data["status"] = "error"
		data["details"] = err.Error()
		c.Render(http.StatusOK, "users-page.tmpl", data)
		return err
	}
	data["users"] = lu
	c.Render(http.StatusOK, "users-page.tmpl", data)

	return nil
}

func (h *UserHandler) GetUsersEdit(c echo.Context) error {
	data := make(map[string]interface{})

	c.Render(http.StatusOK, "users-edit-page.tmpl", data)

	return nil
}

func (h *UserHandler) PostUsersEdit(c echo.Context) error {
	data := make(map[string]interface{})

	name := c.FormValue("name")
	email := c.FormValue("email")

	err := h.service.CreateUser(name, email)
	if err != nil {
		data["status"] = "error"
		data["details"] = err.Error()
		c.Render(http.StatusOK, "users-edit-page.tmpl", data)
		return err
	}

	lu, err := h.service.ListUsers()
	if err != nil {
		data["status"] = "error"
		data["details"] = err.Error()
		c.Render(http.StatusOK, "users-page.tmpl", data)
		return err
	}
	data["users"] = lu
	data["status"] = "ok"
	data["details"] = "User created with success"
	c.Render(http.StatusSeeOther, "users-page.tmpl", data)

	return nil
}
