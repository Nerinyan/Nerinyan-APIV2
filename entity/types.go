package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"time"
)

//==========================================================

// IntArray [0,0,0]
type IntArray []int

// Scan implements the sql.Scanner interface
func (v *IntArray) Scan(src interface{}) error {
	if bytes, ok := src.([]byte); ok {
		return json.Unmarshal(bytes, v)
	}
	return gorm.ErrInvalidData
}

// Value implements the driver.Valuer interface
func (v *IntArray) Value() (driver.Value, error) {
	return json.Marshal(v)
}

//==========================================================

// StringsArray ["",""]
type StringsArray []string

// Scan implements the sql.Scanner interface
func (v *StringsArray) Scan(src interface{}) error {
	if bytes, ok := src.([]byte); ok {
		return json.Unmarshal(bytes, v)
	}
	return gorm.ErrInvalidData
}

// Value implements the driver.Valuer interface
func (v *StringsArray) Value() (driver.Value, error) {
	return json.Marshal(v)
}

//==========================================================

// JsonObject {}
type JsonObject map[string]any

// Scan implements the sql.Scanner interface
func (v *JsonObject) Scan(src interface{}) error {
	if bytes, ok := src.([]byte); ok {
		return json.Unmarshal(bytes, v)
	}
	return gorm.ErrInvalidData
}

// Value implements the driver.Valuer interface
func (v *JsonObject) Value() (driver.Value, error) {
	return json.Marshal(v)
}

//==========================================================

// JsonArray []
type JsonArray []any

// Scan implements the sql.Scanner interface
func (v *JsonArray) Scan(src interface{}) error {
	if bytes, ok := src.([]byte); ok {
		return json.Unmarshal(bytes, v)
	}
	return gorm.ErrInvalidData
}

// Value implements the driver.Valuer interface
func (v *JsonArray) Value() (driver.Value, error) {
	return json.Marshal(v)
}

//==========================================================

// RFC3339 YYYY-MM-DDThh-mm-ssZ
type RFC3339 time.Time

func (t RFC3339) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02T15:04:05Z"))), nil
}

func (t RFC3339) ToTime() time.Time {
	return time.Time(t)
}
