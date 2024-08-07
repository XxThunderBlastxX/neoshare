package handler

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/sujit-baniya/flash"
	"golang.org/x/oauth2"

	"github.com/XxThunderBlastxX/neoshare/cmd/web/page"
	"github.com/XxThunderBlastxX/neoshare/internal/auth"
	"github.com/XxThunderBlastxX/neoshare/internal/model"
	"github.com/XxThunderBlastxX/neoshare/internal/session"
)

type authHandler struct {
	session       *session.Session
	auth          *auth.Authenticator
	authCookieKey string
}

type AuthHandler interface {
	LoginHandler() fiber.Handler
	CallbackHandler() fiber.Handler
	LogoutHandler() fiber.Handler

	LoginView() fiber.Handler
}

func NewAuthHandler(sess *session.Session, auth *auth.Authenticator) AuthHandler {
	return &authHandler{
		session:       sess,
		auth:          auth,
		authCookieKey: "auth_token",
	}
}

func (a *authHandler) LoginView() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		res := flash.Get(ctx)
		if len(res) != 0 {
			var resData model.WebResponse
			resData.ConvertToStruct(res)
			render := adaptor.HTTPHandler(templ.Handler(page.AuthPage(resData)))
			return render(ctx)
		}

		authCookie := ctx.Cookies(a.authCookieKey)
		if authCookie != "" {
			return ctx.Redirect("/dashboard")
		}

		render := adaptor.HTTPHandler(templ.Handler(page.AuthPage()))

		return render(ctx)
	}
}

func (a *authHandler) LoginHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		state, err := generateRandomState()
		if err != nil {
			errRes := model.WebResponse{
				Error:      err.Error(),
				StatusCode: fiber.StatusInternalServerError,
				Success:    false,
			}
			return flash.WithError(ctx, errRes.ConvertToMap()).Redirect("/login")
		}

		sess, _ := a.session.Get(ctx)
		sess.Set("state", state)
		if err := sess.Save(); err != nil {
			errRes := model.WebResponse{
				Error:      err.Error(),
				StatusCode: fiber.StatusInternalServerError,
				Success:    false,
			}
			return flash.WithError(ctx, errRes.ConvertToMap()).Redirect("/login")
		}

		opt := oauth2.SetAuthURLParam("audience", "https://thunder.jp.auth0.com/api/v2/")
		return ctx.Redirect(a.auth.AuthCodeURL(state, opt))
	}
}

func (a *authHandler) CallbackHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		sess, _ := a.session.Get(ctx)
		if ctx.Query("state") != sess.Get("state") {
			errRes := model.WebResponse{
				Error:      "Application state does not match",
				StatusCode: fiber.StatusInternalServerError,
				Success:    false,
			}
			return flash.WithError(ctx, errRes.ConvertToMap()).Redirect("/login")
		}

		code := ctx.Query("code")
		opt := oauth2.SetAuthURLParam("audience", "https://thunder.jp.auth0.com/api/v2/")
		token, err := a.auth.Exchange(ctx.Context(), code, opt)
		if err != nil {
			errRes := model.WebResponse{
				Error:      err.Error(),
				StatusCode: fiber.StatusInternalServerError,
				Success:    false,
			}
			return flash.WithError(ctx, errRes.ConvertToMap()).Redirect("/login")
		}

		_, err = a.auth.VerifyIDToken(ctx.Context(), token)
		if err != nil {
			errRes := model.WebResponse{
				Error:      err.Error(),
				StatusCode: fiber.StatusInternalServerError,
				Success:    false,
			}
			return flash.WithError(ctx, errRes.ConvertToMap()).Redirect("/login")
		}

		c := new(fiber.Cookie)
		c.Name = a.authCookieKey
		c.Value = token.AccessToken
		c.Expires = token.Expiry
		c.Secure = true

		ctx.Cookie(c)

		return ctx.Redirect("/dashboard")
	}
}

func (a *authHandler) LogoutHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.ClearCookie(a.authCookieKey)
		return ctx.Redirect(a.auth.RevokeToken())
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
