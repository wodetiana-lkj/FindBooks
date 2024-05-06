package utils

import (
	"regexp"
	"testing"
)

func TestChineseToArabic(t *testing.T) {
	pattern := `第(.*?)章`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch("萧炎云韵篇", -1)
	if len(matches) > 0 {
		println(matches[0][1])
	}
	println(0)
}
