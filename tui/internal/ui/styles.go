// Project: 20i Stack Manager TUI
// File: styles.go
// Purpose: Shared UI styles, color palette, and helpers.
// Version: 0.1.0
// Updated: 2025-12-28

package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Color palette (semantic names)
var (
	// Status Colors
	ColorRunning = lipgloss.Color("#00ff00") // Green - running
	ColorStopped = lipgloss.Color("#808080") // Gray - stopped
	ColorError   = lipgloss.Color("#ff0000") // Red - error
	ColorWarning = lipgloss.Color("#ffff00") // Yellow - warning
	ColorInfo    = lipgloss.Color("#0000ff") // Blue - info

	// UI Colors
	ColorPrimary   = lipgloss.Color("#7D56F4") // Purple - headers, focus
	ColorSecondary = lipgloss.Color("#04B575") // Green - accents
	ColorAccent    = lipgloss.Color("#0000ff") // Bright blue - accents, selection
	ColorBorder    = lipgloss.Color("#585858") // Dark gray - borders
	ColorText      = lipgloss.Color("#e4e4e4") // Light gray - text
	ColorMuted     = lipgloss.Color("#767676") // Medium gray - hints
	ColorSelection = lipgloss.Color("#585858") // Selection background
	ColorHighlight = lipgloss.Color("#ff00ff") // Magenta - search matches
)

// AdaptiveText is an adaptive light/dark text color used across views
var AdaptiveText = lipgloss.AdaptiveColor{
	Light: "#1a1a1a",
	Dark:  "#f5f5f5",
}

// Header style (set Width dynamically in Update/View)
var HeaderStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(ColorPrimary).
	Padding(0, 1)

// Footer / Help style
var FooterStyle = lipgloss.NewStyle().
	Foreground(ColorMuted).
	Padding(0, 1)

// Panel style for boxed content
var PanelStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(ColorBorder).
	Padding(1, 2)

// Row styles for lists
var (
	RowStyle = lipgloss.NewStyle().Padding(0, 2)

	SelectedRowStyle = RowStyle.Copy().
				Background(ColorSelection).
				Bold(true)

	AlternateRowStyle = RowStyle.Copy().
				Foreground(ColorMuted)
)

// Status / message styles
var (
	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorError).
			Bold(true).
			Padding(0, 1)

	WarningStyle = lipgloss.NewStyle().
			Foreground(ColorWarning).
			Bold(true).
			Padding(0, 1)

	InfoStyle = lipgloss.NewStyle().
			Foreground(ColorInfo).
			Padding(0, 1)
)

// Button styles
var (
	ButtonStyle = lipgloss.NewStyle().
			Foreground(ColorText).
			Background(lipgloss.Color("238")).
			Padding(0, 2).
			Margin(0, 1)

	ActiveButtonStyle = ButtonStyle.Copy().
				Foreground(lipgloss.Color("#FFF")).
				Background(ColorPrimary).
				Bold(true)
)

// Accent and highlight styles
var (
	AccentTextStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true)

	HighlightTextStyle = lipgloss.NewStyle().
				Foreground(ColorHighlight).
				Bold(true)
)

// StatusBadge returns a short, styled badge for a status string.
// If `status` is empty it displays "UNKNOWN" with muted style.
func StatusBadge(status string) string {
	s := strings.ToLower(strings.TrimSpace(status))
	if s == "" {
		s = "unknown"
	}

	style := lipgloss.NewStyle().
		Padding(0, 1).
		Bold(true)

	switch s {
	case "running":
		style = style.Foreground(lipgloss.Color("#000")).Background(ColorRunning)
	case "stopped":
		style = style.Foreground(ColorText).Background(ColorStopped)
	case "error":
		style = style.Foreground(lipgloss.Color("#FFF")).Background(ColorError)
	case "warning":
		style = style.Foreground(lipgloss.Color("#000")).Background(ColorWarning)
	default:
		style = style.Foreground(ColorText).Background(ColorMuted)
	}

	return style.Render(strings.ToUpper(s))
}
