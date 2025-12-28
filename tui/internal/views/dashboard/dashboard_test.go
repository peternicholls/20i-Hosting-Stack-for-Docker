package dashboard

import (
	tea "github.com/charmbracelet/bubbletea"
	"testing"
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

// T146: Key handler tests
func TestDashboardModel_KeyHandler_S_WithPublicHTML(t *testing.T) {
	model := NewModel(nil, "test-project")
	model = model.SetProjectInfo("/test/path", "/test/stack.yml", true)

	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'S'}})

	if updatedModel.rightPanelState != "output" {
		t.Errorf("Expected rightPanelState to be 'output', got '%s'", updatedModel.rightPanelState)
	}

	if cmd == nil {
		t.Error("Expected composeUpCmd to be returned")
	}
}

func TestDashboardModel_KeyHandler_S_WithoutPublicHTML(t *testing.T) {
	model := NewModel(nil, "test-project")
	model = model.SetProjectInfo("/test/path", "/test/stack.yml", false)

	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'S'}})

	if updatedModel.rightPanelState == "output" {
		t.Error("rightPanelState should not change to 'output' when public_html is missing")
	}

	if cmd != nil {
		t.Error("No command should be returned when public_html is missing")
	}

	if updatedModel.lastStatusMsg == "" {
		t.Error("Expected error message when public_html is missing")
	}
}

func TestDashboardModel_KeyHandler_T_WithoutPublicHTML(t *testing.T) {
	model := NewModel(nil, "test-project")
	model = model.SetProjectInfo("/test/path", "/test/stack.yml", false)

	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'T'}})

	if cmd == nil {
		t.Error("Expected installTemplateCmd to be returned")
	}

	// Should trigger template installation, not change state to output
	if updatedModel.rightPanelState == "output" {
		t.Error("rightPanelState should not change to 'output' for template installation")
	}
}

func TestDashboardModel_KeyHandler_T_WithPublicHTML(t *testing.T) {
	model := NewModel(nil, "test-project")
	model = model.SetProjectInfo("/test/path", "/test/stack.yml", true)

	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'T'}})

	if updatedModel.rightPanelState != "output" {
		t.Errorf("Expected rightPanelState to be 'output', got '%s'", updatedModel.rightPanelState)
	}

	if cmd == nil {
		t.Error("Expected composeDownCmd to be returned")
	}
}

func TestDashboardModel_KeyHandler_R(t *testing.T) {
	model := NewModel(nil, "test-project")
	model = model.SetProjectInfo("/test/path", "/test/stack.yml", true)

	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'R'}})

	if updatedModel.rightPanelState != "output" {
		t.Errorf("Expected rightPanelState to be 'output', got '%s'", updatedModel.rightPanelState)
	}

	if cmd == nil {
		t.Error("Expected composeRestartCmd to be returned")
	}
}

func TestDashboardModel_KeyHandler_D_ShowsFirstModal(t *testing.T) {
	model := NewModel(nil, "test-project")
	model = model.SetProjectInfo("/test/path", "/test/stack.yml", true)

	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'D'}})

	if updatedModel.confirmationStage != 1 {
		t.Errorf("Expected confirmationStage to be 1, got %d", updatedModel.confirmationStage)
	}

	if updatedModel.firstInput != "" {
		t.Error("firstInput should be empty initially")
	}
}

func TestDashboardModel_Modal_FirstConfirmation_Yes(t *testing.T) {
	model := NewModel(nil, "test-project")
	model = model.SetProjectInfo("/test/path", "/test/stack.yml", true)
	model.confirmationStage = 1

	// Type 'yes'
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})

	if model.firstInput != "yes" {
		t.Errorf("Expected firstInput to be 'yes', got '%s'", model.firstInput)
	}

	// Press enter
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if updatedModel.confirmationStage != 2 {
		t.Errorf("Expected confirmationStage to be 2, got %d", updatedModel.confirmationStage)
	}

	if updatedModel.secondInput != "" {
		t.Error("secondInput should be empty after advancing to stage 2")
	}
}

func TestDashboardModel_Modal_FirstConfirmation_Wrong(t *testing.T) {
	model := NewModel(nil, "test-project")
	model = model.SetProjectInfo("/test/path", "/test/stack.yml", true)
	model.confirmationStage = 1

	// Type 'no'
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'o'}})

	// Press enter
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if updatedModel.confirmationStage != 1 {
		t.Errorf("Expected confirmationStage to remain 1, got %d", updatedModel.confirmationStage)
	}

	if updatedModel.firstInput != "" {
		t.Error("firstInput should be reset after wrong input")
	}
}

func TestDashboardModel_Modal_SecondConfirmation_Destroy(t *testing.T) {
	model := NewModel(nil, "test-project")
	model = model.SetProjectInfo("/test/path", "/test/stack.yml", true)
	model.confirmationStage = 2

	// Type 'destroy'
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'o'}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

	if model.secondInput != "destroy" {
		t.Errorf("Expected secondInput to be 'destroy', got '%s'", model.secondInput)
	}

	// Press enter
	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if updatedModel.confirmationStage != 0 {
		t.Errorf("Expected confirmationStage to be 0, got %d", updatedModel.confirmationStage)
	}

	if cmd == nil {
		t.Error("Expected composeDestroyCmd to be returned")
	}

	if updatedModel.rightPanelState != "output" {
		t.Errorf("Expected rightPanelState to be 'output', got '%s'", updatedModel.rightPanelState)
	}
}

func TestDashboardModel_Modal_EscapeToCancel(t *testing.T) {
	model := NewModel(nil, "test-project")
	model = model.SetProjectInfo("/test/path", "/test/stack.yml", true)
	model.confirmationStage = 1
	model.firstInput = "some input"

	// Press escape
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEsc})

	if updatedModel.confirmationStage != 0 {
		t.Errorf("Expected confirmationStage to be 0, got %d", updatedModel.confirmationStage)
	}

	if updatedModel.firstInput != "" {
		t.Error("firstInput should be cleared after escape")
	}
}

func TestDashboardModel_Modal_Backspace(t *testing.T) {
	model := NewModel(nil, "test-project")
	model = model.SetProjectInfo("/test/path", "/test/stack.yml", true)
	model.confirmationStage = 1

	// Type 'yes'
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})

	// Press backspace
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyBackspace})

	if updatedModel.firstInput != "ye" {
		t.Errorf("Expected firstInput to be 'ye' after backspace, got '%s'", updatedModel.firstInput)
	}
}

// T146a: Mouse handler tests
// Note: No mouse handlers currently implemented, but adding placeholder test
func TestDashboardModel_MouseHandler_Placeholder(t *testing.T) {
	// Placeholder for future mouse handler tests
	// Currently, the dashboard doesn't handle mouse events
	model := NewModel(nil, "test-project")
	
	// Mouse events should not crash the model
	mouseMsg := tea.MouseMsg{Type: tea.MouseLeft, X: 10, Y: 10}
	_, _ = model.Update(mouseMsg)
	
	// Test passes if no panic occurs
}

