package main

import (
	"encoding/gob"
	"strconv"

	"github.com/XxThunderBlastxX/neoshare/internal/auth"
	"github.com/XxThunderBlastxX/neoshare/internal/routes"
	"github.com/XxThunderBlastxX/neoshare/internal/server"
	"github.com/XxThunderBlastxX/neoshare/internal/session"
)

func main() {
	app := server.New()

	sess := session.New()

	a, _ := auth.New(&app.Config.Auth)

	gob.Register(map[string]interface{}{})

	routes.New(app, a, sess).AuthRouter()

	if err := app.Listen(":" + strconv.Itoa(app.Config.Port)); err != nil {
		panic(err)
	}
}
