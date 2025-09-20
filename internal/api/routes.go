package api

import (
	"github.com/labstack/echo/v4"
)

func (cfg *ApiConfig) RegisterRoutes(e *echo.Echo) {
	e.POST("/api/users", cfg.handlerCreateUser)
	e.POST("/api/login", cfg.handlerLogin)

}
