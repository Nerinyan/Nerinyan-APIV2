package utils

import "fmt"

const (
	KB = 1000
	MB = KB * KB
	GB = KB * KB * KB
	TB = KB * KB * KB * KB
	PB = KB * KB * KB * KB * KB
)

func ToHumanDataSize(bytes int64) string {
	var unit string
	var size float64

	if bytes >= PB {
		unit = "PB"
		size = float64(bytes) / PB
	} else if bytes >= TB {
		unit = "TB"
		size = float64(bytes) / TB
	} else if bytes >= GB {
		unit = "GB"
		size = float64(bytes) / GB
	} else if bytes >= MB {
		unit = "MB"
		size = float64(bytes) / MB
	} else if bytes >= KB {
		unit = "KB"
		size = float64(bytes) / KB
	} else {
		unit = "B"
		size = float64(bytes)
	}

	return fmt.Sprintf("%.2f %s", size, unit)

}
