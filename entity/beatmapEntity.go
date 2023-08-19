package entity

/**
포인터 데이터의경우 nullable 임
*/

func (Beatmaset) TableName() string {
	return "BEATMAPSET"
}

type Beatmaset struct {
	Id             int    `json:"id,omitempty"                   gorm:"primaryKey;column:ID"`
	Artist         string `json:"artist,omitempty"               gorm:"column:ARTIST"`
	ArtistUnicode  string `json:"artist_unicode,omitempty"       gorm:"column:ARTIST_UNICODE"`
	Covers         Covers `json:"covers,omitempty"               gorm:"foreignKey:ID;references:BEATMAPSET_ID"`
	Creator        string `json:"creator,omitempty"              gorm:"column:CREATOR"`
	FavouriteCount int    `json:"favourite_count,omitempty"      gorm:"column:FAVOURITE_COUNT"`
	Hype           *Hype  `json:"hype,omitempty"                 gorm:"embedded;embeddedPrefix:HYPE_"` // nullable
	Nsfw           bool   `json:"nsfw,omitempty"                 gorm:"column:NSFW"`
	Offset         int    `json:"offset,omitempty"               gorm:"column:SET_OFFSET"`
	PlayCount      int    `json:"play_count,omitempty"           gorm:"column:PLAY_COUNT"`
	PreviewUrl     string `json:"preview_url,omitempty"          gorm:"column:PREVIEW_URL"`
	Source         string `json:"source,omitempty"               gorm:"column:SOURCE"`
	Spotlight      bool   `json:"spotlight,omitempty"            gorm:"column:SPOTLIGHT"`
	Status         string `json:"status,omitempty"               gorm:"column:STATUS"` //'graveyard' | 'wip' | 'pending' | 'ranked' | 'approved' | 'qualified' | 'loved'
	Title          string `json:"title,omitempty"                gorm:"column:TITLE"`
	TitleUnicode   string `json:"title_unicode,omitempty"        gorm:"column:TITLE_UNICODE"`
	TrackId        *int   `json:"track_id,omitempty"             gorm:"column:TRACK_ID"`
	UserId         int    `json:"user_id,omitempty"              gorm:"column:USER_ID"`
	Video          bool   `json:"video,omitempty"                gorm:"column:VIDEO"`
	// ========================================================================================================
	Bpm                *float64            `json:"bpm,omitempty"                  gorm:"column:BPM"`
	CanBeHyped         *bool               `json:"can_be_hyped,omitempty"         gorm:"column:CAN_BE_HYPED"`
	DeletedAt          *RFC3339            `json:"deleted_at,omitempty"           gorm:"column:DELETED_AT"`
	DiscussionEnabled  *bool               `json:"discussion_enabled,omitempty"   gorm:"column:DISCUSSION_ENABLED"`
	DiscussionLocked   *bool               `json:"discussion_locked,omitempty"    gorm:"column:DISCUSSION_LOCKED"`
	IsScoreable        *bool               `json:"is_scoreable,omitempty"         gorm:"column:IS_SCOREABLE"`
	LastUpdated        *RFC3339            `json:"last_updated,omitempty"         gorm:"column:LAST_UPDATED"`
	LegacyThreadUrl    *string             `json:"legacy_thread_url,omitempty"    gorm:"column:LEGACY_THREAD_URL"`
	NominationsSummary *Nominations        `json:"nominations_summary,omitempty"  gorm:"embedded;embeddedPrefix:NOMINATIONS_"`
	Ranked             *int                `json:"ranked,omitempty"               gorm:"column:RANKED"`
	RankedDate         *RFC3339            `json:"ranked_date,omitempty"          gorm:"column:RANKED_DATE"`
	Storyboard         *bool               `json:"storyboard,omitempty"           gorm:"column:STORYBOARD"`
	SubmittedDate      *RFC3339            `json:"submitted_date,omitempty"       gorm:"column:SUBMITTED_DATE"`
	Tags               *string             `json:"tags,omitempty"                 gorm:"column:TAGS"`
	Availability       *Availability       `json:"availability,omitempty"         gorm:"embedded;embeddedPrefix:AVAILABILITY_"`
	HasFavourited      *bool               `json:"has_favourited,omitempty"       gorm:"column:HAS_FAVOURITED"`
	Beatmaps           *[]Beatmap          `json:"beatmaps,omitempty"             gorm:"foreignKey:ID;references:BEATMAPSET_ID"`
	PackTags           *StringsArray       `json:"pack_tags,omitempty"            gorm:"column:PACK_TAGS"`
	Converts           *[]Beatmap          `json:"converts,omitempty"             gorm:"-"` // 이거 지원하려면 join 시 추가 쿼리 지원해야함. 가능한지는 봐야할듯
	CurrentNominations *CurrentNominations `json:"current_nominations,omitempty"  gorm:"embedded;embeddedPrefix:CURRENT_NOMINATIONS_"`
	Description        *Description        `json:"description,omitempty"          gorm:"embedded;embeddedPrefix:DESCRIPTION_"`
	Genre              *Genre              `json:"genre,omitempty"                gorm:"embedded;embeddedPrefix:GENRE_"`
	Language           *Language           `json:"language,omitempty"             gorm:"embedded;embeddedPrefix:LANGUAGE_"`
	Ratings            *IntArray           `json:"ratings,omitempty"              gorm:"column:RATINGS"`
	RecentFavourites   *[]BanchoUser       `json:"recent_favourites,omitempty"    gorm:"-"` // 지원 불가능할듯함
	RelatedUsers       *[]BanchoUser       `json:"related_users,omitempty"        gorm:"-"` // 지원 불가능할듯함
	User               *BanchoUser         `json:"user,omitempty"                 gorm:"foreignKey:USER_ID;references:ID"`
}

type CurrentNominations struct {
	BeatmapsetId int          `json:"beatmapset_id"      gorm:"column:CURRENT_NOMINATIONS_BEATMAPSET_ID"`
	Reset        bool         `json:"reset"              gorm:"column:CURRENT_NOMINATIONS_RESET"`
	Rulesets     StringsArray `json:"rulesets,omitempty" gorm:"column:CURRENT_NOMINATIONS_RULESETS"`
	UserId       int          `json:"user_id"            gorm:"column:CURRENT_NOMINATIONS_USER_ID"`
}

type Description struct {
	Bbcode      *string `json:"bbcode"      gorm:"column:DESCRIPTION_BBCODE"`
	Description *string `json:"description" gorm:"column:DESCRIPTION_DESCRIPTION"`
}

type Genre struct {
	Id   *int   `json:"id"   gorm:"column:GENRE_ID"`
	Name string `json:"name" gorm:"column:GENRE_NAME"`
}

type Language struct {
	Id   *int   `json:"id"   gorm:"column:LANGUAGE_ID"`
	Name string `json:"name" gorm:"column:LANGUAGE_NAME"`
}

func (Covers) TableName() string {
	return "BEATMAPSET_COVERS"
}

type Covers struct {
	BeatmapsetId int    `json:"beatmapset_id,omitempty" gorm:"primaryKey;column:BEATMAPSET_ID"`
	Cover        string `json:"cover,omitempty"         gorm:"column:COVER"`
	Cover2X      string `json:"cover@2x,omitempty"      gorm:"column:COVER2X"`
	Card         string `json:"card,omitempty"          gorm:"column:CARD"`
	Card2X       string `json:"card@2x,omitempty"       gorm:"column:CARD2X"`
	List         string `json:"list,omitempty"          gorm:"column:LIST"`
	List2X       string `json:"list@2x,omitempty"       gorm:"column:LIST2X"`
	SlimCover    string `json:"slimcover,omitempty"     gorm:"column:SLIM_COVER"`
	SlimCover2X  string `json:"slimcover@2x,omitempty"  gorm:"column:SLIM_COVER2X"`
}

type Hype struct {
	Current  int `json:"current,omitempty"  gorm:"column:HYPE_CURRENT"`
	Required int `json:"required,omitempty" gorm:"column:HYPE_REQUIRED"`
}

type Nominations struct {
	Current  int `json:"current,omitempty"  gorm:"column:NOMINATIONS_CURRENT"`
	Required int `json:"required,omitempty" gorm:"column:NOMINATIONS_REQUIRED"`
}

type Availability struct {
	DownloadDisabled bool    `json:"download_disabled,omitempty" gorm:"column:AVAILABILITY_DOWNLOAD_DISABLED"`
	MoreInformation  *string `json:"more_information,omitempty"  gorm:"column:AVAILABILITY_MORE_INFORMATION"`
}

func (Beatmap) TableName() string {
	return "BEATMAP"
}

type Beatmap struct {
	Id               int     `json:"id,omitempty"                gorm:"primaryKey;column:ID"`
	BeatmapsetId     int     `json:"beatmapset_id,omitempty"     gorm:"column:BEATMAPSET_ID"`
	DifficultyRating float64 `json:"difficulty_rating,omitempty" gorm:"column:DIFFICULTY_RATING"`
	Mode             string  `json:"mode,omitempty"              gorm:"column:MODE"` //['osu', 'taiko', 'fruits', 'mania']
	Status           string  `json:"status,omitempty"            gorm:"column:STATUS"`
	TotalLength      int     `json:"total_length,omitempty"      gorm:"column:TOTAL_LENGTH"`
	UserId           int     `json:"user_id,omitempty"           gorm:"column:USER_ID"`
	Version          string  `json:"version,omitempty"           gorm:"column:VERSION"`
	//====================================================================================================
	Accuracy      *float64   `json:"accuracy,omitempty"          gorm:"column:ACCURACY"`
	Ar            *float64   `json:"ar,omitempty"                gorm:"column:AR"`
	Bpm           *float64   `json:"bpm,omitempty"               gorm:"column:BPM"`
	Convert       *bool      `json:"convert,omitempty"           gorm:"column:CONVERTED"`
	CountCircles  *int       `json:"count_circles,omitempty"     gorm:"column:COUNT_CIRCLES"`
	CountSliders  *int       `json:"count_sliders,omitempty"     gorm:"column:COUNT_SLIDERS"`
	CountSpinners *int       `json:"count_spinners,omitempty"    gorm:"column:COUNT_SPINNERS"`
	Cs            *float64   `json:"cs,omitempty"                gorm:"column:CS"`
	DeletedAt     *RFC3339   `json:"deleted_at,omitempty"        gorm:"column:DELETED_AT"`
	Drain         *float64   `json:"drain,omitempty"             gorm:"column:DRAIN"`
	HitLength     *int       `json:"hit_length,omitempty"        gorm:"column:HIT_LENGTH"`
	IsScoreable   *bool      `json:"is_scoreable,omitempty"      gorm:"column:IS_SCOREABLE"`
	LastUpdated   *RFC3339   `json:"last_updated,omitempty"      gorm:"column:LAST_UPDATED"`
	ModeInt       *int       `json:"mode_int,omitempty"          gorm:"column:MODE_INT"`
	Passcount     *int       `json:"passcount,omitempty"         gorm:"column:PASSCOUNT"`
	Playcount     *int       `json:"playcount,omitempty"         gorm:"column:PLAYCOUNT"`
	Ranked        *int       `json:"ranked,omitempty"            gorm:"column:RANKED"`
	Url           *string    `json:"url,omitempty"               gorm:"column:URL"`
	Checksum      *string    `json:"checksum,omitempty"          gorm:"column:CHECKSUM"`
	MaxCombo      *int       `json:"max_combo,omitempty"         gorm:"column:MAX_COMBO"`
	Failtimes     *Failtimes `json:"failtimes,omitempty"         gorm:"embedded;embeddedPrefix:FAILTIMES_"`
}

func (Failtimes) TableName() string {
	return "FAILTIMES"
}

type Failtimes struct {
	Fail IntArray `json:"fail" gorm:"column:FAILTIMES_FAIL"`
	Exit IntArray `json:"exit" gorm:"column:FAILTIMES_EXIT"`
}

func (BanchoUser) TableName() string {
	return "BANCHO_USER"
}

type BanchoUser struct {
	Id            int      `json:"id" gorm:"primaryKey;column:ID"`
	AvatarUrl     string   `json:"avatar_url" gorm:"column:AVATAR_URL"`
	CountryCode   string   `json:"country_code" gorm:"column:COUNTRY_CODE"`
	DefaultGroup  *string  `json:"default_group" gorm:"column:DEFAULT_GROUP"`
	IsActive      bool     `json:"is_active" gorm:"column:IS_ACTIVE"`
	IsBot         bool     `json:"is_bot" gorm:"column:IS_BOT"`
	IsDeleted     bool     `json:"is_deleted" gorm:"column:IS_DELETED"`
	IsOnline      bool     `json:"is_online" gorm:"column:IS_ONLINE"`
	IsSupporter   bool     `json:"is_supporter" gorm:"column:IS_SUPPORTER"`
	LastVisit     *RFC3339 `json:"last_visit" gorm:"column:LAST_VISIT"`
	PmFriendsOnly bool     `json:"pm_friends_only" gorm:"column:PM_FRIENDS_ONLY"`
	ProfileColour *string  `json:"profile_colour" gorm:"column:PROFILE_COLOUR"`
	Username      string   `json:"username" gorm:"column:USERNAME"`
}
