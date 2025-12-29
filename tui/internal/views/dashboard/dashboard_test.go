package dashboard

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/peternicholls/20i-stack/tui/internal/docker"
	"github.com/peternicholls/20i-stack/tui/internal/project"
	"github.com/peternicholls/20i-stack/tui/internal/ui"
)

func TestDashboardModel_Init(t *testing.T) {
	model := NewModel(nil, "test-project")
	cmd := model.Init()
	if cmd == nil {
		t.Error("Init should return a command")
	}
}

func TestDashboardModel_Update_WindowSize(t *testing.T) {
	model := NewModel(nil, "test-project")
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	if updatedModel.width != 120 {
		t.Errorf("Expected width 120, got %d", updatedModel.width)
	}
	if updatedModel.height != 40 {
		t.Errorf("Expected height 40, got %d", updatedModel.height)
	}
}

func TestDashboardModel_View_ThreePanelLayout(t *testing.T) {
	model := NewModel(nil, "test-project")
	model.width = 100
	model.height = 30
	model.rightPanelState = "preflight"

	view := model.View()
	if view == "" {
		t.Error("View should return non-empty string")
	}

	// Check that view contains expected panel elements
	if !containsAny(view, "Detecting...", "Pre-flight") {
		t.Error("View should contain preflight panel content")
	}
}

func TestDashboardModel_View_WithProject(t *testing.T) {
	model := NewModel(nil, "test-project")
	model.width = 100
	model.height = 30
	model.project = &project.Project{
		Name:          "test-project",
		Path:          "/home/user/test-project",
		HasPublicHTML: true,
	}
	model.rightPanelState = "preflight"

	view := model.View()

	// Should show project name
	if !containsAny(view, "test-project") {
		t.Error("View should contain project name")
	}

	// Should show path (may be truncated, so check for partial match)
	if !containsAny(view, "/home", "test-project", "user") {
		t.Error("View should contain project path or parts of it")
	}

	// Should show public_html status
	if !containsAny(view, "public_html", "Present") {
		t.Error("View should show public_html status")
	}
}

func TestDashboardModel_View_StatusPanel(t *testing.T) {
	model := NewModel(nil, "test-project")
	model.width = 100
	model.height = 30
	model.rightPanelState = "status"
	model.containers = []docker.Container{
		{
			ID:      "abc123",
			Service: "nginx",
			Status:  docker.StatusRunning,
			Image:   "nginx:latest",
		},
		{
			ID:      "def456",
			Service: "mariadb",
			Status:  docker.StatusRunning,
			Image:   "mariadb:10.11",
		},
	}

	view := model.View()

	// Should show service names
	if !containsAny(view, "nginx", "mariadb") {
		t.Error("View should contain service names")
	}

	// Should show status indicators
	if !containsAny(view, "Running") {
		t.Error("View should show running status")
	}
}

func TestDashboardModel_Update_ProjectDetection(t *testing.T) {
	model := NewModel(nil, "test-project")
	
	msg := projectDetectedMsg{
		project: project.Project{
			Name:          "detected-project",
			Path:          "/home/user/project",
			HasPublicHTML: true,
		},
	}

	updatedModel, _ := model.Update(msg)

	if updatedModel.project == nil {
		t.Fatal("Project should be set after detection")
	}

	if updatedModel.project.Name != "detected-project" {
		t.Errorf("Expected project name 'detected-project', got '%s'", updatedModel.project.Name)
	}
}

func TestDashboardModel_Update_ContainerList(t *testing.T) {
	model := NewModel(nil, "test-project")
	model.rightPanelState = "preflight"

	msg := containerListMsg{
		containers: []docker.Container{
			{ID: "123", Service: "nginx", Status: docker.StatusRunning},
		},
		err: nil,
	}

	updatedModel, _ := model.Update(msg)

	if len(updatedModel.containers) != 1 {
		t.Errorf("Expected 1 container, got %d", len(updatedModel.containers))
	}

	// Should auto-switch to status panel when containers are present
	if updatedModel.rightPanelState != "status" {
		t.Errorf("Expected rightPanelState 'status', got '%s'", updatedModel.rightPanelState)
	}
}

func TestDashboardModel_MouseClick_URLRegion(t *testing.T) {
	mockOpener := &ui.MockURLOpener{}
	
	model := NewModel(nil, "test-project")
	model.urlOpener = mockOpener
	model.rightPanelState = "status"
	model.containers = []docker.Container{
		{Service: "nginx", Status: docker.StatusRunning},
	}

	// Simulate rendering to populate table state
	model.View()

	// Simulate mouse click (this is a basic test - actual click detection is position-based)
	msg := tea.MouseMsg{
		Type:   tea.MouseLeft,
		Button: tea.MouseButtonLeft,
		X:      70,
		Y:      2,
	}

	model.Update(msg)

	// Note: In real usage, URL would be opened if click was on a URL region
	// This test verifies the handler doesn't crash
}

func TestRightPanelState_Transitions(t *testing.T) {
	tests := []struct {
		name           string
		initialState   string
		containers     []docker.Container
		expectedState  string
	}{
		{
			name:          "preflight to status with containers",
			initialState:  "preflight",
			containers:    []docker.Container{{Service: "nginx"}},
			expectedState: "status",
		},
		{
			name:          "status to preflight with no containers",
			initialState:  "status",
			containers:    []docker.Container{},
			expectedState: "preflight",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewModel(nil, "test-project")
			model.rightPanelState = tt.initialState

			msg := containerListMsg{
				containers: tt.containers,
				err:        nil,
			}

			updatedModel, _ := model.Update(msg)

			if updatedModel.rightPanelState != tt.expectedState {
				t.Errorf("Expected state '%s', got '%s'", tt.expectedState, updatedModel.rightPanelState)
			}
		})
	}
}

// Helper function to check if a string contains any of the given substrings
func containsAny(s string, substrings ...string) bool {
	for _, substr := range substrings {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}


// TestDashboardModel_OutputStreaming tests the output streaming functionality
func TestDashboardModel_OutputStreaming(t *testing.T) {
t.Run("stackOutputMsg appends to composeOutput", func(t *testing.T) {
model := NewModel(nil, "test-project")
model.isStreaming = true

msg := stackOutputMsg{Line: "Test output line", IsError: false}
updatedModel, _ := model.Update(msg)

if len(updatedModel.composeOutput) != 1 {
t.Errorf("Expected 1 line in composeOutput, got %d", len(updatedModel.composeOutput))
}

if updatedModel.composeOutput[0] != "Test output line" {
t.Errorf("Expected 'Test output line', got '%s'", updatedModel.composeOutput[0])
}
})

t.Run("completion message transitions to status", func(t *testing.T) {
model := NewModel(nil, "test-project")
model.isStreaming = true
model.rightPanelState = "output"

msg := stackOutputMsg{Line: "[Complete]", IsError: false}
updatedModel, _ := model.Update(msg)

if !updatedModel.streamingComplete {
t.Error("Expected streamingComplete to be true")
}

if updatedModel.isStreaming {
t.Error("Expected isStreaming to be false")
}

// Verify last output contains completion marker
if len(updatedModel.composeOutput) == 0 {
t.Error("Expected composeOutput to contain completion message")
} else if updatedModel.composeOutput[len(updatedModel.composeOutput)-1] != "[Complete]" {
t.Error("Expected last output to be [Complete]")
}
})

t.Run("multiple output messages accumulate", func(t *testing.T) {
model := NewModel(nil, "test-project")
model.isStreaming = true

lines := []string{"Line 1", "Line 2", "Line 3"}
for _, line := range lines {
msg := stackOutputMsg{Line: line, IsError: false}
model, _ = model.Update(msg)
}

if len(model.composeOutput) != 3 {
t.Errorf("Expected 3 lines in composeOutput, got %d", len(model.composeOutput))
}

for i, expected := range lines {
if model.composeOutput[i] != expected {
t.Errorf("Line %d: expected '%s', got '%s'", i, expected, model.composeOutput[i])
}
}
})

t.Run("stackStatusRefreshMsg switches to status panel", func(t *testing.T) {
model := NewModel(nil, "test-project")
model.rightPanelState = "output"

msg := stackStatusRefreshMsg{}
updatedModel, _ := model.Update(msg)

if updatedModel.rightPanelState != "status" {
t.Errorf("Expected rightPanelState 'status', got '%s'", updatedModel.rightPanelState)
}
})

t.Run("complete streaming flow transitions to status", func(t *testing.T) {
// This test demonstrates the requirement:
// "Dashboard transitions to status state after completion message"
model := NewModel(nil, "test-project")
model.rightPanelState = "output"
model.isStreaming = true

// Simulate streaming several lines
lines := []string{
"Container creating",
"Container created",
"Container starting",
"[Complete]",
}

for i, line := range lines {
msg := stackOutputMsg{Line: line, IsError: false}
model, _ = model.Update(msg)

// Before completion, should still be streaming
if i < len(lines)-1 {
if !model.isStreaming {
t.Error("Should still be streaming before completion")
}
}
}

// After [Complete] message
if model.isStreaming {
t.Error("Expected isStreaming to be false after completion")
}
if !model.streamingComplete {
t.Error("Expected streamingComplete to be true after completion")
}

// Verify all lines were captured in order
if len(model.composeOutput) != 4 {
t.Errorf("Expected 4 lines in composeOutput, got %d", len(model.composeOutput))
}

// Then handle the status refresh message that follows
refreshMsg := stackStatusRefreshMsg{}
model, _ = model.Update(refreshMsg)

// Verify transition to status panel
if model.rightPanelState != "status" {
t.Errorf("Expected rightPanelState 'status', got '%s'", model.rightPanelState)
}
})
}

// TestDashboardModel_ComposeStreamStarted tests the composeStreamStartedMsg handler
func TestDashboardModel_ComposeStreamStarted(t *testing.T) {
model := NewModel(nil, "test-project")

// Create a mock channel
mockChan := make(chan string, 5)
mockChan <- "Test line 1"
mockChan <- "Test line 2"
close(mockChan)

msg := composeStreamStartedMsg{channel: mockChan}
updatedModel, cmd := model.Update(msg)

if !updatedModel.isStreaming {
t.Error("Expected isStreaming to be true")
}

if updatedModel.outputChannel == nil {
t.Error("Expected outputChannel to be set")
}

if len(updatedModel.composeOutput) != 0 {
t.Error("Expected composeOutput to be empty (cleared)")
}

if cmd == nil {
t.Error("Expected a command to be returned")
}
}

// TestDashboardModel_KeyHandlers tests the keyboard shortcuts for stack operations
func TestDashboardModel_KeyHandlers(t *testing.T) {
	t.Run("S key starts stack only when public_html exists", func(t *testing.T) {
		model := NewModel(nil, "test-project")
		model.project = &project.Project{
			Name:          "test",
			Path:          "/test",
			HasPublicHTML: true,
		}
		model.stackFile = "/test/docker-compose.yml"
		model.rightPanelState = "preflight"
		
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
		updatedModel, cmd := model.Update(msg)
		
		if updatedModel.rightPanelState != "output" {
			t.Errorf("Expected rightPanelState 'output', got '%s'", updatedModel.rightPanelState)
		}
		
		if cmd == nil {
			t.Error("Expected command to be returned")
		}
		
		if updatedModel.lastStatusMsg != "Starting stack..." {
			t.Errorf("Expected status message 'Starting stack...', got '%s'", updatedModel.lastStatusMsg)
		}
	})
	
	t.Run("S key does nothing when public_html missing", func(t *testing.T) {
		model := NewModel(nil, "test-project")
		model.project = &project.Project{
			Name:          "test",
			Path:          "/test",
			HasPublicHTML: false,
		}
		model.rightPanelState = "preflight"
		
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
		updatedModel, cmd := model.Update(msg)
		
		if updatedModel.rightPanelState != "preflight" {
			t.Errorf("Expected rightPanelState to remain 'preflight', got '%s'", updatedModel.rightPanelState)
		}
		
		if cmd != nil {
			t.Error("Expected no command when public_html missing")
		}
	})
	
	t.Run("S key does nothing when stack already running", func(t *testing.T) {
		model := NewModel(nil, "test-project")
		model.project = &project.Project{
			Name:          "test",
			Path:          "/test",
			HasPublicHTML: true,
		}
		model.containers = []docker.Container{
			{Service: "nginx", Status: docker.StatusRunning},
		}
		model.rightPanelState = "status"
		
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
		updatedModel, cmd := model.Update(msg)
		
		if updatedModel.lastStatusMsg != "Stack is already running" {
			t.Errorf("Expected status message 'Stack is already running', got '%s'", updatedModel.lastStatusMsg)
		}
		
		if cmd != nil {
			t.Error("Expected no command when stack already running")
		}
	})
	
	t.Run("T key installs template when public_html missing", func(t *testing.T) {
		model := NewModel(nil, "test-project")
		model.project = &project.Project{
			Name:          "test",
			Path:          "/test",
			HasPublicHTML: false,
		}
		model.rightPanelState = "preflight"
		
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}}
		updatedModel, cmd := model.Update(msg)
		
		if updatedModel.lastStatusMsg != "Installing template..." {
			t.Errorf("Expected status message 'Installing template...', got '%s'", updatedModel.lastStatusMsg)
		}
		
		if cmd == nil {
			t.Error("Expected command to be returned for template installation")
		}
	})
	
	t.Run("T key stops stack when running", func(t *testing.T) {
		model := NewModel(nil, "test-project")
		model.project = &project.Project{
			Name:          "test",
			Path:          "/test",
			HasPublicHTML: true,
		}
		model.stackFile = "/test/docker-compose.yml"
		model.containers = []docker.Container{
			{Service: "nginx", Status: docker.StatusRunning},
		}
		model.rightPanelState = "status"
		
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}}
		updatedModel, cmd := model.Update(msg)
		
		if updatedModel.rightPanelState != "output" {
			t.Errorf("Expected rightPanelState 'output', got '%s'", updatedModel.rightPanelState)
		}
		
		if cmd == nil {
			t.Error("Expected command to be returned for stack stop")
		}
		
		if updatedModel.lastStatusMsg != "Stopping stack..." {
			t.Errorf("Expected status message 'Stopping stack...', got '%s'", updatedModel.lastStatusMsg)
		}
	})
	
	t.Run("R key triggers restart when stack running", func(t *testing.T) {
		model := NewModel(nil, "test-project")
		model.project = &project.Project{
			Name:          "test",
			Path:          "/test",
			HasPublicHTML: true,
		}
		model.stackFile = "/test/docker-compose.yml"
		model.containers = []docker.Container{
			{Service: "nginx", Status: docker.StatusRunning},
		}
		model.rightPanelState = "status"
		
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}
		updatedModel, cmd := model.Update(msg)
		
		if updatedModel.rightPanelState != "output" {
			t.Errorf("Expected rightPanelState 'output', got '%s'", updatedModel.rightPanelState)
		}
		
		if cmd == nil {
			t.Error("Expected command to be returned for restart")
		}
		
		if updatedModel.lastStatusMsg != "Restarting stack..." {
			t.Errorf("Expected status message 'Restarting stack...', got '%s'", updatedModel.lastStatusMsg)
		}
	})
	
	t.Run("R key refreshes containers when stack not running", func(t *testing.T) {
		model := NewModel(nil, "test-project")
		model.project = &project.Project{
			Name:          "test",
			Path:          "/test",
			HasPublicHTML: true,
		}
		model.containers = []docker.Container{}
		model.rightPanelState = "preflight"
		
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}
		updatedModel, cmd := model.Update(msg)
		
		// Should trigger container list refresh
		if cmd == nil {
			t.Error("Expected command to be returned for refresh")
		}
		
		// State should remain preflight
		if updatedModel.rightPanelState != "preflight" {
			t.Errorf("Expected rightPanelState 'preflight', got '%s'", updatedModel.rightPanelState)
		}
	})
	
	t.Run("D key opens first confirmation modal", func(t *testing.T) {
		model := NewModel(nil, "test-project")
		model.project = &project.Project{
			Name:          "test",
			Path:          "/test",
			HasPublicHTML: true,
		}
		model.rightPanelState = "status"
		
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
		updatedModel, _ := model.Update(msg)
		
		if updatedModel.confirmationStage != 1 {
			t.Errorf("Expected confirmationStage 1, got %d", updatedModel.confirmationStage)
		}
		
		if updatedModel.firstInput != "" {
			t.Errorf("Expected firstInput to be empty, got '%s'", updatedModel.firstInput)
		}
		
		if updatedModel.secondInput != "" {
			t.Errorf("Expected secondInput to be empty, got '%s'", updatedModel.secondInput)
		}
	})
	
	t.Run("D key does nothing when public_html missing", func(t *testing.T) {
		model := NewModel(nil, "test-project")
		model.project = &project.Project{
			Name:          "test",
			Path:          "/test",
			HasPublicHTML: false,
		}
		model.rightPanelState = "preflight"
		
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
		updatedModel, _ := model.Update(msg)
		
		if updatedModel.confirmationStage != 0 {
			t.Errorf("Expected confirmationStage 0, got %d", updatedModel.confirmationStage)
		}
	})
}

// TestDashboardModel_ModalFlow tests the double-confirmation destroy flow
func TestDashboardModel_ModalFlow(t *testing.T) {
	t.Run("Modal stage 1 requires exact 'yes' to advance", func(t *testing.T) {
		model := NewModel(nil, "test-project")
		model.project = &project.Project{
			Name:          "test",
			Path:          "/test",
			HasPublicHTML: true,
		}
		model.confirmationStage = 1
		model.firstInput = ""
		
		// Type 'y'
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
		if model.firstInput != "y" {
			t.Errorf("Expected firstInput 'y', got '%s'", model.firstInput)
		}
		
		// Type 'e'
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
		if model.firstInput != "ye" {
			t.Errorf("Expected firstInput 'ye', got '%s'", model.firstInput)
		}
		
		// Type 's'
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
		if model.firstInput != "yes" {
			t.Errorf("Expected firstInput 'yes', got '%s'", model.firstInput)
		}
		
		// Press Enter
		updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
		
		if updatedModel.confirmationStage != 2 {
			t.Errorf("Expected confirmationStage 2, got %d", updatedModel.confirmationStage)
		}
		
		if updatedModel.secondInput != "" {
			t.Errorf("Expected secondInput to be reset, got '%s'", updatedModel.secondInput)
		}
	})
	
	t.Run("Modal stage 1 rejects incorrect input", func(t *testing.T) {
		model := NewModel(nil, "test-project")
		model.project = &project.Project{
			Name:          "test",
			Path:          "/test",
			HasPublicHTML: true,
		}
		model.confirmationStage = 1
		model.firstInput = "yep"
		
		// Press Enter with wrong input
		updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
		
		if updatedModel.confirmationStage != 1 {
			t.Errorf("Expected to stay in confirmationStage 1, got %d", updatedModel.confirmationStage)
		}
	})
	
	t.Run("Modal stage 2 requires exact 'destroy' to trigger", func(t *testing.T) {
		model := NewModel(nil, "test-project")
		model.project = &project.Project{
			Name:          "test",
			Path:          "/test",
			HasPublicHTML: true,
		}
		model.stackFile = "/test/docker-compose.yml"
		model.confirmationStage = 2
		
		// Type 'destroy'
		for _, r := range "destroy" {
			model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		}
		
		if model.secondInput != "destroy" {
			t.Errorf("Expected secondInput 'destroy', got '%s'", model.secondInput)
		}
		
		// Press Enter
		updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
		
		if updatedModel.confirmationStage != 0 {
			t.Errorf("Expected confirmationStage reset to 0, got %d", updatedModel.confirmationStage)
		}
		
		if updatedModel.rightPanelState != "output" {
			t.Errorf("Expected rightPanelState 'output', got '%s'", updatedModel.rightPanelState)
		}
		
		if cmd == nil {
			t.Error("Expected command to be returned for destroy")
		}
	})
	
	t.Run("Esc cancels modal from stage 1", func(t *testing.T) {
		model := NewModel(nil, "test-project")
		model.confirmationStage = 1
		model.firstInput = "y"
		
		updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
		
		if updatedModel.confirmationStage != 0 {
			t.Errorf("Expected confirmationStage 0 after Esc, got %d", updatedModel.confirmationStage)
		}
		
		if updatedModel.firstInput != "" {
			t.Errorf("Expected firstInput cleared, got '%s'", updatedModel.firstInput)
		}
		
		if updatedModel.secondInput != "" {
			t.Errorf("Expected secondInput cleared, got '%s'", updatedModel.secondInput)
		}
	})
	
	t.Run("Esc cancels modal from stage 2", func(t *testing.T) {
		model := NewModel(nil, "test-project")
		model.confirmationStage = 2
		model.firstInput = "yes"
		model.secondInput = "des"
		
		updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
		
		if updatedModel.confirmationStage != 0 {
			t.Errorf("Expected confirmationStage 0 after Esc, got %d", updatedModel.confirmationStage)
		}
		
		if updatedModel.firstInput != "" || updatedModel.secondInput != "" {
			t.Error("Expected all inputs cleared after Esc")
		}
	})
	
	t.Run("Backspace removes characters", func(t *testing.T) {
		model := NewModel(nil, "test-project")
		model.confirmationStage = 1
		model.firstInput = "yes"
		
		updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
		
		if updatedModel.firstInput != "ye" {
			t.Errorf("Expected firstInput 'ye' after backspace, got '%s'", updatedModel.firstInput)
		}
	})
	
	t.Run("Modal consumes S/T/R keys (does not trigger stack ops)", func(t *testing.T) {
		model := NewModel(nil, "test-project")
		model.project = &project.Project{
			Name:          "test",
			Path:          "/test",
			HasPublicHTML: true,
		}
		model.confirmationStage = 1
		model.rightPanelState = "status"
		
		// Try 'S' key while modal is open
		updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
		
		// Should add 's' to input, not trigger stack start
		if updatedModel.firstInput != "s" {
			t.Errorf("Expected 's' added to firstInput, got '%s'", updatedModel.firstInput)
		}
		
		// Panel state should not change
		if updatedModel.rightPanelState != "status" {
			t.Errorf("Expected rightPanelState 'status', got '%s'", updatedModel.rightPanelState)
		}
		
		// No command should be triggered
		if cmd != nil {
			t.Error("Expected no command while modal is consuming keys")
		}
	})
}

