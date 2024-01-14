package model

import "time"

type Notification struct {
	Title     string
	Refresher string
	Domain    string
	Time      time.Time
	NewIP     string
	OldIP     string
}
