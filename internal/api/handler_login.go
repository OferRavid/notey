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
// Gets user by email provided in the request then generates jwt token.
func (cfg *ApiConfig) handlerLogin(c echo.Context) error {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		// ExpiresInSeconds int64  `json:"expires_in_seconds"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(c.Request().Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"Error": fmt.Sprintf("Couldn't decode parameters: %v", err)})
	}

	duration := time.Hour
	user, err := cfg.DbQueries.GetUserByEmail(c.Request().Context(), params.Email)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"Error": fmt.Sprintf("Incorrect email or password: %v", err)})
	}

	err = auth.CheckPasswordHash(user.HashedPassword, params.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"Error": fmt.Sprintf("Incorrect password: %v", err)})
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

	return c.JSON(
		http.StatusOK,
		response{
			User: User{
				ID:        user.ID,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.CreatedAt,
				Email:     user.Email,
			},
			Token:        token,
			RefreshToken: refreshToken.Token,
		},
	)
}
