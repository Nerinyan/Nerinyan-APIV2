package entity

import (
	"time"
)

type BanchoBeatmapEntity struct {
	BeatmapId        int        `json:"beatmap_id" gorm:"column:BEATMAP_ID"`
	BeatmapsetId     int        `json:"beatmapset_id" gorm:"column:BEATMAPSET_ID"`
	Mode             *string    `json:"mode" gorm:"column:MODE"`
	ModeInt          *int       `json:"mode_int" gorm:"column:MODE_INT"`
	Status           *string    `json:"status" gorm:"column:STATUS"`
	Ranked           *int       `json:"ranked" gorm:"column:RANKED"`
	TotalLength      *int       `json:"total_length" gorm:"column:TOTAL_LENGTH"`
	MaxCombo         *int       `json:"max_combo" gorm:"column:MAX_COMBO"`
	DifficultyRating *float64   `json:"difficulty_rating" gorm:"column:DIFFICULTY_RATING"`
	Version          *string    `json:"version" gorm:"column:VERSION"`
	Accuracy         *float64   `json:"accuracy" gorm:"column:ACCURACY"`
	Ar               *float64   `json:"ar" gorm:"column:AR"`
	Cs               *float64   `json:"cs" gorm:"column:CS"`
	Drain            *float64   `json:"drain" gorm:"column:DRAIN"`
	Bpm              *float64   `json:"bpm" gorm:"column:BPM"`
	Convert          *bool      `json:"convert" gorm:"column:CONVERT"`
	CountCircles     *int       `json:"count_circles" gorm:"column:COUNT_CIRCLES"`
	CountSliders     *int       `json:"count_sliders" gorm:"column:COUNT_SLIDERS"`
	CountSpinners    *int       `json:"count_spinners" gorm:"column:COUNT_SPINNERS"`
	HitLength        *int       `json:"hit_length" gorm:"column:HIT_LENGTH"`
	IsScoreable      *bool      `json:"is_scoreable" gorm:"column:IS_SCOREABLE"`
	LastUpdated      *time.Time `json:"last_updated" gorm:"column:LAST_UPDATED"`
	Passcount        *int       `json:"passcount" gorm:"column:PASSCOUNT"`
	Playcount        *int       `json:"playcount" gorm:"column:PLAYCOUNT"`
	Checksum         *string    `json:"checksum" gorm:"column:CHECKSUM"`
	UserId           *int       `json:"user_id" gorm:"column:USER_ID"`
}

func (BanchoBeatmapEntity) TableName() string {
	return "BEATMAP"
}
