package common

import (
	"github.com/Nerinyan/Nerinyan-APIV2/banchoCrawler"
	"github.com/Nerinyan/Nerinyan-APIV2/src"
	"github.com/labstack/echo/v4"
	"net/http"
	"runtime"
)

func Status(c echo.Context) error {
	return c.JSON(
		http.StatusOK, map[string]interface{}{
			"cpuThreadCount":        runtime.NumCPU(),
			"runningGoroutineCount": runtime.NumGoroutine(),
			"beatmapSetCount":       src.BeatmapSetCount,
			"apiCount":              *banchoCrawler.ApiCount,
			"fileCount":             len(src.FileList),
			"fileSize":              src.FileSizeToString,
		},
	)
}
