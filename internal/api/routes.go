package api

import (
	"github.com/labstack/echo/v4"
)

func (cfg *ApiConfig) RegisterRoutes(e *echo.Echo) {
	e.POST("/api/users", cfg.handlerCreateUser)
	e.POST("/api/login", cfg.handlerLogin)

	e.POST("/api/refresh", cfg.handlerRefreshToken)
	e.POST("/api/revoke", cfg.handlerRevokeToken)
	e.POST("/admin/reset", cfg.handlerReset)
}
