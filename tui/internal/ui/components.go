// Project: 20i Stack Manager TUI
// File: components.go
// Purpose: Small UI components and helpers (StatusIcon, etc.)
// Version: 0.1.0
// Updated: 2025-12-28

package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// StatusIcon returns a short, colored icon for the provided container status.
// The function accepts a status string (e.g., "running", "stopped", "restarting", "error").
// It intentionally accepts a string to avoid importing a domain package (ContainerStatus)
// while that type is being implemented elsewhere. Callers can pass `string(status)`
// or the literal string values defined in the data model.
func StatusIcon(status string) string {
	s := strings.ToLower(strings.TrimSpace(status))
	if s == "" {
		return lipgloss.NewStyle().Foreground(ColorMuted).Render("?")
	}

	style := lipgloss.NewStyle().Bold(true)
	var icon string

	switch s {
	case "running":
		icon = "●" // filled circle
		style = style.Foreground(ColorRunning)
	case "stopped":
		icon = "○" // hollow circle
		style = style.Foreground(ColorStopped)
	case "restarting":
		icon = "⚠" // warning
		style = style.Foreground(ColorWarning)
	case "error":
		icon = "✗" // cross
		style = style.Foreground(ColorError)
	default:
		icon = "?"
		style = style.Foreground(ColorMuted)
	}

	return style.Render(icon)
}

// RenderModal renders a modal dialog box centered on screen.
// T144: Double-confirmation modal component.
func RenderModal(content string, screenWidth, screenHeight int) string {
	modalWidth := min(60, screenWidth-10)
	modalHeight := min(12, screenHeight-6)

	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(ColorWarning).
		Padding(1, 2).
		Width(modalWidth).
		Height(modalHeight).
		Align(lipgloss.Center, lipgloss.Center)

	return modalStyle.Render(content)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
