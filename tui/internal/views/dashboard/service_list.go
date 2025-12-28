// Project: 20i Stack Manager TUI
// File: service_list.go
// Purpose: Service list rendering for dashboard view
// Version: 0.1.0
// Updated: 2025-12-28

package dashboard

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/peternicholls/20i-stack/tui/internal/docker"
	"github.com/peternicholls/20i-stack/tui/internal/ui"
)

// renderServiceList renders the service list panel with status icons and names.
// This is Phase 3 simple rendering (icon + name only, NO stats).
// Phase 5 will extend this to show CPU/memory usage inline.
func renderServiceList(containers []docker.Container, selectedIndex int, width, height int) string {
	if len(containers) == 0 {
		emptyMsg := ui.InfoStyle.Render("No containers found. Press 'r' to refresh.")
		return lipgloss.NewStyle().
			Width(width).
			Height(height).
			Padding(1).
			Render(emptyMsg)
	}

	var rows []string
	for i, c := range containers {
		icon := ui.StatusIcon(string(c.Status))

		// Select appropriate style
		style := ui.RowStyle
		if i == selectedIndex {
			style = ui.SelectedRowStyle
		}

		// Render: icon + space + service name
		row := fmt.Sprintf("%s %s", icon, c.Service)
		rows = append(rows, style.Render(row))
	}

	content := strings.Join(rows, "\n")

	// Wrap in panel with border
	return ui.PanelStyle.
		Width(width - 4).   // Account for border and padding
		Height(height - 2). // Account for border
		Render(content)
}
