package utils

import (
	"fmt"
	"strings"
)

var replacerNotUseFilename = strings.NewReplacer(
	`\`, `＼`,
	`/`, `／`,
	`:`, `：`,
	`*`, `＊`,
	`?`, `？`,
	`"`, `＂`,
	`<`, `＜`,
	`>`, `＞`,
	`|`, `｜`,
)

func BuildOsuFileName(artist, title, creator, version string) string {
	return replacerNotUseFilename.Replace(fmt.Sprintf("%s - %s (%s) [%s].osu", artist, title, creator, version))
}
