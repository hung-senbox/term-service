package helper

import "time"

func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func ValidateDateRange(start, end time.Time) bool {
	return !start.After(end)
}
