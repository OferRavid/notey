package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/OferRavid/notey/internal/auth"
	"github.com/OferRavid/notey/internal/database"
	"github.com/labstack/echo/v4"
)

func (cfg *ApiConfig) handlerRefreshToken(c echo.Context) error {
	cookie, err := c.Cookie(RefreshTokenCookieName)
	if err != nil {
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

	if time.Now().After(refreshToken.ExpiresAt) {
		if refreshToken.RevokedAt.Valid {
			type errorResponse struct {
				Error   string `json:"error"`
				Revoked bool   `json:"revoked"`
			}
			return c.JSON(http.StatusForbidden, errorResponse{
				Error:   "Refresh token expired",
				Revoked: refreshToken.RevokedAt.Valid,
			})
		}
		err = cfg.DbQueries.RevokeRefreshToken(c.Request().Context(), refreshToken.Token)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"Error": "Failed to revoke token"})
		}

		return c.JSON(http.StatusForbidden, echo.Map{"Error": "Refresh token expired"})
	}

	err = cfg.DbQueries.RevokeRefreshToken(c.Request().Context(), refreshToken.Token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"Error": "Failed to revoke old token"})
	}

	newRefreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"Error": "Couldn't create new refresh token"})
	}
	_, err = cfg.DbQueries.CreateRefreshToken(
		c.Request().Context(),
		database.CreateRefreshTokenParams{
			Token:     newRefreshToken,
			UserID:    refreshToken.UserID,
			ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
		},
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"Error": fmt.Sprintf("Failed to create refresh token in database: %v", err)})
	}

	newCookie := &http.Cookie{
		Name:     RefreshTokenCookieName,
		Value:    newRefreshToken,
		Path:     "/",
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HttpOnly: true,
		// Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	c.SetCookie(newCookie)

	jwtToken, err := auth.MakeJWT(refreshToken.UserID, cfg.Secret, time.Hour)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"Error": "Couldn't create access JWT"})
	}

	type response struct {
		Token string `json:"token"`
	}

	return c.JSON(
		http.StatusOK,
		response{
			Token: jwtToken,
		},
	)
}
