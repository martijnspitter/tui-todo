package models

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/martijnspitter/tui-todo/internal/theme"
)

func TestStatus_String(t *testing.T) {
	tests := []struct {
		name     string
		status   Status
		expected string
	}{
		{"Open status", Open, "status.open"},
		{"Doing status", Doing, "status.doing"},
		{"Done status", Done, "status.done"},
		{"Invalid status", Status(999), "status.unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.String(); got != tt.expected {
				t.Errorf("Status.String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestStatus_Color(t *testing.T) {
	tests := []struct {
		name     string
		status   Status
		expected lipgloss.Color
	}{
		{"Open status color", Open, theme.OpenStatusColor},
		{"Doing status color", Doing, theme.DoingStatusColor},
		{"Done status color", Done, theme.DoneStatusColor},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.Color(); got != tt.expected {
				t.Errorf("Status.Color() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestPriority_String(t *testing.T) {
	tests := []struct {
		name     string
		priority Priority
		expected string
	}{
		{"Low priority", Low, "priority.low"},
		{"Medium priority", Medium, "priority.medium"},
		{"High priority", High, "priority.high"},
		{"Major priority", Major, "priority.major"},
		{"Critical priority", Critical, "priority.critical"},
		{"Invalid priority", Priority(999), "priority.unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.priority.String(); got != tt.expected {
				t.Errorf("Priority.String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestPriority_Color(t *testing.T) {
	tests := []struct {
		name     string
		priority Priority
		expected lipgloss.Color
	}{
		{"Low priority color", Low, theme.LowPriorityColor},
		{"Medium priority color", Medium, theme.MediumPriorityColor},
		{"High priority color", High, theme.HighPriorityColor},
		{"Major priority color", Major, theme.MajorPriorityColor},
		{"Critical priority color", Critical, theme.CriticalPriorityColor},
		{"Invalid priority color", Priority(999), theme.MediumPriorityColor},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.priority.Color(); got != tt.expected {
				t.Errorf("Priority.Color() = %q, want %q", got, tt.expected)
			}
		})
	}
}

// Test completeness - verify we've covered all enum values
func TestCompleteness(t *testing.T) {
	// Check that we have tests for all Status values
	for s := Status(0); s <= Done; s++ {
		t.Run("Status completeness "+s.String(), func(t *testing.T) {
			// Just calling String() and Color() verifies they don't panic
			_ = s.String()
			_ = s.Color()
		})
	}

	// Check that we have tests for all Priority values
	for p := Priority(0); p <= Critical; p++ {
		t.Run("Priority completeness "+p.String(), func(t *testing.T) {
			// Just calling String() and Color() verifies they don't panic
			_ = p.String()
			_ = p.Color()
		})
	}
}

// Possible bug test - Done status is using DoingStatusColor
func TestDoneStatus_HasCorrectColor(t *testing.T) {
	// This test will fail with your current implementation
	// because Done is using DoingStatusColor instead of DoneStatusColor
	if Done.Color() != theme.DoneStatusColor {
		t.Errorf("Done status color = %q, want %q", Done.Color(), theme.DoneStatusColor)
	}
}
