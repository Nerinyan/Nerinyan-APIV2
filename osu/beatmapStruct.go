package osu

import "time"

type BeatmapIN struct {
	DifficultyRating float64     `json:"difficulty_rating"`
	Id               int         `json:"id"`
	Mode             *string     `json:"mode"`
	Status           *string     `json:"status"`
	TotalLength      int         `json:"total_length"`
	UserId           int         `json:"user_id"`
	Version          *string     `json:"version"`
	Accuracy         float64     `json:"accuracy"`
	Ar               float64     `json:"ar"`
	BeatmapsetId     int         `json:"beatmapset_id"`
	Bpm              interface{} `json:"bpm"`
	Convert          *bool       `json:"convert"`
	CountCircles     int         `json:"count_circles"`
	CountSliders     int         `json:"count_sliders"`
	CountSpinners    int         `json:"count_spinners"`
	Cs               float64     `json:"cs"`
	DeletedAt        *time.Time  `json:"deleted_at"`
	Drain            float64     `json:"drain"`
	HitLength        int         `json:"hit_length"`
	IsScoreable      *bool       `json:"is_scoreable"`
	LastUpdated      *time.Time  `json:"last_updated"`
	ModeInt          int         `json:"mode_int"`
	Passcount        int         `json:"passcount"`
	Playcount        int         `json:"playcount"`
	Ranked           int         `json:"ranked"`
	Url              *string     `json:"url"`
	Checksum         *string     `json:"checksum"`
	MaxCombo         int         `json:"max_combo"`
}

type BeatmapOUT struct {
	DifficultyRating *float64   `json:"difficulty_rating"`
	Id               *int       `json:"id"`
	Mode             *string    `json:"mode"`
	Status           *string    `json:"status"`
	TotalLength      *int       `json:"total_length"`
	UserId           *int       `json:"user_id"`
	Version          *string    `json:"version"`
	Accuracy         *float64   `json:"accuracy"`
	Ar               *float64   `json:"ar"`
	BeatmapsetId     *int       `json:"beatmapset_id"`
	Bpm              *string    `json:"bpm"`
	Convert          *bool      `json:"convert"`
	CountCircles     *int       `json:"count_circles"`
	CountSliders     *int       `json:"count_sliders"`
	CountSpinners    *int       `json:"count_spinners"`
	Cs               *float64   `json:"cs"`
	DeletedAt        *string    `json:"deleted_at"`
	Drain            *float64   `json:"drain"`
	HitLength        *int       `json:"hit_length"`
	IsScoreable      *bool      `json:"is_scoreable"`
	LastUpdated      *time.Time `json:"last_updated"`
	ModeInt          *int       `json:"mode_int"`
	Passcount        *int       `json:"passcount"`
	Playcount        *int       `json:"playcount"`
	Ranked           *int       `json:"ranked"`
	Url              *string    `json:"url"`
	Checksum         *string    `json:"checksum"`
	MaxCombo         *int       `json:"max_combo"`

	OsuFile string `json:"osu_file"`
}
