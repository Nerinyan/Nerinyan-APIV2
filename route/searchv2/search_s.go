package searchv2

import (
	"fmt"
	"github.com/Nerinyan/Nerinyan-APIV2/banchoCrawler"
	ms "github.com/Nerinyan/Nerinyan-APIV2/db/meilisearch"
	"github.com/labstack/echo/v4"
	"github.com/meilisearch/meilisearch-go"
	"github.com/pterm/pterm"
	"net/http"
	"strconv"
)

// http://127.0.0.1/search/s/123456

func SearchS(c echo.Context) (err error) {
	id, _ := strconv.Atoi(c.Param("setId"))
	if id == 0 {
		return c.JSON(http.StatusOK, struct{}{})
	}

	param := meilisearch.SearchRequest{
		AttributesToRetrieve: nil,
		AttributesToCrop:     nil,
		CropLength:           0,
		CropMarker:           "",
		MatchingStrategy:     "all",
		Filter:               fmt.Sprintf("id = %d", id),
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

	var resp []banchoCrawler.Beatmasets
	if err = ms.HitsMapper(res.Hits, &resp); err != nil {
		return
	}
	for _, beatmasets := range resp {
		if beatmasets.Id == id {
			return c.JSON(http.StatusOK, beatmasets)
		}
	}
	return c.JSON(http.StatusOK, struct{}{})
}
