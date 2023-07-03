package middlewareFunc

import (
	"github.com/Nerinyan/Nerinyan-APIV2/db/mariadb"
	"github.com/Nerinyan/Nerinyan-APIV2/db/mariadb/entity"
	"github.com/labstack/echo/v4"
	"github.com/pterm/pterm"
	"net/http"
	"time"
)

var banList map[string]bool
var logBlocked pterm.PrefixPrinter

func init() {
	logBlocked = pterm.Error
	logBlocked.Prefix.Text = "BLOCKED"
}
func StartHandler() {
	if banList != nil {
		return
	}
	banList = map[string]bool{}

	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for ; ; <-ticker.C {
			var blacklists []entity.BlacklistEntity
			err := mariadb.Mariadb.Where("EXPIRED_AT > ?", time.Now()).Find(&blacklists).Error
			if err != nil {
				pterm.Error.WithShowLineNumber().Println(err)
				continue
			}
			l := map[string]bool{}
			for _, bl := range blacklists {
				l[bl.IPV4] = true
			}
			banList = l
		}
	}()
}

func BlackListHandler() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if banList != nil && banList[c.RealIP()] {
				go logBlocked.Printfln(
					// 2023-05-13T00:01:02Z | GET    | 123.123.123.123 | api.nerinyan.moe | /search
					"%-s | %6s | %15s | %-30s | %-s",
					time.Now().UTC().Format(time.RFC3339), c.Request().Method, c.RealIP(), c.Request().Host, c.Request().URL.Path,
				)

				return c.JSON(http.StatusForbidden, map[string]any{
					"message":    "access denied. contect 'admin@nerinyan.moe'",
					"request_id": c.Response().Header().Get("X-Request-Id"),
					"time":       entity.RFC3339(time.Now().UTC()),
				})
			}
			return next(c)
		}
	}
}
