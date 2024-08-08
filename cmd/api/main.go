package main

import (
	"strconv"

	"github.com/XxThunderBlastxX/neoshare/internal/database"
	"github.com/XxThunderBlastxX/neoshare/internal/routes"
	"github.com/XxThunderBlastxX/neoshare/internal/server"
)

func main() {
	app := server.New()

	_, err := database.ConnectDB(&app.Config.DBConfig)
	if err != nil {
		panic(err)
	}

	r := routes.New(app)

	r.RegisterRoutes()

	if err := app.Listen(":" + strconv.Itoa(app.Config.Port)); err != nil {
		panic(err)
	}
}
