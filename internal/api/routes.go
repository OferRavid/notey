package api

import (
	"github.com/labstack/echo/v4"
)

func (cfg *ApiConfig) RegisterRoutes(e *echo.Echo) {
	e.POST("/api/users", createUser)
	e.POST("/login", login)

}

func createUser(c echo.Context) error {
	return nil
}

func login(c echo.Context) error {
	return nil
}
