package main

import (
	"github.com/XxThunderBlastxX/neoshare/internal/routes"
	"github.com/XxThunderBlastxX/neoshare/internal/server"
	"strconv"
)

func main() {
	app := server.New()

	routes.RegisterRoutes(app)

	if err := app.Listen(":" + strconv.Itoa(app.Config.Port)); err != nil {
		panic(err)
	}
}
