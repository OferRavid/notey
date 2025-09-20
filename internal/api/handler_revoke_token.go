package api

import (
	"net/http"
	"time"

	"github.com/OferRavid/notey/internal/auth"
	"github.com/labstack/echo/v4"
)

func (cfg *ApiConfig) handlerRevokeToken(c echo.Context) error {
	refresh_token, err := auth.GetBearerToken(c.Request().Header)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"Error": "Missing bearer token in headers"})
	}

	refreshToken, err := cfg.DbQueries.GetRefreshTokenByToken(c.Request().Context(), refresh_token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"Error": "Refresh token doesn't exist"})
	}
	if time.Now().After(refreshToken.ExpiresAt) || refreshToken.RevokedAt.Valid {
		return c.JSON(http.StatusUnauthorized, echo.Map{"Error": "Refresh token already expired"})
	}

	err = cfg.DbQueries.UpdateRefreshToken(c.Request().Context(), refreshToken.Token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"Error": "Failed to update record"})
	}

	return c.JSON(http.StatusNoContent, nil)
}
