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
	requestPath := c.Request().URL.Path[len("/app/"):]
	filePath := ""

	switch requestPath {
	case "":
		break
	case "contact":
		filePath = "contact/contact.html"
	case "login":
		filePath = "login/login.html"
	case "400":
		filePath = "bad-request/bad-request.html"
	case "401":
		filePath = "authorization-error/unauthorized.html"
	case "403":
		filePath = "access-denied/access-denied.html"
	case "404":
		filePath = "not-found/not-found.html"
	case "500":
		filePath = "server-error/server-error.html"
	case "privacy-policy":
		filePath = "policy/privacy.html"
	case "terms":
		filePath = "policy/terms.html"
	case "notes":
		filePath = "notes/notes.html"
	case "note":
		filePath = "notes/note.html"
	}
	staticPath := filepath.Join(cfg.StaticDir, filePath)

	return c.File(staticPath)
}
