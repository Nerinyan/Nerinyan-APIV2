package banchoCrawler

import (
	"database/sql"
	"fmt"
	"github.com/Nerinyan/Nerinyan-APIV2/config"
	"github.com/Nerinyan/Nerinyan-APIV2/db"
	"github.com/Nerinyan/Nerinyan-APIV2/osu"
	"github.com/goccy/go-json"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

var api = struct {
	count int
	mutex sync.Mutex
}{}
var pause bool

var ApiCount = &api.count

func apicountAdd() {
	api.mutex.Lock()
	api.count++
	api.mutex.Unlock()
}

func apiCountReset() {
	api.mutex.Lock()
	api.count = 0
	api.mutex.Unlock()
}

func RunGetBeatmapDataASBancho() {

	go func() {
		for {
			time.Sleep(time.Minute)
			if db.Maria.Ping() != nil {
				continue
			}
			apiCountReset()
			go config.Config.Save()
		}
	}()
	go func() { //ALL desc limit 50
		for {
			awaitApiCount()
			time.Sleep(time.Second * 30)
			getUpdatedMapDesc()
		}
	}()
	go func() { //Update Ranked DESC limit 50
		for {
			awaitApiCount()
			time.Sleep(time.Minute)
			getUpdatedMapRanked()
		}
	}()
	go func() { //Update Qualified desc limit 50
		for {
			awaitApiCount()
			time.Sleep(time.Minute)
			getUpdatedMapQualified()
		}
	}()
	go func() { //Update Loved DESC limit 50
		for {
			awaitApiCount()
			time.Sleep(time.Minute)
			getUpdatedMapLoved()
		}
	}()
	go func() { //Update Graveyard asc limit 50
		for {
			awaitApiCount()
			time.Sleep(time.Minute)
			getGraveyardMap()
		}
	}()

	go func() { //ALL asc
		for {
			awaitApiCount()
			getUpdatedMapAsc()
		}
	}()
	pterm.Info.Println("Bancho cron started.")
}
func awaitApiCount() {
	for {
		if api.count < 50 && !pause {
			break
		}
		time.Sleep(time.Millisecond * 500)
	}
}
func ManualUpdateBeatmapSet(id int) {
	var err error
	defer func() {
		if err != nil {
			pterm.Error.Println(err)
		}
	}()
	url := fmt.Sprintf("https://osu.ppy.sh/api/v2/beatmapsets/search?nsfw=true&s=any&q=%d", id)

	var data osu.BeatmapsetsSearch
	if err = stdGETBancho(url, &data); err != nil {
		return
	}
	// updateMapset(&data)

	if err = updateSearchBeatmaps(*data.Beatmapsets); err != nil {
		return
	}
}

func getUpdatedMapRanked() {
	var err error
	defer func() {
		if err != nil {
			pterm.Error.Println(err)
		}
	}()
	url := "https://osu.ppy.sh/api/v2/beatmapsets/search?nsfw=true&s=ranked"

	var data osu.BeatmapsetsSearch

	if err = stdGETBancho(url, &data); err != nil {
		return
	}
	//pterm.Info.Println("getUpdatedMapRanked", string(*utils.ToJsonString(*data.Beatmapsets))[:100])

	if err = updateSearchBeatmaps(*data.Beatmapsets); err != nil {
		return
	}

	return
}

func getUpdatedMapLoved() {
	var err error
	defer func() {
		if err != nil {
			pterm.Error.Println(err)
		}
	}()
	url := "https://osu.ppy.sh/api/v2/beatmapsets/search?nsfw=true&s=loved"

	var data osu.BeatmapsetsSearch
	if err = stdGETBancho(url, &data); err != nil {
		return
	}
	//pterm.Info.Println("getUpdatedMapLoved", string(*utils.ToJsonString(*data.Beatmapsets))[:100])
	if err = updateSearchBeatmaps(*data.Beatmapsets); err != nil {
		return
	}
	return
}
func getUpdatedMapQualified() {
	var err error
	defer func() {
		if err != nil {
			pterm.Error.Println(err)
		}
	}()
	url := "https://osu.ppy.sh/api/v2/beatmapsets/search?nsfw=true&s=qualified"

	var data osu.BeatmapsetsSearch
	if err = stdGETBancho(url, &data); err != nil {
		return
	}
	//pterm.Info.Println("getUpdatedMapQualified", string(*utils.ToJsonString(*data.Beatmapsets))[:100])
	if err = updateSearchBeatmaps(*data.Beatmapsets); err != nil {
		return
	}
	return
}

func getGraveyardMap() {
	var err error
	defer func() {
		if err != nil {
			pterm.Error.Println(err)
		}
	}()
	url := ""
	cs := &config.Config.Osu.BeatmapUpdate.GraveyardAsc.CursorString
	if *cs != "" {
		url = "https://osu.ppy.sh/api/v2/beatmapsets/search?nsfw=true&sort=updated_asc&s=graveyard&cursor_string=" + *cs
	} else {
		url = "https://osu.ppy.sh/api/v2/beatmapsets/search?nsfw=true&sort=updated_asc&s=graveyard"
	}

	var data osu.BeatmapsetsSearch

	err = stdGETBancho(url, &data)
	if err != nil {
		return
	}
	if data.CursorString == "" {
		return
	}
	//pterm.Info.Println("getGraveyardMap", string(*utils.ToJsonString(*data.Beatmapsets))[:100])
	if err = updateSearchBeatmaps(*data.Beatmapsets); err != nil {
		return
	}
	*cs = data.CursorString
	return
}

func getUpdatedMapDesc() {
	var err error
	defer func() {
		if err != nil {
			pterm.Error.Println(err)
		}
	}()
	url := "https://osu.ppy.sh/api/v2/beatmapsets/search?nsfw=true&sort=updated_desc&s=any"

	var data osu.BeatmapsetsSearch

	if err = stdGETBancho(url, &data); err != nil {
		return
	}

	//pterm.Info.Println("getUpdatedMapDesc", string(*utils.ToJsonString(*data.Beatmapsets))[:100])
	if err = updateSearchBeatmaps(*data.Beatmapsets); err != nil {
		return
	}
	if data.CursorString == "" {
		return
	}
	config.Config.Osu.BeatmapUpdate.UpdatedDesc.CursorString = data.CursorString
	return
}

func getUpdatedMapAsc() {
	var err error
	defer func() {
		if err != nil {
			pterm.Error.Println(err)
		}
	}()
	url := ""
	cs := &config.Config.Osu.BeatmapUpdate.UpdatedAsc.CursorString
	if *cs != "" {
		url = "https://osu.ppy.sh/api/v2/beatmapsets/search?nsfw=true&sort=updated_asc&s=any&cursor_string=" + *cs
	} else {
		url = "https://osu.ppy.sh/api/v2/beatmapsets/search?nsfw=true&sort=updated_asc&s=any"
	}

	var data osu.BeatmapsetsSearch

	err = stdGETBancho(url, &data)
	if err != nil {
		return
	}

	//pterm.Info.Println(url, data.CursorString)
	//pterm.Info.Println(data.CursorString, url, string(*utils.ToJsonString(*data.Beatmapsets))[:200])
	if err = updateSearchBeatmaps(*data.Beatmapsets); err != nil {
		return
	}
	*cs = data.CursorString
	return
}

func stdGETBancho(url string, str interface{}) (err error) {
	if config.Config.Log.Crawler {
		pterm.Info.Printfln("%s | %-50s | URL : %s", time.Now().Format("15:04:05.000"), pterm.Yellow("BANCHO CRAWLER"), url)
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return
	}

	req.Header.Add("Authorization", config.Config.Osu.Token.TokenType+" "+config.Config.Osu.Token.AccessToken)
	req.Header.Add("Content-Type", "Application/json")

	res, err := client.Do(req)
	apicountAdd()
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		if res.StatusCode == 401 {
			pause = true
		}
		return errors.New(res.Status)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	return json.Unmarshal(body, &str)
	// return json.NewDecoder(res.Body).Decode(&str)

}

func updateMapset(s *osu.BeatmapSetsIN) {

	//	beatmapset_id,artist,artist_unicode,creator,favourite_count,
	//	hype_current,hype_required,nsfw,play_count,source,
	//	status,title,title_unicode,user_id,video,
	//	availability_download_disabled,availability_more_information,bpm,can_be_hyped,discussion_enabled,
	//	discussion_locked,is_scoreable,last_updated,legacy_thread_url,nominations_summary_current,
	//	nominations_summary_required,ranked,ranked_date,storyboard,submitted_date,
	//	tags,has_favourited,description,genre_id,genre_name,
	//	language_id,language_name,ratings

	r := *s.Ratings

	db.InsertQueueChannel <- db.ExecQueue{
		//DB 큐에 전송
		Query: UpsertBeatmapSet,
		Args: []any{
			s.Id, s.Artist, s.ArtistUnicode, s.Creator, s.FavouriteCount, *s.Hype.Current, *s.Hype.Required, s.Nsfw, s.PlayCount, s.Source, s.Status, s.Title, s.TitleUnicode, s.UserId, s.Video, s.Availability.DownloadDisabled,
			s.Availability.MoreInformation, s.Bpm, s.CanBeHyped, s.DiscussionEnabled, s.DiscussionLocked, s.IsScoreable, s.LastUpdated, s.LegacyThreadUrl, s.NominationsSummary.Current, s.NominationsSummary.Required, s.Ranked, s.RankedDate,
			s.DeletedAt,
			s.Storyboard, s.SubmittedDate, s.Tags, s.HasFavourited, s.Description.Description, s.Genre.Id, s.Genre.Name, s.Language.Id, s.Language.Name,
			fmt.Sprintf("%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d", r[0], r[1], r[2], r[3], r[4], r[5], r[6], r[7], r[8], r[9], r[10]),
		},
	}
	//_, err := db.Maria.Exec(UpsertBeatmapSet, s.Id, s.Artist, s.ArtistUnicode, s.Creator, s.FavouriteCount, s.Hype.Current, s.Hype.Required, s.Nsfw, s.PlayCount, s.Source, s.Status, s.Title, s.TitleUnicode, s.UserId, s.Video, s.Availability.DownloadDisabled, s.Availability.MoreInformation, s.Bpm, s.CanBeHyped, s.DiscussionEnabled, s.DiscussionLocked, s.IsScoreable, s.LastUpdated, s.LegacyThreadUrl, s.NominationsSummary.Current, s.NominationsSummary.Required, s.Ranked, s.RankedDate, s.Storyboard, s.SubmittedDate, s.Tags, s.HasFavourited, s.Description.Description, s.Genre.Id, s.Genre.Name, s.Language.Id, s.Language.Name, fmt.Sprintf("%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d", r[0], r[1], r[2], r[3], r[4], r[5], r[6], r[7], r[8], r[9], r[10]))
	//if err != nil {
	//	log.Error(err)
	//	pterm.Error.Println(err)
	//}
	//커버 이미지 주소
	//if _, err = db.Maria.Exec(fmt.Sprintf(upsertCoverQuery, coverValue), s.Id, s.Covers.Cover, s.Covers.Cover2X, s.Covers.Card, s.Covers.Card2X, s.Covers.List, s.Covers.List2X, s.Covers.Slimcover, s.Covers.Slimcover2X); err != nil {
	//	log.Error(err)
	//	pterm.Error.Println(err)
	//}
	fmt.Println("22222222222222222222222")
	if *s.Beatmaps == nil {
		return
	}

	fmt.Println("333333333333")
	for _, m := range *s.Beatmaps {
		go upsertMap(m)
	}

}

func upsertMap(m osu.BeatmapIN) {
	db.InsertQueueChannel <- db.ExecQueue{
		//DB 큐에 전송
		Query: UpsertBeatmap,
		Args: []any{
			m.Id, m.BeatmapsetId, m.Mode, m.ModeInt, m.Status, m.Ranked, m.TotalLength, m.MaxCombo, m.DifficultyRating, m.Version, m.Accuracy, m.Ar, m.Cs, m.Drain, m.Bpm, m.Convert,
			m.CountCircles, m.CountSliders, m.CountSpinners, m.DeletedAt, m.HitLength, m.IsScoreable, m.LastUpdated, m.Passcount, m.Playcount, m.Checksum, m.UserId,
		},
	}
	//_, err := db.Maria.Exec(UpsertBeatmap, m.Id, m.BeatmapsetId, m.Mode, m.ModeInt, m.Status, m.Ranked, m.TotalLength, m.MaxCombo, m.DifficultyRating, m.Version, m.Accuracy, m.Ar, m.Cs, m.Drain, m.Bpm, m.Convert, m.CountCircles, m.CountSliders, m.CountSpinners, m.DeletedAt, m.HitLength, m.IsScoreable, m.LastUpdated, m.Passcount, m.Playcount, m.Checksum, m.UserId)
	//if err != nil {
	//	log.Error(err)
	//	pterm.Error.Println(err)
	//}

}

const (
	UpsertBeatmap = `/* UPSERT BEATMAP */
	INSERT INTO BEATMAP
		(
			BEATMAP_ID,BEATMAPSET_ID,MODE,MODE_INT,STATUS,	RANKED,TOTAL_LENGTH,MAX_COMBO,DIFFICULTY_RATING,VERSION,
			ACCURACY,AR,CS,DRAIN,BPM,` + "`CONVERT`" + `,COUNT_CIRCLES,COUNT_SLIDERS,COUNT_SPINNERS,DELETED_AT,
			HIT_LENGTH,IS_SCOREABLE,LAST_UPDATED,PASSCOUNT,PLAYCOUNT,	CHECKSUM,USER_ID
		)VALUES(
			?,?,?,?,?,	?,?,?,?,?,
			?,?,?,?,?,	?,?,?,?,?,
			?,?,?,?,?,	?,?
		)ON DUPLICATE KEY UPDATE 
			BEATMAPSET_ID = VALUES(BEATMAPSET_ID), MODE = VALUES(MODE), MODE_INT = VALUES(MODE_INT), STATUS = VALUES(STATUS), 
			RANKED = VALUES(RANKED), TOTAL_LENGTH = VALUES(TOTAL_LENGTH), MAX_COMBO = VALUES(MAX_COMBO), DIFFICULTY_RATING = VALUES(DIFFICULTY_RATING), 
			VERSION = VALUES(VERSION), 	ACCURACY = VALUES(ACCURACY), AR = VALUES(AR), CS = VALUES(CS), DRAIN = VALUES(DRAIN), BPM = VALUES(BPM),` +
		"`CONVERT` = VALUES(`CONVERT`" + `), COUNT_CIRCLES = VALUES(COUNT_CIRCLES), COUNT_SLIDERS = VALUES(COUNT_SLIDERS), 
			COUNT_SPINNERS = VALUES(COUNT_SPINNERS), DELETED_AT = VALUES(DELETED_AT), 	HIT_LENGTH = VALUES(HIT_LENGTH), 
			IS_SCOREABLE = VALUES(IS_SCOREABLE), LAST_UPDATED = VALUES(LAST_UPDATED), PASSCOUNT = VALUES(PASSCOUNT), PLAYCOUNT = VALUES(PLAYCOUNT), 
			CHECKSUM = VALUES(CHECKSUM), USER_ID = VALUES(USER_ID);`

	setUpsert = `/* UPSERT BEATMAPSET */
		INSERT INTO BEATMAPSET (
			BEATMAPSET_ID,ARTIST,ARTIST_UNICODE,CREATOR,FAVOURITE_COUNT,
			NSFW,PLAY_COUNT,SOURCE,
			STATUS,TITLE,TITLE_UNICODE,USER_ID,VIDEO,
			AVAILABILITY_DOWNLOAD_DISABLED,AVAILABILITY_MORE_INFORMATION,BPM,CAN_BE_HYPED,DISCUSSION_ENABLED,
			DISCUSSION_LOCKED,IS_SCOREABLE,LAST_UPDATED,LEGACY_THREAD_URL,NOMINATIONS_SUMMARY_CURRENT,
			NOMINATIONS_SUMMARY_REQUIRED,RANKED,RANKED_DATE,DELETED_AT,STORYBOARD,SUBMITTED_DATE,
			TAGS,HAS_FAVOURITED )
		VALUES %s ON DUPLICATE KEY UPDATE 
			ARTIST = VALUES(ARTIST), ARTIST_UNICODE = VALUES(ARTIST_UNICODE), CREATOR = VALUES(CREATOR), FAVOURITE_COUNT = VALUES(FAVOURITE_COUNT), 
			NSFW = VALUES(NSFW), PLAY_COUNT = VALUES(PLAY_COUNT), SOURCE = VALUES(SOURCE), 
			STATUS = VALUES(STATUS), TITLE = VALUES(TITLE), TITLE_UNICODE = VALUES(TITLE_UNICODE), USER_ID = VALUES(USER_ID), VIDEO = VALUES(VIDEO), 
			AVAILABILITY_DOWNLOAD_DISABLED = VALUES(AVAILABILITY_DOWNLOAD_DISABLED), AVAILABILITY_MORE_INFORMATION = VALUES(AVAILABILITY_MORE_INFORMATION), 
			BPM = VALUES(BPM), CAN_BE_HYPED = VALUES(CAN_BE_HYPED), DISCUSSION_ENABLED = VALUES(DISCUSSION_ENABLED), 
			DISCUSSION_LOCKED = VALUES(DISCUSSION_LOCKED), IS_SCOREABLE = VALUES(IS_SCOREABLE), LAST_UPDATED = VALUES(LAST_UPDATED), 
			LEGACY_THREAD_URL = VALUES(LEGACY_THREAD_URL), NOMINATIONS_SUMMARY_CURRENT = VALUES(NOMINATIONS_SUMMARY_CURRENT), 
			NOMINATIONS_SUMMARY_REQUIRED = VALUES(NOMINATIONS_SUMMARY_REQUIRED), RANKED = VALUES(RANKED), RANKED_DATE = VALUES(RANKED_DATE), DELETED_AT= VALUES(DELETED_AT),
			STORYBOARD = VALUES(STORYBOARD), SUBMITTED_DATE = VALUES(SUBMITTED_DATE), 
			TAGS = VALUES(TAGS), HAS_FAVOURITED = VALUES(HAS_FAVOURITED);`
	setValues = `(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)` //30

	mapUpsert = `/* UPSERT BEATMAP */
		INSERT INTO BEATMAP (	
			BEATMAP_ID,BEATMAPSET_ID,MODE,MODE_INT,STATUS,	RANKED,TOTAL_LENGTH,MAX_COMBO,DIFFICULTY_RATING,VERSION,
			ACCURACY,AR,CS,DRAIN,BPM,` + "`CONVERT`" + `,COUNT_CIRCLES,COUNT_SLIDERS,COUNT_SPINNERS,DELETED_AT,
			HIT_LENGTH,IS_SCOREABLE,LAST_UPDATED,PASSCOUNT,PLAYCOUNT,	CHECKSUM,USER_ID
		)VALUES %s ON DUPLICATE KEY UPDATE 
			BEATMAPSET_ID = VALUES(BEATMAPSET_ID), MODE = VALUES(MODE), MODE_INT = VALUES(MODE_INT), STATUS = VALUES(STATUS), 
			RANKED = VALUES(RANKED), TOTAL_LENGTH = VALUES(TOTAL_LENGTH), MAX_COMBO = VALUES(MAX_COMBO), 
			DIFFICULTY_RATING = VALUES(DIFFICULTY_RATING), VERSION = VALUES(VERSION), 
			ACCURACY = VALUES(ACCURACY), AR = VALUES(AR), CS = VALUES(CS), DRAIN = VALUES(DRAIN), BPM = VALUES(BPM), 
			` + "`CONVERT` = VALUES(`CONVERT`" + `), COUNT_CIRCLES = VALUES(COUNT_CIRCLES), COUNT_SLIDERS = VALUES(COUNT_SLIDERS),
			COUNT_SPINNERS = VALUES(COUNT_SPINNERS), DELETED_AT = VALUES(DELETED_AT), 
			HIT_LENGTH = VALUES(HIT_LENGTH), IS_SCOREABLE = VALUES(IS_SCOREABLE), LAST_UPDATED = VALUES(LAST_UPDATED), 
			PASSCOUNT = VALUES(PASSCOUNT), PLAYCOUNT = VALUES(PLAYCOUNT), 
			CHECKSUM = VALUES(CHECKSUM), USER_ID = VALUES(USER_ID);`
	mapValues = `(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)` //27

	UpsertBeatmapSet = `
/* UPSERT BEATMAPSET */
INSERT INTO BEATMAPSET(
	BEATMAPSET_ID,ARTIST,ARTIST_UNICODE,CREATOR,FAVOURITE_COUNT,
	HYPE_CURRENT,HYPE_REQUIRED,NSFW,PLAY_COUNT,SOURCE,
	STATUS,TITLE,TITLE_UNICODE,USER_ID,VIDEO,
	AVAILABILITY_DOWNLOAD_DISABLED,AVAILABILITY_MORE_INFORMATION,BPM,CAN_BE_HYPED,DISCUSSION_ENABLED,
	DISCUSSION_LOCKED,IS_SCOREABLE,LAST_UPDATED,LEGACY_THREAD_URL,NOMINATIONS_SUMMARY_CURRENT,
	NOMINATIONS_SUMMARY_REQUIRED,RANKED,RANKED_DATE,DELETED_AT,STORYBOARD,SUBMITTED_DATE,
	TAGS,HAS_FAVOURITED,DESCRIPTION,GENRE_ID,GENRE_NAME,
	LANGUAGE_ID,LANGUAGE_NAME,RATINGS
)VALUES(
	?,?,?,?,?,	?,?,?,?,?,
	?,?,?,?,?,	?,?,?,?,?,
	?,?,?,?,?,	?,?,?,?,?,
	?,?,?,?,?,	?,?,?,?
)ON DUPLICATE KEY UPDATE 
	ARTIST= VALUES(ARTIST), ARTIST_UNICODE= VALUES(ARTIST_UNICODE), CREATOR= VALUES(CREATOR), FAVOURITE_COUNT= VALUES(FAVOURITE_COUNT), 
	HYPE_CURRENT= VALUES(HYPE_CURRENT), HYPE_REQUIRED= VALUES(HYPE_REQUIRED), NSFW= VALUES(NSFW), PLAY_COUNT= VALUES(PLAY_COUNT), SOURCE= VALUES(SOURCE), 
	STATUS= VALUES(STATUS), TITLE= VALUES(TITLE), TITLE_UNICODE= VALUES(TITLE_UNICODE), USER_ID= VALUES(USER_ID), VIDEO= VALUES(VIDEO), 
	AVAILABILITY_DOWNLOAD_DISABLED= VALUES(AVAILABILITY_DOWNLOAD_DISABLED), AVAILABILITY_MORE_INFORMATION= VALUES(AVAILABILITY_MORE_INFORMATION), 
	BPM= VALUES(BPM), CAN_BE_HYPED= VALUES(CAN_BE_HYPED), DISCUSSION_ENABLED= VALUES(DISCUSSION_ENABLED), 
	DISCUSSION_LOCKED= VALUES(DISCUSSION_LOCKED), IS_SCOREABLE= VALUES(IS_SCOREABLE), LAST_UPDATED= VALUES(LAST_UPDATED), LEGACY_THREAD_URL= VALUES(LEGACY_THREAD_URL), 
	NOMINATIONS_SUMMARY_CURRENT= VALUES(NOMINATIONS_SUMMARY_CURRENT), 	NOMINATIONS_SUMMARY_REQUIRED= VALUES(NOMINATIONS_SUMMARY_REQUIRED), 
	RANKED= VALUES(RANKED), RANKED_DATE= VALUES(RANKED_DATE), DELETED_AT= VALUES(DELETED_AT), STORYBOARD= VALUES(STORYBOARD), SUBMITTED_DATE= VALUES(SUBMITTED_DATE), 
	TAGS= VALUES(TAGS), HAS_FAVOURITED= VALUES(HAS_FAVOURITED), DESCRIPTION= VALUES(DESCRIPTION), GENRE_ID= VALUES(GENRE_ID), GENRE_NAME= VALUES(GENRE_NAME), 
	LANGUAGE_ID= VALUES(LANGUAGE_ID), LANGUAGE_NAME= VALUES(LANGUAGE_NAME), RATINGS= VALUES(RATINGS)
;
`
)

func buildSqlValues(s string, count int) (r string) {
	var sbuf []string
	for i := 0; i < count; i++ {
		sbuf = append(sbuf, s)
	}
	return strings.Join(sbuf, ",")
}

func updateSearchBeatmaps(data []osu.BeatmapSetsIN) (err error) {
	if data == nil {
		return
	}
	if len(data) < 1 {
		return
	}
	go db.InsertCache(data)
	var (
		setInsertBuf []interface{}
		mapInsertBuf []interface{}
		beatmapSets  []int
		beatmaps     []int
	)

	for _, s := range data {

		beatmapSets = append(beatmapSets, s.Id)
		setInsertBuf = append(
			setInsertBuf, s.Id, s.Artist, s.ArtistUnicode, s.Creator, s.FavouriteCount, s.Nsfw, s.PlayCount, s.Source, s.Status, s.Title, s.TitleUnicode, s.UserId, s.Video, s.Availability.DownloadDisabled, s.Availability.MoreInformation,
			s.Bpm, s.CanBeHyped, s.DiscussionEnabled, s.DiscussionLocked, s.IsScoreable, s.LastUpdated, s.LegacyThreadUrl, s.NominationsSummary.Current, s.NominationsSummary.Required, s.Ranked, s.RankedDate, s.DeletedAt, s.Storyboard,
			s.SubmittedDate,
			s.Tags, s.HasFavourited,
		)
		for _, m := range *s.Beatmaps {
			beatmaps = append(beatmaps, m.Id)
			mapInsertBuf = append(
				mapInsertBuf, m.Id, m.BeatmapsetId, m.Mode, m.ModeInt, m.Status, m.Ranked, m.TotalLength, m.MaxCombo, m.DifficultyRating, m.Version, m.Accuracy, m.Ar, m.Cs, m.Drain, m.Bpm, m.Convert, m.CountCircles, m.CountSliders,
				m.CountSpinners, m.DeletedAt, m.HitLength, m.IsScoreable, m.LastUpdated, m.Passcount, m.Playcount, m.Checksum, m.UserId,
			)
		}
	}
	//맵셋
	db.InsertQueueChannel <- db.ExecQueue{
		//맵셋
		Query: fmt.Sprintf(setUpsert, buildSqlValues(setValues, len(beatmapSets))),
		Args:  setInsertBuf,
	}

	db.InsertQueueChannel <- db.ExecQueue{
		//맵
		Query: fmt.Sprintf(mapUpsert, buildSqlValues(mapValues, len(beatmaps))),
		Args:  mapInsertBuf,
	}
	db.InsertQueueChannel <- db.ExecQueue{
		//삭제된맵 삭제
		Query: `/* SOFT DELETE MAP */ UPDATE BEATMAP SET DELETED_AT = CURRENT_TIMESTAMP WHERE  BEATMAPSET_ID IN  @mapSets AND BEATMAP_ID NOT IN @maps;`,
		Args:  []any{sql.Named("mapSets", beatmapSets), sql.Named("maps", beatmaps)},
	}

	db.InsertQueueChannel <- db.ExecQueue{
		//삭제된맵 삭제
		Query: ` /* SOFT DELETE MAPSET */ UPDATE BEATMAPSET SET DELETED_AT = CURRENT_TIMESTAMP 
                      WHERE  BEATMAPSET_ID IN (
					      SELECT MS.BEATMAPSET_ID FROM BEATMAPSET MS LEFT JOIN BEATMAP M on MS.BEATMAPSET_ID = M.BEATMAPSET_ID
						WHERE TRUE
						  AND MS.BEATMAPSET_ID IN @mapSets
                          AND MS.DELETED_AT IS NULL
						  AND M.DELETED_AT IS NULL
						GROUP BY M.BEATMAPSET_ID HAVING COUNT(1) < 1
);`,
		Args: []any{sql.Named("mapSets", beatmapSets)},
	}
	return
}
