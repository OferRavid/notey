package api

import (
	"github.com/labstack/echo/v4"
)

func (cfg *ApiConfig) RegisterRoutes(e *echo.Echo) {
	// 404 - Not Found error handler
	e.HTTPErrorHandler = handlerNotFoundError

	// Routes to handle users
	e.POST("/api/users", cfg.handlerCreateUser)
	e.POST("/api/login", cfg.handlerLogin)
	e.PUT("/api/users", cfg.handlerUpdateUserData, cfg.Middleware())
	e.DELETE("/api/users/:userID", cfg.handlerDeleteUser, cfg.Middleware())

	// Routes to hanlde refresh token
	e.POST("/api/refresh", cfg.handlerRefreshToken)
	e.POST("/api/revoke", cfg.handlerRevokeToken)

	// Routes to handle notes
	e.GET("/api/notes", cfg.handlerRetrieveNotes, cfg.Middleware())
	e.GET("/api/notes/:noteID", cfg.handlerGetNoteByID, cfg.Middleware())
	e.POST("/api/notes", cfg.handlerCreateNote, cfg.Middleware())
	e.PUT("/api/notes/:noteID", cfg.handlerUpdateNote, cfg.Middleware())
	e.DELETE("/api/notes/:noteID", cfg.handlerDeleteNote, cfg.Middleware())

	// Admin only routes
	e.POST("/admin/reset", cfg.handlerReset)
	e.GET("/admin/metrics", cfg.handlerMetrics)
}
