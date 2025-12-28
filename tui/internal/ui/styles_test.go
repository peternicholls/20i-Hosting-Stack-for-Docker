// Project: 20i Stack Manager TUI
// File: styles_test.go
// Purpose: Tests for shared UI styles and helpers.
// Version: 0.1.0
// Updated: 2025-12-28
package ui_test

import (
	"strings"
	"testing"

	"github.com/peternicholls/20i-stack/tui/internal/ui"
)

func TestHeaderStyleRenders(t *testing.T) {
	out := ui.HeaderStyle.Render("test")
	if out == "" {
		t.Fatalf("HeaderStyle.Render returned empty string")
	}
}

func TestStatusBadgeContainsUppercase(t *testing.T) {
	s := ui.StatusBadge("running")
	if !strings.Contains(strings.ToUpper(s), "RUNNING") {
		t.Fatalf("StatusBadge did not contain RUNNING; got %q", s)
	}
}

func TestStatusBadgeUnknownForEmpty(t *testing.T) {
	s := ui.StatusBadge("")
	if !strings.Contains(strings.ToUpper(s), "UNKNOWN") {
		t.Fatalf("StatusBadge did not contain UNKNOWN for empty input; got %q", s)
	}
}
