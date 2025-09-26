package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/OferRavid/notey/internal/auth"
	"github.com/OferRavid/notey/internal/database"
	"github.com/labstack/echo/v4"
)

type parameters struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Create new user in database.
// Gets email and password from request and hashes the password.
// Creates new user record in database using email and hashedPassword.
func (cfg *ApiConfig) handlerCreateUser(c echo.Context) error {
	hashedPassword, email, statusCode, err := getHashedPasswordAndEmail(c)
	if err != nil {
		return c.JSON(statusCode, echo.Map{"Error": err})
	}

	user, err := cfg.DbQueries.CreateUser(
		c.Request().Context(),
		database.CreateUserParams{
			Email:          email,
			HashedPassword: hashedPassword,
		},
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"Error": "Couldn't create user"})
	}

	return c.JSON(
		http.StatusCreated,
		User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	)
}

// Handles updating user's email address and/or password.
// Gets email and password from request and hashes the password.
// Uses the token given to find the user_id of the user to update and makes the update.
func (cfg *ApiConfig) handlerUpdateUserData(c echo.Context) error {
	token, err := auth.GetBearerToken(c.Request().Header)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"Error": "Missing or malformed token"})
	}

	user_id, err := auth.ValidateJWT(token, cfg.Secret)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"Error": "Invalid bearerToken for user"})
	}

	hashedPassword, email, statusCode, err := getHashedPasswordAndEmail(c)
	if err != nil {
		return c.JSON(statusCode, echo.Map{"Error": err})
	}

	user, err := cfg.DbQueries.UpdateUser(
		c.Request().Context(),
		database.UpdateUserParams{
			Email:          email,
			HashedPassword: hashedPassword,
			ID:             user_id,
		},
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"Error": "Couldn't update user"})
	}

	return c.JSON(http.StatusOK,
		User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	)
}

// Decodes email and password from request hashes the password then returns the email and hashed password.
func getHashedPasswordAndEmail(c echo.Context) (string, string, int, error) {
	decoder := json.NewDecoder(c.Request().Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		return "", "", http.StatusInternalServerError, errors.New("couldn't decode parameters")
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		return "", "", http.StatusBadRequest, errors.New("couldn't create hashed password")
	}

	return hashedPassword, params.Email, 0, nil
}
