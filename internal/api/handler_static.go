package api

import (
	"path/filepath"

	"github.com/labstack/echo/v4"
)

func (cfg *ApiConfig) ServeStaticFiles(c echo.Context) error {
	// Increment the hit counter
	cfg.FileserverHits.Add(1)

	// Strip the prefix /app and serve the requested file
	filePath := filepath.Join(cfg.FilepathRoot, c.Request().URL.Path[len("/app/"):])
	staticPath := filepath.Join(cfg.StaticDir, filePath)
	return c.File(staticPath)
}
