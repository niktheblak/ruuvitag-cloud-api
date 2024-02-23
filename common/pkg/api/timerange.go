package api

import (
	"fmt"
	"time"
)

func ParseTimeRange(fromStr, toStr string) (from time.Time, to time.Time, err error) {
	if fromStr != "" {
		from, err = time.Parse("2006-01-02", fromStr)
	}
	if err != nil {
		return
	}
	if toStr != "" {
		to, err = time.Parse("2006-01-02", toStr)
	}
	if err != nil {
		return
	}
	if !from.IsZero() && !to.IsZero() && from == to {
		to = to.AddDate(0, 0, 1)
	}
	if to.IsZero() || to.After(time.Now()) {
		to = time.Now().UTC()
	}
	if from.After(to) {
		err = fmt.Errorf("from timestamp cannot be after to timestamp")
	}
	return
}
