package models

import "time"

type Status int

func (s Status) String() string {
	switch s {
	case Open:
		return "Open"
	case Doing:
		return "Doing"
	case Done:
		return "Done"
	case Archived:
		return "Archived"
	default:
		return "Unknown"
	}
}

const (
	Open Status = iota
	Doing
	Done
	Archived
)

type Priority int

func (p Priority) String() string {
	switch p {
	case Low:
		return "Low"
	case Medium:
		return "Medium"
	case High:
		return "High"
	default:
		return "Unknown"
	}
}

const (
	Low Priority = iota
	Medium
	High
)

type Todo struct {
	ID          int64
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DueDate     *time.Time
	Status      Status
	Priority    Priority
	Tags        []string
}
