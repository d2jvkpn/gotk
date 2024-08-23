package gotk

import (
	// "fmt"
	"math"
	"regexp"
	"strings"
)

var (
	_RE_FirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	_RE_AllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

// https://stackoverflow.com/posts/56616250/revisions
func ToSnakeCase(str string) string {
	snake := _RE_FirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = _RE_AllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func Round3(value float64) float64 {
	return math.Round(value*1e3) / 1e3
}
