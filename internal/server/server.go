package server

import (
	"github.com/gofiber/fiber/v2"

	"github.com/XxThunderBlastxX/neoshare/internal/auth"
	"github.com/XxThunderBlastxX/neoshare/internal/config"
	"github.com/XxThunderBlastxX/neoshare/internal/session"
)

type Server struct {
	*fiber.App

	Config *config.AppConfig

	Session *session.Session

	Authenticator *auth.Authenticator
}

func New() *Server {
	c := config.New()

	a, err := auth.New(&c.Auth)
	if err != nil {
		panic(err)
	}

	server := &Server{
		App: fiber.New(fiber.Config{
			AppName: "NeoShare",
		}),
		Config:        c,
		Session:       session.New(),
		Authenticator: a,
	}

	return server
}
