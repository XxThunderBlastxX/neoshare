package handler

import (
	"log"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/sujit-baniya/flash"
	"golang.org/x/oauth2"

	"github.com/XxThunderBlastxX/neoshare/cmd/web/page"
	"github.com/XxThunderBlastxX/neoshare/internal/auth"
	"github.com/XxThunderBlastxX/neoshare/internal/model"
	"github.com/XxThunderBlastxX/neoshare/internal/session"
	"github.com/XxThunderBlastxX/neoshare/internal/utils"
)

type authHandler struct {
	session         *session.Session
	auth            *auth.Authenticator
	authCodeOptions oauth2.AuthCodeOption
	authCookieKey   string
}

type AuthHandler interface {
	LoginHandler() fiber.Handler
	CallbackHandler() fiber.Handler
	LogoutHandler() fiber.Handler

	LoginView() fiber.Handler
	LogoutCallbackHandler() fiber.Handler
}

func NewAuthHandler(sess *session.Session, auth *auth.Authenticator, authAudience, authCookieKey string) AuthHandler {
	return &authHandler{
		session:         sess,
		auth:            auth,
		authCookieKey:   authCookieKey,
		authCodeOptions: oauth2.SetAuthURLParam("audience", authAudience),
	}
}

func (a *authHandler) LoginView() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Checks if there is any flash message as error
		res := flash.Get(ctx)
		if len(res) != 0 {
			var resData model.WebResponse
			resData.ConvertToStruct(res)
			render := adaptor.HTTPHandler(templ.Handler(page.AuthPage(resData)))
			return render(ctx)
		}

		// Checks if the user is already authenticated then redirects to the dashboard
		authCookie := ctx.Cookies(a.authCookieKey)
		if authCookie != "" {
			return ctx.Redirect("/dashboard")
		}

		// Rendering the login view page
		render := adaptor.HTTPHandler(templ.Handler(page.AuthPage()))
		return render(ctx)
	}
}

func (a *authHandler) LoginHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		state, err := utils.GenerateRandomState()
		if err != nil {
			errRes := model.WebResponse{
				Message:    err.Error(),
				StatusCode: fiber.StatusInternalServerError,
				Success:    false,
			}
			return flash.WithError(ctx, errRes.ConvertToMap()).Redirect("/login")
		}

		// Generate PKCE code verifier and challenge
		codeVerifier := utils.GenerateCodeVerifier()
		codeChallenge := utils.GenerateCodeChallenge(codeVerifier)

		sess, _ := a.session.Get(ctx)
		sess.Set("state", state)
		sess.Set("code_verifier", codeVerifier)
		if err := sess.Save(); err != nil {
			errRes := model.WebResponse{
				Message:    err.Error(),
				StatusCode: fiber.StatusInternalServerError,
				Success:    false,
			}
			return flash.WithError(ctx, errRes.ConvertToMap()).Redirect("/login")
		}

		// Add PKCE parameters to the authorization URL
		opts := []oauth2.AuthCodeOption{
			oauth2.SetAuthURLParam("code_challenge", codeChallenge),
			oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		}
		opts = append(opts, a.authCodeOptions)

		// Redirects to the OAuth2 provider consent page
		return ctx.Redirect(a.auth.AuthCodeURL(state, opts...))
	}
}

func (a *authHandler) CallbackHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		sess, _ := a.session.Get(ctx)
		if ctx.Query("state") != sess.Get("state") {
			log.Println("Application state does not match")
			errRes := model.WebResponse{
				Message:    "Application state does not match",
				StatusCode: fiber.StatusInternalServerError,
				Success:    false,
			}
			return flash.WithError(ctx, errRes.ConvertToMap()).Redirect("/login")
		}

		code := ctx.Query("code")
		codeVerifier := sess.Get("code_verifier").(string)

		// Add PKCE code verifier to the token exchange
		opts := []oauth2.AuthCodeOption{
			oauth2.SetAuthURLParam("code_verifier", codeVerifier),
		}
		opts = append(opts, a.authCodeOptions)

		// Converting the authorization code to token
		token, err := a.auth.Exchange(ctx.Context(), code, opts...)
		if err != nil {
			log.Printf("Failed to exchange the code: %v", err)
			errRes := model.WebResponse{
				Message:    err.Error(),
				StatusCode: fiber.StatusInternalServerError,
				Success:    false,
			}
			return flash.WithError(ctx, errRes.ConvertToMap()).Redirect("/login")
		}

		// Verifies the ID token
		_, err = a.auth.VerifyIDToken(ctx.Context(), token)
		if err != nil {
			log.Printf("Failed to verify the ID token: %v", err)
			errRes := model.WebResponse{
				Message:    err.Error(),
				StatusCode: fiber.StatusInternalServerError,
				Success:    false,
			}
			return flash.WithError(ctx, errRes.ConvertToMap()).Redirect("/login")
		}

		// Creates a new cookie with the access token
		c := new(fiber.Cookie)
		c.Name = a.authCookieKey
		c.Value = token.AccessToken
		c.Expires = token.Expiry
		c.Secure = true

		// Sets the cookie in the response
		ctx.Cookie(c)

		// Redirects to the dashboard
		return ctx.Redirect("/dashboard")
	}
}

func (a *authHandler) LogoutHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authCookie := ctx.Cookies(a.authCookieKey)
		return ctx.Redirect(a.auth.LogoutURL(authCookie))
	}
}

func (a *authHandler) LogoutCallbackHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Deletes the cookie
		ctx.ClearCookie() // TODO: Check if this is the correct way to delete the cookie

		// Redirects to the login page
		return ctx.Redirect("/login")
	}
}
