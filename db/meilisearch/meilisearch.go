package meilisearch

import (
	"github.com/Nerinyan/Nerinyan-APIV2/config"
	"github.com/meilisearch/meilisearch-go"
	"github.com/pterm/pterm"
	"time"
)

var meiliSearch *meilisearch.Client

const INDEX_BEATMAP_SET = "BEATMAP_SET"

func InitMeiliSearch() {
	meiliSearch = meilisearch.NewClient(meilisearch.ClientConfig{
		Host:    config.Config.MeiliSearch.Host,
		APIKey:  config.Config.MeiliSearch.APIKey,
		Timeout: time.Second * time.Duration(config.Config.MeiliSearch.Timeout),
	})
	if meiliSearch.IsHealthy() {
		pterm.Success.Println("Connected meiliSearch")
	}

	res, err := GetIndexBeatMapSet().UpdateSortableAttributes(&[]string{
		"ranked_date",
		"artist",
		"favourite_count",
		"play_count",
		"last_updated",
		"title",
		"difficulty_rating",
	})
	if err != nil {
		pterm.Error.WithShowLineNumber().Println(err)
	}
	pterm.Info.Println(res)
	res, err = GetIndexBeatMapSet().UpdateFilterableAttributes(&[]string{
		"id",
		"ranked",
		"nsfw",
		"video",
		"storyboard",
		// =========================
		"beatmaps.ar",
		"beatmaps.accuracy",
		"beatmaps.total_length",
		"beatmaps.bpm",
		"beatmaps.cs",
		"beatmaps.difficulty_rating",
		"beatmaps.total_length",
		"beatmaps.drain",
		"beatmaps.mode",
		"beatmaps.drain",
		"beatmaps.id",
		"beatmaps.mode_int",
		"beatmaps.max_combo",
	})
	if err != nil {
		pterm.Error.WithShowLineNumber().Println(err)
	}
	pterm.Info.Println(res)

}

func GetIndexBeatMapSet() *meilisearch.Index {
	return meiliSearch.Index(INDEX_BEATMAP_SET)
}
