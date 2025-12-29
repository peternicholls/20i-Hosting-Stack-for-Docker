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

// RenderConfirmationModal renders a double-confirmation modal for destructive operations.
// It supports two stages:
//   - Stage 1: User must type "yes" to continue
//   - Stage 2: User must type "destroy" to confirm
//
// Parameters:
//   - stage: Current confirmation stage (1 or 2)
//   - currentInput: The current user input for this stage
//   - width: Terminal width for centering
//   - height: Terminal height for centering
//
// Returns:
//   - Rendered modal as a string
func RenderConfirmationModal(stage int, currentInput string, width, height int) string {
	var title, prompt, progress string
	
	switch stage {
	case 1:
		title = "⚠️  Destroy stack?"
		prompt = "Type 'yes' to continue"
		progress = "Step 1/2"
	case 2:
		title = "🔴 Are you SURE?"
		prompt = "Type 'destroy' to confirm"
		progress = "Step 2/2"
	default:
		return ""
	}
	
	// Modal styling
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("9")). // Red
		Align(lipgloss.Center).
		Padding(0, 1)
	
	promptStyle := lipgloss.NewStyle().
		Foreground(ColorMuted).
		Align(lipgloss.Center).
		Padding(0, 1)
	
	progressStyle := lipgloss.NewStyle().
		Foreground(ColorMuted).
		Align(lipgloss.Center).
		Italic(true).
		Padding(0, 1)
	
	inputStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("11")). // Yellow
		Align(lipgloss.Center).
		Padding(0, 1)
	
	hintStyle := lipgloss.NewStyle().
		Foreground(ColorMuted).
		Align(lipgloss.Center).
		Italic(true).
		Padding(0, 1)
	
	modalContentStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("9")). // Red
		Padding(1, 2).
		Width(50)
	
	// Build modal content
	var lines []string
	lines = append(lines, titleStyle.Render(title))
	lines = append(lines, "") // Empty line
	lines = append(lines, promptStyle.Render(prompt))
	lines = append(lines, "") // Empty line
	
	// Show current input with cursor
	inputDisplay := currentInput + "█"
	lines = append(lines, inputStyle.Render(inputDisplay))
	lines = append(lines, "") // Empty line
	
	// Show progress and help
	lines = append(lines, progressStyle.Render(progress))
	lines = append(lines, hintStyle.Render("Esc: cancel"))
	
	content := lipgloss.JoinVertical(lipgloss.Center, lines...)
	modal := modalContentStyle.Render(content)
	
	// Center the modal on screen
	centered := lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		modal,
	)
	
	return centered
}
