package pkg

import "strconv"

func Max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
func TruncateString(str string, num int) string {
	truncated := str
	if num <= 3 {
		return str
	}

	if len(str) > num {
		if num > 3 {
			num -= 3
		}

		skipped := len(str) - num
		truncated = str[0:skipped] + "_"
	}

	return truncated
}

func CastIntToStr(a int) string {
	str := strconv.Itoa(a)

	return str
}
