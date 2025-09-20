package api

import (
	"encoding/json"
	"fmt"
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
// Gets email and password from request and hashes the password
// Creates new user record in database using email and hashedPassword
func (cfg *ApiConfig) handlerCreateUser(c echo.Context) error {
	hashedPassword, email, err := getHashedPasswordAndEmail(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"Error": err})
	}

	user, err := cfg.DbQueries.CreateUser(
		c.Request().Context(),
		database.CreateUserParams{
			Email:          email,
			HashedPassword: hashedPassword,
		},
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Couldn't create user:; %v", err))
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

// Decodes email and password from request hashes the password then returns the email and hashed password.
func getHashedPasswordAndEmail(c echo.Context) (string, string, error) {
	decoder := json.NewDecoder(c.Request().Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		return "", "", c.JSON(http.StatusInternalServerError, echo.Map{"Error": fmt.Sprintf("Couldn't decode parameters: %v", err)})
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		return "", "", c.JSON(http.StatusInternalServerError, echo.Map{"Error": fmt.Sprintf("Couldn't create hashed password: %v", err)})
	}

	return hashedPassword, params.Email, nil
}
