package entity

import (
	"fmt"
	"time"
)

type RFC3339 time.Time

func (t RFC3339) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02T15:04:05Z"))), nil
}

func (t RFC3339) ToTime() time.Time {
	return time.Time(t)
}
