package db

import (
	"github.com/Nerinyan/Nerinyan-APIV2/osu"
	"github.com/Nerinyan/Nerinyan-APIV2/utils"
	"github.com/pterm/pterm"
	"github.com/surgebase/porter2"
	"regexp"
	"strconv"
	"strings"
)

func init() {
	cacheChannel = make(chan []osu.BeatmapSetsIN, 100)
	go func() {
		for ins := range cacheChannel {
			insertStringIndex(ins)
		}
	}()

}

var (
	regexpReplace, _ = regexp.Compile(`[^0-9A-z]|[\[\]]`)
	cacheChannel     chan []osu.BeatmapSetsIN
)

func InsertCache(data []osu.BeatmapSetsIN) {
	cacheChannel <- data

}

type insertData struct {
	Strbuf  []string
	Artist  []row
	Creator []row
	Title   []row
	Tags    []row
	Other   []row
}
type row struct {
	KEY          []string
	BeatmapsetId int
}

func insertStringIndex(data []osu.BeatmapSetsIN) {
	defer func() {
		err, e := recover().(error)
		if e {
			pterm.Error.Println(err)
		}
	}()
	var insertDataa insertData
	//pterm.Info.Println(unsafe.Pointer(&insertDataa), string(*utils.ToJsonString(insertDataa)))

	for _, in := range data {

		artist := splitStringUnique(*in.Artist)
		artist = append(artist, findRepeats(artist)...)

		creator := splitStringUnique(*in.Creator)
		creator = append(creator, findRepeats(creator)...)

		title := splitStringUnique(*in.Title)
		title = append(title, findRepeats(title)...)

		tags := splitStringUnique(*in.Tags)
		tags = append(tags, findRepeats(tags)...)

		insertDataa.Artist = append(insertDataa.Artist, row{
			KEY:          artist,
			BeatmapsetId: in.Id,
		})
		insertDataa.Creator = append(insertDataa.Creator, row{
			KEY:          creator,
			BeatmapsetId: in.Id,
		})
		insertDataa.Title = append(insertDataa.Title, row{
			KEY:          title,
			BeatmapsetId: in.Id,
		})
		insertDataa.Tags = append(insertDataa.Tags, row{
			KEY:          tags,
			BeatmapsetId: in.Id,
		})
		insertDataa.Strbuf = append(insertDataa.Strbuf, artist...)
		insertDataa.Strbuf = append(insertDataa.Strbuf, creator...)
		insertDataa.Strbuf = append(insertDataa.Strbuf, title...)
		insertDataa.Strbuf = append(insertDataa.Strbuf, tags...)

		insertDataa.Strbuf = append(insertDataa.Strbuf, strconv.Itoa(in.Id))
		for _, beatmapIN := range *in.Beatmaps {
			other := splitStringUnique(*beatmapIN.Version)
			other = append(other, findRepeats(other)...)

			insertDataa.Strbuf = append(insertDataa.Strbuf, other...)
			insertDataa.Strbuf = append(insertDataa.Strbuf, strconv.Itoa(beatmapIN.Id))
			insertDataa.Other = append(insertDataa.Other, row{
				KEY:          []string{strconv.Itoa(beatmapIN.Id)},
				BeatmapsetId: beatmapIN.BeatmapsetId,
			})
			insertDataa.Other = append(insertDataa.Other, row{
				KEY:          []string{strconv.Itoa(beatmapIN.BeatmapsetId)},
				BeatmapsetId: beatmapIN.BeatmapsetId,
			})
			insertDataa.Other = append(insertDataa.Other, row{
				KEY:          other,
				BeatmapsetId: beatmapIN.BeatmapsetId,
			})
		}

	}
	if idata := utils.MakeArrayUniqueInterface(&insertDataa.Strbuf); len(idata) > 0 {
		AddInsertQueue("/* INSERT SEARCH_CACHE_STRING_INDEX */ INSERT INTO SEARCH_CACHE_STRING_INDEX (STRING) VALUES "+
			utils.StringRepeatJoin("(?)", ",", len(idata))+" ON DUPLICATE KEY UPDATE TMP = 1;",
			idata...,
		)
	}

	if idata := toIndexKV(insertDataa.Artist); len(idata) > 0 {
		AddInsertQueue("/* INSERT SEARCH_CACHE_ARTIST */ INSERT INTO SEARCH_CACHE_ARTIST (INDEX_KEY,BEATMAPSET_ID) VALUES "+
			utils.StringRepeatJoin("((SELECT ID FROM SEARCH_CACHE_STRING_INDEX WHERE `STRING` = ?), ?)", ",", len(idata)/2)+
			" ON DUPLICATE KEY UPDATE TMP = 1;", idata...,
		)
	}

	if idata := toIndexKV(insertDataa.Title); len(idata) > 0 {
		AddInsertQueue("/* INSERT SEARCH_CACHE_TITLE */ INSERT INTO SEARCH_CACHE_TITLE (INDEX_KEY,BEATMAPSET_ID) VALUES "+
			utils.StringRepeatJoin("((SELECT ID FROM SEARCH_CACHE_STRING_INDEX WHERE `STRING` = ?), ?)", ",", len(idata)/2)+
			" ON DUPLICATE KEY UPDATE TMP = 1;", idata...,
		)
	}

	if idata := toIndexKV(insertDataa.Creator); len(idata) > 0 {
		AddInsertQueue("/* INSERT SEARCH_CACHE_CREATOR */ INSERT INTO SEARCH_CACHE_CREATOR (INDEX_KEY,BEATMAPSET_ID) VALUES "+
			utils.StringRepeatJoin("((SELECT ID FROM SEARCH_CACHE_STRING_INDEX WHERE `STRING` = ?), ?)", ",", len(idata)/2)+
			" ON DUPLICATE KEY UPDATE TMP = 1;", idata...,
		)
	}
	if idata := toIndexKV(insertDataa.Tags); len(idata) > 0 {
		AddInsertQueue("/* INSERT SEARCH_CACHE_TAG */ INSERT INTO SEARCH_CACHE_TAG (INDEX_KEY,BEATMAPSET_ID) VALUES "+
			utils.StringRepeatJoin("((SELECT ID FROM SEARCH_CACHE_STRING_INDEX WHERE `STRING` = ?), ?)", ",", len(idata)/2)+
			" ON DUPLICATE KEY UPDATE TMP = 1;", idata...,
		)
	}
	if idata := toIndexKV(insertDataa.Other); len(idata) > 0 {
		AddInsertQueue("/* INSERT SEARCH_CACHE_OTHER */ INSERT INTO SEARCH_CACHE_OTHER (INDEX_KEY,BEATMAPSET_ID) VALUES "+
			utils.StringRepeatJoin("((SELECT ID FROM SEARCH_CACHE_STRING_INDEX WHERE `STRING` = ?), ?)", ",", len(idata)/2)+
			" ON DUPLICATE KEY UPDATE TMP = 1;", idata...,
		)
	}
}

func toIndexKV(data []row) (AA []interface{}) {
	for _, A := range data {

		for _, K := range A.KEY {
			AA = append(AA, K, A.BeatmapsetId)
		}
	}
	return
}

func splitString(input string) (ss []string) {
	for _, s := range strings.Split(strings.ToLower(regexpReplace.ReplaceAllString(input, " ")), " ") {
		if s == "" || s == " " {
			continue
		}
		ss = append(ss, s, porter2.Stem(s))
	}
	return
}

func splitStringUnique(input string) (ss []string) {
	ip := splitString(input)
	return utils.MakeStringArrayUniqueAndCheckLength(&ip, 254)
}
func findRepeats(s []string) (res []string) {
	for _, s2 := range s {
		res = append(res, findRepeat(s2)...)
	}
	return
}

func findRepeat(s string) []string {
	var result []string
	checkMap := make(map[string]int)
	length := len(s)
	for i := 0; i < length; i++ {
		for j := i + 2; j <= length; j++ {
			substr := s[i:j]
			v, ok := checkMap[substr]
			if ok {
				if v == 1 {
					result = append(result, substr)
				}
				checkMap[substr]++
			} else {
				checkMap[substr] = 1
			}
		}
	}

	// Remove smaller substrings if they're part of a larger substring
	length = len(result)
	for i := 0; i < length; i++ {
		for j := i + 1; j < length; j++ {
			if strings.Contains(result[i], result[j]) {
				result[j] = result[i]
			} else if strings.Contains(result[j], result[i]) {
				result[i] = result[j]
			}
		}
	}
	// Deduplicate the result
	keys := make(map[string]bool)
	var list []string
	for _, entry := range result {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
