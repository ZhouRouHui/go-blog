package types

import (
	"goblog/pkg/logger"
	"strconv"
)

// Int64ToString 将 int64 转为 string
func Int64ToString(num int64) string {
	return strconv.FormatInt(num, 10)
}

// StringToInt 字符串转int
func StringToInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		logger.LogError(err)
	}
	return i
}
