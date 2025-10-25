package api

import (
	"errors"
	"net/http"

	"github.com/OferRavid/notey/internal/auth"
	"github.com/labstack/echo/v4"
)

// Middleware validates JWT and sets user_id in Echo context.
func (cfg *ApiConfig) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := auth.GetBearerToken(c)
			if token == "" || err != nil {
				if errors.As(err, &auth.ErrNoAuthHeaderIncluded) {
					return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
				}
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid Authorization header format"})
			}

			user_id, err := auth.ValidateJWT(token, cfg.Secret)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid token"})
			}

			c.Set("user_id", user_id)
			return next(c)
		}
	}
}
