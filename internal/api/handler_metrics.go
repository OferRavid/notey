package api

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/OferRavid/notey/internal/auth"
	"github.com/labstack/echo/v4"
)

func (cfg *ApiConfig) ServeStaticFiles(c echo.Context) error {
	// Increment the hit counter
	cfg.FileserverHits.Add(1)

	// Strip the prefix /app and serve the requested file
	filePath := filepath.Join(cfg.FilepathRoot, c.Request().URL.Path[len("/app/"):])
	return c.File(filePath)
}

func (cfg *ApiConfig) handlerMetrics(c echo.Context) error {
	htmlResponse := fmt.Sprintf(
		`
<html>
  <body>
    <h1>Welcome, %s Admin</h1>
	<p>%s has been visited %d times!</p>
  </body>
</html>`,
		auth.TokenTypeAccess,
		auth.TokenTypeAccess,
		cfg.FileserverHits.Load(),
	)

	return c.HTML(http.StatusOK, htmlResponse)
}
