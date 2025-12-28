package dashboard

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/peternicholls/20i-stack/tui/internal/docker"
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

// TestRefreshTimer_StartsWhenContainersLoaded tests that the refresh timer starts when containers are loaded
func TestRefreshTimer_StartsWhenContainersLoaded(t *testing.T) {
	model := NewModel(nil, "test-project")
	
	// Initially, refresh should not be active
	if model.refreshActive {
		t.Error("Refresh should not be active initially")
	}
	
	// Simulate loading containers
	containers := []docker.Container{
		{ID: "c1", Service: "web", Status: docker.StatusRunning},
	}
	msg := containerListMsg{containers: containers, err: nil}
	
	updatedModel, cmd := model.Update(msg)
	
	// After loading containers, refresh should be active
	if !updatedModel.refreshActive {
		t.Error("Refresh should be active after loading containers")
	}
	
	// Should return a command to start the timer
	if cmd == nil {
		t.Error("Should return a command to start refresh timer")
	}
}

// TestRefreshTimer_DoesNotStartWithNoContainers tests that timer doesn't start when there are no containers
func TestRefreshTimer_DoesNotStartWithNoContainers(t *testing.T) {
	model := NewModel(nil, "test-project")
	
	// Simulate loading empty container list
	msg := containerListMsg{containers: []docker.Container{}, err: nil}
	
	updatedModel, _ := model.Update(msg)
	
	// Refresh should not be active with no containers
	if updatedModel.refreshActive {
		t.Error("Refresh should not be active with no containers")
	}
}

// TestRefreshTimer_CancelsOnViewSwitch tests that timer cancels when switching views
func TestRefreshTimer_CancelsOnViewSwitch(t *testing.T) {
	model := NewModel(nil, "test-project")
	model.refreshActive = true // Simulate active refresh
	
	testKeys := []string{"esc", "?", "p"}
	
	for _, key := range testKeys {
		m := model
		m.refreshActive = true
		
		keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
		updatedModel, _ := m.Update(keyMsg)
		
		if updatedModel.refreshActive {
			t.Errorf("Refresh should be cancelled on key '%s'", key)
		}
	}
}

// TestRefreshTimer_StopsWhenStackStops tests that timer stops when stack has no containers
func TestRefreshTimer_StopsWhenStackStops(t *testing.T) {
	model := NewModel(nil, "test-project")
	model.refreshActive = true
	model.containers = []docker.Container{
		{ID: "c1", Service: "web", Status: docker.StatusRunning},
	}
	
	// Simulate refresh with empty container list (stack stopped)
	msg := stackContainersMsg{containers: []docker.Container{}, err: nil}
	
	updatedModel, _ := model.Update(msg)
	
	// Refresh should be stopped
	if updatedModel.refreshActive {
		t.Error("Refresh should stop when stack has no containers")
	}
}

// TestRefreshTimer_TickTriggersRefresh tests that refresh tick triggers container refresh
func TestRefreshTimer_TickTriggersRefresh(t *testing.T) {
	model := NewModel(nil, "test-project")
	model.refreshActive = true
	model.containers = []docker.Container{
		{ID: "c1", Service: "web", Status: docker.StatusRunning},
	}
	
	// Simulate a refresh tick
	msg := refreshTickMsg{}
	
	_, cmd := model.Update(msg)
	
	// Should return a command (batch of refresh + next tick)
	if cmd == nil {
		t.Error("Refresh tick should return a command")
	}
}

// TestRefreshTimer_NoTickWhenInactive tests that tick is ignored when refresh is inactive
func TestRefreshTimer_NoTickWhenInactive(t *testing.T) {
	model := NewModel(nil, "test-project")
	model.refreshActive = false
	
	// Simulate a refresh tick when inactive
	msg := refreshTickMsg{}
	
	_, cmd := model.Update(msg)
	
	// Should not return a command when inactive
	if cmd != nil {
		t.Error("Refresh tick should not return command when inactive")
	}
}

// TestStackContainersMsg_UpdatesContainers tests that stackContainersMsg updates container list
func TestStackContainersMsg_UpdatesContainers(t *testing.T) {
	model := NewModel(nil, "test-project")
	model.containers = []docker.Container{
		{ID: "c1", Service: "web", Status: docker.StatusRunning, CPUPercent: 0},
	}
	
	// Simulate refreshed container data with CPU stats
	updatedContainers := []docker.Container{
		{ID: "c1", Service: "web", Status: docker.StatusRunning, CPUPercent: 25.5},
	}
	msg := stackContainersMsg{containers: updatedContainers, err: nil}
	
	updatedModel, _ := model.Update(msg)
	
	// Container list should be updated
	if len(updatedModel.containers) != 1 {
		t.Errorf("Expected 1 container, got %d", len(updatedModel.containers))
	}
	
	if updatedModel.containers[0].CPUPercent != 25.5 {
		t.Errorf("Expected CPU 25.5%%, got %.1f%%", updatedModel.containers[0].CPUPercent)
	}
}

// TestStackContainersMsg_HandlesError tests that stackContainersMsg handles errors gracefully
func TestStackContainersMsg_HandlesError(t *testing.T) {
	model := NewModel(nil, "test-project")
	model.refreshActive = true
	
	// Simulate error during refresh
	msg := stackContainersMsg{containers: nil, err: docker.ErrDaemonUnreachable}
	
	updatedModel, _ := model.Update(msg)
	
	// Should keep refresh active even on error
	if !updatedModel.refreshActive {
		t.Error("Refresh should remain active even on error")
	}
	
	// Error should be recorded
	if updatedModel.lastError == nil {
		t.Error("Error should be recorded in lastError")
	}
}
