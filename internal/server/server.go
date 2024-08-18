package server

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"

	"github.com/XxThunderBlastxX/neoshare/internal/auth"
	"github.com/XxThunderBlastxX/neoshare/internal/config"
	"github.com/XxThunderBlastxX/neoshare/internal/database"
	"github.com/XxThunderBlastxX/neoshare/internal/session"
	"github.com/XxThunderBlastxX/neoshare/internal/utils"
)

type Server struct {
	*fiber.App

	Config *config.AppConfig

	Favicon []byte

	Session *session.Session

	Authenticator *auth.Authenticator

	DB *sql.DB
}

func New() *Server {
	c := config.New()

	a, err := auth.New(&c.Auth)
	if err != nil {
		panic(err)
	}

	db := database.MustConnectDB(&c.DBConfig)

	icon, _ := utils.GetFavicon(c.FaviconURI)

	server := &Server{
		App: fiber.New(fiber.Config{
			AppName:   "NeoShare",
			BodyLimit: 10 * 1024 * 1024 * 1024,
		}),
		Config:        c,
		Favicon:       icon,
		Session:       session.New(),
		Authenticator: a,
		DB:            db,
	}

	return server
}
