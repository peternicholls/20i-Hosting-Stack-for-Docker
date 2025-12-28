// Project: 20i Stack Manager TUI
// File: dashboard.go
// Purpose: Dashboard model for container lifecycle view
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

// Model represents the dashboard view state for Phase 3.
type Model struct {
	containers    []docker.Container
	selectedIndex int
	projectName   string
	dockerClient  *docker.Client
	width         int
	height        int
	lastError     error
	lastStatusMsg string
}

// NewModel creates a dashboard model with required dependencies.
func NewModel(client *docker.Client, projectName string) Model {
	return Model{
		dockerClient: client,
		projectName:  projectName,
		width:        80,
		height:       24,
	}
}

// Init loads the initial container list.
func (m Model) Init() tea.Cmd {
	return loadContainersCmd(m.dockerClient, m.projectName)
}

// Update handles incoming messages for the dashboard.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		// Navigation keys (T051)
		switch msg.String() {
		case "up", "k":
			if m.selectedIndex > 0 {
				m.selectedIndex--
			}
			return m, nil
		case "down", "j":
			if m.selectedIndex < len(m.containers)-1 {
				m.selectedIndex++
			}
			return m, nil

		// Container action keys (T053-T054)
		case "s":
			// Toggle start/stop for selected container
			if m.selectedIndex >= 0 && m.selectedIndex < len(m.containers) {
				container := m.containers[m.selectedIndex]
				action := "start"
				if container.Status == docker.StatusRunning {
					action = "stop"
				}
				return m, containerActionCmd(m.dockerClient, container.ID, action, container.Service)
			}
			return m, nil

		case "r":
			// Restart selected container
			if m.selectedIndex >= 0 && m.selectedIndex < len(m.containers) {
				container := m.containers[m.selectedIndex]
				return m, containerActionCmd(m.dockerClient, container.ID, "restart", container.Service)
			}
			return m, nil
		}

	case containerActionResultMsg:
		// Handle action results (T061)
		m.lastStatusMsg = msg.message
		if msg.success {
			// Refresh container list after successful action
			return m, loadContainersCmd(m.dockerClient, m.projectName)
		}
		return m, nil

	case containerListMsg:
		if msg.err != nil {
			m.lastError = msg.err
			return m, nil
		}

		m.containers = msg.containers
		if m.selectedIndex >= len(m.containers) {
			m.selectedIndex = clampIndex(len(m.containers) - 1)
		}
		m.lastError = nil
	}

	return m, nil
}

// View renders the dashboard view with 2-panel layout (Phase 3: service list + status messages).
// Phase 5 will expand to 3-panel layout (service list + detail panel + status messages).
func (m Model) View() string {
	// Calculate panel dimensions
	// Service list: 30% width, full height minus header/footer
	// Status panel: 70% width, full height minus header/footer

	serviceListWidth := max(20, m.width*30/100)
	statusPanelWidth := m.width - serviceListWidth
	panelHeight := m.height - 4 // Reserve space for header (1) and footer (2)

	// Render service list (left panel, 30%)
	serviceList := renderServiceList(m.containers, m.selectedIndex, serviceListWidth, panelHeight)

	// Render status messages panel (right panel, 70%)
	statusPanel := m.renderStatusPanel(statusPanelWidth, panelHeight)

	// Join horizontally
	mainContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		serviceList,
		statusPanel,
	)

	// Render footer
	footer := m.renderFooter()

	// Join vertically
	return lipgloss.JoinVertical(
		lipgloss.Left,
		mainContent,
		footer,
	)
}

// renderStatusPanel renders the status messages panel (Phase 3: simple error/success display).
func (m Model) renderStatusPanel(width, height int) string {
	var content string

	if m.lastError != nil {
		// Show error message
		content = "❌ Error: " + m.lastError.Error()
	} else if len(m.containers) == 0 {
		// No containers loaded yet
		content = "ℹ Loading containers..."
	} else {
		// Show helpful hint
		content = fmt.Sprintf("Selected: %s\n\nPress 's' to start/stop\nPress 'r' to restart",
			m.getSelectedContainerName())
	}

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.ColorBorder).
		Width(width-4).
		Height(height-2).
		Padding(1, 2).
		Render(content)
}

// renderFooter renders the keyboard shortcuts footer.
func (m Model) renderFooter() string {
	shortcuts := "↑↓/k/j:navigate  s:start/stop  r:restart  S:stop-all  R:restart-all  D:destroy  q:quit"
	return ui.FooterStyle.Width(m.width).Render(shortcuts)
}

// getSelectedContainerName returns the name of the currently selected container.
func (m Model) getSelectedContainerName() string {
	if m.selectedIndex >= 0 && m.selectedIndex < len(m.containers) {
		return m.containers[m.selectedIndex].Service
	}
	return "(none)"
}

func loadContainersCmd(client *docker.Client, projectName string) tea.Cmd {
	return func() tea.Msg {
		containers, err := client.ListContainers(projectName)
		return containerListMsg{containers: containers, err: err}
	}
}

func clampIndex(index int) int {
	if index < 0 {
		return 0
	}
	return index
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type containerListMsg struct {
	containers []docker.Container
	err        error
}

// containerActionResultMsg is returned after a container action command completes.
type containerActionResultMsg struct {
	success bool
	message string
	err     error
}

// containerActionCmd creates a tea.Cmd that performs a container action asynchronously.
// This is the generic command function per ADR-004 (T058).
func containerActionCmd(client *docker.Client, containerID, action, serviceName string) tea.Cmd {
	return func() tea.Msg {
		var err error

		switch action {
		case "start":
			err = client.StartContainer(containerID)
		case "stop":
			err = client.StopContainer(containerID, 10)
		case "restart":
			err = client.RestartContainer(containerID, 10)
		default:
			return containerActionResultMsg{
				success: false,
				message: fmt.Sprintf("❌ Invalid action: %s", action),
				err:     fmt.Errorf("invalid action: %s", action),
			}
		}

		if err != nil {
			return containerActionResultMsg{
				success: false,
				message: formatDockerError(err, action, serviceName),
				err:     err,
			}
		}

		return containerActionResultMsg{
			success: true,
			message: fmt.Sprintf("✅ Container '%s' %sed successfully", serviceName, action),
		}
	}
}

// formatDockerError formats Docker errors in a user-friendly way (T065).
func formatDockerError(err error, action, containerName string) string {
	if err == nil {
		return ""
	}

	errStr := err.Error()

	// Check for known Docker error patterns
	if strings.Contains(errStr, "port is already allocated") || strings.Contains(errStr, "address already in use") {
		return fmt.Sprintf("❌ Port conflict: Cannot %s '%s'. Port already in use.", action, containerName)
	}

	if strings.Contains(errStr, "timeout") || strings.Contains(errStr, "context deadline exceeded") {
		return fmt.Sprintf("❌ Timeout: Container '%s' took too long to %s. Try again.", containerName, action)
	}

	if strings.Contains(errStr, "No such container") || strings.Contains(errStr, "not found") {
		return fmt.Sprintf("❌ Container '%s' not found. It may have been removed. Press 'r' to refresh.", containerName)
	}

	if strings.Contains(errStr, "permission denied") {
		return "❌ Permission denied. Add your user to the docker group."
	}

	// Generic fallback
	return fmt.Sprintf("❌ Failed to %s '%s': %s", action, containerName, err)
}
