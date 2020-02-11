package server

import (
	"strconv"
)

func ParseLimit(limitStr string) int {
	var limit int64
	if limitStr != "" {
		limit, _ = strconv.ParseInt(limitStr, 10, 64)
	}
	if limit <= 0 {
		limit = 20
	}
	return int(limit)
}
