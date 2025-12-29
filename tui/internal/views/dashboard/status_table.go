// Project: 20i Stack Manager TUI
// File: status_table.go
// Purpose: Service status table with CPU bars and clickable URLs
// Version: 0.1.0
// Updated: 2025-12-28

package dashboard

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/peternicholls/20i-stack/tui/internal/docker"
	"github.com/peternicholls/20i-stack/tui/internal/ui"
)

// URLRegion represents a clickable region in the rendered table.
// It stores the URL and its position (row, column range) for click detection.
type URLRegion struct {
	URL      string // The URL to open when clicked
	Row      int    // Row number in the rendered table
	ColStart int    // Starting column of the clickable region
	ColEnd   int    // Ending column of the clickable region
}

// StatusTableState tracks the rendered state for click detection.
// It maintains a list of clickable URL regions and panel dimensions.
type StatusTableState struct {
	URLRegions []URLRegion // Clickable URL regions in the table
	Width      int         // Table width in characters
	Height     int         // Table height in lines
}

// renderStatusTable renders a table showing container status information.
// The table includes columns for Service, Status, Image, URL/Port, and CPU%.
// URLs for web services (nginx, phpmyadmin) are highlighted and clickable.
//
// Parameters:
//   - containers: List of Docker containers to display
//   - width: Table width in characters
//   - height: Table height in lines
//
// Returns:
//   - rendered: The rendered table as a string
//   - state: Updated state with URL regions for click detection
func renderStatusTable(containers []docker.Container, width, height int) (string, StatusTableState) {
	state := StatusTableState{
		URLRegions: []URLRegion{},
		Width:      width,
		Height:     height,
	}

	if len(containers) == 0 {
		return lipgloss.NewStyle().
			Width(width - 4).
			Height(height - 4).
			Padding(2).
			Render("No containers running"), state
	}

	// Build table
	var rows []string

	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(ui.ColorPrimary).
		Padding(0, 1)

	header := headerStyle.Render(fmt.Sprintf("%-20s %-12s %-30s %-25s %-10s",
		"Service",
		"Status",
		"Image",
		"URL/Port",
		"CPU%",
	))
	rows = append(rows, header)

	// Separator
	rows = append(rows, strings.Repeat("─", width-4))

	// Data rows
	for i, container := range containers {
		row := i + 2 // +2 for header and separator

		// Service name
		serviceName := truncateString(container.Service, 20)

		// Status badge
		statusBadge := getStatusBadge(container.Status)

		// Image name (truncate)
		imageName := truncateString(container.Image, 30)

		// URL/Port with click detection
		urlText, urlStart, urlEnd := formatURLField(container, 25)
		if urlStart >= 0 {
			state.URLRegions = append(state.URLRegions, URLRegion{
				URL:      extractURL(container),
				Row:      row,
				ColStart: 64 + urlStart, // Offset for previous columns
				ColEnd:   64 + urlEnd,
			})
		}

		// CPU bar (placeholder for now, will be populated with real stats)
		cpuBar := renderCPUBar(0.0, 10)

		rowText := fmt.Sprintf("%-20s %-12s %-30s %-25s %-10s",
			serviceName,
			statusBadge,
			imageName,
			urlText,
			cpuBar,
		)

		rows = append(rows, rowText)
	}

	content := strings.Join(rows, "\n")

	rendered := lipgloss.NewStyle().
		Padding(1, 2).
		Render(content)

	return rendered, state
}

// renderCPUBar renders a CPU usage bar using block characters.
// Uses ▓ for filled blocks and ░ for empty blocks.
// The bar is colored based on CPU usage: green (< 50%), yellow (50-80%), red (> 80%).
//
// Parameters:
//   - cpuPercent: CPU usage percentage (0-100+)
//   - width: Bar width in characters
//
// Returns:
//   - Rendered and colored CPU bar
func renderCPUBar(cpuPercent float64, width int) string {
	if width <= 0 {
		return ""
	}

	filled := int((cpuPercent / 100.0) * float64(width))
	if filled > width {
		filled = width
	}

	bar := strings.Repeat("▓", filled) + strings.Repeat("░", width-filled)
	
	// Color based on usage
	style := lipgloss.NewStyle()
	if cpuPercent > 80 {
		style = style.Foreground(ui.ColorError)
	} else if cpuPercent > 50 {
		style = style.Foreground(ui.ColorWarning)
	} else {
		style = style.Foreground(ui.ColorRunning)
	}

	return style.Render(bar)
}

// getStatusBadge returns a styled status badge for the container status.
// The badge is colored based on status: green (running), gray (stopped),
// yellow (restarting), red (error).
//
// Parameters:
//   - status: Container status from docker.ContainerStatus
//
// Returns:
//   - Styled status badge string
func getStatusBadge(status docker.ContainerStatus) string {
	var style lipgloss.Style
	var text string

	switch status {
	case docker.StatusRunning:
		style = lipgloss.NewStyle().Foreground(ui.ColorRunning).Bold(true)
		text = "Running"
	case docker.StatusStopped:
		style = lipgloss.NewStyle().Foreground(ui.ColorStopped)
		text = "Stopped"
	case docker.StatusRestarting:
		style = lipgloss.NewStyle().Foreground(ui.ColorWarning).Bold(true)
		text = "Restarting"
	case docker.StatusError:
		style = lipgloss.NewStyle().Foreground(ui.ColorError).Bold(true)
		text = "Error"
	default:
		style = lipgloss.NewStyle().Foreground(ui.ColorMuted)
		text = "Unknown"
	}

	return style.Render(text)
}

// formatURLField formats the URL/Port field and returns the text, start, and end positions
// for click detection. Web URLs (http/https) are highlighted with color and underline.
//
// Parameters:
//   - container: Docker container to extract URL from
//   - maxWidth: Maximum width for the field
//
// Returns:
//   - text: Formatted and styled URL text
//   - start: Starting column for click detection (-1 if no URL)
//   - end: Ending column for click detection (-1 if no URL)
func formatURLField(container docker.Container, maxWidth int) (string, int, int) {
	url := extractURL(container)
	if url == "" {
		return truncateString("N/A", maxWidth), -1, -1
	}

	// Highlight URLs for web services
	isWebURL := strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")

	urlText := truncateString(url, maxWidth)
	// Store visual length before applying any styling (ANSI codes)
	visualLength := len(urlText)

	if isWebURL {
		// Highlight in blue/underline
		styled := lipgloss.NewStyle().
			Foreground(ui.ColorAccent).
			Underline(true).
			Render(urlText)
		return styled, 0, visualLength
	}

	return urlText, 0, visualLength
}

// extractURL extracts the URL or port information from a container.
// For nginx and phpmyadmin, it generates HTTP URLs.
// For mariadb/mysql, it returns the port without HTTP.
// For apache, it returns "internal" (proxied via nginx).
//
// Note: This is a placeholder implementation. Actual port detection
// would come from container port mappings (US2 requirement).
//
// Parameters:
//   - container: Docker container to extract URL from
//
// Returns:
//   - URL string or empty if no URL is applicable
func extractURL(container docker.Container) string {
	serviceName := strings.ToLower(container.Service)

	// These are placeholder implementations - actual port detection
	// would come from container port mappings (US2 requirement)
	switch {
	case strings.Contains(serviceName, "nginx"):
		return "http://localhost:80"
	case strings.Contains(serviceName, "phpmyadmin"):
		return "http://localhost:8081"
	case strings.Contains(serviceName, "mariadb") || strings.Contains(serviceName, "mysql"):
		return "localhost:3306"
	case strings.Contains(serviceName, "apache"):
		return "internal"
	default:
		return ""
	}
}

// truncateString truncates a string to maxLen, adding "..." if needed.
//
// Parameters:
//   - s: String to truncate
//   - maxLen: Maximum length (including "..." if truncated)
//
// Returns:
//   - Truncated string
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen < 3 {
		return "..."
	}
	return s[:maxLen-3] + "..."
}

// handleURLClick processes a mouse click and opens the URL if clicked on a URL region.
// The URL is opened in a goroutine to avoid blocking the UI.
//
// Parameters:
//   - msg: Mouse message from Bubble Tea
//   - state: Status table state containing URL regions
//   - urlOpener: URL opener implementation (for testing abstraction)
//
// Returns:
//   - tea.Cmd (always nil currently)
func handleURLClick(msg tea.MouseMsg, state StatusTableState, urlOpener ui.URLOpener) tea.Cmd {
	// Convert mouse position to row/column
	row := msg.Y
	col := msg.X

	// Check if click is within any URL region
	for _, region := range state.URLRegions {
		if row == region.Row && col >= region.ColStart && col <= region.ColEnd {
			// Open URL in browser
			go urlOpener.OpenURL(region.URL)
			return nil
		}
	}

	return nil
}
