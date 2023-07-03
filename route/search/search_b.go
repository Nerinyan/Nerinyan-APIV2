package search

import (
	"github.com/Nerinyan/Nerinyan-APIV2/db/mariadb"
	"github.com/Nerinyan/Nerinyan-APIV2/db/mariadb/entity"
	"github.com/labstack/echo/v4"
	"github.com/pterm/pterm"
	"net/http"
)

func SearchB(c echo.Context) (err error) { // http://127.0.0.1/search/s/123456
	query := mariadb.Mariadb.Select("M.*").Table("BEATMAP M").
		Where("M.DELETED_AT IS NULL").
		Where("M.BEATMAP_ID = ?", c.Param("mapId"))

	var maps entity.BanchoBeatmapEntity
	if err = query.Find(&maps).Error; err != nil {
		pterm.Error.Println(err)
		c.Error(err)
		return
	}
	return c.JSON(http.StatusOK, maps)
}
