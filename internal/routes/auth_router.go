package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/XxThunderBlastxX/neoshare/internal/handler"
)

func (r *Router) AuthRouter() {
	api := r.app.Group("/api")
	view := r.app.Group("/")

	authHandler := handler.NewAuthHandler(r.sessionStore, r.authenticator)

	// API Routes
	api.Get("/login", authHandler.LoginHandler())
	api.Get("/callback", authHandler.CallbackHandler())
	api.Get("/logout", authHandler.LogoutHandler())

	// View Routes
	view.Get("/login", authHandler.LoginView())
	view.Get("/kuchbhi", r.middleware.VerifyToken(), func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello, World!")
	})
}
