package routes

import (
	"github.com/XxThunderBlastxX/neoshare/internal/handler"
)

func (r *Router) AuthRouter() {
	api := r.app.Group("/api")
	view := r.app.Group("/")

	authHandler := handler.NewAuthHandler(
		r.sessionStore,
		r.authenticator,
		r.app.Config.Auth.Audience,
		r.authCookieKey,
	)

	// API Routes
	api.Get("/login", authHandler.LoginHandler())
	api.Get("/callback", authHandler.CallbackHandler())
	api.Get("/logout", authHandler.LogoutHandler())
	api.Get("/logout/callback", authHandler.LogoutCallbackHandler())

	// View Routes
	view.Get("/login", authHandler.LoginView())
}
