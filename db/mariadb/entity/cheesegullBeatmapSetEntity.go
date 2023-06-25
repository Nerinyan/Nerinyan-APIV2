package entity

import (
	"gorm.io/gorm"
	"strconv"
	"time"
)

type CheesegullBeatmapSetEntity struct {
	SetId            int                       `json:"SetID" gorm:"column:BEATMAPSET_ID"`
	ChildrenBeatmaps []CheesegoolBeatmapEntity `json:"ChildrenBeatmaps" gorm:"foreignKey:BEATMAPSET_ID;references:BEATMAPSET_ID"`
	RankedStatus     *int                      `json:"RankedStatus" gorm:"column:RANKED"`
	ApprovedDate     *string                   `json:"ApprovedDate" gorm:"column:"`
	LastUpdate       *string                   `json:"LastUpdate" gorm:"column:LAST_UPDATED"`             // 비트맵 업데이트된 시간
	LastChecked      *time.Time                `json:"LastChecked" gorm:"column:SYSTEM_UPDATE_TIMESTAMP"` // 크롤링한 시간
	Artist           *string                   `json:"Artist" gorm:"column:ARTIST"`
	Title            *string                   `json:"Title" gorm:"column:TITLE"`
	Creator          *string                   `json:"Creator" gorm:"column:CREATOR"`
	CreatorId        *string                   `json:"CreatorID" gorm:"column:USER_ID"`
	Source           *string                   `json:"Source" gorm:"column:SOURCE"`
	Tags             *string                   `json:"Tags" gorm:"column:TAGS"`
	HasVideo         *bool                     `json:"HasVideo" gorm:"column:VIDEO"`
	GenreId          *string                   `json:"-" gorm:"column:GENRE_ID"`
	Genre            *int                      `json:"Genre" gorm:"-"`
	LanguageId       *string                   `json:"-" gorm:"column:LANGUAGE_ID"`
	Language         *int                      `json:"Language" gorm:"-"`
	Favourites       *int                      `json:"Favourites" gorm:"column:FAVOURITE_COUNT"`
}

func (CheesegullBeatmapSetEntity) TableName() string {
	return "BEATMAPSET"
}

func (v *CheesegullBeatmapSetEntity) AfterFind(tx *gorm.DB) (err error) {
	if v.GenreId != nil {
		*v.Genre, _ = strconv.Atoi(*v.GenreId)
	}

	if v.LanguageId != nil {
		*v.Language, _ = strconv.Atoi(*v.LanguageId)
	}

	return
}
