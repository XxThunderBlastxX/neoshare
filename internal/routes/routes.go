package routes

import (
	"github.com/XxThunderBlastxX/neoshare/cmd/web"
	pages "github.com/XxThunderBlastxX/neoshare/cmd/web/page"
	"github.com/XxThunderBlastxX/neoshare/internal/server"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"net/http"
)

func RegisterRoutes(app *server.Server) {
	app.Use("/assets", filesystem.New(filesystem.Config{
		Root:       http.FS(web.Files),
		PathPrefix: "assets",
		Browse:     false,
	}))

	app.App.Get("/", func(c *fiber.Ctx) error {
		render := adaptor.HTTPHandler(templ.Handler(pages.HomePage()))

		return render(c)
	})
}
