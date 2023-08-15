package banchoCrawler

import "time"

type BeatmapSetSearch struct {
	Beatmapsets []Beatmasets `json:"beatmapsets,omitempty"`
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
type Beatmasets struct {
	Artist        *string `json:"artist,omitempty"`
	ArtistUnicode *string `json:"artist_unicode,omitempty"`
	Covers        *struct {
		Cover       *string `json:"cover,omitempty"`
		Cover2X     *string `json:"cover@2x,omitempty"`
		Card        *string `json:"card,omitempty"`
		Card2X      *string `json:"card@2x,omitempty"`
		List        *string `json:"list,omitempty"`
		List2X      *string `json:"list@2x,omitempty"`
		Slimcover   *string `json:"slimcover,omitempty"`
		Slimcover2X *string `json:"slimcover@2x,omitempty"`
	} `json:"covers,omitempty"`
	Creator        *string `json:"creator,omitempty"`
	FavouriteCount int     `json:"favourite_count,omitempty"`
	Hype           *struct {
		Current  *int `json:"current,omitempty"`
		Required *int `json:"required,omitempty"`
	} `json:"hype,omitempty"`
	Id                 int        `json:"id,omitempty"`
	Nsfw               *bool      `json:"nsfw,omitempty"`
	Offset             *int       `json:"offset,omitempty"`
	PlayCount          *int       `json:"play_count,omitempty"`
	PreviewUrl         *string    `json:"preview_url,omitempty"`
	Source             *string    `json:"source,omitempty"`
	Spotlight          *bool      `json:"spotlight,omitempty"`
	Status             *string    `json:"status,omitempty"`
	Title              *string    `json:"title,omitempty"`
	TitleUnicode       *string    `json:"title_unicode,omitempty"`
	TrackId            *int       `json:"track_id,omitempty"`
	UserId             *int       `json:"user_id,omitempty"`
	Video              *bool      `json:"video,omitempty"`
	Bpm                *float64   `json:"bpm,omitempty"`
	CanBeHyped         *bool      `json:"can_be_hyped,omitempty"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty"`
	DiscussionEnabled  *bool      `json:"discussion_enabled,omitempty"`
	DiscussionLocked   *bool      `json:"discussion_locked,omitempty"`
	IsScoreable        *bool      `json:"is_scoreable,omitempty"`
	LastUpdated        *time.Time `json:"last_updated,omitempty"`
	LegacyThreadUrl    *string    `json:"legacy_thread_url,omitempty"`
	NominationsSummary *struct {
		Current  *int `json:"current,omitempty"`
		Required *int `json:"required,omitempty"`
	} `json:"nominations_summary,omitempty"`
	Ranked        *int       `json:"ranked,omitempty"`
	RankedDate    *time.Time `json:"ranked_date,omitempty"`
	Storyboard    *bool      `json:"storyboard,omitempty"`
	SubmittedDate *time.Time `json:"submitted_date,omitempty"`
	Tags          *string    `json:"tags,omitempty"`
	Availability  *struct {
		DownloadDisabled *bool   `json:"download_disabled,omitempty"`
		MoreInformation  *string `json:"more_information,omitempty"`
	} `json:"availability,omitempty"`
	HasFavourited *bool      `json:"has_favourited,omitempty"`
	Beatmaps      []Beatmaps `json:"beatmaps,omitempty"`
	PackTags      []string   `json:"pack_tags,omitempty"`
}
type Beatmaps struct {
	BeatmapsetId     *int       `json:"beatmapset_id,omitempty"`
	DifficultyRating *float64   `json:"difficulty_rating,omitempty"`
	Id               int        `json:"id,omitempty"`
	Mode             *string    `json:"mode,omitempty"`
	Status           *string    `json:"status,omitempty"`
	TotalLength      *int       `json:"total_length,omitempty"`
	UserId           *int       `json:"user_id,omitempty"`
	Version          *string    `json:"version,omitempty"`
	Accuracy         *float64   `json:"accuracy,omitempty"`
	Ar               *float64   `json:"ar,omitempty"`
	Bpm              *float64   `json:"bpm,omitempty"`
	Convert          *bool      `json:"convert,omitempty"`
	CountCircles     *int       `json:"count_circles,omitempty"`
	CountSliders     *int       `json:"count_sliders,omitempty"`
	CountSpinners    *int       `json:"count_spinners,omitempty"`
	Cs               *float64   `json:"cs,omitempty"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
	Drain            *float64   `json:"drain,omitempty"`
	HitLength        *int       `json:"hit_length,omitempty"`
	IsScoreable      *bool      `json:"is_scoreable,omitempty"`
	LastUpdated      *time.Time `json:"last_updated,omitempty"`
	ModeInt          *int       `json:"mode_int,omitempty"`
	Passcount        *int       `json:"passcount,omitempty"`
	Playcount        *int       `json:"playcount,omitempty"`
	Ranked           *int       `json:"ranked,omitempty"`
	Url              *string    `json:"url,omitempty"`
	Checksum         *string    `json:"checksum,omitempty"`
	MaxCombo         *int       `json:"max_combo,omitempty"`
}
