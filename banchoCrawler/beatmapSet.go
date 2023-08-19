package banchoCrawler

import (
	"github.com/Nerinyan/Nerinyan-APIV2/entity"
)

type BeatmapSetSearch struct {
	Beatmapsets []entity.Beatmaset `json:"beatmapsets,omitempty"`
	Search      *struct {
		Sort *string `json:"sort,omitempty"`
	} `json:"search,omitempty"`
	RecommendedDifficulty *float64     `json:"recommended_difficulty,omitempty"`
	Error                 *interface{} `json:"error,omitempty"`
	Total                 *int         `json:"total,omitempty"`
	Cursor                *struct {
		ApprovedDate *string `json:"approved_date,omitempty"`
		Id           *string `json:"id,omitempty"`
	} `json:"cursor,omitempty"`
	CursorString *string `json:"cursor_string,omitempty"`
}
