package common

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func Health(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}
