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

	// Column widths (used for formatting and click region calculation)
	const (
		colService = 20
		colStatus  = 12
		colImage   = 30
		colURL     = 25
		colCPU     = 10
	)

	// Header (apply styling to entire row after formatting)
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(ui.ColorPrimary)

	headerText := fmt.Sprintf("%-*s %-*s %-*s %-*s %-*s",
		colService, "Service",
		colStatus, "Status",
		colImage, "Image",
		colURL, "URL/Port",
		colCPU, "CPU%",
	)
	rows = append(rows, headerStyle.Render(headerText))

	// Separator
	rows = append(rows, strings.Repeat("─", width-4))

	// Data rows
	for i, container := range containers {
		row := i + 2 // +2 for header and separator

		// Service name (plain text)
		serviceName := truncateString(container.Service, colService)

		// Status badge (plain text, styling applied later)
		statusText := getStatusText(container.Status)

		// Image name (plain text)
		imageName := truncateString(container.Image, colImage)

		// URL/Port (plain text for now, will style entire row)
		url := extractURL(container)
		urlText := truncateString(url, colURL)

		// Track URL region if present
		if url != "" {
			// Calculate column offset: service + status + image + spacing
			urlColStart := colService + colStatus + colImage + 2
			state.URLRegions = append(state.URLRegions, URLRegion{
				URL:      url,
				Row:      row,
				ColStart: urlColStart,
				ColEnd:   urlColStart + len(urlText),
			})
		}

		// CPU bar (plain text)
		cpuBar := renderCPUBarPlain(0.0, colCPU)

		// Format row with plain text
		rowText := fmt.Sprintf("%-*s %-*s %-*s %-*s %-*s",
			colService, serviceName,
			colStatus, statusText,
			colImage, imageName,
			colURL, urlText,
			colCPU, cpuBar,
		)

		// Apply styling to specific cells by wrapping the entire row
		styledRow := styleTableRow(rowText, container.Status, url != "")

		rows = append(rows, styledRow)
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
	bar := renderCPUBarPlain(cpuPercent, width)

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

// renderCPUBarPlain renders a plain text CPU bar without styling.
func renderCPUBarPlain(cpuPercent float64, width int) string {
	if width <= 0 {
		return ""
	}

	filled := int((cpuPercent / 100.0) * float64(width))
	if filled > width {
		filled = width
	}

	return strings.Repeat("▓", filled) + strings.Repeat("░", width-filled)
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
	text := getStatusText(status)

	var style lipgloss.Style
	switch status {
	case docker.StatusRunning:
		style = lipgloss.NewStyle().Foreground(ui.ColorRunning).Bold(true)
	case docker.StatusStopped:
		style = lipgloss.NewStyle().Foreground(ui.ColorStopped)
	case docker.StatusRestarting:
		style = lipgloss.NewStyle().Foreground(ui.ColorWarning).Bold(true)
	case docker.StatusError:
		style = lipgloss.NewStyle().Foreground(ui.ColorError).Bold(true)
	default:
		style = lipgloss.NewStyle().Foreground(ui.ColorMuted)
	}

	return style.Render(text)
}

// getStatusText returns plain text status without styling.
func getStatusText(status docker.ContainerStatus) string {
	switch status {
	case docker.StatusRunning:
		return "Running"
	case docker.StatusStopped:
		return "Stopped"
	case docker.StatusRestarting:
		return "Restarting"
	case docker.StatusError:
		return "Error"
	default:
		return "Unknown"
	}
}

// styleTableRow applies styling to a formatted table row while preserving alignment.
// This is applied after the row is formatted to avoid ANSI code length issues.
func styleTableRow(rowText string, status docker.ContainerStatus, hasURL bool) string {
	// For now, return plain text to avoid alignment issues
	// Individual cell styling can be added back using a more sophisticated approach
	return rowText
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
// Returns a tea.Cmd that opens the URL asynchronously and can report errors.
//
// Parameters:
//   - msg: Mouse message from Bubble Tea
//   - state: Status table state containing URL regions
//   - urlOpener: URL opener implementation (for testing abstraction)
//
// Returns:
//   - tea.Cmd that opens the URL, or nil if no URL was clicked
func handleURLClick(msg tea.MouseMsg, state StatusTableState, urlOpener ui.URLOpener) tea.Cmd {
	// Convert mouse position to row/column
	row := msg.Y
	col := msg.X

	// Check if click is within any URL region
	for _, region := range state.URLRegions {
		if row == region.Row && col >= region.ColStart && col <= region.ColEnd {
			// Return a command that opens the URL
			return func() tea.Msg {
				err := urlOpener.OpenURL(region.URL)
				if err != nil {
					// Return error message that could be handled by the model
					return urlOpenErrorMsg{url: region.URL, err: err}
				}
				return urlOpenedMsg{url: region.URL}
			}
		}
	}

	return nil
}

// urlOpenedMsg is sent when a URL is successfully opened.
type urlOpenedMsg struct {
	url string
}

// urlOpenErrorMsg is sent when opening a URL fails.
type urlOpenErrorMsg struct {
	url string
	err error
}
