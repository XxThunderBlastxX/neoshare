package middleware

import (
	"fmt"
	"log"

	contribJwt "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/sujit-baniya/flash"

	"github.com/XxThunderBlastxX/neoshare/internal/model"
)

func (m *Middleware) VerifyToken() fiber.Handler {
	jwksUri, _ := m.auth.JwksUri()

	return contribJwt.New(contribJwt.Config{
		Filter:     nil,
		JWKSetURLs: []string{jwksUri},
		SuccessHandler: func(ctx *fiber.Ctx) error {
			return ctx.Next()
		},
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			log.Println(err)
			ctx.ClearCookie(m.authCookieKey)
			errRes := model.WebResponse{
				Message:    "Authentication failed",
				StatusCode: fiber.StatusUnauthorized,
				Success:    false,
			}
			return flash.WithError(ctx, errRes.ConvertToMap()).Redirect("/login")
		},
		TokenLookup: fmt.Sprintf("cookie:%s", m.authCookieKey),
	})
}
