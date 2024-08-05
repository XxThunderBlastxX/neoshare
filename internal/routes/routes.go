package routes

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"

	"github.com/XxThunderBlastxX/neoshare/cmd/web"
	"github.com/XxThunderBlastxX/neoshare/internal/auth"
	"github.com/XxThunderBlastxX/neoshare/internal/middleware"
	"github.com/XxThunderBlastxX/neoshare/internal/server"
	"github.com/XxThunderBlastxX/neoshare/internal/service"
	"github.com/XxThunderBlastxX/neoshare/internal/session"
)

type Router struct {
	app           *server.Server
	authenticator *auth.Authenticator
	sessionStore  *session.Session
	middleware    *middleware.Middleware
	s3service     service.S3Service
}

func New(app *server.Server) *Router {
	app.Use("/assets", filesystem.New(filesystem.Config{
		Root:       http.FS(web.Files),
		PathPrefix: "assets",
		Browse:     false,
	}))

	// TODO: Redirect to login page or dashboard.templ page
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Redirect("/dashboard")
	})

	return &Router{
		app:           app,
		authenticator: app.Authenticator,
		sessionStore:  app.Session,
		middleware:    middleware.New(app.Session, app.Authenticator),
		s3service:     service.New(&app.Config.S3Config),
	}
}

func (r *Router) RegisterRoutes() {
	r.AuthRouter()
	r.DashboardRouter()
}
