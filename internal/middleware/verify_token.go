package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/sujit-baniya/flash"

	"github.com/XxThunderBlastxX/neoshare/internal/model"
)

func (m *Middleware) VerifyToken() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Get the JWT token from the cookie
		tokenString := ctx.Cookies(m.authCookieKey)
		if tokenString == "" {
			log.Println("No JWT token found in the cookie")
			res := model.WebResponse{
				Success:    false,
				StatusCode: fiber.StatusUnauthorized,
				Message:    "Unauthorized access",
			}
			return flash.WithError(ctx, res.ConvertToMap()).Redirect("/login")
		}

		ok, err := m.auth.VerifyUserInfo(tokenString)
		if err != nil {
			log.Println("Error verifying Access token: ", err)
			res := model.WebResponse{
				Success:    false,
				StatusCode: fiber.StatusUnauthorized,
				Message:    "Unauthorized access",
			}
			return flash.WithError(ctx, res.ConvertToMap()).Redirect("/login")
		}

		if !ok {
			log.Println("Invalid token found")
			res := model.WebResponse{
				Success:    false,
				StatusCode: fiber.StatusUnauthorized,
				Message:    "Unauthorized access",
			}
			return flash.WithError(ctx, res.ConvertToMap()).Redirect("/login")
		}

		return ctx.Next()
	}
}
