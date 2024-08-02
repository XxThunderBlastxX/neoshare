package middleware

import (
	"github.com/XxThunderBlastxX/neoshare/internal/auth"
	"github.com/XxThunderBlastxX/neoshare/internal/session"
)

type Middleware struct {
	session *session.Session
	auth    *auth.Authenticator
}

func New(s *session.Session, a *auth.Authenticator) *Middleware {
	return &Middleware{
		session: s,
		auth:    a,
	}
}
