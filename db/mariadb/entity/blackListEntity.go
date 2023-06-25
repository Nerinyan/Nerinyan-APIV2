package entity

import (
	"time"
)

type BlacklistEntity struct {
	IPV4      string    `gorm:"column:IPV4"`
	ExpiredAt time.Time `gorm:"column:EXPIRED_AT"`
}

func (BlacklistEntity) TableName() string {
	return "BLACKLIST"
}
