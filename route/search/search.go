package search

import (
	"database/sql"
	"github.com/Nerinyan/Nerinyan-APIV2/db"
	"github.com/Nerinyan/Nerinyan-APIV2/db/mariadb/entity"
	"github.com/Nerinyan/Nerinyan-APIV2/utils"
	"github.com/labstack/echo/v4"
	"github.com/pterm/pterm"
	"net/http"
	"strings"
)

func splitString(input string) (ss []string) {
	for _, s := range strings.Split(strings.ToLower(regexpReplace.ReplaceAllString(input, " ")), " ") {
		s = strings.TrimSpace(s)
		if s == "" || s == " " {
			continue
		}
		//ss = append(ss, s, porter2.Stemmer.Stem(s))
		ss = append(ss, s)
	}
	return
}

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
	text := splitString(params.Text)
	text = utils.MakeArrayUnique(&text)

	setQuery := db.Gorm.Select("MS.*").Table("BEATMAPSET MS").Where("MS.DELETED_AT IS NULL")
	//===============================================================================================
	// 맵셋 조건
	if params.Ranked != "all" {
		setQuery.Where("MS.RANKED IN ?", utils.NotInMapFindAllAppendDefault(ranked, utils.SplitTrimLower(params.Ranked, ",")))
	}
	if !params.getNsfw() {
		setQuery.Where("MS.NSFW = ?", false)
	}
	if params.getVideo() {
		setQuery.Where("MS.VIDEO = ?", true)
	}
	if params.getStoryboard() {
		setQuery.Where("MS.STORYBOARD = ?", params.getStoryboard())
	}
	if optionB != 0xFFFFFFFF && optionB&(1<<6) > 0 && len(text) > 0 {
		setQuery.Where("MS.BEATMAPSET_ID IN @BEATMAPSET_ID", sql.Named("BEATMAPSET_ID", text))
	}
	//===============================================================================================
	// 맵 조건
	mapQuery := db.Gorm.Select("M.BEATMAPSET_ID").Table("BEATMAP M").Where("M.DELETED_AT IS NULL")
	useMap := false
	if params.TotalLength.Min != 0 {
		mapQuery.Where("? <= M.TOTAL_LENGTH", params.TotalLength.Min)
		useMap = true
	}
	if params.TotalLength.Max != 0 {
		mapQuery.Where("M.TOTAL_LENGTH <= ?", params.TotalLength.Max)
		useMap = true
	}
	if params.MaxCombo.Min != 0 {
		mapQuery.Where("? <= M.MAX_COMBO", params.MaxCombo.Min)
		useMap = true
	}
	if params.MaxCombo.Max != 0 {
		mapQuery.Where("M.MAX_COMBO <= ?", params.MaxCombo.Max)
		useMap = true
	}
	if params.DifficultyRating.Min != 0 {
		mapQuery.Where("? <= M.DIFFICULTY_RATING", params.DifficultyRating.Min)
		useMap = true
	}
	if params.DifficultyRating.Max != 0 {
		mapQuery.Where("M.DIFFICULTY_RATING <= ?", params.DifficultyRating.Max)
		useMap = true
	}
	if params.Accuracy.Min != 0 {
		mapQuery.Where("? <= M.ACCURACY", params.Accuracy.Min)
		useMap = true
	}
	if params.Accuracy.Max != 0 {
		mapQuery.Where("M.ACCURACY <= ?", params.Accuracy.Max)
		useMap = true
	}
	if params.Ar.Min != 0 {
		mapQuery.Where("? <= M.AR", params.Ar.Min)
		useMap = true
	}
	if params.Ar.Max != 0 {
		mapQuery.Where("M.AR <= ?", params.Ar.Max)
		useMap = true
	}
	if params.Cs.Min != 0 {
		mapQuery.Where("? <= M.CS", params.Cs.Min)
		useMap = true
	}
	if params.Cs.Max != 0 {
		mapQuery.Where("M.CS <= ?", params.Cs.Max)
		useMap = true
	}
	if params.Drain.Min != 0 {
		mapQuery.Where("? <= M.DRAIN", params.Drain.Min)
		useMap = true
	}
	if params.Drain.Max != 0 {
		mapQuery.Where("M.DRAIN <= ?", params.Drain.Max)
		useMap = true
	}
	if params.Bpm.Min != 0 {
		mapQuery.Where("? <= M.BPM", params.Bpm.Min)
		useMap = true
	}
	if params.Bpm.Max != 0 {
		mapQuery.Where("M.BPM <= ?", params.Bpm.Max)
		useMap = true
	}
	if params.Mode != "all" && params.Mode != "" {
		mapQuery.Where("M.MODE_INT IN ?", utils.NotInMapFindAllDefault(mode, utils.SplitTrimLower(params.Mode, ",")))
		useMap = true
	}
	if optionB != 0xFFFFFFFF && optionB&(1<<4) > 0 && len(text) > 0 {
		mapQuery.Where("M.CHECKSUM IN @CHECKSUM", sql.Named("CHECKSUM", text))
		useMap = true
	}
	if optionB != 0xFFFFFFFF && optionB&(1<<5) > 0 && len(text) > 0 {
		mapQuery.Where("M.BEATMAP_ID IN @BEATMAP_ID", sql.Named("BEATMAP_ID", text))
		useMap = true
	}

	if useMap {
		setQuery.Where("MS.BEATMAPSET_ID IN (?)", mapQuery)
	}

	if len(text) > 0 && (optionB&0b1111 > 0 || optionB == 0xFFFFFFFF) {
		setQuery.Where(
			`
			MS.BEATMAPSET_ID IN ( 
				SELECT BEATMAPSET_ID FROM (
							  SELECT BEATMAPSET_ID FROM SEARCH_CACHE_ARTIST  WHERE @ARTIST  AND INDEX_KEY IN ( SELECT ID FROM SEARCH_CACHE_STRING_INDEX WHERE STRING IN @TEXT )
					UNION ALL SELECT BEATMAPSET_ID FROM SEARCH_CACHE_CREATOR WHERE @CREATOR AND INDEX_KEY IN ( SELECT ID FROM SEARCH_CACHE_STRING_INDEX WHERE STRING IN @TEXT )
					UNION ALL SELECT BEATMAPSET_ID FROM SEARCH_CACHE_TAG     WHERE @TAG     AND INDEX_KEY IN ( SELECT ID FROM SEARCH_CACHE_STRING_INDEX WHERE STRING IN @TEXT )
					UNION ALL SELECT BEATMAPSET_ID FROM SEARCH_CACHE_TITLE   WHERE @TITLE   AND INDEX_KEY IN ( SELECT ID FROM SEARCH_CACHE_STRING_INDEX WHERE STRING IN @TEXT )
					UNION ALL SELECT BEATMAPSET_ID FROM SEARCH_CACHE_OTHER   WHERE @OTHER   AND INDEX_KEY IN ( SELECT ID FROM SEARCH_CACHE_STRING_INDEX WHERE STRING IN @TEXT )
				) C group by BEATMAPSET_ID HAVING COUNT(1) >= @LEN
			)
`,
			sql.Named("TEXT", text),
			sql.Named("ARTIST", optionB&(1<<0) > 0),
			sql.Named("CREATOR", optionB&(1<<1) > 0),
			sql.Named("TAG", optionB&(1<<2) > 0),
			sql.Named("TITLE", optionB&(1<<3) > 0),
			sql.Named("OTHER", optionB == 0xFFFFFFFF),
			sql.Named("LEN", len(text)),
		)
	}

	// 조건 order, join, page
	setQuery.Order(
		utils.NotInMapFindDefault(orderBy, params.Sort),
	).Limit(
		params.getPageSize(),
	).Offset(
		utils.Multiply(params.getPage(), params.getPageSize()),
	).Preload("Beatmaps", "DELETED_AT IS NULL") // 이게 있어야 gorm join이 작동함

	var sets []entity.BanchoBeatmapSetEntity
	if err = setQuery.Find(&sets).Error; err != nil {
		pterm.Error.Println(err)
		c.Error(err)
		return
	}
	return c.JSON(http.StatusOK, sets)
	//===============================================================================================

}
