package models

import (
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/martijnspitter/tui-todo/internal/theme"
)

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

func (s Status) Color() lipgloss.Color {
	switch s {
	case Open:
		return theme.OpenStatusColor
	case Doing:
		return theme.DoingStatusColor
	case Done:
		return theme.DoneStatusColor
	default:
		return theme.OpenStatusColor
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

func (p Priority) Color() lipgloss.Color {
	switch p {
	case Low:
		return theme.LowPriorityColor
	case Medium:
		return theme.MediumPriorityColor
	case High:
		return theme.HighPriorityColor
	case Major:
		return theme.MajorPriorityColor
	case Critical:
		return theme.CriticalPriorityColor
	default:
		return theme.MediumPriorityColor
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
