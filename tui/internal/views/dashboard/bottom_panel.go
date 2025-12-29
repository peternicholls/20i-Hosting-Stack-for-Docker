// Project: 20i Stack Manager TUI
// File: bottom_panel.go
// Purpose: Bottom panel rendering for dashboard (commands and status messages)
// Version: 0.1.0
// Updated: 2025-12-28

package dashboard

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/peternicholls/20i-stack/tui/internal/ui"
)

// renderBottomPanel renders the bottom panel showing available commands and status messages.
// The panel is exactly 3 lines tall:
// - Line 1: Available commands based on current state
// - Line 2: Status/error messages
// - Line 3: General help (quit command)
//
// Parameters:
//   - rightPanelState: Current state of the right panel ("preflight", "output", or "status")
//   - statusMsg: Current status message to display
//   - width: Panel width in characters
//
// Returns:
//   - Rendered panel content (3 lines)
func renderBottomPanel(rightPanelState string, statusMsg string, width int) string {
	var lines []string

	// Line 1: State-specific commands
	commands := getAvailableCommands(rightPanelState)
	commandLine := lipgloss.NewStyle().
		Foreground(ui.ColorMuted).
		Width(width).
		Render(commands)
	lines = append(lines, commandLine)

	// Line 2: Status messages
	statusLine := ""
	if statusMsg != "" {
		statusLine = statusMsg
	}
	statusLineStyled := lipgloss.NewStyle().
		Width(width).
		Render(statusLine)
	lines = append(lines, statusLineStyled)

	// Line 3: General help
	helpLine := lipgloss.NewStyle().
		Foreground(ui.ColorMuted).
		Width(width).
		Render("q: quit")
	lines = append(lines, helpLine)

	return strings.Join(lines, "\n")
}

// getAvailableCommands returns the command help text based on the current right panel state.
//
// Parameters:
//   - rightPanelState: Current state ("preflight", "output", or "status")
//
// Returns:
//   - Command help string appropriate for the current state
func getAvailableCommands(rightPanelState string) string {
	switch rightPanelState {
	case "preflight":
		return "s: start stack  t: install template  r: refresh"
	case "output":
		return "Streaming output... (wait for completion)"
	case "status":
		return "t: stop stack  r: restart  d: destroy  Click URL to open"
	default:
		return "r: refresh"
	}
}
