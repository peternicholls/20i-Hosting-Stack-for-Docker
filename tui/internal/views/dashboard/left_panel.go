// Project: 20i Stack Manager TUI
// File: left_panel.go
// Purpose: Left panel rendering for dashboard (project info and stack status)
// Version: 0.1.0
// Updated: 2025-12-28

package dashboard

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/peternicholls/20i-stack/tui/internal/project"
	"github.com/peternicholls/20i-stack/tui/internal/ui"
)

// renderLeftPanel renders the left panel showing project information and stack status.
// It displays:
// - Project name with status icon (✅ ready / ⚠️ missing public_html / 🔄 detecting)
// - Project path (truncated to fit)
// - Stack status (Not Running / Running / Starting)
// - public_html directory status
//
// Parameters:
//   - proj: Detected project information (nil if detection in progress)
//   - stackRunning: Whether the Docker stack is currently running
//   - width: Panel width in characters
//   - height: Panel height in lines
//
// Returns:
//   - Rendered panel content with border and styling
func renderLeftPanel(proj *project.Project, stackRunning bool, width, height int) string {
	var lines []string

	// Status icon based on project state
	statusIcon := "🔄" // detecting
	if proj != nil {
		if proj.HasPublicHTML {
			statusIcon = "✅" // ready
		} else {
			statusIcon = "⚠️" // warning - missing public_html
		}
	}

	// Project name
	projectName := "Detecting..."
	if proj != nil {
		projectName = proj.Name
	}
	lines = append(lines, lipgloss.NewStyle().Bold(true).Render(statusIcon+" "+projectName))
	lines = append(lines, "")

	// Project path
	if proj != nil {
		pathLabel := lipgloss.NewStyle().Foreground(ui.ColorMuted).Render("Path:")
		lines = append(lines, pathLabel)
		lines = append(lines, truncatePath(proj.Path, width-4))
		lines = append(lines, "")
	}

	// Stack status
	stackStatusLabel := lipgloss.NewStyle().Foreground(ui.ColorMuted).Render("Stack:")
	lines = append(lines, stackStatusLabel)

	stackStatus := "Not Running"
	stackStyle := lipgloss.NewStyle().Foreground(ui.ColorStopped)
	if stackRunning {
		stackStatus = "Running"
		stackStyle = lipgloss.NewStyle().Foreground(ui.ColorRunning)
	}
	lines = append(lines, stackStyle.Render(stackStatus))

	// public_html status if project detected
	if proj != nil {
		lines = append(lines, "")
		htmlLabel := lipgloss.NewStyle().Foreground(ui.ColorMuted).Render("public_html:")
		lines = append(lines, htmlLabel)

		htmlStatus := "Missing"
		htmlStyle := lipgloss.NewStyle().Foreground(ui.ColorWarning)
		if proj.HasPublicHTML {
			htmlStatus = "Present"
			htmlStyle = lipgloss.NewStyle().Foreground(ui.ColorRunning)
		}
		lines = append(lines, htmlStyle.Render(htmlStatus))
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

// truncatePath shortens a path to fit within maxWidth by replacing middle parts with "...".
// If the path already fits, it is returned unchanged.
// If maxWidth is too small (< 10), only "..." is returned.
//
// Parameters:
//   - path: The file path to truncate
//   - maxWidth: Maximum width in characters
//
// Returns:
//   - Truncated path string
func truncatePath(path string, maxWidth int) string {
	if len(path) <= maxWidth {
		return path
	}

	if maxWidth < 10 {
		return "..."
	}

	// Show start and end of path
	prefixLen := (maxWidth - 3) / 2
	suffixLen := maxWidth - 3 - prefixLen

	return path[:prefixLen] + "..." + path[len(path)-suffixLen:]
}
