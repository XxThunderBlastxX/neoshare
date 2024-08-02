package routes

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/XxThunderBlastxX/neoshare/cmd/web"
	pages "github.com/XxThunderBlastxX/neoshare/cmd/web/page"
	"github.com/XxThunderBlastxX/neoshare/internal/auth"
	"github.com/XxThunderBlastxX/neoshare/internal/server"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/session"
	"net/http"
)

func RegisterRoutes(app *server.Server, authenticator *auth.Authenticator, sessionStore *session.Store) {
	app.Use("/assets", filesystem.New(filesystem.Config{
		Root:       http.FS(web.Files),
		PathPrefix: "assets",
		Browse:     false,
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, X-Auth-Token",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		render := adaptor.HTTPHandler(templ.Handler(pages.HomePage()))

		return render(c)
	})

	app.Get("/login", func(c *fiber.Ctx) error {
		state, err := generateRandomState()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		sess, _ := sessionStore.Get(c)

		sess.Set("state", state)
		if err := sess.Save(); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return c.Redirect(authenticator.AuthCodeURL(state))
	})

	app.Get("/callback", func(ctx *fiber.Ctx) error {
		sess, _ := sessionStore.Get(ctx)
		if ctx.Query("state") != sess.Get("state") {
			return ctx.Status(http.StatusUnauthorized).SendString("Invalid state.")
		}

		// Exchange an authorization code for a token.
		token, err := authenticator.Exchange(ctx.Context(), ctx.Query("code"))
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).SendString(err.Error())
		}

		_, err = authenticator.VerifyIDToken(ctx.Context(), token)
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).SendString(err.Error())
		}

		fmt.Println(token.AccessToken)

		sess.Set("access_token", token.AccessToken)
		if err := sess.Save(); err != nil {
			return ctx.Status(http.StatusInternalServerError).SendString(err.Error())
		}

		return ctx.Redirect("/")
	})
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}
