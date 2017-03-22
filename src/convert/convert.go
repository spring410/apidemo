package convert

import (
	"strconv"
)

//This is to wrap strconv to be a simple package.

func StringToInt(s string) (int, error) {
	return strconv.Atoi(s)
}

func StringToIntBase16(s string) (int64, error) {
	return strconv.ParseInt(s, 16, 64)
}

func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func IntToString(i int) string {
	return strconv.Itoa(i)
}
