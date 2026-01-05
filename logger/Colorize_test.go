package logger

import "testing"

func TestColorize(t *testing.T) {
	tests := []struct {
		name     string
		color    string
		text     string
		expected string
	}{
		{"red text", Red, "error", Red + "error" + Reset},
		{"green text", Green, "success", Green + "success" + Reset},
		{"yellow text", Yellow, "warning", Yellow + "warning" + Reset},
		{"blue text", Blue, "info", Blue + "info" + Reset},
		{"empty text", Green, "", Green + "" + Reset},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Colorize(tt.color, tt.text)
			if result != tt.expected {
				t.Errorf("Colorize(%q, %q) = %q, want %q", tt.color, tt.text, result, tt.expected)
			}
		})
	}
}

func TestColorConstants(t *testing.T) {
	// Verify ANSI escape codes are correct
	if Red != "\033[31m" {
		t.Errorf("Red = %q, want %q", Red, "\033[31m")
	}
	if Green != "\033[32m" {
		t.Errorf("Green = %q, want %q", Green, "\033[32m")
	}
	if Yellow != "\033[33m" {
		t.Errorf("Yellow = %q, want %q", Yellow, "\033[33m")
	}
	if Blue != "\033[34m" {
		t.Errorf("Blue = %q, want %q", Blue, "\033[34m")
	}
	if Reset != "\033[0m" {
		t.Errorf("Reset = %q, want %q", Reset, "\033[0m")
	}
}
