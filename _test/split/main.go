package main

import (
	"fmt"
	"regexp"
	"strings"
)

func main() {
	input := `"kimi" no`
	fmt.Println(strings.Join(SplitInput(input), ", ")) // Outputs: [hello world this is a quoted string this is also quoted this too]
}

func SplitInput(input string) []string {
	// The regex matches any string surrounded by quotes or backticks,
	// or any non-space character.
	re := regexp.MustCompile(`"([^"]*)"|'([^']*)'|\S+`)
	matches := re.FindAllString(input, -1)

	result := make([]string, len(matches))
	for i, match := range matches {
		// Trim any surrounding quotes or backticks from the match.
		result[i] = strings.Trim(match, "`'\"")
	}

	return result
}
