package meilisearch

import "github.com/goccy/go-json"

func HitsMapper[T any](hits []interface{}, target *T) (err error) {
	b, err := json.Marshal(hits)
	if err != nil {
		return
	}
	return json.Unmarshal(b, target)
}
