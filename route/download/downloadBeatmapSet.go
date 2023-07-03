package download

import (
	"bytes"
	"fmt"
	"github.com/Nerinyan/Nerinyan-APIV2/banchoCrawler"
	"github.com/Nerinyan/Nerinyan-APIV2/bodyStruct"
	"github.com/Nerinyan/Nerinyan-APIV2/config"
	"github.com/Nerinyan/Nerinyan-APIV2/db/mariadb"
	"github.com/Nerinyan/Nerinyan-APIV2/db/mariadb/entity"
	"github.com/Nerinyan/Nerinyan-APIV2/logger"
	"github.com/Nerinyan/Nerinyan-APIV2/src"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var cannotUseFilename, _ = regexp.Compile(`([\\/:*?"<>|])`)

type downloadBeatmapSet_requestBody struct {
	NoVideo       bool `query:"noVideo"`
	NoVideo2      bool `query:"nv"`
	NoBg          bool `query:"noBg"`
	NoBg2         bool `query:"nb"`
	NoHitsound    bool `query:"noHitsound"`
	NoHitsound2   bool `query:"nh"`
	NoStoryboard  bool `query:"noStoryboard"`
	NoStoryboard2 bool `query:"nsb"`
	MapId         int  `param:"mapId"`
	SetId         int  `param:"setId"`
}

var downloadCount int
var mutex = sync.Mutex{}

func isLimitedDownload() bool {
	mutex.Lock()
	defer mutex.Unlock()
	return downloadCount > 80
}
func init() {
	ticker := time.NewTicker(time.Minute * 10)
	go func() {
		for range ticker.C {
			mutex.Lock()
			downloadCount = 0
			mutex.Unlock()
			pterm.Info.Println("reset download count")
		}
	}()

}
func DownloadBeatmapSet(c echo.Context) (err error) {
	var request downloadBeatmapSet_requestBody
	err = c.Bind(&request)
	if err != nil {
		pterm.Error.WithShowLineNumber().Println(err)
		return
	}
	request.NoVideo = request.NoVideo || request.NoVideo2

	request.NoBg = request.NoBg || request.NoBg2
	request.NoHitsound = request.NoHitsound || request.NoHitsound2
	request.NoStoryboard = request.NoStoryboard || request.NoStoryboard2

	redirecturl := fmt.Sprintf("https://subapi.nerinyan.moe/d/%d", request.SetId)
	// NoBg 또는 NoHitsound 요청시 서브API로 리다이렉트
	// NoBG NoHitsound NoStoryboard
	var subParams = []string{}
	if request.NoBg == true {
		subParams = append(subParams, "nb=1")
	}
	if request.NoHitsound == true {
		subParams = append(subParams, "nh=1")
	}
	if request.NoStoryboard == true {
		subParams = append(subParams, "nsb=1")
	}
	if request.NoVideo == true {
		subParams = append(subParams, "nv=1")
	}
	if len(subParams) > 0 {
		redirecturl += "?" + strings.Join(subParams, "&")
	}

	if request.NoBg == true || request.NoHitsound == true || request.NoStoryboard == true {
		return c.Redirect(http.StatusPermanentRedirect, redirecturl)
	}
	var beatmap4d entity.BanchoBeatmaSetForDownloaderEntity
	if request.SetId != 0 {
		err = mariadb.Mariadb.Find(&beatmap4d, request.SetId).Error
	} else if request.MapId != 0 {
		err = mariadb.Mariadb.Find(&beatmap4d, "BEATMAPSET_ID = (SELECT BEATMAPSET_ID FROM BEATMAP WHERE BEATMAP_ID = ?)", request.MapId).Error
	} else {
		return errors.New("set id & map id not found")
	}
	if err != nil {
		return errors.New(err.Error())
	}
	pterm.Info.Println(beatmap4d)
	go banchoCrawler.ManualUpdateBeatmapSet(beatmap4d.BeatmapsetId)

	if beatmap4d.AvailabilityDownloadDisabled {
		return errors.New("download is disabled from bancho.")
	}

	url := fmt.Sprintf("https://osu.ppy.sh/api/v2/beatmapsets/%d/download", request.SetId)
	if beatmap4d.Video && request.NoVideo {
		beatmap4d.BeatmapsetId *= -1
		beatmap4d.Title += " [no video]"
		url += "?noVideo=1"
	}

	serverFileName := fmt.Sprintf("%s/%d.osz", config.Config.TargetDir, beatmap4d.BeatmapsetId)
	realFilename := cannotUseFilename.ReplaceAllString(fmt.Sprintf("%d %s - %s.osz", beatmap4d.BeatmapsetId, beatmap4d.Artist, beatmap4d.Title), "_")
	c.Response().Header().Set("FileName", realFilename)
	if src.FileList[beatmap4d.BeatmapsetId].Unix() >= beatmap4d.LastUpdated.ToTime().Unix() { // 맵이 최신인경우
		c.Response().Header().Set("Content-Type", "application/x-osu-beatmap-archive")
		c.Response().Header().Set("Content-Source", "nerinyan.moe")
		return c.Attachment(serverFileName, realFilename)
	}

	//==========================================
	//=        비트맵 파일이 서버에 없는경우        =
	//==========================================
	var req *http.Request
	var res *http.Response
	var client *http.Client
	if !isLimitedDownload() {
		client = &http.Client{}
		req, err = http.NewRequest("GET", url, nil)

		if err != nil {
			return errors.WithStack(err)
		}
		req.Header.Add("Authorization", config.Config.Osu.Token.TokenType+" "+config.Config.Osu.Token.AccessToken)

		res, err = client.Do(req)
		if err != nil {
			return errors.New("Bancho request Build Error")
		}
		if res.StatusCode != http.StatusOK {
			pterm.Error.Println("Bancho request Error. :" + res.Status)
			res.Body.Close()
		}
		mutex.Lock()
		downloadCount++
		mutex.Unlock()
	}

	if res == nil || res.StatusCode != http.StatusOK {
		client = &http.Client{}
		req, err = http.NewRequest("GET", fmt.Sprintf("https://beatconnect.io/b/%d", beatmap4d.BeatmapsetId), nil)

		if err != nil {
			return errors.New("beatconnect request Build Error")
		}

		res, err = client.Do(req)
		if err != nil {
			return errors.New("beatconnect request Build Error")
		}
		if res.StatusCode != http.StatusOK {
			res.Body.Close()
			return errors.New("beatconnect request Error. :" + res.Status)
		}
		defer res.Body.Close()
	}

	pterm.Info.Println("beatmapSet Download Source", req.Host)
	pterm.Info.Println("beatmapSet Downloading at", serverFileName)

	cLen, _ := strconv.Atoi(res.Header.Get("Content-Length"))
	c.Response().Header().Set("Content-Source", req.Host)
	c.Response().Header().Set("Content-Length", res.Header.Get("Content-Length"))
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, realFilename))
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
					err = c.JSON(
						http.StatusInternalServerError, logger.Error(
							c, &bodyStruct.ErrorStruct{
								Code:      "DownloadBeatmapSet-010",
								Path:      c.Path(),
								RequestId: c.Response().Header().Get("X-Request-ID"),
								Error:     ee,
								Message:   "write response stream error",
							},
						),
					)
				}
			}

		}

		if e == io.EOF {
			break
		} else if e != nil { //에러처리
			return c.JSON(
				http.StatusInternalServerError, logger.Error(
					c, &bodyStruct.ErrorStruct{
						Code:      "DownloadBeatmapSet-011",
						Path:      c.Path(),
						RequestId: c.Response().Header().Get("X-Request-ID"),
						Error:     e,
						Message:   "fail to read Bancho Stream",
					},
				),
			)
		}
	}
	if cLen == buf.Len() {
		return saveLocal(&buf, serverFileName, beatmap4d.BeatmapsetId)
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
