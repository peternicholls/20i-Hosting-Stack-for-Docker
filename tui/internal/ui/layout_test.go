// Project: 20i Stack Manager TUI
// File: layout_test.go
// Purpose: Unit tests for layout sizing helpers
// Version: 0.1.0
// Updated: 2025-12-28

package ui

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestCalculatePanelWidths(t *testing.T) {
	tests := []struct {
		name         string
		totalWidth   int
		leftContent  string
		rightContent string
		wantLeft     int
		wantRight    int
	}{
		{
			name:         "non-positive width",
			totalWidth:   0,
			leftContent:  "left",
			rightContent: "right",
			wantLeft:     0,
			wantRight:    0,
		},
		{
			name:         "default split with small content",
			totalWidth:   100,
			leftContent:  "left",
			rightContent: "right",
			wantLeft:     28,
			wantRight:    68,
		},
		{
			name:         "left content expands panel",
			totalWidth:   30,
			leftContent:  "left content",
			rightContent: "right",
			wantLeft:     12,
			wantRight:    14,
		},
		{
			name:         "right content expands panel",
			totalWidth:   30,
			leftContent:  "left",
			rightContent: "right content is long",
			wantLeft:     5,
			wantRight:    21,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotLeft, gotRight := calculatePanelWidths(tc.totalWidth, tc.leftContent, tc.rightContent)
			if gotLeft != tc.wantLeft || gotRight != tc.wantRight {
				t.Fatalf("calculatePanelWidths(%d, %q, %q) = (%d, %d); want (%d, %d)",
					tc.totalWidth, tc.leftContent, tc.rightContent, gotLeft, gotRight, tc.wantLeft, tc.wantRight)
			}
		})
	}
}

func TestCalculatePanelWidthsContentMinimums(t *testing.T) {
	leftContent := "left content is very long"
	rightContent := "right content is even longer"

	gotLeft, gotRight := calculatePanelWidths(10, leftContent, rightContent)
	if gotLeft < len(leftContent) || gotRight < len(rightContent) {
		t.Fatalf("calculatePanelWidths enforces content minimums; got (%d, %d)", gotLeft, gotRight)
	}
}

func TestCalculatePanelHeights(t *testing.T) {
	tests := []struct {
		name          string
		totalHeight   int
		headerContent string
		footerContent string
		wantHeader    int
		wantContent   int
		wantFooter    int
	}{
		{
			name:          "non-positive height",
			totalHeight:   0,
			headerContent: "header",
			footerContent: "footer",
			wantHeader:    0,
			wantContent:   0,
			wantFooter:    0,
		},
		{
			name:          "normal layout",
			totalHeight:   10,
			headerContent: "header\nline",
			footerContent: "footer",
			wantHeader:    2,
			wantContent:   7,
			wantFooter:    1,
		},
		{
			name:          "content clamped at zero",
			totalHeight:   2,
			headerContent: "h1\nh2\nh3",
			footerContent: "f1",
			wantHeader:    3,
			wantContent:   0,
			wantFooter:    1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotHeader, gotContent, gotFooter := calculatePanelHeights(tc.totalHeight, tc.headerContent, tc.footerContent)
			if gotHeader != tc.wantHeader || gotContent != tc.wantContent || gotFooter != tc.wantFooter {
				t.Fatalf("calculatePanelHeights(%d, %q, %q) = (%d, %d, %d); want (%d, %d, %d)",
					tc.totalHeight, tc.headerContent, tc.footerContent,
					gotHeader, gotContent, gotFooter,
					tc.wantHeader, tc.wantContent, tc.wantFooter)
			}
		})
	}
}

func TestCalculatePanelWidthsUnicode(t *testing.T) {
	left := "左側"     // Japanese characters
	right := "право" // Cyrillic
	gotLeft, gotRight := calculatePanelWidths(40, left, right)
	if gotLeft < lipgloss.Width(left) || gotRight < lipgloss.Width(right) {
		t.Fatalf("calculatePanelWidths enforces content minima for unicode; got (%d, %d)", gotLeft, gotRight)
	}
}

func TestCalculatePanelWidthsContentTooWide(t *testing.T) {
	left := "LLLLLLLLLLLLLLLLLLLLLLLLLLLLL"
	right := "RRRRRRRRRRRRRRRRRRRRRRRRRRR"
	gotLeft, gotRight := calculatePanelWidths(10, left, right)
	if gotLeft < 0 || gotRight < 0 {
		t.Fatalf("calculatePanelWidths should not return negative widths; got (%d, %d)", gotLeft, gotRight)
	}
	if gotLeft < lipgloss.Width(left) || gotRight < lipgloss.Width(right) {
		t.Fatalf("calculatePanelWidths enforces content minima when too wide; got (%d, %d)", gotLeft, gotRight)
	}
}
