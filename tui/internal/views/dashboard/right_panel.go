// Project: 20i Stack Manager TUI
// File: right_panel.go
// Purpose: Right panel rendering with state-based content switching
// Version: 0.1.0
// Updated: 2025-12-28

package dashboard

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/peternicholls/20i-stack/tui/internal/docker"
	"github.com/peternicholls/20i-stack/tui/internal/project"
	"github.com/peternicholls/20i-stack/tui/internal/ui"
)

// Cached styles for output panel to improve performance during frequent renders
var (
	completeStyle = lipgloss.NewStyle().Foreground(ui.ColorRunning).Bold(true)
	errorStyle    = lipgloss.NewStyle().Foreground(ui.ColorError)
)

// renderRightPanel renders the right panel content based on the current state.
// States:
// - "preflight": Show public_html status and template installation option
// - "output": Show streaming compose output
// - "status": Show stack status table with container information
//
// Parameters:
//   - state: Current right panel state
//   - proj: Detected project information
//   - containers: List of Docker containers
//   - composeOutput: Lines of Docker Compose output
//   - tableState: Status table state (updated when rendering status table)
//   - width: Panel width in characters
//   - height: Panel height in lines
//
// Returns:
//   - Rendered panel content
func renderRightPanel(state string, proj *project.Project, containers []docker.Container, composeOutput []string, tableState *StatusTableState, width, height int) string {
	switch state {
	case "preflight":
		return renderPreflightPanel(proj, width, height)
	case "output":
		return renderOutputPanel(composeOutput, width, height)
	case "status":
		rendered, newTableState := renderStatusTable(containers, width, height)
		*tableState = newTableState
		return rendered
	default:
		return renderPreflightPanel(proj, width, height)
	}
}

// renderPreflightPanel shows pre-flight checks and setup options.
// Displays:
// - Project detection status
// - public_html directory status
// - Docker availability
// - Template installation option (if public_html is missing)
// - Ready to start prompt (if public_html is present)
//
// Parameters:
//   - proj: Detected project information (nil if detection in progress)
//   - width: Panel width in characters
//   - height: Panel height in lines
//
// Returns:
//   - Rendered panel content with border
func renderPreflightPanel(proj *project.Project, width, height int) string {
	var lines []string

	lines = append(lines, lipgloss.NewStyle().Bold(true).Render("Pre-flight Check"))
	lines = append(lines, "")

	if proj == nil {
		lines = append(lines, "Detecting project...")
		content := strings.Join(lines, "\n")
		return lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ui.ColorBorder).
			Width(width - 2).
			Height(height - 2).
			Padding(1).
			Render(content)
	}

	// public_html status
	lines = append(lines, lipgloss.NewStyle().Foreground(ui.ColorMuted).Render("✓ public_html directory:"))
	if proj.HasPublicHTML {
		lines = append(lines, "  "+lipgloss.NewStyle().Foreground(ui.ColorRunning).Render("✅ Present"))
	} else {
		lines = append(lines, "  "+lipgloss.NewStyle().Foreground(ui.ColorWarning).Render("⚠️  Missing"))
		lines = append(lines, "")
		lines = append(lines, lipgloss.NewStyle().Foreground(ui.ColorMuted).Render("You can install a template to create public_html:"))
		lines = append(lines, "")
		lines = append(lines, "  "+lipgloss.NewStyle().Bold(true).Render("Press Enter to install default template"))
	}

	lines = append(lines, "")
	lines = append(lines, "")

	// Docker availability
	lines = append(lines, lipgloss.NewStyle().Foreground(ui.ColorMuted).Render("✓ Docker:"))
	lines = append(lines, "  "+lipgloss.NewStyle().Foreground(ui.ColorRunning).Render("✅ Available"))

	lines = append(lines, "")
	lines = append(lines, "")

	// Next steps
	if proj.HasPublicHTML {
		lines = append(lines, lipgloss.NewStyle().Bold(true).Render("Ready to start!"))
		lines = append(lines, "")
		lines = append(lines, "Press 's' to start the stack")
	} else {
		lines = append(lines, lipgloss.NewStyle().Foreground(ui.ColorWarning).Render("Install a template to continue"))
	}

	content := strings.Join(lines, "\n")

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.ColorBorder).
		Width(width - 2).
		Height(height - 2).
		Padding(1).
		Render(content)
}

// renderOutputPanel shows streaming compose output.
// Displays the most recent lines of Docker Compose command output,
// truncating long lines to fit the panel width.
//
// Parameters:
//   - composeOutput: Lines of output from Docker Compose commands
//   - width: Panel width in characters
//   - height: Panel height in lines
//
// Returns:
//   - Rendered panel content with border
func renderOutputPanel(composeOutput []string, width, height int) string {
	var lines []string

	lines = append(lines, lipgloss.NewStyle().Bold(true).Render("Stack Output"))
	lines = append(lines, strings.Repeat("─", width-6))
	lines = append(lines, "")

	if len(composeOutput) == 0 {
		lines = append(lines, lipgloss.NewStyle().Foreground(ui.ColorMuted).Render("No output yet..."))
	} else {
		// Show last N lines to fit in the panel
		maxLines := height - 8 // Account for border, title, padding
		startIdx := 0
		if len(composeOutput) > maxLines {
			startIdx = len(composeOutput) - maxLines
		}

		for _, line := range composeOutput[startIdx:] {
			// Check styling conditions before truncation
			isComplete := line == "[Complete]"
			isError := strings.HasPrefix(line, "ERROR:")
			
			// Truncate long lines to fit panel width
			maxWidth := width - 6
			if maxWidth < 3 {
				maxWidth = 3 // Minimum width for "..."
			}
			displayLine := line
			if len(line) > maxWidth {
				if maxWidth <= 3 {
					displayLine = "..."
				} else {
					displayLine = line[:maxWidth-3] + "..."
				}
			}
			
			// Apply styling based on original line content
			if isComplete {
				displayLine = completeStyle.Render(displayLine)
			} else if isError {
				displayLine = errorStyle.Render(displayLine)
			}
			
			lines = append(lines, displayLine)
		}
	}

	content := strings.Join(lines, "\n")

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.ColorBorder).
		Width(width - 2).
		Height(height - 2).
		Padding(1).
		Render(content)
}
