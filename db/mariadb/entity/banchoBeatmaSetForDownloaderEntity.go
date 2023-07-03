package entity

type BanchoBeatmaSetForDownloaderEntity struct {
	BeatmapsetId                 int     `gorm:"column:BEATMAPSET_ID;primaryKey"`
	Artist                       string  `gorm:"column:ARTIST"`
	Title                        string  `gorm:"column:TITLE"`
	LastUpdated                  RFC3339 `gorm:"column:LAST_UPDATED"`
	Video                        bool    `gorm:"column:VIDEO"`
	AvailabilityDownloadDisabled bool    `gorm:"column:AVAILABILITY_DOWNLOAD_DISABLED"` // 조회용
}

func (v BanchoBeatmaSetForDownloaderEntity) TableName() string {
	return "BEATMAPSET"
}
