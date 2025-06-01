package models

import "time"

type Tag struct {
	Name        string
	Description string
	ID          int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
