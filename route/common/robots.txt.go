package common

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func Robots(c echo.Context) error {
	return c.NoContent(http.StatusInternalServerError)
}
