package utils

import (
	"strings"
)

func SplitTrim(s, sep string) (res []string) {
	ss := strings.Split(s, sep)
	for _, i := range ss {
		res = append(res, strings.TrimSpace(i))
	}
	return
}
func SplitTrimUpper(s, sep string) (res []string) {
	ss := strings.Split(strings.ToUpper(s), sep)
	for _, i := range ss {
		res = append(res, strings.TrimSpace(i))
	}
	return
}
func SplitTrimLower(s, sep string) (res []string) {
	ss := strings.Split(strings.ToLower(s), sep)
	for _, i := range ss {
		res = append(res, strings.TrimSpace(i))
	}
	return
}
func SplitTrimLowerUniqueLen(s, sep string) int {
	inResult := make(map[string]bool)
	var result []string
	for _, str := range SplitTrimLower(s, sep) {
		if str == "" {
			continue
		}
		if _, ok := inResult[str]; !ok {
			inResult[str] = true
			result = append(result, str)
		}
	}
	return len(result)
}

func NotInMapFindDefault[T string | []int](m map[string]T, key string) (val T) {
	if v, ok := m[key]; ok {
		return v
	} else {
		return m["default"]
	}
}

func NotInMapFindAllDefault[T any](m map[string]T, keys []string) (val []T) {
	for _, key := range keys {
		if v, ok := m[key]; ok {
			val = append(val, v)
		}
	}
	if len(val) < 1 {
		val = append(val, m["default"])
	}
	return
}

func NotInMapFindAllAppendDefault[T any](m map[string][]T, keys []string) (val []T) {
	for _, key := range keys {
		if v, ok := m[key]; ok {
			val = append(val, v...)
		}
	}
	if len(val) < 1 {
		val = append(val, m["default"]...)
	}
	return
}

func TrimLower(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func TrimUpper(s string) string {
	return strings.ToUpper(strings.TrimSpace(s))
}
