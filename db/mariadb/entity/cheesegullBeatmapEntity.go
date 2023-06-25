package entity

type CheesegoolBeatmapEntity struct {
	BeatmapId        *int     `json:"BeatmapID" gorm:"column:BEATMAP_ID"`
	ParentSetId      *int     `json:"ParentSetID" gorm:"column:BEATMAPSET_ID"`
	DiffName         *string  `json:"DiffName" gorm:"column:VERSION"`
	FileMd5          *string  `json:"FileMD5" gorm:"column:CHECKSUM"`
	Mode             *int     `json:"Mode" gorm:"column:MODE_INT"`
	Bpm              *float64 `json:"BPM" gorm:"column:BPM"`
	Ar               *float64 `json:"AR" gorm:"column:AR"`
	Od               *float64 `json:"OD" gorm:"column:ACCURACY"`
	Cs               *float64 `json:"CS" gorm:"column:CS"`
	Hp               *float64 `json:"HP" gorm:"column:DRAIN"`
	TotalLength      *int     `json:"TotalLength" gorm:"column:TOTAL_LENGTH"`
	HitLength        *int     `json:"HitLength" gorm:"column:HIT_LENGTH"`
	Playcount        *int     `json:"Playcount" gorm:"column:PLAYCOUNT"`
	Passcount        *int     `json:"Passcount" gorm:"column:PASSCOUNT"`
	MaxCombo         *int     `json:"MaxCombo" gorm:"column:MAX_COMBO"`
	DifficultyRating *float64 `json:"DifficultyRating" gorm:"column:DIFFICULTY_RATING"`
}

func (CheesegoolBeatmapEntity) TableName() string {
	return "BEATMAP"
}
