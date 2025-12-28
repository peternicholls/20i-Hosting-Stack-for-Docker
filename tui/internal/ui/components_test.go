// Project: 20i Stack Manager TUI
// File: components_test.go
// Purpose: Unit tests for UI components
// Version: 0.1.0
// Updated: 2025-12-28

package ui_test

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/peternicholls/20i-stack/tui/internal/ui"
)

func TestStatusIconGlyphs(t *testing.T) {
	// Ensure truecolor profile so lipgloss emits color on supported terminals during tests
	lipgloss.SetColorProfile(termenv.TrueColor)

	tests := []struct {
		status string
		glyph  string
	}{
		{"running", "●"},
		{"stopped", "○"},
		{"restarting", "⚠"},
		{"error", "✗"},
		{"unknown", "?"},
		{"", "?"},
	}

	for _, tc := range tests {
		out := ui.StatusIcon(tc.status)
		if !strings.Contains(out, tc.glyph) {
			t.Fatalf("StatusIcon(%q) = %q; want glyph %q", tc.status, out, tc.glyph)
		}
	}
}
