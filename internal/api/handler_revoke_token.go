package api

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (cfg *ApiConfig) handlerRevokeToken(c echo.Context) error {
	cookie, err := c.Cookie(RefreshTokenCookieName)
	if err != nil {
		// No cookie, likely expired or missing
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "missing refresh token cookie"})
	}
	refresh_token := cookie.Value

	refreshToken, err := cfg.DbQueries.GetRefreshTokenByToken(c.Request().Context(), refresh_token)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"Error": "Refresh token doesn't exist"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"Error": "Bad response from database"})
	}
	if refreshToken.RevokedAt.Valid {
		return c.JSON(http.StatusForbidden, echo.Map{"Error": "Refresh token already revoked"})
	}

	err = cfg.DbQueries.RevokeRefreshToken(c.Request().Context(), refreshToken.Token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"Error": "Failed to revoke token"})
	}

	return c.JSON(http.StatusNoContent, nil)
}
