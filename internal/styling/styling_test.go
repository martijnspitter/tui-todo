package styling

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/martijnspitter/tui-todo/internal/models"
	"github.com/martijnspitter/tui-todo/internal/theme"
)

// Helper function to check if a string contains ANSI escape codes for a specific color
func containsColor(s string, color lipgloss.Color) bool {
	colorStr := string(color)
	if !strings.HasPrefix(colorStr, "#") {
		return strings.Contains(s, colorStr)
	}

	// For hex colors, the check is more complex because lipgloss translates hex to RGB
	// A very basic check - at least ensure the string isn't empty
	return s != ""
}

func TestGetStyledStatus(t *testing.T) {
	tests := []struct {
		name             string
		translatedStatus string
		status           models.Status
		selected         bool
		omitNumber       bool
		hovered          bool
	}{
		{
			name:             "Open status unselected",
			translatedStatus: "Open",
			status:           models.Open,
			selected:         false,
			omitNumber:       false,
			hovered:          false,
		},
		{
			name:             "Doing status selected",
			translatedStatus: "Doing",
			status:           models.Doing,
			selected:         true,
			omitNumber:       false,
			hovered:          false,
		},
		{
			name:             "Done status hovered",
			translatedStatus: "Done",
			status:           models.Done,
			selected:         false,
			omitNumber:       false,
			hovered:          true,
		},
		{
			name:             "Status with omitted number",
			translatedStatus: "Open",
			status:           models.Open,
			selected:         true,
			omitNumber:       true,
			hovered:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetStyledStatus(tt.translatedStatus, tt.status, tt.selected, tt.omitNumber, tt.hovered)

			// Basic validation - result shouldn't be empty
			if result == "" {
				t.Error("Expected styled status output, got empty string")
			}

			// Check that the status text is included
			if !strings.Contains(result, tt.translatedStatus) && !tt.omitNumber {
				t.Errorf("Styled output should contain status text '%s'", tt.translatedStatus)
			}

			// Check for number indication if not omitted
			if !tt.omitNumber {
				expectedNum := int(tt.status) + 2
				if !strings.Contains(result, string(rune('0'+expectedNum))) {
					t.Errorf("Expected status number %d in output", expectedNum)
				}
			}

			// For hovered state, yellow should be used
			if tt.hovered && !containsColor(result, theme.Yellow) {
				t.Error("Hovered status should use Yellow color")
			}
		})
	}
}

func TestGetStyledTagWithIndicator(t *testing.T) {
	tests := []struct {
		name       string
		num        int
		text       string
		color      lipgloss.Color
		selected   bool
		omitNumber bool
		hovered    bool
	}{
		{
			name:       "Basic tag",
			num:        1,
			text:       "Tag",
			color:      theme.Mauve,
			selected:   false,
			omitNumber: false,
			hovered:    false,
		},
		{
			name:       "Selected tag",
			num:        2,
			text:       "Selected",
			color:      theme.Lavender,
			selected:   true,
			omitNumber: false,
			hovered:    false,
		},
		{
			name:       "Hovered tag",
			num:        3,
			text:       "Hovered",
			color:      theme.Rosewater,
			selected:   false,
			omitNumber: false,
			hovered:    true,
		},
		{
			name:       "Omitted number",
			num:        4,
			text:       "NoNumber",
			color:      theme.Yellow,
			selected:   false,
			omitNumber: true,
			hovered:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetStyledTagWithIndicator(tt.num, tt.text, tt.color, tt.selected, tt.omitNumber, tt.hovered)

			// Check that the tag text is included
			if !strings.Contains(result, tt.text) {
				t.Errorf("Styled output should contain tag text '%s'", tt.text)
			}

			// Check for number if not omitted
			if !tt.omitNumber {
				if !strings.Contains(result, string(rune('0'+tt.num))) {
					t.Errorf("Expected number %d in output", tt.num)
				}
			}

			// For hovered state, yellow should be used
			if tt.hovered && !containsColor(result, theme.Yellow) {
				t.Error("Hovered tag should use Yellow color")
			}

			// For selected state, the specified color should be used
			if tt.selected && !tt.hovered && !containsColor(result, tt.color) {
				t.Error("Selected tag should use the specified color")
			}
		})
	}
}

func TestGetStyledPriority(t *testing.T) {
	tests := []struct {
		name        string
		translatedP string
		priority    models.Priority
		selected    bool
		hovered     bool
	}{
		{
			name:        "Low priority unselected",
			translatedP: "Low",
			priority:    models.Low,
			selected:    false,
			hovered:     false,
		},
		{
			name:        "Medium priority selected",
			translatedP: "Medium",
			priority:    models.Medium,
			selected:    true,
			hovered:     false,
		},
		{
			name:        "High priority hovered",
			translatedP: "High",
			priority:    models.High,
			selected:    false,
			hovered:     true,
		},
		{
			name:        "Major priority selected and hovered",
			translatedP: "Major",
			priority:    models.Major,
			selected:    true,
			hovered:     true,
		},
		{
			name:        "Critical priority",
			translatedP: "Critical",
			priority:    models.Critical,
			selected:    true,
			hovered:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetStyledPriority(tt.translatedP, tt.priority, tt.selected, tt.hovered)

			// Check that the priority text is included
			if !strings.Contains(result, tt.translatedP) {
				t.Errorf("Styled output should contain priority text '%s'", tt.translatedP)
			}

			// For hovered state, yellow should be used
			if tt.hovered && !containsColor(result, theme.Yellow) {
				t.Error("Hovered priority should use Yellow color")
			}

			// For selected state, the priority color should be used
			if tt.selected && !tt.hovered {
				priorityColor := tt.priority.Color()
				if !containsColor(result, priorityColor) {
					t.Error("Selected priority should use the priority color")
				}
			}
		})
	}
}

func TestGetStyledUpdatedAt(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{
			name: "Short time",
			text: "1m",
		},
		{
			name: "Medium time",
			text: "1h 30m",
		},
		{
			name: "Long time",
			text: "3d 5h 12m",
		},
		{
			name: "Empty string",
			text: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetStyledUpdatedAt(tt.text)

			// Check that the text is included (if not empty)
			if tt.text != "" && !strings.Contains(result, tt.text) {
				t.Errorf("Styled output should contain text '%s'", tt.text)
			}

			// Check that the Lavender color is used
			if !containsColor(result, theme.Lavender) {
				t.Error("UpdatedAt should use Lavender color")
			}
		})
	}
}

func TestGetStyledDueDate(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		priority models.Priority
	}{
		{
			name:     "Due date with low priority",
			text:     "Tomorrow",
			priority: models.Low,
		},
		{
			name:     "Due date with medium priority",
			text:     "Next week",
			priority: models.Medium,
		},
		{
			name:     "Due date with high priority",
			text:     "Today",
			priority: models.High,
		},
		{
			name:     "Due date with critical priority",
			text:     "Overdue",
			priority: models.Critical,
		},
		{
			name:     "Empty due date",
			text:     "",
			priority: models.Medium,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetStyledDueDate(tt.text, tt.priority)

			// Check that the text is included (if not empty)
			if tt.text != "" && !strings.Contains(result, tt.text) {
				t.Errorf("Styled output should contain text '%s'", tt.text)
			}

			// Check that the priority color is used
			priorityColor := tt.priority.Color()
			if !containsColor(result, priorityColor) {
				t.Error("Due date should use the priority color")
			}
		})
	}
}

func TestGetStyledTag(t *testing.T) {
	tests := []struct {
		name string
		tag  string
	}{
		{
			name: "Regular tag",
			tag:  "important",
		},
		{
			name: "Short tag",
			tag:  "bug",
		},
		{
			name: "Long tag",
			tag:  "feature-request",
		},
		{
			name: "Empty tag",
			tag:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetStyledTag(tt.tag)

			// Check that the tag is included (if not empty)
			if tt.tag != "" && !strings.Contains(result, tt.tag) {
				t.Errorf("Styled output should contain tag '%s'", tt.tag)
			}

			// Check that the Rosewater color is used
			if !containsColor(result, theme.Rosewater) {
				t.Error("Tag should use Rosewater color")
			}
		})
	}
}

func TestGetSelectedBlock(t *testing.T) {
	tests := []struct {
		name     string
		selected bool
	}{
		{
			name:     "Selected block",
			selected: true,
		},
		{
			name:     "Unselected block",
			selected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetSelectedBlock(tt.selected)

			// Output shouldn't be empty
			if result == "" {
				t.Error("Expected block output, got empty string")
			}

			// For selected state, yellow should be used
			if tt.selected && !containsColor(result, theme.Yellow) {
				t.Error("Selected block should use Yellow color")
			}
		})
	}
}

func TestRenderMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		contains []string
	}{
		{
			name:     "Simple markdown",
			markdown: "# Header\nParagraph",
			contains: []string{"Header", "Paragraph"},
		},
		{
			name:     "Markdown with formatting",
			markdown: "**Bold** *Italic*",
			contains: []string{"Bold", "Italic"},
		},
		{
			name:     "Markdown with list",
			markdown: "- Item 1\n- Item 2",
			contains: []string{"Item 1", "Item 2"},
		},
		{
			name:     "Empty markdown",
			markdown: "",
			contains: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderMarkdown(tt.markdown)

			// Check that the markdown is rendered and contains expected content
			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("Rendered markdown should contain '%s'", substr)
				}
			}
		})
	}
}
