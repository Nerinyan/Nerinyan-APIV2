package banchoCrawler

import (
	"database/sql"
	"fmt"
	"github.com/Nerinyan/Nerinyan-APIV2/db"
	"github.com/pterm/pterm"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

var (
	regPackDownloadLink, _ = regexp.Compile(`(?:<a href=")(.+?)(?:"(\s+?|.+?)class="beatmap-pack-download__link">)`)
	regMapSetId, _         = regexp.Compile(`(?:<a href=")(?:https://osu[.]ppy[.]sh/beatmapsets/)([0-9]+?)(?:"(\s+?|.+?)class="beatmap-pack-items__link">)`)
	regPackName, _         = regexp.Compile(`(?:<div class="beatmap-pack__name">)(.*)(?:</div>)`)
	regPackId, _           = regexp.Compile(`(?:https://osu[.]ppy[.]sh/beatmaps/packs/)([0-9]+?)(?:")`)
	regPackDate, _         = regexp.Compile(`(?:<span class="beatmap-pack__date">)(.*)(?:</span>)`)
	regCreator, _          = regexp.Compile(`(?:<span class="beatmap-pack__author beatmap-pack__author--bold">)(.*)(?:</span>)`)
	regLastPage, _         = regexp.Compile(`(?:<a class="pagination-v2__link" href=".+?">)([0-9]+?)(?:</a>)`)
)

var packType = []string{
	"standard",
	"chart",
	"theme",
	"artist",
}

func UpdateAllPackList() {
	pterm.Info.Println("beatmap pack cron started.")
	defer pterm.Info.Println("beatmap pack cron end.")
	for i := 0; i < len(packType); i++ {
		var dbPackId string
		page := 0

		row := db.Maria.QueryRow(`SELECT PACK_ID FROM BEATMAP_PACK WHERE TYPE = ? ORDER BY PACK_ID DESC LIMIT 1`, packType[i])
		if err := row.Err(); err != nil {
			if err != sql.ErrNoRows {
				log.Println(err)
				pterm.Error.Println(err)
				return
			}
		} else {
			if err = row.Scan(&dbPackId); err != nil {
				if err != sql.ErrNoRows {
					log.Println(err)
					pterm.Error.Println(err)
					return
				}
			}
		}

		bodyString := fetchPacks(packType[i], strconv.Itoa(page))
		m := regLastPage.FindAllStringSubmatch(bodyString, -1)
		for _, sm := range m {
			if len(sm) > 1 {
				t1, _ := strconv.Atoi(sm[1])
				if page < t1 {
					page = t1
				}
			}
		}
		packId, packName, packCreator, packDate := parseBody(bodyString)
		upsertPack(packType[i], packId, packName, packCreator, packDate)
		for j := 1; j <= page; j++ {
			time.Sleep(time.Millisecond * 500)
			bodyString = fetchPacks(packType[i], strconv.Itoa(j))
			packId, packName, packCreator, packDate = parseBody(bodyString)
			upsertPack(packType[i], packId, packName, packCreator, packDate)

		}
	}
}

const upsertPackSql = `
INSERT INTOBEATMAP_PACK(PACK_ID,TYPE,NAME,CREATOR,DATE)
VALUES %s
ON DUPLICATE KEY UPDATE 
TYPE = VALUES(TYPE),NAME = VALUES(NAME),CREATOR = VALUES(CREATOR),DATE = VALUES(DATE);`
const upsertPackSqlRows = `(?,?,?,?,?)`

func upsertPack(packType string, packId, packName, packCreator, packDate []string) {
	sizePackName := len(packName)
	sizePackId := len(packId)
	sizePackDate := len(packDate)
	sizePackCreator := len(packCreator)
	if sizePackName&sizePackId&sizePackDate&sizePackCreator != sizePackCreator {
		pterm.Error.Println("pack data parse Errpr")
		pterm.Error.Println(sizePackName, sizePackId, sizePackDate, sizePackCreator)
		return
	}
	var buf []interface{}
	for i := 0; i < sizePackId; i++ {
		buf = append(buf, packId[i], packType, packName[i], packCreator[i], packDate[i])
	}

	for j := 0; j < sizePackId; j++ {
		_, err := db.Maria.Exec(fmt.Sprintf(upsertPackSql, buildSqlValues(upsertPackSqlRows, sizePackId)), buf...)
		if err != nil {
			log.Println(err)
			pterm.Error.Println(err)
			return
		}
	}
}

func parseBody(bodyString string) (packId, packName, packCreator, packDate []string) {
	m := regPackName.FindAllStringSubmatch(bodyString, -1)
	for _, sm := range m {
		if len(sm) > 1 {
			packName = append(packName, sm[1])
		} else {
			packName = append(packName, "")
		}
	}
	m = regPackId.FindAllStringSubmatch(bodyString, -1)
	for _, sm := range m {
		if len(sm) > 1 {
			packId = append(packId, sm[1])
		} else {
			packId = append(packId, "-1")
		}
	}
	m = regPackDate.FindAllStringSubmatch(bodyString, -1)
	for _, sm := range m {
		if len(sm) > 1 {
			packDate = append(packDate, sm[1])
		} else {
			packDate = append(packDate, "")
		}
	}
	m = regCreator.FindAllStringSubmatch(bodyString, -1)
	for _, sm := range m {
		if len(sm) > 1 {
			packCreator = append(packCreator, sm[1])
		} else {
			packCreator = append(packCreator, "")
		}
	}
	return
}
func fetchPacks(Type, page string) string {
	res, err := http.Get(fmt.Sprintf("https://osu.ppy.sh/beatmaps/packs?type=%s&page=%s", Type, page))
	if err != nil || res.StatusCode != http.StatusOK {
		if res != nil {
			panic(res.Status)
		}
		panic(err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	return string(body)
}
