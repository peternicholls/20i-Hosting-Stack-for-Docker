// Project: 20i Stack Manager TUI
// File: dashboard.go
// Purpose: Dashboard model for three-panel layout with state-based rendering
// Version: 0.1.0
// Updated: 2025-12-28

// Package dashboard provides a three-panel TUI dashboard for managing Docker stack lifecycle.
//
// The dashboard consists of:
// - Left panel (25%): Project information and stack status
// - Right panel (75%): Dynamic content based on state (preflight/output/status)
// - Bottom panel (3 lines): Available commands and status messages
//
// The right panel switches between three states:
// - "preflight": Pre-flight checks and template installation
// - "output": Streaming Docker Compose output
// - "status": Live container status table with clickable URLs
package dashboard

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/peternicholls/20i-stack/tui/internal/docker"
	"github.com/peternicholls/20i-stack/tui/internal/project"
	"github.com/peternicholls/20i-stack/tui/internal/stack"
	"github.com/peternicholls/20i-stack/tui/internal/ui"
)

// Status refresh delay after compose completion.
// Allows time for containers to fully start before querying their status.
const statusRefreshDelay = 2 * time.Second

// DashboardModel represents the dashboard view state with three-panel layout.
// It manages project detection, container status, and dynamic right panel rendering.
type DashboardModel struct {
	// Project information
	project *project.Project

	// Container list
	containers []docker.Container

	// Right panel state: "preflight" | "output" | "status"
	rightPanelState string

	// Compose output for streaming display
	composeOutput []string
	
	// Streaming state
	isStreaming       bool
	streamingComplete bool
	outputChannel     <-chan string

	// Status table state for URL click detection
	tableState StatusTableState

	// URL opener (injectable for testing)
	urlOpener ui.URLOpener

	// Legacy fields for compatibility
	selectedIndex int
	projectName   string
	dockerClient  *docker.Client
	stackFile     string // Path to docker-compose.yml
	codeDir       string // Project code directory
	width         int
	height        int
	lastError     error
	lastStatusMsg string
}

// Model is a legacy type alias for backward compatibility.
type Model = DashboardModel

// NewModel creates a new DashboardModel with the specified Docker client and project name.
// The model is initialized with default dimensions (80x24) and the "preflight" panel state.
//
// Parameters:
//   - client: Docker client for container operations (can be nil for testing)
//   - projectName: Name of the Docker Compose project to monitor
//
// Returns:
//   - A new DashboardModel ready for initialization via Init()
func NewModel(client *docker.Client, projectName string) DashboardModel {
	// Detect stack environment
	stackEnv, err := stack.DetectStackEnv()
	stackFile := ""
	if err == nil && stackEnv != nil {
		stackFile = stackEnv.StackFile
	}
	
	return DashboardModel{
		dockerClient:    client,
		projectName:     projectName,
		stackFile:       stackFile,
		codeDir:         "", // Will be detected from project
		width:           80,
		height:          24,
		rightPanelState: "preflight",
		urlOpener:       &ui.DefaultURLOpener{},
	}
}

// Init triggers async project detection and initial container list load.
// This method is called once when the dashboard model is first initialized.
// It returns a batch command that runs both project detection and container listing in parallel.
//
// Returns:
//   - tea.Cmd that performs async initialization tasks
func (m DashboardModel) Init() tea.Cmd {
	return tea.Batch(
		detectProjectCmd(),
		loadContainersCmd(m.dockerClient, m.projectName),
	)
}

// Update handles incoming messages for the dashboard.
// This is the main event handler that processes:
// - Window resize events (tea.WindowSizeMsg)
// - Mouse clicks for URL opening (tea.MouseMsg)
// - Keyboard input for navigation and actions (tea.KeyMsg)
// - Project detection results (projectDetectedMsg)
// - Container list updates (containerListMsg)
// - Container action results (containerActionResultMsg)
//
// The Update method never blocks - all I/O operations are performed asynchronously via tea.Cmd.
//
// Parameters:
//   - msg: The incoming Bubble Tea message to process
//
// Returns:
//   - Updated model state and optional command to execute
func (m DashboardModel) Update(msg tea.Msg) (DashboardModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.MouseMsg:
		// Handle URL clicks in status table
		if msg.Type == tea.MouseLeft && msg.Button == tea.MouseButtonLeft && m.rightPanelState == "status" {
			return m, handleURLClick(msg, m.tableState, m.urlOpener)
		}
		return m, nil

	case tea.KeyMsg:
		// State-specific key handling
		switch msg.String() {
		case "s":
			// Start/stop stack based on state
			if m.rightPanelState == "preflight" && m.project != nil && m.project.HasPublicHTML {
				// Switch to output mode and start stack
				m.rightPanelState = "output"
				m.lastStatusMsg = "Starting stack..."
				
				// Use project path as code directory
				codeDir := m.codeDir
				if codeDir == "" && m.project != nil {
					codeDir = m.project.Path
				}
				
				return m, startComposeUpCmd(m.stackFile, codeDir)
			} else if m.rightPanelState == "status" {
				// Stop stack
				m.rightPanelState = "output"
				m.lastStatusMsg = "Stopping stack..."
				return m, nil // TODO: Add stack stop command
			}
			return m, nil

		case "r":
			// Refresh - reload containers
			return m, loadContainersCmd(m.dockerClient, m.projectName)

		case "enter":
			// Install template in preflight mode
			if m.rightPanelState == "preflight" && m.project != nil && !m.project.HasPublicHTML {
				m.lastStatusMsg = "Installing template..."
				return m, nil // TODO: Add template install command
			}
			return m, nil

		// Legacy navigation keys for compatibility
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
		}

	case projectDetectedMsg:
		m.project = &msg.project
		// Set code directory from project path
		if m.codeDir == "" && m.project != nil {
			m.codeDir = m.project.Path
		}
		// Auto-switch to status if stack is running
		if len(m.containers) > 0 {
			m.rightPanelState = "status"
		}
		return m, nil

	case containerActionResultMsg:
		// Handle action results (legacy compatibility)
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
		
		// Auto-switch panel state based on containers, but preserve "output" state
		// Only transition between "preflight" and "status"
		if m.rightPanelState != "output" {
			if len(m.containers) > 0 && m.rightPanelState == "preflight" {
				m.rightPanelState = "status"
			} else if len(m.containers) == 0 && m.rightPanelState == "status" {
				m.rightPanelState = "preflight"
			}
		}

		if m.selectedIndex >= len(m.containers) {
			m.selectedIndex = clampIndex(len(m.containers) - 1)
		}
		m.lastError = nil
		return m, nil

	case urlOpenedMsg:
		// URL was successfully opened, update status message
		m.lastStatusMsg = "Opened URL: " + msg.url
		return m, nil

	case urlOpenErrorMsg:
		// URL opening failed, show error message
		m.lastStatusMsg = "Failed to open URL: " + msg.err.Error()
		return m, nil
	
	case stackOutputMsg:
		// Handle streaming output from compose operations
		m.composeOutput = append(m.composeOutput, msg.Line)
		
		// Detect critical errors that should stop streaming
		if strings.HasPrefix(msg.Line, "ERROR: Failed to start command") ||
			strings.HasPrefix(msg.Line, "ERROR: Failed to create") {
			m.streamingComplete = true
			m.isStreaming = false
			m.outputChannel = nil
			m.lastStatusMsg = "Compose operation failed"
			return m, nil
		}
		
		// Detect completion
		if msg.Line == "[Complete]" {
			m.streamingComplete = true
			m.isStreaming = false
			m.outputChannel = nil
			// Trigger status refresh and switch to status panel after delay
			// Delay allows containers to fully start before querying status
			return m, tea.Batch(
				loadContainersCmd(m.dockerClient, m.projectName),
				tea.Tick(statusRefreshDelay, func(t time.Time) tea.Msg {
					return stackStatusRefreshMsg{}
				}),
			)
		}
		
		// Continue reading from the channel if streaming
		if m.isStreaming && m.outputChannel != nil {
			return m, waitForNextLineCmd(m.outputChannel)
		}
		
		return m, nil
	
	case stackStatusRefreshMsg:
		// Switch to status panel after streaming completes
		m.rightPanelState = "status"
		m.lastStatusMsg = "Stack started successfully"
		return m, nil
	
	case composeStreamStartedMsg:
		// Store the channel and start reading from it
		m.outputChannel = msg.channel
		m.isStreaming = true
		m.composeOutput = []string{} // Clear previous output
		m.streamingComplete = false
		return m, waitForNextLineCmd(m.outputChannel)
	}

	return m, nil
}

// stackStatusRefreshMsg is sent to trigger a switch to status panel
type stackStatusRefreshMsg struct{}


// View renders the three-panel dashboard layout.
// Layout:
// - Left panel (25% width): Project info, stack status
// - Right panel (75% width): Dynamic content based on rightPanelState
// - Bottom panel (3 lines): Commands and status messages
//
// The view is fully responsive and handles terminal resize events.
// All rendering uses Lipgloss for styling - no raw ANSI codes.
//
// Returns:
//   - String representation of the complete dashboard UI
func (m DashboardModel) View() string {
	// Calculate panel dimensions
	leftWidth := m.width * 25 / 100
	rightWidth := m.width - leftWidth
	bottomHeight := 3
	mainHeight := m.height - bottomHeight

	// Check if stack is running
	stackRunning := len(m.containers) > 0

	// Render left panel
	leftPanel := renderLeftPanel(m.project, stackRunning, leftWidth, mainHeight)

	// Render right panel
	rightPanel := renderRightPanel(
		m.rightPanelState,
		m.project,
		m.containers,
		m.composeOutput,
		&m.tableState,
		rightWidth,
		mainHeight,
	)

	// Join horizontally
	mainContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPanel,
		rightPanel,
	)

	// Render bottom panel
	bottomPanel := renderBottomPanel(m.rightPanelState, m.lastStatusMsg, m.width)

	// Join vertically
	return lipgloss.JoinVertical(
		lipgloss.Left,
		mainContent,
		bottomPanel,
	)
}

func loadContainersCmd(client *docker.Client, projectName string) tea.Cmd {
	return func() tea.Msg {
		// Handle nil client (allowed for testing)
		if client == nil {
			return containerListMsg{containers: []docker.Container{}, err: nil}
		}
		containers, err := client.ListContainers(projectName)
		return containerListMsg{containers: containers, err: err}
	}
}

// waitForNextLineCmd creates a recursive command that continues reading from the channel.
// This allows non-blocking subscription to the output stream.
func waitForNextLineCmd(outputChan <-chan string) tea.Cmd {
	return func() tea.Msg {
		// Guard against a nil channel, which would block forever on receive.
		if outputChan == nil {
			return stackOutputMsg{
				Line:    "ERROR: compose output channel is nil",
				IsError: true,
			}
		}

		// Use a timeout to avoid blocking indefinitely if the channel is never closed or written to.
		select {
		case line, ok := <-outputChan:
			if !ok {
				// Channel closed
				return stackOutputMsg{Line: "[Complete]", IsError: false}
			}
			return stackOutputMsg{
				Line:    line,
				IsError: strings.HasPrefix(line, "ERROR:"),
			}
		case <-time.After(30 * time.Second):
			return stackOutputMsg{
				Line:    "ERROR: timed out waiting for compose output",
				IsError: true,
			}
		}
	}
}

// startComposeUpCmd starts a compose up operation with streaming output.
// It returns the output channel and a command to start reading from it.
func startComposeUpCmd(stackFile, codeDir string) tea.Cmd {
	return func() tea.Msg {
		outputChan, err := stack.ComposeUpStreaming(stackFile, codeDir)
		if err != nil {
			return stackOutputMsg{
				Line:    fmt.Sprintf("ERROR: Failed to start compose: %v", err),
				IsError: true,
			}
		}
		
		// Return a message that includes the channel
		return composeStreamStartedMsg{channel: outputChan}
	}
}

// composeStreamStartedMsg is sent when compose streaming begins.
type composeStreamStartedMsg struct {
	channel <-chan string
}


// detectProjectCmd triggers async project detection.
func detectProjectCmd() tea.Cmd {
	return func() tea.Msg {
		proj, err := project.DetectProject()
		if err != nil {
			return projectDetectedMsg{
				project: project.Project{
					Name:          "unknown",
					Path:          "",
					HasPublicHTML: false,
				},
			}
		}
		return projectDetectedMsg{project: *proj}
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

// projectDetectedMsg is sent when project detection completes.
type projectDetectedMsg struct {
	project project.Project
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

// stackOutputMsg sent when streaming compose command output.
type stackOutputMsg struct {
	// Line of output from compose command
	Line string
	// True if from stderr
	IsError bool
}
