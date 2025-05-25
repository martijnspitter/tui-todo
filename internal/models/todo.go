package models

import (
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/martijnspitter/tui-todo/internal/theme"
	"github.com/martijnspitter/tui-todo/internal/utils"
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
	case Blocked:
		return "status.blocked"
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
	case Blocked:
		return theme.BlockedStatusColor
	default:
		return theme.OpenStatusColor
	}
}

const (
	Open Status = iota
	Doing
	Done
	Blocked
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
	TimeSpent   int64      // Total time spent in seconds
	TimeStarted *time.Time // When the task was last set to Doing status
}

// FormatTimeSpent returns a human-readable format of the time spent on this todo
func (t *Todo) FormatTimeSpent() string {
	// Calculate total seconds including current session
	totalSeconds := t.GetTotalSeconds()

	return utils.FormatTime(totalSeconds)
}

// IsCurrentlyTracking returns whether this todo is actively tracking time
func (t *Todo) isCurrentlyTracking() bool {
	return t.Status == Doing && t.TimeStarted != nil
}

// GetTotalSeconds returns the total seconds spent including any current tracking session
func (t *Todo) GetTotalSeconds() int64 {
	if !t.isCurrentlyTracking() {
		return t.TimeSpent
	}

	currentSessionSeconds := int64(time.Since(*t.TimeStarted).Seconds())
	return t.TimeSpent + currentSessionSeconds
}
