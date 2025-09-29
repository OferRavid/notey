package api

import (
	"fmt"
	"net/http"

	"github.com/OferRavid/notey/internal/auth"
	"github.com/labstack/echo/v4"
)

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
