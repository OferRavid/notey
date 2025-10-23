package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Redirects to the 404 page for unmatched routes.
func handlerNotFoundError(err error, c echo.Context) {
	he, ok := err.(*echo.HTTPError)
	if ok && he.Code == http.StatusNotFound {
		// Log the error for server-side debugging
		c.Logger().Errorf("404 Not Found: %s", c.Request().RequestURI)

		// Render the custom 404 page
		c.File("static/not-found/not-found.html")
	}
	// For other errors, use the default Echo handler
	c.Echo().DefaultHTTPErrorHandler(err, c)
}
