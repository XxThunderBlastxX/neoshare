package session

import (
	"github.com/gofiber/fiber/v2/middleware/session"
)

type Session struct {
	*session.Store
}

func New() *Session {
	sess := session.New(session.Config{
		CookieSecure:   true,
		CookieHTTPOnly: false,
	})
	sess.RegisterType(map[string]any{})

	return &Session{Store: sess}
}
