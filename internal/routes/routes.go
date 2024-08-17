package routes

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"

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
	// Setting up the logging middleware
	r.app.Use(r.middleware.StyledLogger(r.app.Config.AppEnv))

	// Serve static files
	r.app.Static("/static", "./static", fiber.Static{
		Browse: false,
	})

	// Setting favicon for the application
	r.app.Use(favicon.New(favicon.Config{
		Data:         r.app.Favicon,
		CacheControl: "public, max-age=31536000",
	}))

	// Setting up the rate limiter middleware
	r.app.Use(r.middleware.RateLimiter())

	// Redirecting the root path to the login page
	r.app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Redirect("/login")
	})

	// Registering the routes
	r.AuthRouter()
	r.DashboardRouter()
}
