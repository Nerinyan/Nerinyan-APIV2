package route

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/Nerinyan/Nerinyan-APIV2/Logger"
	"github.com/Nerinyan/Nerinyan-APIV2/banchoCroller"
	"github.com/Nerinyan/Nerinyan-APIV2/bodyStruct"
	"github.com/Nerinyan/Nerinyan-APIV2/config"
	"github.com/Nerinyan/Nerinyan-APIV2/db"
	"github.com/Nerinyan/Nerinyan-APIV2/src"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type downloadBeatmapSet_requestBody struct {
	NoVideo     	bool `query:"noVideo"`
	NoVideo2    	bool `query:"nv"`
	NoBg  	    	bool `query:"noBg"`
	NoBg2  	    	bool `query:"nb"`
	NoHitsound  	bool `query:"noHitsound"`
	NoHitsound2 	bool `query:"nh"`
	NoStoryboard	bool `query:"noStoryboard"`
	NoStoryboard2	bool `query:"nsb"`
	MapId       	int  `param:"mapId"`
	SetId       	int  `param:"setId"`
}

func DownloadBeatmapSet(c echo.Context) (err error) {
	var request downloadBeatmapSet_requestBody
	err = c.Bind(&request)
	if err != nil {

		return c.JSON(http.StatusInternalServerError, logger.Error(c, &bodyStruct.ErrorStruct{
			Code:        "DownloadBeatmapSet-001",
			Path:        c.Path(),
			RequestId:   c.Response().Header().Get("X-Request-ID"),
			Error:       err,
			Message:     "request parse error",
			RequestData: request,
		}))
	}
	request.NoVideo = request.NoVideo || request.NoVideo2

	request.NoBg = request.NoBg || request.NoBg2
	request.NoHitsound = request.NoHitsound || request.NoHitsound2
	request.NoStoryboard = request.NoStoryboard || request.NoStoryboard2

	redirecturl := fmt.Sprintf("https://subapi.nerinyan.moe/d/%d", request.SetId)
	// NoBg 또는 NoHitsound 요청시 서브API로 리다이렉트
	// NoBG NoHitsound NoStoryboard
	if request.NoBg == true && request.NoHitsound == true && request.NoStoryboard == true {
		redirecturl = redirecturl + "?nb=1&nh=1&nsb=1"
		if request.NoVideo == true {
			redirecturl = redirecturl + "&nv=1"
		}
	// NoBG NoHitsound
	} else if request.NoBg == true && request.NoHitsound == true && request.NoStoryboard == false {
		redirecturl = redirecturl + "?nb=1&nh=1"
		if request.NoVideo == true {
			redirecturl = redirecturl + "&nv=1"
		}
	// NoBG NoStoryboard
	} else if request.NoBg == true && request.NoHitsound == false && request.NoStoryboard == true {
		redirecturl = redirecturl + "?nb=1&nsb=1"
		if request.NoVideo == true {
			redirecturl = redirecturl + "&nv=1"
		}
	// NoHitsound NoStoryboard
	} else if request.NoBg == false && request.NoHitsound == true && request.NoStoryboard == true {
		redirecturl = redirecturl + "?nh=1&nsb=1"
		if request.NoVideo == true {
			redirecturl = redirecturl + "&nv=1"
		}
	// NoBG 
	} else if request.NoBg == true && request.NoHitsound == false && request.NoStoryboard == false {
		redirecturl = redirecturl + "?nb=1"
		if request.NoVideo == true {
			redirecturl = redirecturl + "&nv=1"
		}
	// NoHitsound
	} else if request.NoBg == false && request.NoHitsound == true && request.NoStoryboard == false {
		redirecturl = redirecturl + "?nh=1"
		if request.NoVideo == true {
			redirecturl = redirecturl + "&nv=1"
		}
	// NoStoryboard
	} else if request.NoBg == false && request.NoHitsound == false && request.NoStoryboard == true {
		redirecturl = redirecturl + "?nsb=1"
		if request.NoVideo == true {
			redirecturl = redirecturl + "&nv=1"
		}
	}

	if request.NoBg == true || request.NoHitsound == true || request.NoStoryboard == true {
		return c.Redirect(http.StatusPermanentRedirect, redirecturl)
	}

	var row *sql.Row
	if request.SetId != 0 {
		go banchoCroller.ManualUpdateBeatmapSet(request.SetId)
		row = db.Maria.QueryRow(`SELECT BEATMAPSET_ID,ARTIST,TITLE,LAST_UPDATED,VIDEO FROM BEATMAPSET WHERE BEATMAPSET_ID = ?`, request.SetId)
	} else if request.MapId != 0 {
		row = db.Maria.QueryRow(`SELECT BEATMAPSET_ID,ARTIST,TITLE,LAST_UPDATED,VIDEO FROM BEATMAPSET WHERE BEATMAPSET_ID = (SELECT BEATMAPSET_ID FROM BEATMAP WHERE BEATMAP_ID = ?);`, request.MapId)
	} else {
		return c.JSON(http.StatusInternalServerError, logger.Error(c, &bodyStruct.ErrorStruct{
			Code:        "DownloadBeatmapSet-001-1",
			Path:        c.Path(),
			RequestId:   c.Response().Header().Get("X-Request-ID"),
			Error:       nil,
			Message:     "set id & map id not found",
			RequestData: request,
		}))
	}

	if err = row.Err(); err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, logger.Error(c, &bodyStruct.ErrorStruct{
				Code:        "DownloadBeatmapSet-002",
				Path:        c.Path(),
				RequestId:   c.Response().Header().Get("X-Request-ID"),
				Error:       err,
				Message:     "not in database",
				RequestData: request,
			}))
		}
		return c.JSON(http.StatusInternalServerError, logger.Error(c, &bodyStruct.ErrorStruct{
			Code:        "DownloadBeatmapSet-003",
			Path:        c.Path(),
			RequestId:   c.Response().Header().Get("X-Request-ID"),
			Error:       err,
			Message:     "database Query error",
			RequestData: request,
		}))
	}

	var a struct {
		Id          string
		Artist      string
		Title       string
		LastUpdated string
		Video       bool
	}

	if err = row.Scan(&a.Id, &a.Artist, &a.Title, &a.LastUpdated, &a.Video); err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, logger.Error(c, &bodyStruct.ErrorStruct{
				Code:        "DownloadBeatmapSet-004",
				Path:        c.Path(),
				RequestId:   c.Response().Header().Get("X-Request-ID"),
				Error:       err,
				Message:     "not in database",
				RequestData: request,
			}))

		}
		return c.JSON(http.StatusInternalServerError, logger.Error(c, &bodyStruct.ErrorStruct{
			Code:        "DownloadBeatmapSet-005",
			Path:        c.Path(),
			RequestId:   c.Response().Header().Get("X-Request-ID"),
			Error:       err,
			Message:     "database Query error",
			RequestData: request,
		}))
	}

	lu, err := time.Parse("2006-01-02 15:04:05", a.LastUpdated)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, logger.Error(c, &bodyStruct.ErrorStruct{
			Code:        "DownloadBeatmapSet-006",
			Path:        c.Path(),
			RequestId:   c.Response().Header().Get("X-Request-ID"),
			Error:       err,
			Message:     "time Parse error",
			RequestData: request,
		}))
	}

	url := fmt.Sprintf("https://osu.ppy.sh/api/v2/beatmapsets/%d/download", request.SetId)
	if a.Video && request.NoVideo {
		request.SetId *= -1
		a.Title += " [no video]"
		url += "?noVideo=1"
	}

	serverFileName := fmt.Sprintf("%s/%d.osz", config.Config.TargetDir, request.SetId)

	if src.FileList[request.SetId].Unix() >= lu.Unix() { // 맵이 최신인경우
		c.Response().Header().Set("Content-Type", "application/x-osu-beatmap-archive")
		return c.Attachment(serverFileName, fmt.Sprintf("%s %s - %s.osz", a.Id, a.Artist, a.Title))
	}

	//==========================================
	//=        비트맵 파일이 서버에 없는경우        =
	//==========================================

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, logger.Error(c, &bodyStruct.ErrorStruct{
			Code:      "DownloadBeatmapSet-007",
			Path:      c.Path(),
			RequestId: c.Response().Header().Get("X-Request-ID"),
			Error:     err,
			Message:   "Bancho request Build Error",
		}))
	}
	req.Header.Add("Authorization", config.Config.Osu.Token.TokenType+" "+config.Config.Osu.Token.AccessToken)

	res, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, logger.Error(c, &bodyStruct.ErrorStruct{
			Code:      "DownloadBeatmapSet-008",
			Path:      c.Path(),
			RequestId: c.Response().Header().Get("X-Request-ID"),
			Error:     err,
			Message:   "Bancho request Build Error",
		}))
	}
	defer res.Body.Close()

	if res.Header.Get("Content-Type") != "application/x-osu-beatmap-archive" {
		return c.JSON(http.StatusInternalServerError, logger.Error(c, &bodyStruct.ErrorStruct{
			Code:      "DownloadBeatmapSet-012",
			Path:      c.Path(),
			RequestId: c.Response().Header().Get("X-Request-ID"),
			Error:     errors.New(http.StatusText(res.StatusCode)),
			Message:   "Bancho request Error. :" + res.Status,
		}))
	}

	if res.StatusCode != http.StatusOK {
		return c.JSON(http.StatusInternalServerError, logger.Error(c, &bodyStruct.ErrorStruct{
			Code:      "DownloadBeatmapSet-009",
			Path:      c.Path(),
			RequestId: c.Response().Header().Get("X-Request-ID"),
			Error:     errors.New(http.StatusText(res.StatusCode)),
			Message:   "Bancho request Error. :" + res.Status,
		}))
	}

	pterm.Info.Println("beatmapSet Downloading at", serverFileName)

	cLen, _ := strconv.Atoi(res.Header.Get("Content-Length"))
	c.Response().Header().Set("Content-Length", res.Header.Get("Content-Length"))
	c.Response().Header().Set("Content-Disposition", res.Header.Get("Content-Disposition"))
	c.Response().Header().Set("Content-Type", "application/x-osu-beatmap-archive")

	var buf bytes.Buffer

	//http://localhost/d/1573058
	//http://localhost/d/1469677

	defer c.Response().Flush()
	clientError := false
	for i := 0; i < cLen; { // 읽을 데이터 사이즈 체크
		var b = make([]byte, 64000) // 바이트 배열
		n, e := res.Body.Read(b)    // 반쵸 스트림에서 64k 읽어서 바이트 배열 b 에 넣음

		i += n // 현재까지 읽은 바이트
		if n > 0 {
			buf.Write(b[:n]) // 서버에 저장할 파일 버퍼에 쓴다
			if !clientError {
				if _, ee := c.Response().Write(b[:n]); ee != nil {
					clientError = true // write response 에러 발생시
					err = c.JSON(http.StatusInternalServerError, logger.Error(c, &bodyStruct.ErrorStruct{
						Code:      "DownloadBeatmapSet-010",
						Path:      c.Path(),
						RequestId: c.Response().Header().Get("X-Request-ID"),
						Error:     ee,
						Message:   "write response stream error",
					}))
				}
			}

		}

		if e == io.EOF {
			break
		} else if e != nil { //에러처리
			return c.JSON(http.StatusInternalServerError, logger.Error(c, &bodyStruct.ErrorStruct{
				Code:      "DownloadBeatmapSet-011",
				Path:      c.Path(),
				RequestId: c.Response().Header().Get("X-Request-ID"),
				Error:     e,
				Message:   "fail to read Bancho Stream",
			}))
		}
	}
	if cLen == buf.Len() {
		return saveLocal(&buf, serverFileName, request.SetId)
	}
	errMsg := fmt.Sprintf("filesize not match: bancho response bytes : %d | downloaded bytes : %d", cLen, buf.Len())
	pterm.Error.Printfln(errMsg)
	return errors.New(errMsg)

}

func saveLocal(data *bytes.Buffer, path string, id int) (err error) {
	tmp := path + ".down"
	file, err := os.Create(tmp)
	if err != nil {
		return
	}
	if file == nil {
		return errors.New("")
	}
	_, err = file.Write(data.Bytes())
	if err != nil {
		return
	}
	file.Close()

	if _, err = os.Stat(path); !os.IsNotExist(err) {
		err = os.Remove(path)
		if err != nil {
			return
		}
	}
	err = os.Rename(tmp, path)
	if err != nil {
		return
	}

	src.FileList[id] = time.Now()
	pterm.Info.Println("beatmapSet Downloading Finished", path)
	return
}

//
