package middleware

import (
	"github.com/XxThunderBlastxX/neoshare/internal/auth"
	"github.com/XxThunderBlastxX/neoshare/internal/session"
)

type Middleware struct {
	session       *session.Session
	auth          *auth.Authenticator
	authDomain    string
	authCookieKey string
}

func New(s *session.Session, a *auth.Authenticator, authCookieKey, authDomain string) *Middleware {
	return &Middleware{
		session:       s,
		auth:          a,
		authCookieKey: authCookieKey,
		authDomain:    authDomain,
	}
}
