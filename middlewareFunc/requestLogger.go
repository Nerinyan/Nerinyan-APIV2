package middlewareFunc

import (
	"github.com/labstack/echo/v4"
	"github.com/pterm/pterm"
	"time"
)

func RequestConsolLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			go pterm.Info.Printfln(
				// 2023-05-13T00:01:02Z | GET    | 123.123.123.123 | api.nerinyan.moe | /search
				"%-s | %6s | %15s | %-30s | %-s",
				time.Now().UTC().Format(time.RFC3339), c.Request().Method, c.RealIP(), c.Request().Host, c.Request().URL.Path,
			)
			return next(c)
		}
	}
}
