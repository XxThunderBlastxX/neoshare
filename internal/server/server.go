package server

import (
	"github.com/XxThunderBlastxX/neoshare/internal/config"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	*fiber.App

	Config *config.AppConfig
}

func New() *Server {
	server := &Server{
		App: fiber.New(fiber.Config{
			AppName: "NeoShare",
		}),
		Config: config.New(),
	}

	return server
}
