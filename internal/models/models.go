package models

import "time"

type User struct {
	Name     string
	Birthday time.Time
	ChatID   int64
}
