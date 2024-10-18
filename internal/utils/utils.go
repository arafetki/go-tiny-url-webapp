package utils

import "time"

func IsExpired(t time.Time) bool {
	return time.Now().After(t)
}
