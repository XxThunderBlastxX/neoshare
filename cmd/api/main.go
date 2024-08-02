package main

import (
	"github.com/XxThunderBlastxX/neoshare/internal/auth"
	"github.com/XxThunderBlastxX/neoshare/internal/routes"
	"github.com/XxThunderBlastxX/neoshare/internal/server"
	"github.com/gofiber/fiber/v2/middleware/session"
	"strconv"
)

func main() {
	app := server.New()

	sessionStore := session.New()

	a, _ := auth.New(&app.Config.Auth)

	routes.RegisterRoutes(app, a, sessionStore)

	if err := app.Listen(":" + strconv.Itoa(app.Config.Port)); err != nil {
		panic(err)
	}
}
