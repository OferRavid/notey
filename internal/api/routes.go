package api

import (
	"github.com/labstack/echo/v4"
)

func (cfg *ApiConfig) RegisterRoutes(e *echo.Echo) {
	// Routes to handle users
	e.POST("/api/users", cfg.handlerCreateUser)
	e.POST("/api/login", cfg.handlerLogin)

	// Routes to hanlde refresh token
	e.POST("/api/refresh", cfg.handlerRefreshToken)
	e.POST("/api/revoke", cfg.handlerRevokeToken)

	// Admin only routes
	e.POST("/admin/reset", cfg.handlerReset)
	e.GET("/admin/metrics", cfg.handlerMetrics)
}
