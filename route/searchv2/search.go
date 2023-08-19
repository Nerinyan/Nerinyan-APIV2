package searchv2

import (
	"fmt"
	ms "github.com/Nerinyan/Nerinyan-APIV2/db/ms"
	"github.com/Nerinyan/Nerinyan-APIV2/utils"
	"github.com/labstack/echo/v4"
	"github.com/meilisearch/meilisearch-go"
	"github.com/pterm/pterm"
	"strconv"
	"strings"
)

func Search(c echo.Context) (err error) {
	var params SearchQuery
	err = c.Bind(&params)
	if err != nil {
		pterm.Error.WithShowLineNumber().Println(err.Error())
		pterm.Error.WithShowLineNumber().Println("URI", c.Request().RequestURI)
	}
	err = params.parseB64()
	if err != nil {
		pterm.Error.WithShowLineNumber().Println(err.Error())
		return
	}
	optionB := params.parseOption()
	_ = optionB

	param := meilisearch.SearchRequest{
		AttributesToRetrieve:  nil,
		AttributesToCrop:      nil,
		CropLength:            0,
		CropMarker:            "",
		AttributesToHighlight: nil,
		HighlightPreTag:       "",
		HighlightPostTag:      "",
		MatchingStrategy:      "all",
		Filter:                nil,
		ShowMatchesPosition:   false,
		Facets:                nil,
		PlaceholderSearch:     false,
		Sort:                  []string{parseOrder(params.Sort)},
		HitsPerPage:           int64(params.getPageSize()),
		Page:                  int64(params.getPage()),
	}
	var filter []string //c.QueryParam("hq")
	//===============================================================================================
	// 맵셋 조건
	if params.Ranked != "all" {
		r := utils.NotInMapFindAllAppendDefault(ranked, utils.SplitTrimLower(params.Ranked, ","))
		filter = append(filter, fmt.Sprintf("ranked IN [%s]", intSliceToCSV(r)))
	}
	if !params.getNsfw() {
		filter = append(filter, "nsfw = false")
	}
	if params.getVideo() {
		filter = append(filter, "video = true")
	}
	if params.getStoryboard() {
		filter = append(filter, fmt.Sprintf("storyboard = %t", params.getStoryboard()))
	}

	if params.Mode != "all" && params.Mode != "" {
		r := utils.NotInMapFindAllDefault(mode, utils.SplitTrimLower(params.Mode, ","))
		filter = append(filter, fmt.Sprintf("beatmaps.mode_int IN [%s]", intSliceToCSV(r)))
	}
	//================================================================================================================================
	if params.TotalLength.Min != 0 {
		filter = append(filter, fmt.Sprintf("beatmaps.total_length >= %f", params.TotalLength.Min))
	}
	if params.TotalLength.Max != 0 {
		filter = append(filter, fmt.Sprintf("beatmaps.total_length <= %f", params.TotalLength.Max))
	}
	//================================================================================================================================
	if params.MaxCombo.Min != 0 {
		filter = append(filter, fmt.Sprintf("beatmaps.max_combo >= %f", params.MaxCombo.Min))
	}
	if params.MaxCombo.Max != 0 {
		filter = append(filter, fmt.Sprintf("beatmaps.max_combo <= %f", params.MaxCombo.Max))
	}
	//================================================================================================================================
	if params.DifficultyRating.Min != 0 {
		filter = append(filter, fmt.Sprintf("beatmaps.difficulty_rating >= %f", params.DifficultyRating.Min))
	}
	if params.DifficultyRating.Max != 0 {
		filter = append(filter, fmt.Sprintf("beatmaps.difficulty_rating <= %f", params.DifficultyRating.Max))
	}
	//================================================================================================================================
	if params.Accuracy.Min != 0 {
		filter = append(filter, fmt.Sprintf("beatmaps.accuracy >= %f", params.Accuracy.Min))
	}
	if params.Accuracy.Max != 0 {
		filter = append(filter, fmt.Sprintf("beatmaps.accuracy <= %f", params.Accuracy.Max))
	}
	//================================================================================================================================
	if params.Ar.Min != 0 {
		filter = append(filter, fmt.Sprintf("beatmaps.ar >= %f", params.Ar.Min))
	}
	if params.Ar.Max != 0 {
		filter = append(filter, fmt.Sprintf("beatmaps.ar <= %f", params.Ar.Max))
	}
	//================================================================================================================================
	if params.Cs.Min != 0 {
		filter = append(filter, fmt.Sprintf("beatmaps.cs >= %f", params.Cs.Min))
	}
	if params.Cs.Max != 0 {
		filter = append(filter, fmt.Sprintf("beatmaps.cs <= %f", params.Cs.Max))
	}
	//================================================================================================================================
	if params.Drain.Min != 0 {
		filter = append(filter, fmt.Sprintf("beatmaps.drain >= %f", params.Drain.Min))
	}
	if params.Drain.Max != 0 {
		filter = append(filter, fmt.Sprintf("beatmaps.drain <= %f", params.Drain.Max))
	}
	//================================================================================================================================
	if params.Bpm.Min != 0 {
		filter = append(filter, fmt.Sprintf("beatmaps.bpm >= %f", params.Bpm.Min))
	}
	if params.Bpm.Max != 0 {
		filter = append(filter, fmt.Sprintf("beatmaps.bpm <= %f", params.Bpm.Max))
	}
	//================================================================================================================================

	if optionB != 0xFFFFFFFF && optionB&(1<<5) > 0 && len(params.Text) > 0 {
		id, _ := strconv.Atoi(params.Text)
		if id != 0 {
			filter = append(filter, fmt.Sprintf("beatmaps.id = %d", id))
		}
	}
	if optionB != 0xFFFFFFFF && optionB&(1<<6) > 0 && len(params.Text) > 0 {
		id, _ := strconv.Atoi(params.Text)
		if id != 0 {
			filter = append(filter, fmt.Sprintf("id = %d", id))
		}
	}
	//================================================================================================================================
	param.Filter = filter
	res, err := ms.GetIndexBeatMapSet().Search(params.Text, &param)
	if err != nil {
		return err
	}
	pterm.Info.Printfln("query: '%s' sort: '%s' Filter: '%v'", res.Query, params.Sort, param.Filter)
	return c.JSON(200, res.Hits)
}

func parseOrder(order string) (res string) {
	defer func() {
		if res == "" {
			res = orderBy["default"]
			return
		}
	}()
	res = orderBy[strings.ToLower(order)]
	return
}

func intSliceToCSV(nums []int) string {
	strs := make([]string, len(nums))
	for i, n := range nums {
		strs[i] = strconv.Itoa(n)
	}
	return strings.Join(strs, ",")
}
