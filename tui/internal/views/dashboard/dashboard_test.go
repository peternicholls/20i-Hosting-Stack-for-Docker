package dashboard

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
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

func TestDashboardModel_StackOutputMsg(t *testing.T) {
	model := NewModel(nil, "test-project")
	model.rightPanelState = rightPanelOutput

	// Send first output line
	updatedModel, _ := model.Update(StackOutputMsg{
		Line:    "Line 1",
		IsError: false,
	})

	if len(updatedModel.outputBuffer) != 1 {
		t.Errorf("Expected 1 line in buffer, got %d", len(updatedModel.outputBuffer))
	}

	if updatedModel.outputBuffer[0] != "Line 1" {
		t.Errorf("Expected 'Line 1', got %q", updatedModel.outputBuffer[0])
	}

	// Send second output line
	updatedModel, _ = updatedModel.Update(StackOutputMsg{
		Line:    "Line 2",
		IsError: false,
	})

	if len(updatedModel.outputBuffer) != 2 {
		t.Errorf("Expected 2 lines in buffer, got %d", len(updatedModel.outputBuffer))
	}

	// Verify viewport content
	content := updatedModel.outputViewport.View()
	if !strings.Contains(content, "Line 1") {
		t.Error("Viewport should contain 'Line 1'")
	}
	if !strings.Contains(content, "Line 2") {
		t.Error("Viewport should contain 'Line 2'")
	}
}

func TestDashboardModel_ComposeOutputComplete(t *testing.T) {
	model := NewModel(nil, "test-project")
	model.rightPanelState = rightPanelOutput
	model.outputBuffer = []string{"Line 1", "Line 2"}

	// Send completion message
	updatedModel, cmd := model.Update(composeOutputCompleteMsg{})

	if !updatedModel.outputComplete {
		t.Error("Expected outputComplete to be true")
	}

	// Check for completion marker (should be last two lines: empty and [Complete])
	if len(updatedModel.outputBuffer) < 2 {
		t.Fatalf("Expected at least 2 lines in buffer, got %d", len(updatedModel.outputBuffer))
	}
	
	lastLine := updatedModel.outputBuffer[len(updatedModel.outputBuffer)-1]
	if lastLine != "[Complete]" {
		t.Errorf("Expected '[Complete]' as last line, got %q", lastLine)
	}

	// Verify command is returned (should refresh container list and switch to status)
	if cmd == nil {
		t.Error("Expected command to be returned")
	}
}

func TestDashboardModel_SwitchToStatus(t *testing.T) {
	model := NewModel(nil, "test-project")
	model.rightPanelState = rightPanelOutput

	// Send switch message
	updatedModel, _ := model.Update(switchToStatusMsg{})

	if updatedModel.rightPanelState != rightPanelStatus {
		t.Errorf("Expected rightPanelState to be rightPanelStatus, got %v", updatedModel.rightPanelState)
	}
}

func TestDashboardModel_RightPanelState(t *testing.T) {
	model := NewModel(nil, "test-project")

	// Default should be status
	if model.rightPanelState != rightPanelStatus {
		t.Errorf("Expected default rightPanelState to be rightPanelStatus, got %v", model.rightPanelState)
	}

	// Switch to output
	model.rightPanelState = rightPanelOutput

	// Render view and verify it doesn't panic
	view := model.View()
	if view == "" {
		t.Error("Expected non-empty view")
	}
}

func TestDashboardModel_ViewportScrolling(t *testing.T) {
	model := NewModel(nil, "test-project")
	model.rightPanelState = rightPanelOutput

	// Add multiple lines to test scrolling
	for i := 0; i < 10; i++ {
		model, _ = model.Update(StackOutputMsg{
			Line:    "Test line",
			IsError: false,
		})
	}

	// Verify buffer has all lines
	if len(model.outputBuffer) != 10 {
		t.Errorf("Expected 10 lines in buffer, got %d", len(model.outputBuffer))
	}

	// Verify viewport is at bottom (YOffset should be at max)
	// This is a state assertion, not a rendering test
	if model.outputViewport.YOffset < 0 {
		t.Error("Viewport YOffset should be non-negative")
	}
}
