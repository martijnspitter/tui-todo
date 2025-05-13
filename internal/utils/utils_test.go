package utils

import (
	"fmt"
	"strings"
	"testing"

	"pgregory.net/rapid"
)

// Traditional example-based tests
func TestFormatTime(t *testing.T) {
	testCases := []struct {
		name           string
		inputSeconds   int64
		expectedOutput string
	}{
		{
			name:           "seconds only",
			inputSeconds:   45,
			expectedOutput: "45s",
		},
		{
			name:           "minutes and seconds",
			inputSeconds:   125,
			expectedOutput: "2m 5s",
		},
		{
			name:           "just minutes",
			inputSeconds:   60,
			expectedOutput: "1m 0s",
		},
		{
			name:           "hours, minutes, and seconds",
			inputSeconds:   3725,
			expectedOutput: "1h 2m 5s",
		},
		{
			name:           "zero seconds",
			inputSeconds:   0,
			expectedOutput: "0s",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := FormatTime(tc.inputSeconds)
			if result != tc.expectedOutput {
				t.Errorf("FormatTime(%d) = %s; want %s", tc.inputSeconds, result, tc.expectedOutput)
			}
		})
	}
}

// Property-based tests using rapid
func TestFormatTimeProperties(t *testing.T) {
	// Property: The formatted time should be consistent with the input seconds
	t.Run("consistency", rapid.MakeCheck(func(t *rapid.T) {
		// Generate random non-negative seconds
		seconds := rapid.Int64Range(0, 100000).Draw(t, "seconds")

		formatted := FormatTime(seconds)

		// Parse the formatted string back to seconds for verification
		var h, m, s int64

		if strings.Contains(formatted, "h") {
			_, err := fmt.Sscanf(formatted, "%dh %dm %ds", &h, &m, &s)
			if err != nil {
				t.Fatalf("Failed to parse formatted time: %s", formatted)
			}
		} else if strings.Contains(formatted, "m") {
			_, err := fmt.Sscanf(formatted, "%dm %ds", &m, &s)
			if err != nil {
				t.Fatalf("Failed to parse formatted time: %s", formatted)
			}
		} else {
			_, err := fmt.Sscanf(formatted, "%ds", &s)
			if err != nil {
				t.Fatalf("Failed to parse formatted time: %s", formatted)
			}
		}

		// Calculate total seconds from parsed values
		parsedSeconds := h*3600 + m*60 + s

		// Verify the result matches the original input
		if parsedSeconds != seconds {
			t.Errorf("FormatTime(%d) resulted in %s, which parses back to %d seconds",
				seconds, formatted, parsedSeconds)
		}
	}))

	// Property: The format of the output string should follow the expected pattern
	t.Run("format", rapid.MakeCheck(func(t *rapid.T) {
		seconds := rapid.Int64Range(0, 100000).Draw(t, "seconds")

		formatted := FormatTime(seconds)

		// Expected format checks
		hours := seconds / 3600
		minutes := (seconds % 3600) / 60
		secs := seconds % 60

		if hours > 0 {
			expected := fmt.Sprintf("%dh %dm %ds", hours, minutes, secs)
			if formatted != expected {
				t.Errorf("Expected format %s for %d seconds, but got %s",
					expected, seconds, formatted)
			}
		} else if minutes > 0 {
			expected := fmt.Sprintf("%dm %ds", minutes, secs)
			if formatted != expected {
				t.Errorf("Expected format %s for %d seconds, but got %s",
					expected, seconds, formatted)
			}
		} else {
			expected := fmt.Sprintf("%ds", secs)
			if formatted != expected {
				t.Errorf("Expected format %s for %d seconds, but got %s",
					expected, seconds, formatted)
			}
		}
	}))
}
