package route

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Nerinyan/Nerinyan-APIV2/bodyStruct"
	"github.com/Nerinyan/Nerinyan-APIV2/db"
	"github.com/Nerinyan/Nerinyan-APIV2/entity"
	"github.com/Nerinyan/Nerinyan-APIV2/logger"
	"github.com/Nerinyan/Nerinyan-APIV2/osu"
	"github.com/Nerinyan/Nerinyan-APIV2/src"
	"github.com/Nerinyan/Nerinyan-APIV2/utils"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func DerefString(s *string) string {
	if s != nil {
		return *s
	}

	return ""
}

func (s *SearchQuery) parsePage() {
	// 에러 발생시 int value = 0
	page, _ := strconv.Atoi(s.Page)
	pageSize, _ := strconv.Atoi(s.PageSize)
	if page < 0 {
		page = 0
	}
	if pageSize < 10 {
		pageSize = 50
	}
	if pageSize > 1000 {
		pageSize = 1000
	}
	if page == 0 && pageSize == 0 {
		s.Page = "LIMIT 50"
	} else if page != 0 && pageSize == 0 {
		s.Page = fmt.Sprintf("LIMIT %d,50", page*50)
	} else if page == 0 && pageSize != 0 {
		s.Page = fmt.Sprintf("LIMIT 0,%d", pageSize)
	} else if page != 0 && pageSize != 0 {
		s.Page = fmt.Sprintf("LIMIT %d,%d", page*pageSize, pageSize)
	}
	return
}

func (s *SearchQuery) parseNsfw() {
	ss := strings.ToLower(s.Nsfw)

	switch ss {
	case "1", "true", "all":
		s.Nsfw = "all"
	case "0", "false":
		s.Nsfw = "0"

	default:
		if s.OptionB&(searchOption["cks"]|searchOption["m"]|searchOption["s"]) > 0 {
			s.Nsfw = "all"
		} else {
			s.Nsfw = "0"
		}
	}
	return
}
func (s *SearchQuery) parseOption() {
	ss := strings.ToLower(s.Option)
	if ss == "" {
		s.OptionB |= 0xFFFFFFFF
		return
	}
	for _, s2 := range strings.Split(ss, ",") {
		s.OptionB |= searchOption[s2]
	}
	if s.OptionB == 0 {
		s.OptionB = 0xFFFFFFFF
	}
	return
}

func (s *SearchQuery) parseExtra() {
	s.Extra = strings.ToLower(strings.TrimSpace(s.Extra))
	if s.Storyboard != "1" && s.Storyboard != "all" {
		if strings.Contains(s.Extra, "storyboard") {
			s.Storyboard = "1"
		} else {
			s.Storyboard = "all"
		}
	}
	if s.Video != "1" && s.Video != "all" {
		if strings.Contains(s.Extra, "video") {
			s.Video = "1"
		} else {
			s.Video = "all"
		}
	}
}

func (v *minMax) minAsString() string {
	return fmt.Sprintf("%.1f", v.Min)
}
func (v *minMax) maxAsString() string {
	return fmt.Sprintf("%.1f", v.Max)
}
func (v *minMax) minMaxAsQuery() (query string) {
	if v == nil || (v.Min == 0 && v.Max == 0) {
		return
	} else if v.Min != 0 && v.Max == 0 {
		return fmt.Sprintf(">= %.1f", v.Min)
	} else if v.Min == 0 && v.Max != 0 {
		return fmt.Sprintf("<= %.1f", v.Max)
	}
	return fmt.Sprintf("BETWEEN %.1f AND %.1f", v.Min, v.Max)
}
func (v *minMax) minMaxIsNil() (isNil bool) {
	return v == nil || (v.Min == 0 && v.Max == 0)
}
func (s *SearchQuery) parseQuery() {
	s.parseOption()
	s.parseNsfw()
	s.parsePage()
	s.parseExtra()

}
func (s *SearchQuery) parseRankedStatus() (status []int) {
	for _, st := range strings.Split(s.Ranked, ",") {
		st = strings.ToLower(strings.TrimSpace(st))
		status = append(status, ranked[st]...)
	}

	if len(status) < 1 {
		status = ranked["default"]
	}

	return utils.MakeArrayUnique(&status)
}

func (s *SearchQuery) queryBuilder() (qs string, args []interface{}) {
	s.parseQuery()

	var query bytes.Buffer
	var setAnd []string // 맵셋 	AND 문
	var mapAnd []string // 맵	AND 문
	//var textSearchQuery []string
	text := splitString(s.Text)
	text = utils.MakeArrayUnique(&text)

	if s.Ranked != "all" && s.Ranked != "any" {
		setAnd = append(setAnd, "RANKED IN @ranked")
		args = append(args, sql.Named("ranked", s.parseRankedStatus()))
	}
	if s.Nsfw != "all" {
		setAnd = append(setAnd, "NSFW = "+s.Nsfw)
	}
	if s.Video != "all" {
		setAnd = append(setAnd, "VIDEO = "+s.Video)
	}
	if s.Storyboard != "all" {
		setAnd = append(setAnd, "STORYBOARD = "+s.Storyboard)
	}
	if s.Mode != "all" && s.Mode != "" {
		mapAnd = append(mapAnd, "MODE_INT = @modes")
		args = append(args, sql.Named("modes", mode[s.Mode]))
	}
	if !s.TotalLength.minMaxIsNil() {
		mapAnd = append(mapAnd, `TOTAL_LENGTH `+s.TotalLength.minMaxAsQuery())
	}
	if !s.MaxCombo.minMaxIsNil() {
		mapAnd = append(mapAnd, `MAX_COMBO `+s.MaxCombo.minMaxAsQuery())
	}
	if !s.DifficultyRating.minMaxIsNil() {
		mapAnd = append(mapAnd, `DIFFICULTY_RATING `+s.DifficultyRating.minMaxAsQuery())
	}
	if !s.Accuracy.minMaxIsNil() {
		mapAnd = append(mapAnd, `ACCURACY `+s.Accuracy.minMaxAsQuery())
	}
	if !s.Ar.minMaxIsNil() {
		mapAnd = append(mapAnd, `AR `+s.Ar.minMaxAsQuery())
	}
	if !s.Cs.minMaxIsNil() {
		mapAnd = append(mapAnd, `CS `+s.Cs.minMaxAsQuery())
	}
	if !s.Drain.minMaxIsNil() {
		mapAnd = append(mapAnd, `DRAIN `+s.Drain.minMaxAsQuery())
	}
	if !s.Bpm.minMaxIsNil() {
		mapAnd = append(mapAnd, `BPM `+s.Bpm.minMaxAsQuery())
	}

	if len(mapAnd) > 0 { // beatmapset_id IN ()
		mapAnd = append([]string{"BEATMAPSET_ID IN (MAPSET.BEATMAPSET_ID)"}, mapAnd...)
		setAnd = append(setAnd, "BEATMAPSET_ID IN (SELECT BEATMAPSET_ID FROM BEATMAP WHERE "+strings.Join(mapAnd, " AND ")+" )")
	}

	query.WriteString(
		`SELECT 
    MAPSET.BEATMAPSET_ID, ARTIST, ARTIST_UNICODE, CREATOR, FAVOURITE_COUNT,
    HYPE_CURRENT, HYPE_REQUIRED, NSFW, PLAY_COUNT, SOURCE, STATUS,
    TITLE, TITLE_UNICODE, USER_ID, VIDEO, AVAILABILITY_DOWNLOAD_DISABLED,
    AVAILABILITY_MORE_INFORMATION, BPM, CAN_BE_HYPED, DISCUSSION_ENABLED,
    DISCUSSION_LOCKED, IS_SCOREABLE, LAST_UPDATED, LEGACY_THREAD_URL,
    NOMINATIONS_SUMMARY_CURRENT, NOMINATIONS_SUMMARY_REQUIRED, RANKED,
    RANKED_DATE, STORYBOARD, SUBMITTED_DATE, TAGS, HAS_FAVOURITED,
    DESCRIPTION, GENRE_ID, GENRE_NAME, LANGUAGE_ID, LANGUAGE_NAME, RATINGS
FROM `,
	)

	if len(text) > 0 {
		args = append(args, sql.Named("text", text))
		var textQuery1 []string
		var textQuery2 []string
		{
			query.WriteString("(\n")
			query.WriteString("    SELECT MAPSET.* FROM (\n")
		}
		if s.OptionB == 0xFFFFFFFF || s.OptionB&0b00001111 > 0 {
			var query bytes.Buffer
			query.WriteString("        WITH SCSI AS (SELECT ID FROM SEARCH_CACHE_STRING_INDEX WHERE STRING IN @text )\n")
			query.WriteString("        SELECT BEATMAPSET_ID FROM (\n")
			if s.OptionB&(1<<0) > 0 {
				textQuery1 = append(textQuery1, `SELECT BEATMAPSET_ID, INDEX_KEY FROM SEARCH_CACHE_ARTIST  WHERE INDEX_KEY IN ( SELECT ID FROM SCSI )`) // a
			}
			if s.OptionB&(1<<1) > 0 {
				textQuery1 = append(textQuery1, `SELECT BEATMAPSET_ID, INDEX_KEY FROM SEARCH_CACHE_CREATOR WHERE INDEX_KEY IN ( SELECT ID FROM SCSI )`) // c
			}
			if s.OptionB&(1<<2) > 0 {
				textQuery1 = append(textQuery1, `SELECT BEATMAPSET_ID, INDEX_KEY FROM SEARCH_CACHE_TAG     WHERE INDEX_KEY IN ( SELECT ID FROM SCSI )`) // tg
			}
			if s.OptionB&(1<<3) > 0 {
				textQuery1 = append(textQuery1, `SELECT BEATMAPSET_ID, INDEX_KEY FROM SEARCH_CACHE_TITLE   WHERE INDEX_KEY IN ( SELECT ID FROM SCSI )`) // t
			}
			if s.OptionB == 0xFFFFFFFF {
				textQuery1 = append(textQuery1, `SELECT BEATMAPSET_ID, INDEX_KEY FROM SEARCH_CACHE_OTHER   WHERE INDEX_KEY IN ( SELECT ID FROM SCSI )`) // all
			}
			query.WriteString("                      ")
			query.WriteString(strings.Join(textQuery1, "\n            UNION ALL ") + "\n")
			query.WriteString("        ) TEXTINDEX GROUP BY BEATMAPSET_ID HAVING COUNT(DISTINCT INDEX_KEY) = ( SELECT COUNT(1) FROM SCSI )")
			textQuery2 = append(textQuery2, query.String())
		}
		if s.OptionB&0b01110000 > 0 {
			if s.OptionB&(1<<4) > 1 {
				textQuery2 = append(textQuery2, `SELECT BEATMAPSET_ID FROM BEATMAP    WHERE DELETED_AT IS NULL AND CHECKSUM      IN @text`) // cks
			}
			if s.OptionB&(1<<5) > 1 {
				textQuery2 = append(textQuery2, `SELECT BEATMAPSET_ID FROM BEATMAP    WHERE DELETED_AT IS NULL AND BEATMAP_ID    IN @text`) // m
			}
			if s.OptionB&(1<<6) > 1 {
				textQuery2 = append(textQuery2, `SELECT BEATMAPSET_ID FROM BEATMAPSET WHERE BEATMAPSET_ID IN @text`) // s
			}

		}
		{
			if len(textQuery1) < 1 {
				query.WriteString("                  ")
			}
			query.WriteString(strings.Join(textQuery2, "\n        UNION ALL ") + "\n")
			query.WriteString("    )  TEXTINDEX\n")
			query.WriteString("    LEFT JOIN BEATMAPSET MAPSET ON TEXTINDEX.BEATMAPSET_ID = MAPSET.BEATMAPSET_ID\n")
			query.WriteString("    GROUP BY TEXTINDEX.BEATMAPSET_ID  ) MAPSET\n")
		}
	} else {
		query.WriteString("BEATMAPSET AS MAPSET \n")
	}
	if len(setAnd) > 0 { // SELECT * FROM Beatmapset WHERE ranked in (4,2,1) AND nsfw = 1 ...
		query.WriteString("WHERE " + strings.Join(setAnd, " AND ") + "\n")
	}

	query.WriteString("ORDER BY " + utils.TernaryOperator(orderBy[s.Sort] == "", orderBy["default"], orderBy[s.Sort]) + "\n")
	query.WriteString(s.Page + ";")
	qs = query.String()

	return

}

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

	var sq SearchQuery

	err = c.Bind(&sq)
	if sq.B64 != "" {
		b6, err := base64.StdEncoding.DecodeString(sq.B64)
		if err != nil {
			return c.JSON(
				http.StatusInternalServerError, logger.Error(
					c, &bodyStruct.ErrorStruct{
						Code:    "SEARCH-000-0",
						Error:   err,
						Message: "request parm 'b64' base64 decode fail.",
					},
				),
			)
		}
		err = json.Unmarshal(b6, &sq)
		if err != nil {
			return c.JSON(
				http.StatusInternalServerError, logger.Error(
					c, &bodyStruct.ErrorStruct{
						Code:    "SEARCH-000-1",
						Error:   err,
						Message: "request parm 'b64' json parse fail.",
					},
				),
			)
		}
	}

	if err != nil {
		return c.JSON(
			http.StatusInternalServerError, logger.Error(
				c, &bodyStruct.ErrorStruct{
					Code:    "SEARCH-001",
					Error:   err,
					Message: "request parse error",
				},
			),
		)
	}

	if sq.Option == "" && banchoMapRegex.MatchString(sq.Text) {
		sq.Option = "s"
		sq.Text = banchoMapRegex.FindStringSubmatch(sq.Text)[1]
	}

	q, args := sq.queryBuilder()

	pterm.Info.Println(q)
	rows, err := db.Gorm.Raw(q, args...).Rows()

	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(
				http.StatusNotFound, logger.Error(
					c, &bodyStruct.ErrorStruct{
						Code:    "SEARCH-003",
						Error:   err,
						Message: "not in database",
					},
				),
			)

		}
		return c.JSON(
			http.StatusInternalServerError, logger.Error(
				c, &bodyStruct.ErrorStruct{
					Code:    "SEARCH-004",
					Error:   err,
					Message: "database Query error",
				},
			),
		)

	}
	defer rows.Close()
	var sets []osu.BeatmapSetsOUT
	var index = map[int]int{}
	var mapids []int

	for rows.Next() {
		var set osu.BeatmapSetsOUT

		err = rows.Scan(
			// beatmapset_id, artist, artist_unicode, creator, favourite_count, hype_current,
			//hype_required, nsfw, play_count, source, status, title, title_unicode, user_id,
			//video, availability_download_disabled, availability_more_information, bpm, can_be_hyped,
			//discussion_enabled, discussion_locked, is_scoreable, last_updated, legacy_thread_url,
			//nominations_summary_current, nominations_summary_required, ranked, ranked_date, storyboard,
			//submitted_date, tags, has_favourited, description, genre_id, genre_name, language_id, language_name, ratings

			&set.Id, &set.Artist, &set.ArtistUnicode, &set.Creator, &set.FavouriteCount, &set.Hype.Current, &set.Hype.Required, &set.Nsfw, &set.PlayCount, &set.Source, &set.Status, &set.Title, &set.TitleUnicode, &set.UserId, &set.Video,
			&set.Availability.DownloadDisabled, &set.Availability.MoreInformation, &set.Bpm, &set.CanBeHyped, &set.DiscussionEnabled, &set.DiscussionLocked, &set.IsScoreable, &set.LastUpdated, &set.LegacyThreadUrl,
			&set.NominationsSummary.Current, &set.NominationsSummary.Required, &set.Ranked, &set.RankedDate, &set.Storyboard, &set.SubmittedDate, &set.Tags, &set.HasFavourited, &set.Description.Description, &set.Genre.Id, &set.Genre.Name,
			&set.Language.Id, &set.Language.Name, &set.RatingsString,
		)
		if err != nil {
			return c.JSON(
				http.StatusInternalServerError, logger.Error(
					c, &bodyStruct.ErrorStruct{
						Code:    "SEARCH-005",
						Error:   err,
						Message: "database Query scan error",
					},
				),
			)
		}

		lu, err := time.Parse("2006-01-02 15:04:05", *set.LastUpdated)
		if err != nil {
			return c.JSON(
				http.StatusInternalServerError, logger.Error(
					c, &bodyStruct.ErrorStruct{
						Code:    "SEARCH-005-1",
						Error:   err,
						Message: "time Parse error",
					},
				),
			)
		}
		set.Cache.Video = src.FileList[*set.Id].Unix() >= lu.Unix()
		set.Cache.NoVideo = src.FileList[(*set.Id)*-1].Unix() >= lu.Unix()

		index[*set.Id] = len(sets)
		mapids = append(mapids, *set.Id)
		sets = append(sets, set)

	}

	if len(sets) < 1 {
		return c.JSON(
			http.StatusNotFound, logger.Error(
				c, &bodyStruct.ErrorStruct{
					Code:    "SEARCH-006",
					Error:   errors.New(http.StatusText(http.StatusNotFound)),
					Message: "not in database",
				},
			),
		)
	}

	mapids = utils.MakeArrayUnique(&mapids)

	rows, err = db.Gorm.Raw(
		"SELECT BEATMAP_ID, BEATMAPSET_ID, MODE, MODE_INT, STATUS, RANKED, TOTAL_LENGTH, MAX_COMBO, DIFFICULTY_RATING, VERSION, ACCURACY, AR, CS, DRAIN, BPM, "+
			"`CONVERT`, "+
			"COUNT_CIRCLES, COUNT_SLIDERS, COUNT_SPINNERS, DELETED_AT, HIT_LENGTH, IS_SCOREABLE, LAST_UPDATED, PASSCOUNT, PLAYCOUNT, CHECKSUM, "+
			"USER_ID FROM BEATMAP WHERE DELETED_AT IS NULL AND BEATMAPSET_ID IN @setId ORDER BY DIFFICULTY_RATING;", sql.Named("setId", mapids),
	).Rows()
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(
				http.StatusNotFound, logger.Error(
					c, &bodyStruct.ErrorStruct{
						Code:    "SEARCH-007",
						Error:   err,
						Message: "not in database",
					},
				),
			)
		}
		return c.JSON(
			http.StatusInternalServerError, logger.Error(
				c, &bodyStruct.ErrorStruct{
					Code:    "SEARCH-008",
					Error:   err,
					Message: "database Query error",
				},
			),
		)
	}
	defer rows.Close()
	var tmp []int
	for rows.Next() {
		var Map osu.BeatmapOUT
		err = rows.Scan(
			//beatmap_id, beatmapset_id, mode, mode_int, status, ranked, total_length, max_combo, difficulty_rating,
			//version, accuracy, ar, cs, drain, bpm, convert, count_circles, count_sliders, count_spinners, deleted_at,
			//hit_length, is_scoreable, last_updated, passcount, playcount, checksum, user_id
			&Map.Id, &Map.BeatmapsetId, &Map.Mode, &Map.ModeInt, &Map.Status, &Map.Ranked, &Map.TotalLength, &Map.MaxCombo, &Map.DifficultyRating, &Map.Version, &Map.Accuracy, &Map.Ar, &Map.Cs, &Map.Drain, &Map.Bpm, &Map.Convert,
			&Map.CountCircles, &Map.CountSliders, &Map.CountSpinners, &Map.DeletedAt, &Map.HitLength, &Map.IsScoreable, &Map.LastUpdated, &Map.Passcount, &Map.Playcount, &Map.Checksum, &Map.UserId,
		)
		if err != nil {
			return c.JSON(
				http.StatusInternalServerError, logger.Error(
					c, &bodyStruct.ErrorStruct{
						Code:    "SEARCH-009",
						Error:   err,
						Message: "database Query scan error",
					},
				),
			)
		}
		tmp = append(tmp, *Map.BeatmapsetId)
		OsuFile := fmt.Sprintf(
			"%s - %s (%s) [%s].osu", NotAllowedString.Replace(DerefString(sets[index[*Map.BeatmapsetId]].Artist)), NotAllowedString.Replace(DerefString(sets[index[*Map.BeatmapsetId]].Title)),
			NotAllowedString.Replace(DerefString(sets[index[*Map.BeatmapsetId]].Creator)), NotAllowedString.Replace(maniaKeyRegex.ReplaceAllString(DerefString(Map.Version), "")),
		)
		Map.OsuFile = OsuFile
		sets[index[*Map.BeatmapsetId]].Beatmaps = append(sets[index[*Map.BeatmapsetId]].Beatmaps, Map)

	}
	return c.JSON(http.StatusOK, sets)

}
func Search2(c echo.Context) (err error) {
	var params SearchQuery
	_ = c.Bind(&params)

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
	if len(text) > 0 {
		query.Where(
			`
			MS.BEATMAPSET_ID IN ( 
				SELECT BEATMAPSET_ID FROM (
							  SELECT BEATMAPSET_ID FROM SEARCH_CACHE_ARTIST  WHERE @ARTIST  INDEX_KEY IN ( SELECT ID FROM SEARCH_CACHE_STRING_INDEX WHERE STRING IN @TEXT )
					UNION ALL SELECT BEATMAPSET_ID FROM SEARCH_CACHE_CREATOR WHERE @CREATOR INDEX_KEY IN ( SELECT ID FROM SEARCH_CACHE_STRING_INDEX WHERE STRING IN @TEXT )
					UNION ALL SELECT BEATMAPSET_ID FROM SEARCH_CACHE_TAG     WHERE @TAG     INDEX_KEY IN ( SELECT ID FROM SEARCH_CACHE_STRING_INDEX WHERE STRING IN @TEXT )
					UNION ALL SELECT BEATMAPSET_ID FROM SEARCH_CACHE_TITLE   WHERE @TITLE   INDEX_KEY IN ( SELECT ID FROM SEARCH_CACHE_STRING_INDEX WHERE STRING IN @TEXT )
					UNION ALL SELECT BEATMAPSET_ID FROM SEARCH_CACHE_OTHER   WHERE @OTHER   INDEX_KEY IN ( SELECT ID FROM SEARCH_CACHE_STRING_INDEX WHERE STRING IN @TEXT )
				) C group by BEATMAPSET_ID HAVING COUNT(1) >= @LEN
			)
`,
			sql.Named("TEXT", text),
			sql.Named("ARTIST", params.OptionB&(1<<0) > 0),
			sql.Named("CREATOR", params.OptionB&(1<<1) > 0),
			sql.Named("TAG", params.OptionB&(1<<2) > 0),
			sql.Named("TITLE", params.OptionB&(1<<3) > 0),
			sql.Named("OTHER", params.OptionB == 0xFFFFFFFF),
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
	).Preload("Beatmaps") //

	var sets []entity.BanchoBeatmapSetEntity
	if err = query.Find(&sets).Error; err != nil {
		pterm.Error.Println(err)
		c.Error(err)
		return
	}
	return c.JSON(http.StatusOK, sets)
	//===============================================================================================

}
