package handler

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"golang.org/x/oauth2"

	"github.com/XxThunderBlastxX/neoshare/cmd/web/page"
	"github.com/XxThunderBlastxX/neoshare/internal/auth"
	"github.com/XxThunderBlastxX/neoshare/internal/session"
)

type authHandler struct {
	session *session.Session
	auth    *auth.Authenticator
}

type AuthHandler interface {
	LoginHandler() fiber.Handler
	CallbackHandler() fiber.Handler
	LogoutHandler() fiber.Handler

	LoginView() fiber.Handler
}

func NewAuthHandler(sess *session.Session, auth *auth.Authenticator) AuthHandler {
	return &authHandler{
		session: sess,
		auth:    auth,
	}
}

func (a *authHandler) LoginView() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		render := adaptor.HTTPHandler(templ.Handler(page.HomePage()))

		return render(ctx)
	}
}

func (a *authHandler) LoginHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		state, err := generateRandomState()
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		sess, _ := a.session.Get(ctx)
		sess.Set("state", state)
		if err := sess.Save(); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		opt := oauth2.SetAuthURLParam("audience", "https://thunder.jp.auth0.com/api/v2/")
		return ctx.Redirect(a.auth.AuthCodeURL(state, opt))
	}
}

func (a *authHandler) CallbackHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		sess, _ := a.session.Get(ctx)
		if ctx.Query("state") != sess.Get("state") {
			return ctx.Status(fiber.StatusUnauthorized).SendString("Invalid state.")
		}

		code := ctx.Query("code")
		opt := oauth2.SetAuthURLParam("audience", "https://thunder.jp.auth0.com/api/v2/")
		token, err := a.auth.Exchange(ctx.Context(), code, opt)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		_, err = a.auth.VerifyIDToken(ctx.Context(), token)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		c := new(fiber.Cookie)
		c.Name = "auth_token"
		c.Value = token.AccessToken
		c.Expires = token.Expiry

		ctx.Cookie(c)

		return ctx.Redirect("/")
	}
}

func (a *authHandler) LogoutHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return nil
	}
}

// TODO: Refactor this to a separate package
func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}
