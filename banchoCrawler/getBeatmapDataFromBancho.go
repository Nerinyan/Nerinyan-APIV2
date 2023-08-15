package banchoCrawler

import (
	"fmt"
	"github.com/Nerinyan/Nerinyan-APIV2/config"
	"github.com/Nerinyan/Nerinyan-APIV2/db/meilisearch"
	"github.com/goccy/go-json"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"io"
	"net/http"
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
			pterm.Error.WithShowLineNumber().Println(err)
		}
	}()
	url := fmt.Sprintf("https://osu.ppy.sh/api/v2/beatmapsets/search?nsfw=true&s=any&q=%d", id)

	var data BeatmapSetSearch
	if err = stdGETBancho(url, &data); err != nil {
		return
	}
	// updateMapset(&data)

	if err = updateSearchBeatmaps(data.Beatmapsets); err != nil {
		return
	}
}

func getUpdatedMapRanked() {
	var err error
	defer func() {
		if err != nil {
			pterm.Error.WithShowLineNumber().Println(err)
		}
	}()
	url := "https://osu.ppy.sh/api/v2/beatmapsets/search?nsfw=true&s=ranked"

	var data BeatmapSetSearch

	if err = stdGETBancho(url, &data); err != nil {
		return
	}
	//pterm.Info.Println("getUpdatedMapRanked", string(*utils.ToJsonString(*data.Beatmapsets))[:100])

	if err = updateSearchBeatmaps(data.Beatmapsets); err != nil {
		return
	}

	return
}

func getUpdatedMapLoved() {
	var err error
	defer func() {
		if err != nil {
			pterm.Error.WithShowLineNumber().Println(err)
		}
	}()
	url := "https://osu.ppy.sh/api/v2/beatmapsets/search?nsfw=true&s=loved"

	var data BeatmapSetSearch
	if err = stdGETBancho(url, &data); err != nil {
		return
	}
	//pterm.Info.Println("getUpdatedMapLoved", string(*utils.ToJsonString(*data.Beatmapsets))[:100])
	if err = updateSearchBeatmaps(data.Beatmapsets); err != nil {
		return
	}
	return
}
func getUpdatedMapQualified() {
	var err error
	defer func() {
		if err != nil {
			pterm.Error.WithShowLineNumber().Println(err)
		}
	}()
	url := "https://osu.ppy.sh/api/v2/beatmapsets/search?nsfw=true&s=qualified"

	var data BeatmapSetSearch
	if err = stdGETBancho(url, &data); err != nil {
		return
	}
	//pterm.Info.Println("getUpdatedMapQualified", string(*utils.ToJsonString(*data.Beatmapsets))[:100])
	if err = updateSearchBeatmaps(data.Beatmapsets); err != nil {
		return
	}
	return
}

func getGraveyardMap() {
	var err error
	defer func() {
		if err != nil {
			pterm.Error.WithShowLineNumber().Println(err)
		}
	}()
	url := ""
	cs := &config.Config.Osu.BeatmapUpdate.GraveyardAsc.CursorString
	if *cs != "" {
		url = "https://osu.ppy.sh/api/v2/beatmapsets/search?nsfw=true&sort=updated_asc&s=graveyard&cursor_string=" + *cs
	} else {
		url = "https://osu.ppy.sh/api/v2/beatmapsets/search?nsfw=true&sort=updated_asc&s=graveyard"
	}

	var data BeatmapSetSearch

	err = stdGETBancho(url, &data)
	if err != nil {
		return
	}
	if *data.CursorString == "" {
		return
	}
	//pterm.Info.Println("getGraveyardMap", string(*utils.ToJsonString(*data.Beatmapsets))[:100])
	if err = updateSearchBeatmaps(data.Beatmapsets); err != nil {
		return
	}
	cs = data.CursorString
	return
}

func getUpdatedMapDesc() {
	var err error
	defer func() {
		if err != nil {
			pterm.Error.WithShowLineNumber().Println(err)
		}
	}()
	url := "https://osu.ppy.sh/api/v2/beatmapsets/search?nsfw=true&sort=updated_desc&s=any"

	var data BeatmapSetSearch

	if err = stdGETBancho(url, &data); err != nil {
		return
	}

	//pterm.Info.Println("getUpdatedMapDesc", string(*utils.ToJsonString(*data.Beatmapsets))[:100])
	if err = updateSearchBeatmaps(data.Beatmapsets); err != nil {
		return
	}
	if *data.CursorString == "" {
		return
	}
	config.Config.Osu.BeatmapUpdate.UpdatedDesc.CursorString = *data.CursorString
	return
}

func getUpdatedMapAsc() {
	var err error
	defer func() {
		if err != nil {
			pterm.Error.WithShowLineNumber().Println(err)
		}
	}()
	url := ""
	cs := &config.Config.Osu.BeatmapUpdate.UpdatedAsc.CursorString
	if *cs != "" {
		url = "https://osu.ppy.sh/api/v2/beatmapsets/search?nsfw=true&sort=updated_asc&s=any&cursor_string=" + *cs
	} else {
		url = "https://osu.ppy.sh/api/v2/beatmapsets/search?nsfw=true&sort=updated_asc&s=any"
	}

	var data BeatmapSetSearch

	err = stdGETBancho(url, &data)
	if err != nil {
		return
	}

	//pterm.Info.Println(url, data.CursorString)
	//pterm.Info.Println(data.CursorString, url, string(*utils.ToJsonString(*data.Beatmapsets))[:200])
	if err = updateSearchBeatmaps(data.Beatmapsets); err != nil {
		return
	}
	cs = data.CursorString
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
	//pterm.Info.Println("X-Ratelimit-Remaining", res.Header.Get("X-Ratelimit-Remaining"))
	return json.Unmarshal(body, &str)
	// return json.NewDecoder(res.Body).Decode(&str)

}

func updateSearchBeatmaps(data []Beatmasets) (err error) {

	if data == nil {
		return
	}
	if len(data) < 1 {
		return
	}
	_, err = meilisearch.GetIndexBeatMapSet().UpdateDocuments(data, "id")
	if err != nil {
		pterm.Error.WithShowLineNumber().Println(err)
		return err
	}
	//pterm.Info.Println(res)
	return
}
