package src

import (
	"github.com/Nerinyan/Nerinyan-APIV2/config"
	"github.com/Nerinyan/Nerinyan-APIV2/utils"
	"github.com/pterm/pterm"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type FileIndex map[int]time.Time

var FileList = make(FileIndex)
var FileSizeToString = totalFileSize()
var fileSize uint64

const goos = runtime.GOOS

func StartIndex() {
	FileListUpdate()
	go func() {
		time.Sleep(time.Second * 60 * 5)
		for {
			FileListUpdate()
			time.Sleep(time.Second * 60 * 5)
		}
	}()

}
func FileListUpdate() {
	var err error

	checkDir()
	dirs, err := os.ReadDir(config.Config.TargetDir)
	if err != nil {
		return
	}

	tmp := make(FileIndex)
	fileSize = 0
	for _, dir := range dirs {
		fi, err := dir.Info()
		if err != nil || dir.IsDir() {
			continue
		}
		if sid, err := strconv.Atoi(strings.Replace(dir.Name(), ".osz", "", -1)); err == nil {
			tmp[sid] = fi.ModTime()
			fileSize += uint64(fi.Size())
		}
	}
	FileSizeToString = totalFileSize()
	FileList = tmp
	pterm.Info.Printfln(
		"%s File List Indexing : %s files [%s]",
		time.Now().Format("2006-01-02 15:04:05"),
		pterm.LightYellow(strconv.Itoa(len(FileList))),
		pterm.LightYellow(totalFileSize()),
	)

}

func totalFileSize() (s string) {
	return utils.ToHumanDataSize(fileSize)
}

func checkDir() {
	if _, e := os.Stat(config.Config.TargetDir); os.IsNotExist(e) {
		err := os.MkdirAll(config.Config.TargetDir, 666)
		if err != nil {
			pterm.Error.Println(err)
			panic(err)
		}
	}
}
