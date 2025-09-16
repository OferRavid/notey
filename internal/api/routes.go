package api

import (
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	e.POST("/api/users", createUsers, Middleware())
	e.POST("/login", login)
}

func createUsers(c echo.Context) error {
	return nil
}

func login(c echo.Context) error {
	return nil
}
