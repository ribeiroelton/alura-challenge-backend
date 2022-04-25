package ui

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ribeiroelton/alura-challenge-backend/internal/core/domain/ports/api"
	"github.com/ribeiroelton/alura-challenge-backend/pkg/logger"
)

//UserHandlerConfig config struct used as param on NewUserHandler
type UserHandlerConfig struct {
	Service api.User
	Log     logger.Logger
	Srv     *echo.Echo
}

//UserHandler handles and requests related to User api interface
type UserHandler struct {
	service api.User
	log     logger.Logger
}

//NewUserHandler creates and UserHandler
func NewUserHandler(c *UserHandlerConfig) {
	h := &UserHandler{
		service: c.Service,
		log:     c.Log,
	}

	c.Srv.GET("/users", h.GetUsers)
	c.Srv.GET("/users-edit", h.GetUsersEdit)
	c.Srv.POST("/users-edit", h.PostUsersEdit)
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
