package api

import (
	"errors"
	"net/http"

	"github.com/OferRavid/notes-app/internal/auth"
	"github.com/labstack/echo/v4"
)

var jwtSecret = []byte("your_secret_key")

// Middleware validates JWT and sets user_id in Echo context.
func Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := auth.GetBearerToken(c.Request().Header)
			if token == "" || err != nil {
				if errors.As(err, &auth.ErrNoAuthHeaderIncluded) {
					return c.JSON(http.StatusUnauthorized, echo.Map{"error": "missing Authorization header"})
				}
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid Authorization header format"})
			}

			user_id, err := auth.ValidateJWT(token, string(jwtSecret))
			if err != nil {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid token"})
			}

			c.Set("user_id", user_id)
			return next(c)
		}
	}
}
