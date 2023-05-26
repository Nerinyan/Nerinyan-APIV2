package osu

import "time"

type BeatmapSetsIN struct {
	Artist        *string `json:"artist"`
	ArtistUnicode *string `json:"artist_unicode"`
	Covers        *struct {
		Cover       *string `json:"cover"`
		Cover2X     *string `json:"cover@2x"`
		Card        *string `json:"card"`
		Card2X      *string `json:"card@2x"`
		List        *string `json:"list"`
		List2X      *string `json:"list@2x"`
		Slimcover   *string `json:"slimcover"`
		Slimcover2X *string `json:"slimcover@2x"`
	} `json:"covers"`
	Creator        *string `json:"creator"`
	FavouriteCount int     `json:"favourite_count"`
	Hype           *struct {
		Current  *int `json:"current"`
		Required *int `json:"required"`
	} `json:"hype"`
	Id           int     `json:"id"`
	Nsfw         bool    `json:"nsfw"`
	PlayCount    int     `json:"play_count"`
	PreviewUrl   *string `json:"preview_url"`
	Source       *string `json:"source"`
	Status       *string `json:"status"`
	Title        *string `json:"title"`
	TitleUnicode *string `json:"title_unicode"`
	UserId       int     `json:"user_id"`
	Video        bool    `json:"video"`
	Availability *struct {
		DownloadDisabled *bool   `json:"download_disabled"`
		MoreInformation  *string `json:"more_information"`
	} `json:"availability"`
	Bpm                float64    `json:"bpm"`
	CanBeHyped         bool       `json:"can_be_hyped"`
	DiscussionEnabled  bool       `json:"discussion_enabled"`
	DiscussionLocked   bool       `json:"discussion_locked"`
	IsScoreable        bool       `json:"is_scoreable"`
	LastUpdated        *time.Time `json:"last_updated"`
	LegacyThreadUrl    *string    `json:"legacy_thread_url"`
	NominationsSummary *struct {
		Current  int `json:"current"`
		Required int `json:"required"`
	} `json:"nominations_summary"`
	Ranked        int          `json:"ranked"`
	RankedDate    *time.Time   `json:"ranked_date"`
	DeletedAt     *time.Time   `json:"deleted_at"`
	Storyboard    bool         `json:"storyboard"`
	SubmittedDate *time.Time   `json:"submitted_date"`
	Tags          *string      `json:"tags"`
	HasFavourited bool         `json:"has_favourited"`
	Beatmaps      *[]BeatmapIN `json:"beatmaps"`

	Description *struct {
		Description *string `json:"description"`
	} `json:"description"`
	Genre *struct {
		Id   int     `json:"id"`
		Name *string `json:"name"`
	} `json:"genre"`
	Language *struct {
		Id   int     `json:"id"`
		Name *string `json:"name"`
	} `json:"language"`
	Ratings       *[]int `json:"ratings"`
	RatingsString *[]int `json:"ratings_string"`
	User          *struct {
		AvatarUrl     *string      `json:"avatar_url"`
		CountryCode   *string      `json:"country_code"`
		DefaultGroup  *string      `json:"default_group"`
		Id            int          `json:"id"`
		IsActive      bool         `json:"is_active"`
		IsBot         bool         `json:"is_bot"`
		IsDeleted     bool         `json:"is_deleted"`
		IsOnline      bool         `json:"is_online"`
		IsSupporter   bool         `json:"is_supporter"`
		LastVisit     *time.Time   `json:"last_visit"`
		PmFriendsOnly bool         `json:"pm_friends_only"`
		ProfileColour *interface{} `json:"profile_colour"`
		Username      *string      `json:"username"`
	} `json:"user"`
}

type BeatmapSetsOUT struct {
	Artist         *string `json:"artist"`
	ArtistUnicode  *string `json:"artist_unicode"`
	Creator        *string `json:"creator"`
	FavouriteCount *int    `json:"favourite_count"`
	Hype           struct {
		Current  *int `json:"current"`
		Required *int `json:"required"`
	} `json:"hype"`
	Id           *int    `json:"id"`
	Nsfw         *bool   `json:"nsfw"`
	PlayCount    *int    `json:"play_count"`
	PreviewUrl   *string `json:"preview_url"`
	Source       *string `json:"source"`
	Status       *string `json:"status"`
	Title        *string `json:"title"`
	TitleUnicode *string `json:"title_unicode"`
	UserId       *int    `json:"user_id"`
	Video        *bool   `json:"video"`
	Availability struct {
		DownloadDisabled *bool   `json:"download_disabled"`
		MoreInformation  *string `json:"more_information"`
	} `json:"availability"`
	Bpm                *string    `json:"bpm"`
	CanBeHyped         *bool      `json:"can_be_hyped"`
	DiscussionEnabled  *bool      `json:"discussion_enabled"`
	DiscussionLocked   *bool      `json:"discussion_locked"`
	IsScoreable        *bool      `json:"is_scoreable"`
	LastUpdated        *time.Time `json:"last_updated"`
	LegacyThreadUrl    *string    `json:"legacy_thread_url"`
	NominationsSummary struct {
		Current  *int `json:"current"`
		Required *int `json:"required"`
	} `json:"nominations_summary"`
	Ranked        int          `json:"ranked"`
	RankedDate    *time.Time   `json:"ranked_date"`
	Storyboard    *bool        `json:"storyboard"`
	SubmittedDate *time.Time   `json:"submitted_date"`
	Tags          *string      `json:"tags"`
	HasFavourited *bool        `json:"has_favourited"`
	Beatmaps      []BeatmapOUT `json:"beatmaps"`

	Description struct {
		Description *string `json:"description"`
	} `json:"description"`
	Genre struct {
		Id   *int    `json:"id"`
		Name *string `json:"name"`
	} `json:"genre"`
	Language struct {
		Id   *int    `json:"id"`
		Name *string `json:"name"`
	} `json:"language"`
	RatingsString *string `json:"ratings_string"`

	Cache struct {
		Video   bool `json:"video"`
		NoVideo bool `json:"noVideo"`
	} `json:"cache"`
}
