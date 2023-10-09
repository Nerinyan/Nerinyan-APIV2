package osu

type BeatmapsetsSearch struct {
	Beatmapsets *[]BeatmapSetsIN `json:"beatmapsets"`
	Search      *struct {
		Sort *string `json:"sort"`
	} `json:"search"`
	//Cursor      *struct {
	//	LastUpdate *int64 `json:"last_update"`
	//	Id         *string `json:"_id"`
	//} `json:"cursor"`

	CursorString string `json:"cursor_string"`

	RecommendedDifficulty float64      `json:"recommended_difficulty"`
	Error                 *interface{} `json:"error"`
	Total                 int          `json:"total"`
}
