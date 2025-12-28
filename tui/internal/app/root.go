// Project: 20i Stack Manager TUI
// File: root.go
// Purpose: RootModel coordinate view routing, Docker client initialization, and global shortcuts
// Version: 0.1.0
// Updated: 2025-12-28

package app

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/peternicholls/20i-stack/tui/internal/docker"
)

// RootModel coordinates the entire application state and view routing.
type RootModel struct {
	// Docker client initialized once at startup
	dockerClient *docker.Client

	// Current active view state
	activeView string // "dashboard", "help", "projects", "logs"

	// Terminal dimensions
	width  int
	height int

	// Error state for display
	lastError error
}

// NewRootModel creates a new RootModel with initialized Docker client
func NewRootModel(ctx context.Context) (*RootModel, error) {
	cli, err := docker.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	return &RootModel{
		dockerClient: cli,
		activeView:   "dashboard",
		width:        80,
		height:       24,
	}, nil
}

// Init initializes the root model and returns initial commands
func (m *RootModel) Init() tea.Cmd {
	// Could return commands to fetch initial data, but for now just return nil
	return nil
}

// Update handles all incoming messages and routes to appropriate handlers
func (m *RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "?":
			m.activeView = "help"
			m.lastError = nil // Clear error on view change
		case "p":
			m.activeView = "projects"
			m.lastError = nil // Clear error on view change
		case "esc":
			if m.activeView != "dashboard" {
				m.activeView = "dashboard"
				m.lastError = nil // Clear error on view change
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case ErrorMsg:
		m.lastError = msg.Err

	case SuccessMsg:
		// Clear error on success
		m.lastError = nil

	case ProjectSwitchMsg:
		m.activeView = "dashboard"
		m.lastError = nil // Clear error on view change
		// TODO: Notify dashboard of project switch
	}

	return m, nil
}

// View renders the current view based on activeView state
func (m *RootModel) View() string {
	switch m.activeView {
	case "help":
		return m.renderHelpView()
	case "projects":
		return m.renderProjectsView()
	default:
		return m.renderDashboardView()
	}
}

func (m *RootModel) renderDashboardView() string {
	// TODO: Render dashboard view
	if m.lastError != nil {
		return "ERROR: " + m.lastError.Error()
	}
	return "Dashboard View\nPress '?' for help, 'p' for projects, 'q' to quit"
}

func (m *RootModel) renderHelpView() string {
	return `Help - Keyboard Shortcuts

Global:
  q, ctrl+c    Quit
  ?            Show this help
  p            Projects
  esc          Back to dashboard

Dashboard:
  s            Start/stop container
  r            Restart container
  S            Stop all services
  R            Restart all services
  D            Destroy stack
  l            View logs
  Enter        Show container details
  Tab          Next focus

Press 'esc' to return to dashboard`
}

func (m *RootModel) renderProjectsView() string {
	return `Projects View
  
Select a project to switch context
Press 'esc' to return`
}
