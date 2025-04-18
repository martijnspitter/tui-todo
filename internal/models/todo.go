package models

import "time"

type Status int

func (s Status) String() string {
	switch s {
	case Open:
		return "status.open"
	case Doing:
		return "status.doing"
	case Done:
		return "status.done"
	default:
		return "status.unknown"
	}
}

const (
	Open Status = iota
	Doing
	Done
)

type Priority int

func (p Priority) String() string {
	switch p {
	case Low:
		return "priority.low"
	case Medium:
		return "priority.medium"
	case High:
		return "priority.high"
	case Major:
		return "priority.major"
	case Critical:
		return "priority.critical"
	default:
		return "priority.unknown"
	}
}

const (
	Low Priority = iota
	Medium
	High
	Major
	Critical
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
	Archived    bool
}
