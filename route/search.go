package route

import (
	"database/sql"
	"github.com/Nerinyan/Nerinyan-APIV2/db"
	"github.com/Nerinyan/Nerinyan-APIV2/entity"
	"github.com/Nerinyan/Nerinyan-APIV2/utils"
	"github.com/labstack/echo/v4"
	"github.com/pterm/pterm"
	"net/http"
	"regexp"
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

var banchoMapRegex, _ = regexp.Compile(`(?:https://osu[.]ppy[.]sh/beatmapsets/)(\d+?)(?:\D|$)`)
var maniaKeyRegex, _ = regexp.Compile(`(\[[0-9]K\] )`)
var NotAllowedString = strings.NewReplacer("\\", "", "/", "", "|", "", ":", "", "?", "", "*", "", "<", "", ">", "", "\"", "")

func Search(c echo.Context) (err error) {
	var params SearchQuery
	err = c.Bind(&params)
	if err != nil {
		pterm.Error.WithShowLineNumber().Println(err.Error())
		return
	}
	params.parseB64()
	query := db.Gorm.Select(
		"MS.*",
	).Table(
		"BEATMAPSET MS",
	).Joins(
		"INNER JOIN BEATMAP M ON MS.BEATMAPSET_ID = M.BEATMAPSET_ID AND M.DELETED_AT IS NULL",
	)
	//===============================================================================================
	// 조건처리
	if params.Ranked != "all" && params.Ranked != "" {
		query.Where("MS.RANKED IN ?", utils.NotInMapFindDefault(ranked, params.Ranked))
	}
	if !params.getNsfw() {
		query.Where("MS.NSFW = ?", params.getNsfw())
	}
	if params.getVideo() {
		query.Where("MS.VIDEO = ?", params.getVideo())
	}
	if params.getStoryboard() {
		query.Where("MS.STORYBOARD = ?", params.getStoryboard())
	}
	if params.TotalLength.Min != 0 {
		query.Where("M.TOTAL_LENGTH >= ?", params.TotalLength.Min)
	}
	if params.TotalLength.Max != 0 {
		query.Where("M.TOTAL_LENGTH <= ?", params.TotalLength.Max)
	}
	if params.MaxCombo.Min != 0 {
		query.Where("M.MAX_COMBO >= ?", params.MaxCombo.Min)
	}
	if params.MaxCombo.Max != 0 {
		query.Where("M.MAX_COMBO <= ?", params.MaxCombo.Max)
	}
	if params.DifficultyRating.Min != 0 {
		query.Where("M.DIFFICULTY_RATING >= ?", params.DifficultyRating.Min)
	}
	if params.DifficultyRating.Max != 0 {
		query.Where("M.DIFFICULTY_RATING <= ?", params.DifficultyRating.Max)
	}
	if params.Accuracy.Max != 0 {
		query.Where("M.ACCURACY >= ?", params.Accuracy.Max)
	}
	if params.Accuracy.Min != 0 {
		query.Where("M.ACCURACY <= ?", params.Accuracy.Min)
	}
	if params.Ar.Max != 0 {
		query.Where("M.AR >= ?", params.Ar.Max)
	}
	if params.Ar.Min != 0 {
		query.Where("M.AR <= ?", params.Ar.Min)
	}
	if params.Cs.Max != 0 {
		query.Where("M.CS >= ?", params.Cs.Max)
	}
	if params.Cs.Min != 0 {
		query.Where("M.CS <= ?", params.Cs.Min)
	}
	if params.Drain.Max != 0 {
		query.Where("M.DRAIN >= ?", params.Drain.Max)
	}
	if params.Drain.Min != 0 {
		query.Where("M.DRAIN <= ?", params.Drain.Min)
	}
	if params.Bpm.Max != 0 {
		query.Where("M.BPM >= ?", params.Bpm.Max)
	}
	if params.Bpm.Min != 0 {
		query.Where("M.BPM <= ?", params.Bpm.Min)
	}
	if params.Mode != "all" && params.Mode != "" {
		query.Where("M.MODE_INT IN ?", utils.NotInMapFindAllDefault(mode, utils.SplitTrimLower(params.Mode, ",")))
	}
	text := splitString(params.Text)
	text = utils.MakeArrayUnique(&text)
	optionB := params.parseOption()
	if len(text) > 0 && (optionB&0b1111 > 0 || optionB == 0xFFFFFFFF) {
		query.Where(
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
	query.Order(
		utils.NotInMapFindDefault(orderBy, params.Sort),
	).Limit(
		params.getPageSize(),
	).Offset(
		utils.Multiply(params.getPage(), params.getPageSize()),
	).Preload("Beatmaps") // 이게 있어야 gorm join이 작동함

	var sets []entity.BanchoBeatmapSetEntity
	if err = query.Find(&sets).Error; err != nil {
		pterm.Error.Println(err)
		c.Error(err)
		return
	}
	return c.JSON(http.StatusOK, sets)
	//===============================================================================================

}
