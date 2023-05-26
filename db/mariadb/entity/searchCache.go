package entity

type SearchCache struct {
	Text         string `gorm:"column:TEXT;primaryKey"`
	BeatmapSetID int    `gorm:"column:BEATMAPSET_ID;primaryKey"`
	SearchOption int    `gorm:"column:SEARCH_OPTION;default:0"`
}

func (SearchCache) TableName() string {
	return "SEARCH_CACHE"
}
