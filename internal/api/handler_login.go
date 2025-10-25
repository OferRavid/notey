package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/OferRavid/notey/internal/auth"
	"github.com/OferRavid/notey/internal/database"
	"github.com/labstack/echo/v4"
)

// Login current user.
func (cfg *ApiConfig) handlerLogin(c echo.Context) error {
	type parameters struct {
		Username string `json:"username"`
		Password string `json:"password"`
		// ExpiresInSeconds int64  `json:"expires_in_seconds"`
	}
	type response struct {
		User
		Token string `json:"token"`
	}

	decoder := json.NewDecoder(c.Request().Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"Error": fmt.Sprintf("Couldn't decode parameters: %v", err)})
	}

	duration := time.Hour
	user, err := cfg.DbQueries.GetUserByUsername(c.Request().Context(), params.Username)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"Error": fmt.Sprintf("Incorrect username: %v", err)})
	}

	err = auth.CheckPasswordHash(user.HashedPassword, params.Password)
	if err != nil {
		return c.JSON(http.StatusForbidden, echo.Map{"Error": fmt.Sprintf("Incorrect password: %v", err)})
	}

	token, err := auth.MakeJWT(user.ID, cfg.Secret, duration)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"Error": fmt.Sprintf("Couldn't create access JWT: %v", err)})
	}

	refresh_token, err := auth.MakeRefreshToken()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"Error": fmt.Sprintf("Couldn't create refresh token: %v", err)})
	}
	refreshToken, err := cfg.DbQueries.CreateRefreshToken(
		c.Request().Context(),
		database.CreateRefreshTokenParams{
			Token:     refresh_token,
			UserID:    user.ID,
			ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
		},
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"Error": fmt.Sprintf("Failed to create refresh token in database: %v", err)})
	}

	cookie := &http.Cookie{
		Name:     RefreshTokenCookieName,
		Value:    refreshToken.Token,
		Path:     "/",
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HttpOnly: true,
		// Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	c.SetCookie(cookie)

	return c.JSON(
		http.StatusOK,
		response{
			User: User{
				ID:        user.ID,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.CreatedAt,
				Username:  user.Username,
				Email:     user.Email,
			},
			Token: token,
		},
	)
}
