package utils

import (
	"fmt"
	"time"
)

func Truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length]
}

func GetAge(t time.Time) string {
	d := time.Since(t)
	if int(d.Hours()) < 1 {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	} else if int(d.Hours()) < 24 {
		return fmt.Sprintf("%dh", int(d.Hours()))
	} else {
		return fmt.Sprintf("%dh", int(d.Hours())/24)
	}
}
