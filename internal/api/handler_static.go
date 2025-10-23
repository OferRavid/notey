package api

import (
	"path/filepath"

	"github.com/labstack/echo/v4"
)

func (cfg *ApiConfig) ServeStaticFiles(c echo.Context) error {
	// Increment the hit counter
	cfg.FileserverHits.Add(1)
	cfg.PageVisitsGauge.Inc()

	// Strip the prefix /app and serve the requested file
	filePath := c.Request().URL.Path[len("/app/"):]
	staticPath := filepath.Join(cfg.StaticDir, filePath)
	return c.File(staticPath)
}
