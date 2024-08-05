package routes

import (
	"net/http"

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
	return &Router{
		app:           app,
		authenticator: app.Authenticator,
		sessionStore:  app.Session,
		middleware:    middleware.New(app.Session, app.Authenticator),
		s3service:     service.New(&app.Config.S3Config),
	}
}

func (r *Router) RegisterRoutes() {
	r.app.Use("/assets", filesystem.New(filesystem.Config{
		Root:       http.FS(web.Files),
		PathPrefix: "assets",
		Browse:     false,
	}))

	r.AuthRouter()

	r.app.Use(r.middleware.VerifyToken())

	r.DashboardRouter()
}
