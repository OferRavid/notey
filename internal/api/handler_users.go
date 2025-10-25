package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/OferRavid/notey/internal/auth"
	"github.com/OferRavid/notey/internal/database"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type parameters struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Create new user in database.
func (cfg *ApiConfig) handlerCreateUser(c echo.Context) error {
	hashedPassword, email, username, statusCode, err := getHashedPasswordAndEmail(c)
	if err != nil {
		return c.JSON(statusCode, echo.Map{"Error": err})
	}

	user, err := cfg.DbQueries.CreateUser(
		c.Request().Context(),
		database.CreateUserParams{
			Username:       username,
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
			Username:  user.Username,
			Email:     user.Email,
		},
	)
}

// Handles updating username and/or password.
func (cfg *ApiConfig) handlerUpdateUserData(c echo.Context) error {
	user_id := c.Get("user_id").(uuid.UUID)

	hashedPassword, email, username, statusCode, err := getHashedPasswordAndEmail(c)
	if err != nil {
		return c.JSON(statusCode, echo.Map{"Error": err})
	}

	userFromEmail, err := cfg.DbQueries.GetUserByEmail(c.Request().Context(), email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"Error": "Couldn't retrieve user"})
	}

	if user_id != userFromEmail.ID {
		return c.JSON(http.StatusForbidden, echo.Map{"Error": "Unauthorized to edit user params"})
	}

	user, err := cfg.DbQueries.UpdateUser(
		c.Request().Context(),
		database.UpdateUserParams{
			Username:       username,
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
			Username:  user.Username,
			Email:     user.Email,
		},
	)
}

// Handles user removal.
func (cfg *ApiConfig) handlerDeleteUser(c echo.Context) error {
	user_id := c.Get("user_id").(uuid.UUID)

	userIDStr := c.Param("userID")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"Error": "Failed to parse userID"})
	}

	if userID != user_id {
		return c.JSON(http.StatusForbidden, echo.Map{"Error": "Unauthorized to remove user"})
	}

	err = cfg.DbQueries.RemoveUser(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"Error": "Failed to remove user"})
	}

	return c.JSON(http.StatusNoContent, nil)
}

// Decodes email and password from request hashes the password then returns the email and hashed password.
func getHashedPasswordAndEmail(c echo.Context) (string, string, string, int, error) {
	decoder := json.NewDecoder(c.Request().Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		return "", "", "", http.StatusInternalServerError, errors.New("couldn't decode parameters")
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		return "", "", "", http.StatusBadRequest, errors.New("couldn't create hashed password")
	}

	return hashedPassword, params.Email, params.Username, 0, nil
}
