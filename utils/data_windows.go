package utils

import "fmt"

const (
	KiB = 1024
	MiB = KiB * KiB
	GiB = KiB * KiB * KiB
	TiB = KiB * KiB * KiB * KiB
	PiB = KiB * KiB * KiB * KiB * KiB
)

func ToHumanDataSize(bytes uint64) string {
	var unit string
	var size float64

	if bytes >= PiB {
		unit = "PiB"
		size = float64(bytes) / PiB
	} else if bytes >= TiB {
		unit = "TiB"
		size = float64(bytes) / TiB
	} else if bytes >= GiB {
		unit = "GiB"
		size = float64(bytes) / GiB
	} else if bytes >= MiB {
		unit = "MiB"
		size = float64(bytes) / MiB
	} else if bytes >= KiB {
		unit = "KiB"
		size = float64(bytes) / KiB
	} else {
		unit = "B"
		size = float64(bytes)
	}

	return fmt.Sprintf("%.2f%s", size, unit)
}
