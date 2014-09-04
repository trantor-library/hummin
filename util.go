package main

import (
	"strings"
)

func hasSubString(str string, substr string) bool {
	matches := strings.Count(str, substr)
	return matches >= 1

}
