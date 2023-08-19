package ms

import (
	"fmt"
	"github.com/Nerinyan/Nerinyan-APIV2/entity"
	"github.com/meilisearch/meilisearch-go"
	"strconv"
)

func FindBeatmapSetById(id string) (res *entity.Beatmaset, err error) {
	idi, _ := strconv.Atoi(id)
	return FindBeatmapSetByIdInt(idi)
}
func FindBeatmapSetByIdInt(id int) (res *entity.Beatmaset, err error) {
	searchRes, err := GetIndexBeatMapSet().Search("", &meilisearch.SearchRequest{
		AttributesToRetrieve: nil,
		AttributesToCrop:     nil,
		CropLength:           0,
		CropMarker:           "",
		MatchingStrategy:     "all",
		Filter:               fmt.Sprintf("id = %d", id),
		ShowMatchesPosition:  false,
	})
	if err != nil {
		return
	}
	var sets []entity.Beatmaset
	err = HitsMapper(searchRes.Hits, &sets)
	if err != nil {
		return
	}
	for _, set := range sets {
		if set.Id == id {
			res = &set
			return
		}
	}
	return
}
