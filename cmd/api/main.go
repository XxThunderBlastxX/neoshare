package main

import (
	"strconv"

	"github.com/XxThunderBlastxX/neoshare/internal/routes"
	"github.com/XxThunderBlastxX/neoshare/internal/server"
)

func main() {
	app := server.New()

	r := routes.New(app)

	r.RegisterRoutes()

	if err := app.Listen(":" + strconv.Itoa(app.Config.Port)); err != nil {
		panic(err)
	}
}
