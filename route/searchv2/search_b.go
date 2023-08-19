package searchv2

import (
	"fmt"
	ms "github.com/Nerinyan/Nerinyan-APIV2/db/ms"
	"github.com/Nerinyan/Nerinyan-APIV2/entity"
	"github.com/labstack/echo/v4"
	"github.com/meilisearch/meilisearch-go"
	"github.com/pterm/pterm"
	"net/http"
	"strconv"
)

func SearchB(c echo.Context) (err error) { // http://127.0.0.1/search/b/123456
	id, _ := strconv.Atoi(c.Param("mapId"))
	if id == 0 {
		return c.JSON(http.StatusOK, struct{}{})
	}

	param := meilisearch.SearchRequest{
		AttributesToRetrieve: []string{"beatmaps"},
		AttributesToCrop:     nil,
		CropLength:           0,
		CropMarker:           "",
		MatchingStrategy:     "all",
		Filter:               fmt.Sprintf("beatmaps.id = %d", id),
		ShowMatchesPosition:  false,
		Facets:               nil,
		PlaceholderSearch:    false,
	}

	pterm.Info.Printfln("Filter: '%v'", param.Filter)
	res, err := ms.GetIndexBeatMapSet().Search("", &param)
	if err != nil {
		return err
	}

	if len(res.Hits) < 1 {
		return c.JSON(http.StatusOK, struct{}{})
	}

	var resp []entity.Beatmaset
	if err = ms.HitsMapper(res.Hits, &resp); err != nil {
		return
	}
	for _, beatmasets := range resp {
		for _, beatmap := range beatmasets.Beatmaps {
			if beatmap.Id == id {
				return c.JSON(http.StatusOK, beatmap)
			}
		}
	}
	return c.JSON(http.StatusOK, struct{}{})
}
