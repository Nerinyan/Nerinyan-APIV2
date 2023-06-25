package entity

type ServerCache struct {
	Key   string `gorm:"column:KEY;primaryKey"`
	Value string `gorm:"column:VALUE"`
}

func (ServerCache) TableName() string {
	return "SERVER_CACHE"
}
