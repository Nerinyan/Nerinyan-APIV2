package route

import (
	"github.com/Nerinyan/Nerinyan-APIV2/banchoCroller"
	"github.com/Nerinyan/Nerinyan-APIV2/src"
	"github.com/labstack/echo/v4"
	"net/http"
	"runtime"
)

func Status(c echo.Context) error {
	return c.JSON(
		http.StatusOK, map[string]interface{}{
			"CpuThreadCount":        runtime.NumCPU(),
			"RunningGoroutineCount": runtime.NumGoroutine(),
			"apiCount":              *banchoCroller.ApiCount,
			"fileCount":             len(src.FileList),
			"fileSize":              src.FileSizeToString,
		},
	)
}
