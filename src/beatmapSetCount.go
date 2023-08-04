package src

import (
	"github.com/Nerinyan/Nerinyan-APIV2/db/mariadb"
	"github.com/pterm/pterm"
	"time"
)

var BeatmapSetCount int64 = 0

func StartBeatmapSetCounter() {
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		for ; ; <-ticker.C {
			err := mariadb.Mariadb.Table("BEATMAPSET").Select("COUNT(*)").Count(&BeatmapSetCount).Error
			if err != nil {
				pterm.Error.WithShowLineNumber().Println(err)
				continue
			}

		}
	}()
}
