package main

import (
	"github.com/XxThunderBlastxX/neoshare/internal/server"
	"strconv"
)

func main() {
	app := server.New()

	if err := app.Listen(":" + strconv.Itoa(app.Config.Port)); err != nil {
		panic(err)
	}
}
