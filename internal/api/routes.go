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

	// Routes to handle notes
	e.GET("/notes", cfg.handlerRetrieveNotes, cfg.Middleware())
	e.GET("/notes:noteID", cfg.handlerGetNoteByID, cfg.Middleware())
	e.POST("/notes", cfg.handlerCreateNote, cfg.Middleware())
	// e.PUT("/notes/:id", cfg.handlerUpdateNote, cfg.Middleware())
	e.DELETE("/notes/:noteID", cfg.handlerDeleteNote, cfg.Middleware())

	// Admin only routes
	e.POST("/admin/reset", cfg.handlerReset)
	e.GET("/admin/metrics", cfg.handlerMetrics)
}
