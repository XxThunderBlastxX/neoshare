package server

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"

	"github.com/XxThunderBlastxX/neoshare/internal/auth"
	"github.com/XxThunderBlastxX/neoshare/internal/config"
	"github.com/XxThunderBlastxX/neoshare/internal/database"
	"github.com/XxThunderBlastxX/neoshare/internal/session"
)

type Server struct {
	*fiber.App

	Config *config.AppConfig

	Session *session.Session

	Authenticator *auth.Authenticator

	Db *sql.DB
}

func New() *Server {
	c := config.New()

	a, err := auth.New(&c.Auth)
	if err != nil {
		panic(err)
	}

	db, err := database.ConnectDB(&c.DBConfig)
	if err != nil {
		panic(err)
	}

	server := &Server{
		App: fiber.New(fiber.Config{
			AppName:   "NeoShare",
			BodyLimit: 10 * 1024 * 1024 * 1024,
		}),
		Config:        c,
		Session:       session.New(),
		Authenticator: a,
		Db:            db,
	}

	return server
}
