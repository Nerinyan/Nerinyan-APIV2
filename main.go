package main

import (
	"github.com/Nerinyan/Nerinyan-APIV2/banchoCrawler"
	"github.com/Nerinyan/Nerinyan-APIV2/config"
	"github.com/Nerinyan/Nerinyan-APIV2/db"
	"github.com/Nerinyan/Nerinyan-APIV2/logger"
	"github.com/Nerinyan/Nerinyan-APIV2/middlewareFunc"
	"github.com/Nerinyan/Nerinyan-APIV2/route/common"
	"github.com/Nerinyan/Nerinyan-APIV2/route/download"
	"github.com/Nerinyan/Nerinyan-APIV2/route/search"
	"github.com/Nerinyan/Nerinyan-APIV2/src"
	"github.com/Nerinyan/Nerinyan-APIV2/webhook"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pterm/pterm"
	"log"
	"net/http"
	"time"
)

// TODO DOING DB 테이블 없으면 자동으로 생성하게
// TODO DOING 헤더로 프론트인지 api 인지 구분할수있게
// TODO DOING 서버간 비트맵파일 해시값 비교해서 서로 다른경우 둘다 서버에서 삭제.
// TODO DOING 서버끼리 서로 비트맵파일 동기화 시킬수 있게
// TODO DOING 반쵸 비트맵 다운로드 제한 10분간 약 200건 10분 정지. (429 too many request) => 10분 내 100건 봇 감지 알고리즘
// TODO DOING 서버 자체적으로 10분당 150건 이내만 다운로드 가능하게 셋팅
// TODO DOING /status에 들어갈 상태값 추가.

func init() {
	ch := make(chan struct{})
	config.LoadConfig()
	src.StartIndex()
	db.ConnectRDBMS()
	middlewareFunc.StartHandler()
	go banchoCrawler.LoadBancho(ch)
	_ = <-ch

	if config.Config.Debug {
		//go banchoCrawler.UpdateAllPackList()
	} else {
		go banchoCrawler.RunGetBeatmapDataASBancho()
	}

}

func main() {
	e := echo.New()
	e.HideBanner = true
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		pterm.Error.WithShowLineNumber().Printfln("%+v", err)
		_ = c.JSON(
			http.StatusInternalServerError, map[string]interface{}{
				"error":      err,
				"request_id": c.Response().Header().Get("X-Request-Id"),
				"time":       time.Now(),
			},
		)

	}

	e.Renderer = &download.Renderer

	go func() {
		for {
			<-logger.Ch
			e.Logger.SetOutput(log.Writer())
			pterm.Info.Println("UPDATED ECHO LOGGER.")
		}
	}()

	e.Pre(
		//필수 우선순
		middleware.Recover(),
		middleware.RequestID(),
		middleware.RemoveTrailingSlash(),
		middleware.Logger(),
		middlewareFunc.RequestConsolLogger(),
		middleware.RemoveTrailingSlash(),

		//1차 필터
		middlewareFunc.BlackListHandler(), // 1분주기 갱신
		middleware.CORSWithConfig(middleware.CORSConfig{AllowOrigins: []string{"*"}, AllowMethods: []string{echo.GET, echo.HEAD, echo.POST}}),

		//2차 필터
		middleware.RateLimiter(
			middleware.NewRateLimiterMemoryStoreWithConfig(
				middleware.RateLimiterMemoryStoreConfig{
					Rate:      200,
					Burst:     1000,
					ExpiresIn: time.Minute,
				},
			),
		),

		//middleware.RateLimiterWithConfig(middleWareFunc.RateLimiterConfig),

	)

	// docs ============================================================================================================
	e.GET("/", common.Root)

	// 서버상태 체크용 ====================================================================================================
	e.GET("/health", common.Health)
	e.GET("/robots.txt", common.Robots)
	e.GET("/status", common.Status)

	// 맵 파일 다운로드 ===================================================================================================
	e.GET("/d/:setId", download.DownloadBeatmapSet, download.Embed)
	e.GET("/beatmap/:mapId", download.DownloadBeatmapSet)
	e.GET("/beatmapset/:setId", download.DownloadBeatmapSet)
	//TODO 맵아이디, 맵셋아이디 지원

	// 비트맵 BG  =========================================================================================================
	e.GET(
		"/bg/:setId", func(c echo.Context) error {
			redirectUrl := "https://subapi.nerinyan.moe/bg/" + c.Param("setId")
			return c.Redirect(http.StatusPermanentRedirect, redirectUrl)
		},
	)

	// 비트맵 리스트 검색용 ================================================================================================
	e.GET("/search", search.Search)
	e.POST("/search", search.Search)

	// 개발중 || 테스트중 ===================================================================================================

	// ====================================================================================================================
	pterm.Info.Println("ECHO STARTED AT", config.Config.Port)
	webhook.DiscordInfoStartUP()
	e.Logger.Fatal(e.Start(":" + config.Config.Port))

}
