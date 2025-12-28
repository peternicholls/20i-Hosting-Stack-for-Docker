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
	"github.com/peternicholls/20i-stack/tui/internal/project"
	"github.com/peternicholls/20i-stack/tui/internal/stack"
	"github.com/peternicholls/20i-stack/tui/internal/ui"
)

// Model represents the dashboard view state for Phase 3.
type Model struct {
	containers         []docker.Container
	selectedIndex      int
	projectName        string
	dockerClient       *docker.Client
	width              int
	height             int
	lastError          error
	lastStatusMsg      string
	projectPath        string // Absolute path to the project
	stackFile          string // Path to docker-compose.yml
	hasPublicHTML      bool   // Whether public_html exists
	rightPanelState    string // "status" or "output"
	outputLines        []string // Stack operation output lines
	confirmationStage  int    // 0 (none) | 1 (first) | 2 (second)
	firstInput         string // Input for first confirmation
	secondInput        string // Input for second confirmation
}

// NewModel creates a dashboard model with required dependencies.
func NewModel(client *docker.Client, projectName string) Model {
	return Model{
		dockerClient:      client,
		projectName:       projectName,
		width:             80,
		height:            24,
		rightPanelState:   "status",
		confirmationStage: 0,
	}
}

// SetProjectInfo configures project-specific information for stack operations.
func (m Model) SetProjectInfo(projectPath, stackFile string, hasPublicHTML bool) Model {
	m.projectPath = projectPath
	m.stackFile = stackFile
	m.hasPublicHTML = hasPublicHTML
	return m
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
		// Handle modal input first
		if m.confirmationStage > 0 {
			return m.handleModalInput(msg)
		}

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

		// Stack operation keys (T140-T143)
		case "S":
			// T140: Start stack (ComposeUp)
			if !m.hasPublicHTML {
				m.lastStatusMsg = "❌ Cannot start stack: public_html directory not found"
				return m, nil
			}
			m.rightPanelState = "output"
			m.outputLines = []string{"Starting stack..."}
			return m, composeUpCmd(m.stackFile, m.projectPath)

		case "T":
			// T141: Template installation or stop stack
			if !m.hasPublicHTML {
				// Install template
				return m, installTemplateCmd(m.projectPath)
			}
			// Stop stack (ComposeDown)
			m.rightPanelState = "output"
			m.outputLines = []string{"Stopping stack..."}
			return m, composeDownCmd(m.stackFile, m.projectPath)

		case "R":
			// T142: Restart stack
			m.rightPanelState = "output"
			m.outputLines = []string{"Restarting stack..."}
			return m, composeRestartCmd(m.stackFile, m.projectPath)

		case "D":
			// T143: Show first destroy confirmation
			m.confirmationStage = 1
			m.firstInput = ""
			m.secondInput = ""
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

	case stackOperationResultMsg:
		// Handle stack operation results
		if msg.err != nil {
			m.outputLines = append(m.outputLines, fmt.Sprintf("❌ Error: %s", msg.err))
			m.lastStatusMsg = fmt.Sprintf("Stack operation failed: %s", msg.err)
		} else {
			m.outputLines = append(m.outputLines, msg.output)
			m.lastStatusMsg = "✅ Stack operation completed successfully"
		}
		// Refresh container list after stack operation
		return m, loadContainersCmd(m.dockerClient, m.projectName)

	case templateInstallResultMsg:
		// Handle template installation result
		if msg.err != nil {
			m.lastStatusMsg = fmt.Sprintf("❌ Template installation failed: %s", msg.err)
		} else {
			m.lastStatusMsg = "✅ Template installed successfully"
			m.hasPublicHTML = true
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
	// If modal is active, render modal overlay
	if m.confirmationStage > 0 {
		return m.renderWithModal()
	}

	// Calculate panel dimensions
	// Service list: 30% width, full height minus header/footer
	// Status panel: 70% width, full height minus header/footer

	serviceListWidth := max(20, m.width*30/100)
	statusPanelWidth := m.width - serviceListWidth
	panelHeight := m.height - 4 // Reserve space for header (1) and footer (2)

	// Render service list (left panel, 30%)
	serviceList := renderServiceList(m.containers, m.selectedIndex, serviceListWidth, panelHeight)

	// Render status messages panel (right panel, 70%)
	var rightPanel string
	if m.rightPanelState == "output" {
		rightPanel = m.renderOutputPanel(statusPanelWidth, panelHeight)
	} else {
		rightPanel = m.renderStatusPanel(statusPanelWidth, panelHeight)
	}

	// Join horizontally
	mainContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		serviceList,
		rightPanel,
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
	shortcuts := "↑↓/k/j:navigate  s:start/stop  r:restart  S:stack-start  T:template/stop  R:stack-restart  D:destroy  q:quit"
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

// renderOutputPanel renders the stack operation output panel.
func (m Model) renderOutputPanel(width, height int) string {
	content := strings.Join(m.outputLines, "\n")
	
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.ColorBorder).
		Width(width-4).
		Height(height-2).
		Padding(1, 2).
		Render(content)
}

// renderWithModal renders the view with a modal overlay.
func (m Model) renderWithModal() string {
	// First render the base view
	serviceListWidth := max(20, m.width*30/100)
	statusPanelWidth := m.width - serviceListWidth
	panelHeight := m.height - 4

	serviceList := renderServiceList(m.containers, m.selectedIndex, serviceListWidth, panelHeight)
	
	var rightPanel string
	if m.rightPanelState == "output" {
		rightPanel = m.renderOutputPanel(statusPanelWidth, panelHeight)
	} else {
		rightPanel = m.renderStatusPanel(statusPanelWidth, panelHeight)
	}

	background := lipgloss.JoinHorizontal(
		lipgloss.Top,
		serviceList,
		rightPanel,
	)

	footer := m.renderFooter()
	
	fullBackground := lipgloss.JoinVertical(
		lipgloss.Left,
		background,
		footer,
	)

	// Render modal content
	var modalContent string
	if m.confirmationStage == 1 {
		modalContent = fmt.Sprintf("⚠️  Destroy stack? Type 'yes' to continue\n\nStep 1/2\n\n> %s\n\nPress Esc to cancel", m.firstInput)
	} else if m.confirmationStage == 2 {
		modalContent = fmt.Sprintf("🔴 Are you SURE? Type 'destroy' to confirm\n\nStep 2/2\n\n> %s\n\nPress Esc to cancel", m.secondInput)
	}

	modal := ui.RenderModal(modalContent, m.width, m.height)

	// Place modal over background
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		modal,
		lipgloss.WithWhitespaceBackground(lipgloss.Color("0")),
	) + "\n" + fullBackground
}

// handleModalInput handles keyboard input when a modal is active.
func (m Model) handleModalInput(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Cancel and close modal
		m.confirmationStage = 0
		m.firstInput = ""
		m.secondInput = ""
		return m, nil

	case "enter":
		if m.confirmationStage == 1 {
			// Check first confirmation
			if strings.ToLower(strings.TrimSpace(m.firstInput)) == "yes" {
				m.confirmationStage = 2
				m.secondInput = ""
				return m, nil
			}
			// Wrong input, reset
			m.firstInput = ""
			return m, nil
		} else if m.confirmationStage == 2 {
			// Check second confirmation
			if strings.ToLower(strings.TrimSpace(m.secondInput)) == "destroy" {
				// Trigger destroy
				m.confirmationStage = 0
				m.firstInput = ""
				m.secondInput = ""
				m.rightPanelState = "output"
				m.outputLines = []string{"Destroying stack (removing volumes)..."}
				return m, composeDestroyCmd(m.stackFile, m.projectPath)
			}
			// Wrong input, reset
			m.secondInput = ""
			return m, nil
		}

	case "backspace":
		if m.confirmationStage == 1 && len(m.firstInput) > 0 {
			m.firstInput = m.firstInput[:len(m.firstInput)-1]
		} else if m.confirmationStage == 2 && len(m.secondInput) > 0 {
			m.secondInput = m.secondInput[:len(m.secondInput)-1]
		}
		return m, nil

	default:
		// Add character to current input
		if len(msg.String()) == 1 {
			if m.confirmationStage == 1 {
				m.firstInput += msg.String()
			} else if m.confirmationStage == 2 {
				m.secondInput += msg.String()
			}
		}
		return m, nil
	}

	return m, nil
}

// Message types for stack operations
type stackOperationResultMsg struct {
	output string
	err    error
}

type templateInstallResultMsg struct {
	err error
}

// composeUpCmd creates a tea.Cmd that starts the stack.
func composeUpCmd(stackFile, projectPath string) tea.Cmd {
	return func() tea.Msg {
		result := stack.ComposeUp(stackFile, projectPath)
		if result.Error != nil {
			return stackOperationResultMsg{
				output: result.Output,
				err:    result.Error,
			}
		}
		return stackOperationResultMsg{
			output: "✅ Stack started successfully\n" + result.Output,
			err:    nil,
		}
	}
}

// composeDownCmd creates a tea.Cmd that stops the stack.
func composeDownCmd(stackFile, projectPath string) tea.Cmd {
	return func() tea.Msg {
		result := stack.ComposeDown(stackFile, projectPath)
		if result.Error != nil {
			return stackOperationResultMsg{
				output: result.Output,
				err:    result.Error,
			}
		}
		return stackOperationResultMsg{
			output: "✅ Stack stopped successfully\n" + result.Output,
			err:    nil,
		}
	}
}

// composeRestartCmd creates a tea.Cmd that restarts the stack.
func composeRestartCmd(stackFile, projectPath string) tea.Cmd {
	return func() tea.Msg {
		result := stack.ComposeRestart(stackFile, projectPath)
		if result.Error != nil {
			return stackOperationResultMsg{
				output: result.Output,
				err:    result.Error,
			}
		}
		return stackOperationResultMsg{
			output: "✅ Stack restarted successfully\n" + result.Output,
			err:    nil,
		}
	}
}

// composeDestroyCmd creates a tea.Cmd that destroys the stack (including volumes).
func composeDestroyCmd(stackFile, projectPath string) tea.Cmd {
	return func() tea.Msg {
		result := stack.ComposeDestroy(stackFile, projectPath)
		if result.Error != nil {
			return stackOperationResultMsg{
				output: result.Output,
				err:    result.Error,
			}
		}
		return stackOperationResultMsg{
			output: "✅ Stack destroyed successfully\n" + result.Output,
			err:    nil,
		}
	}
}

// installTemplateCmd creates a tea.Cmd that installs the template.
func installTemplateCmd(projectPath string) tea.Cmd {
	return func() tea.Msg {
		err := project.InstallTemplate(projectPath)
		return templateInstallResultMsg{err: err}
	}
}
