package api

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (cfg *ApiConfig) handlerReset(c echo.Context) error {
	if cfg.Platform != "dev" {
		return c.JSON(http.StatusForbidden, echo.Map{"msg": "Forbidden"})
	}

	if err := cfg.DbQueries.DeleteUsers(c.Request().Context()); err != nil {
		log.Printf("failed to delete users: %s", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"Error": err})
	}

	return c.JSON(http.StatusOK, nil)
}
