package utils

import (
	"encoding/json"
	"fmt"
	"github.com/pterm/pterm"
	"runtime"
	"strings"
)

func ToJsonIndentString(i any) (str *[]byte) {
	b, err := json.MarshalIndent(i, "", "    ")
	if err != nil {
		pterm.Error.WithShowLineNumber().Println(err)
		return
	}
	return &b
}
func ToJsonString(i any) (str *[]byte) {
	b, err := json.Marshal(i)
	if err != nil {
		pterm.Error.WithShowLineNumber().Println(err)
		return
	}
	return &b
}

func StringRepeatArray(s string, count int) (arr []string) {
	for i := 0; i < count; i++ {
		arr = append(arr, s)
	}
	return
}

func TernaryOperator[V any](tf bool, t, f V) V {
	if tf {
		return t
	}
	return f
}

func MakeArrayUnique[T comparable](array *[]T) (res []T) {
	keys := make(map[T]struct{})
	for _, s := range *array {
		keys[s] = struct{}{}
	}
	for i := range keys {
		res = append(res, i)
	}
	return
}
func MakeStringArrayUniqueAndCheckLength(array *[]string, limit int) (res []string) {

	keys := make(map[string]struct{})
	for _, s := range *array {
		keys[s] = struct{}{}
	}
	for i := range keys {
		if len(i) <= limit {
			res = append(res, i)

		}
	}
	return
}

func MakeArrayUniqueInterface[T comparable](array *[]T) (res []interface{}) {

	keys := make(map[T]struct{})
	for _, s := range *array {
		keys[s] = struct{}{}
	}
	for i := range keys {
		res = append(res, i)
	}
	return
}
func StringRepeatJoin(str, sep string, count int) string {
	var repeatBuf []string
	for i := 0; i < count; i++ {
		repeatBuf = append(repeatBuf, str)
	}
	return strings.Join(repeatBuf, sep)

}
func GetFileLine() string {
	_, file, line, _ := runtime.Caller(1)
	return fmt.Sprintf("%s:%d ", file, line)
}
