package db

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"regexp"
	"time"
)

type ExecQueue struct {
	Ch    chan error
	Query string
	Args  []any
}

var InsertQueueChannel = make(chan ExecQueue, 10)
var queryNameRegex, _ = regexp.Compile("^(/[*])(.+?)([*]/)")

func AddInsertQueue(query string, args ...any) {
	st := time.Now().UnixMilli()
	InsertQueueChannel <- ExecQueue{
		Query: query,
		Args:  args,
	}
	et := time.Now().UnixMilli() - st
	if et > 3000 {
		pterm.Warning.Println("SLOW QUERY WARN", et, "ms\n", query)
	}
}

func init() {
	go func() {
		for {

			for ins := range InsertQueueChannel {
				st := time.Now().UnixMilli()
				result := Gorm.Exec(ins.Query, ins.Args...)
				err := result.Error
				ins := ins

				go func() { //로그용
					tagSize := len(ins.Query)
					if tagSize > 100 {
						tagSize = 100
					}
					queryName := fmt.Sprintf("%s | %-50s | ", time.Now().Format("15:04:05.000"), pterm.Yellow(queryNameRegex.FindString(ins.Query[:tagSize])))

					if err == nil && result.RowsAffected > 0 {
						pterm.Info.Printfln(queryName+"%5d ROWS | %6dms |", result.RowsAffected, time.Now().UnixMilli()-st)
					}
					if ins.Ch != nil {
						if err != nil {
							ins.Ch <- errors.New(queryName + err.Error() + ins.Query)
						} else {
							ins.Ch <- nil
						}
					} else if err != nil {
						pterm.Error.Println(queryName+err.Error(), ins.Query)
					}
				}()
			}
		}
	}()
}
