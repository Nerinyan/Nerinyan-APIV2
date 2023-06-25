package search

import (
	"github.com/Nerinyan/Nerinyan-APIV2/db"
	"github.com/Nerinyan/Nerinyan-APIV2/db/mariadb/entity"
	"github.com/labstack/echo/v4"
	"github.com/pterm/pterm"
	"net/http"
)

// http://127.0.0.1/search/s/123456

func SearchS(c echo.Context) (err error) {
	query := db.Gorm.Select("MS.*").Table("BEATMAPSET MS").
		Where("MS.DELETED_AT IS NULL").
		Where("MS.BEATMAPSET_ID = ?", c.Param("setId")).
		Preload("Beatmaps", "DELETED_AT IS NULL") // 이게 있어야 gorm join이 작동함

	var sets entity.BanchoBeatmapSetEntity
	if err = query.Find(&sets).Error; err != nil {
		pterm.Error.Println(err)
		c.Error(err)
		return
	}
	return c.JSON(http.StatusOK, sets)
}
