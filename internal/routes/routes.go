package routes

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2/middleware/filesystem"

	"github.com/XxThunderBlastxX/neoshare/cmd/web"
	"github.com/XxThunderBlastxX/neoshare/internal/auth"
	"github.com/XxThunderBlastxX/neoshare/internal/middleware"
	"github.com/XxThunderBlastxX/neoshare/internal/repository"
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
	fileService   service.FileService
	authCookieKey string
}

func New(app *server.Server) *Router {
	authCookieKey := "auth_token"
	s3Service := service.New(&app.Config.S3Config)
	fileService := service.NewFileService(context.Background(), repository.New(app.Db), app.Db, s3Service)

	return &Router{
		app:           app,
		authenticator: app.Authenticator,
		sessionStore:  app.Session,
		middleware:    middleware.New(app.Session, app.Authenticator, authCookieKey, app.Config.Auth.Domain),
		s3service:     s3Service,
		fileService:   fileService,
		authCookieKey: authCookieKey,
	}
}

func (r *Router) RegisterRoutes() {
	r.app.Use("/assets", filesystem.New(filesystem.Config{
		Root:       http.FS(web.Files),
		PathPrefix: "assets",
		Browse:     false,
	}))

	r.app.Use(r.middleware.StyledLogger(r.app.Config.AppEnv))

	r.AuthRouter()
	r.DashboardRouter()
}
