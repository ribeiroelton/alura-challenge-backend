package ui

import (
	"embed"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ribeiroelton/alura-challenge-backend/config"
	"github.com/ribeiroelton/alura-challenge-backend/internal/core/domain/ports/spi"
)

//go:embed static
var static embed.FS

type ServerConfig struct {
	Config *config.Config
	Log    spi.Logger
}

type Server struct {
	config *config.Config
	log    spi.Logger
	Srv    *echo.Echo
}

func NewServer(c ServerConfig) *Server {
	return &Server{
		config: c.Config,
		Srv:    echo.New(),
		log:    c.Log,
	}
}

func (u *Server) StartServer() error {
	body := middleware.BodyLimitConfig{
		Limit: "10M",
	}

	cors := middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
		AllowMethods: []string{"POST", "GET", "HEAD", "OPTIONS"},
	}

	u.Srv.Use(middleware.BodyLimitWithConfig(body))
	u.Srv.Use(middleware.CORSWithConfig(cors))
	u.Srv.Use(middleware.Logger())
	u.Srv.Use(middleware.Recover())

	u.Srv.StaticFS("/", static)

	renderer, err := newRender()
	if err != nil {
		return err
	}

	u.Srv.Renderer = renderer

	u.Srv.HideBanner = true
	if err := u.Srv.Start(u.config.ServerAddress); err != nil {
		return err
	}

	return nil
}
