package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Logout current user.
func (cfg *ApiConfig) handlerLogout(c echo.Context) error {

	user_id := c.Get("user_id").(uuid.UUID)
	err := cfg.DbQueries.DeleteTokensByUserID(c.Request().Context(), user_id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"Error": "Failed to delete refresh tokens"})
	}

	return c.JSON(http.StatusNoContent, nil)
}
