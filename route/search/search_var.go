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
	maniaKeyRegex, _ = regexp.Compile(`(\[[0-9]K\] )`)

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
		"title_asc":  "title:asc",
		"title":      "title:asc",
		"title asc":  "title:asc",
		"title_desc": "title:desc",
		"title desc": "title:desc",

		"artist_desc": "artist:desc",
		"artist desc": "artist:desc",
		"artist_asc":  "artist:asc",
		"artist":      "artist:asc",
		"artist asc":  "artist:asc",

		"difficulty_rating":      "difficulty_rating:asc",
		"difficulty_rating asc":  "difficulty_rating:asc",
		"difficulty_rating_asc":  "difficulty_rating:asc",
		"difficulty_rating desc": "difficulty_rating:desc",
		"difficulty_rating_desc": "difficulty_rating:desc",

		"favourites_asc":       "favourite_count:asc",
		"favourite_count":      "favourite_count:asc",
		"favourite_count asc":  "favourite_count:asc",
		"favourites_desc":      "favourite_count:desc",
		"favourite_count desc": "favourite_count:desc",

		"plays_asc":       "play_count:asc",
		"play_count":      "play_count:asc",
		"play_count asc":  "play_count:asc",
		"plays_desc":      "play_count:desc",
		"play_count desc": "play_count:desc",
		"play_count_desc": "play_count:desc",

		"ranked_asc":       "ranked_date:asc",
		"ranked_date":      "ranked_date:asc",
		"ranked_date asc":  "ranked_date:asc",
		"ranked_desc":      "ranked_date:desc",
		"ranked_date desc": "ranked_date:desc",

		"last_updated":      "last_updated:asc",
		"last_updated asc":  "last_updated:asc",
		"last_updated desc": "last_updated:desc",
		"updated_asc":       "last_updated:asc",
		"updated_desc":      "last_updated:desc",

		"default": "ranked_date:desc",
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
		"mapid":    1 << 5, // 32
		"m":        1 << 5,
		"setid":    1 << 6, // 64
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
