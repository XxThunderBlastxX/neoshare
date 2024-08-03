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

	r := routes.New(app, a, sess)

	r.AuthRouter()
	r.S3Router()

	if err := app.Listen(":" + strconv.Itoa(app.Config.Port)); err != nil {
		panic(err)
	}
}
