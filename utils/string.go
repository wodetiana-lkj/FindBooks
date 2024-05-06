package utils

import (
	"strconv"
	"strings"
)

var chineseToArabicMap = map[rune]int{
	'零': 0,
	'一': 1,
	'二': 2,
	'三': 3,
	'四': 4,
	'五': 5,
	'六': 6,
	'七': 7,
	'八': 8,
	'九': 9,
}

var chineseToArabicUnitMap = map[rune]int{
	'十': 10,
	'百': 100,
	'千': 1000,
	'万': 10000,
	'亿': 100000000,
}

func Concat(arr ...string) string {
	var builder strings.Builder
	for i := range arr {
		builder.WriteString(arr[i])
	}
	return builder.String()
}

/**
 * 中文转阿拉伯数字 连续单位无法判断
 */
func ChineseToArabic(chinese string) int {
	var total int = 0
	var val int
	var unit = 1
	var tmp int
	for _, char := range chinese {
		val = chineseToArabicMap[char]
		// 单位值
		if val == 0 && char != 38646 {
			unit = chineseToArabicUnitMap[char]
			if tmp == 0 {
				tmp = 1
			}
			tmp *= unit
			total += tmp
			tmp = 0
		} else {
			tmp += val
		}
	}

	if tmp != 0 {
		total += tmp
	}
	return total
}

func IsNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	if err == nil {
		return true
	}
	_, err = strconv.ParseFloat(s, 64)
	if err == nil {
		return true
	}
	return false
}
