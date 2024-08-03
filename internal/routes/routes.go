package routes

import (
	"github.com/XxThunderBlastxX/neoshare/internal/service"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"

	"github.com/XxThunderBlastxX/neoshare/cmd/web"
	"github.com/XxThunderBlastxX/neoshare/internal/auth"
	"github.com/XxThunderBlastxX/neoshare/internal/middleware"
	"github.com/XxThunderBlastxX/neoshare/internal/server"
	"github.com/XxThunderBlastxX/neoshare/internal/session"
)

type Router struct {
	app           *server.Server
	authenticator *auth.Authenticator
	sessionStore  *session.Session
	middleware    *middleware.Middleware
	s3service     service.S3Service
}

func New(app *server.Server, authenticator *auth.Authenticator, sessionStore *session.Session) *Router {
	app.Use("/assets", filesystem.New(filesystem.Config{
		Root:       http.FS(web.Files),
		PathPrefix: "assets",
		Browse:     false,
	}))

	// TODO: Redirect to login page or dashboard page
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Redirect("/login")
	})

	return &Router{
		app:           app,
		authenticator: authenticator,
		sessionStore:  sessionStore,
		middleware:    middleware.New(sessionStore, authenticator),
		s3service:     service.New(&app.Config.S3Config),
	}
}
