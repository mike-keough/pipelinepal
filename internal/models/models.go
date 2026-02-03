package models

import "time"

type Lead struct {
	ID             int64
	Name           string
	Phone          string
	Email          string
	Source         string
	Kind           string // lead|buyer|seller
	Status         string
	Notes          string
	CreatedAt      time.Time
	LastContactAt  *time.Time
	NextFollowUpAt *time.Time
}
