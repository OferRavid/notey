package api

import (
	"net/http"

	"github.com/OferRavid/notes-app/internal/auth"
	"github.com/labstack/echo/v4"
)

var jwtSecret = []byte("your_secret_key")

func Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")
			if token == "" {
				return c.JSON(http.StatusUnauthorized, "missing token")
			}
			_, err := auth.ValidateJWT(token, string(jwtSecret))
			if err != nil {
				return c.JSON(http.StatusUnauthorized, "invalid token")
			}
			return next(c)
		}
	}
}
