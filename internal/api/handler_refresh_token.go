package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/OferRavid/notey/internal/auth"
	"github.com/labstack/echo/v4"
)

func (cfg *ApiConfig) handlerRefreshToken(c echo.Context) error {
	refresh_token, err := auth.GetBearerToken(c.Request().Header)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"Error": "Missing bearer token in headers"})
	}

	refreshToken, err := cfg.DbQueries.GetRefreshTokenByToken(c.Request().Context(), refresh_token)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"Error": "Refresh token doesn't exist"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"Error": "Bad response from database"})
	}

	type errorResponse struct {
		Error   string `json:"error"`
		Revoked bool   `json:"revoked"`
	}
	if time.Now().After(refreshToken.ExpiresAt) || refreshToken.RevokedAt.Valid {
		return c.JSON(http.StatusUnauthorized, errorResponse{
			Error:   "Refresh token expired",
			Revoked: refreshToken.RevokedAt.Valid,
		})
	}

	token, err := auth.MakeJWT(refreshToken.UserID, cfg.Secret, time.Hour)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"Error": "Couldn't create access JWT"})
	}

	type response struct {
		Token string `json:"token"`
	}

	return c.JSON(
		http.StatusOK,
		response{
			Token: token,
		},
	)
}
