package common

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func Root(c echo.Context) error {
	return c.Redirect(http.StatusPermanentRedirect, `https://nerinyan.stoplight.io/docs/nerinyan-api`)
}
