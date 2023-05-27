package search

import (
	"encoding/base64"
	"github.com/Nerinyan/Nerinyan-APIV2/utils"
	"github.com/goccy/go-json"
	"github.com/pterm/pterm"
	"regexp"
	"strings"
)

var (
	banchoMapRegex, _ = regexp.Compile(`(?:https://osu[.]ppy[.]sh/beatmapsets/)(\d+?)(?:\D|$)`)
	maniaKeyRegex, _  = regexp.Compile(`(\[[0-9]K\] )`)

	regexpReplace, _    = regexp.Compile(`[^0-9A-z]|[\[\]]`)
	regexpByteString, _ = regexp.Compile(`^((0x[\da-fA-F]{1,2})|([\da-fA-F]{1,2})|(1[0-2][0-7]))$`)
	mode                = map[string]int{
		"0": 0, "o": 0, "std": 0, "entity": 0, "entity!": 0, "standard": 0,
		"1": 1, "t": 1, "taiko": 1, "entity!taiko": 1,
		"2": 2, "c": 2, "ctb": 2, "catch": 2, "entity!catch": 2,
		"3": 3, "m": 3, "mania": 3, "entity!mania": 3,
		"default": 0,
	}
	ranked = map[string][]int{
		"ranked":    {1, 2},
		"qualified": {3},
		"loved":     {4},
		"pending":   {0},
		"wip":       {-1},
		"graveyard": {-2},
		"unranked":  {0, -1, -2},
		"-2":        {-2},
		"-1":        {-1},
		"0":         {0},
		"1":         {1},
		"2":         {2},
		"3":         {3},
		"4":         {4},
		"default":   {4, 2, 1},
	}
	orderBy = map[string]string{
		"artist_desc": "ARTIST DESC",
		"artist desc": "ARTIST DESC",
		"artist_asc":  "ARTIST",
		"artist":      "ARTIST",
		"artist asc":  "ARTIST",

		"favourites_asc":       "FAVOURITE_COUNT",
		"favourite_count":      "FAVOURITE_COUNT",
		"favourite_count asc":  "FAVOURITE_COUNT",
		"favourites_desc":      "FAVOURITE_COUNT DESC",
		"favourite_count desc": "FAVOURITE_COUNT DESC",

		"plays_asc":       "PLAY_COUNT",
		"play_count":      "PLAY_COUNT",
		"play_count asc":  "PLAY_COUNT",
		"plays_desc":      "PLAY_COUNT DESC",
		"play_count desc": "PLAY_COUNT DESC",

		"ranked_asc":       "RANKED_DATE",
		"ranked_date":      "RANKED_DATE",
		"ranked_date asc":  "RANKED_DATE",
		"ranked_desc":      "RANKED_DATE DESC",
		"ranked_date desc": "RANKED_DATE DESC",

		"last_updated":      "LAST_UPDATED",
		"last_updated asc":  "LAST_UPDATED",
		"last_updated desc": "LAST_UPDATED DESC",
		"updated_asc":       "LAST_UPDATED",
		"updated_desc":      "LAST_UPDATED DESC",

		"title_asc":  "TITLE",
		"title":      "TITLE",
		"title asc":  "TITLE",
		"title_desc": "TITLE DESC",
		"title desc": "TITLE DESC",
		"default":    "RANKED_DATE DESC",
	}
	searchOption = map[string]uint32{
		"artist":   1 << 0, // 1
		"a":        1 << 0,
		"creator":  1 << 1, // 2
		"c":        1 << 1,
		"tag":      1 << 2, // 4
		"tg":       1 << 2,
		"title":    1 << 3, // 8
		"t":        1 << 3,
		"checksum": 1 << 4, // 16
		"cks":      1 << 4,
		"mapId":    1 << 5, // 32
		"m":        1 << 5,
		"setId":    1 << 6, // 64
		"s":        1 << 6,
		"default":  0xFFFF, // all
	}
)

type minMax struct {
	Min float32 `json:"min"`
	Max float32 `json:"max"`
}

type SearchQuery struct {
	// global
	Extra string `query:"e" json:"extra"` // 스토리보드 비디오.

	// set
	Ranked     string `query:"s" json:"ranked"`      // 랭크상태   set.ranked
	Nsfw       bool   `query:"nsfw" json:"nsfw"`     // R18       set.nsfw
	Video      bool   `query:"v" json:"video"`       // 비디오     set.video
	Storyboard bool   `query:"sb" json:"storyboard"` // 스토리보드  set.storyboard

	// map
	Mode             string `query:"m" json:"m"`      // 게임모드			    map.mode_int
	TotalLength      minMax `json:"totalLength"`      // 플레이시간			map.totalLength
	MaxCombo         minMax `json:"maxCombo"`         // 콤보				map.maxCombo
	DifficultyRating minMax `json:"difficultyRating"` // 난이도				map.difficultyRating
	Accuracy         minMax `json:"accuracy"`         // od					map.accuracy
	Ar               minMax `json:"ar"`               // ar					map.ar
	Cs               minMax `json:"cs"`               // cs					map.cs
	Drain            minMax `json:"drain"`            // hp					map.drain
	Bpm              minMax `json:"bpm"`              // bpm				map.bpm

	// query
	Sort       string   `query:"sort" json:"sort"`   // 정렬	  order by
	Page       int      `query:"p" json:"page"`      // 페이지 limit
	PageSize   int      `query:"ps" json:"pageSize"` // 페이지 당 크기
	Text       string   `query:"q" json:"query"`     // 문자열 검색
	ParsedText []string `json:"-"`                   // 문자열 검색 파싱 내부 사용용
	Option     string   `query:"option" json:"option"`
	OptionB    uint32   `json:"-"`    //artist 1,creator 2,tags 4 ,title 8
	B64        string   `query:"b64"` // body
}

//	func (v *SearchQuery) getVideo() (allow bool) {
//		sb := v.Video
//		switch sb.(type) {
//		case bool:
//			allow = sb.(bool)
//		case string:
//			allow, _ = strconv.ParseBool(sb.(string))
//			allow = allow || (sb.(string) == "all")
//		}
//		return allow || strings.Contains(utils.TrimLower(v.Extra), "video")
//
// }
func (v *SearchQuery) getVideo() (allow bool) {
	return v.Video || strings.Contains(utils.TrimLower(v.Extra), "video")
}

//	func (v *SearchQuery) getStoryboard() (allow bool) {
//		sb := v.Storyboard
//		switch sb.(type) {
//		case bool:
//			allow = sb.(bool)
//		case string:
//			allow, _ = strconv.ParseBool(sb.(string))
//			allow = allow || (sb.(string) == "all")
//		}
//		return allow || strings.Contains(utils.TrimLower(v.Extra), "storyboard")
//	}
func (v *SearchQuery) getStoryboard() (allow bool) {
	return v.Storyboard || strings.Contains(utils.TrimLower(v.Extra), "storyboard")
}

func (v *SearchQuery) getPage() (page int) {
	return utils.IntMin(utils.ToInt(v.Page), 0)
}

func (v *SearchQuery) getPageSize() int {
	return utils.IntMinMaxDefault(utils.ToInt(v.PageSize), 1, 1000, 50)
}

//	func (v *SearchQuery) getNsfw() (allow bool) {
//		if v.Nsfw != nil {
//			if n, ok := (v.Nsfw).(bool); ok {
//				return n
//			}
//			if n, ok := (v.Nsfw).(string); ok {
//				allow, _ = strconv.ParseBool(n)
//				allow = allow || n == "all"
//				return
//			}
//		}
//		return
//	}
func (v *SearchQuery) getNsfw() (allow bool) {
	return v.Nsfw
}

func (v *SearchQuery) parseOption() uint32 {
	ss := strings.ToLower(v.Option)
	if ss == "" {
		v.OptionB |= 0xFFFFFFFF
		return v.OptionB
	}
	for _, s2 := range strings.Split(ss, ",") {
		v.OptionB |= searchOption[s2]
	}
	if v.OptionB == 0 {
		v.OptionB = 0xFFFFFFFF
	}
	return v.OptionB
}

func (v *SearchQuery) parseB64() (err error) {
	if v.B64 != "" {
		b6, err := base64.StdEncoding.DecodeString(v.B64)
		if err != nil {
			pterm.Error.WithShowLineNumber().Println(err.Error())
			return err
		}
		err = json.Unmarshal(b6, &v)
		if err != nil {
			pterm.Error.WithShowLineNumber().Println(err.Error())
			return err
		}
	}
	return
}
