package main

import (
	"fmt"
	"strings"
	"time"
)

func findRepeats(s string) []string {
	var result []string
	checkMap := make(map[string]int)
	length := len(s)
	for i := 0; i < length; i++ {
		for j := i + 2; j <= length; j++ {
			substr := s[i:j]
			v, ok := checkMap[substr]
			if ok {
				if v == 1 {
					result = append(result, substr)
				}
				checkMap[substr]++
			} else {
				checkMap[substr] = 1
			}
		}
	}

	// Remove smaller substrings if they're part of a larger substring
	length = len(result)
	for i := 0; i < length; i++ {
		for j := i + 1; j < length; j++ {
			if strings.Contains(result[i], result[j]) {
				result[j] = result[i]
			} else if strings.Contains(result[j], result[i]) {
				result[i] = result[j]
			}
		}
	}
	// Deduplicate the result
	keys := make(map[string]bool)
	var list []string
	for _, entry := range result {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func main() {
	st := time.Now().UnixMilli()
	for i := 0; i < 10000; i++ {
		findRepeats("kiwi")           // Output: ["kiwi"]
		findRepeats("owoowo")         // Output: ["owo"]
		findRepeats("yukiyukiowoowo") // Output: ["yuki" "owo"]
	}
	et := time.Now().UnixMilli()
	fmt.Println(et-st, "ms")
	//fmt.Println(findRepeats("kiwi"))           // Output: ["kiwi"]
	//fmt.Println(findRepeats("owoowo"))         // Output: ["owo"]
	//fmt.Println(findRepeats("yukiyukiowoowo")) // Output: ["yuki" "owo"]
}
