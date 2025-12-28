package dashboard

import (
	"errors"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/peternicholls/20i-stack/tui/internal/docker"
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

func TestDashboardModel_ErrorDisplay(t *testing.T) {
	t.Run("error message is set on container action failure", func(t *testing.T) {
		model := NewModel(nil, "test-project")
		
		// Simulate a container action failure
		msg := containerActionResultMsg{
			success: false,
			message: "Failed to start container",
			err:     errors.New("docker: Error response from daemon: port is already allocated"),
		}
		
		updatedModel, cmd := model.Update(msg)
		
		// Error message should be set
		if updatedModel.errorMsg == "" {
			t.Error("Expected errorMsg to be set after container action failure")
		}
		
		// Should contain user-friendly message
		if len(updatedModel.errorMsg) == 0 {
			t.Error("Expected non-empty error message")
		}
		
		// Timer command should be returned
		if cmd == nil {
			t.Error("Expected timer command to be returned")
		}
	})

	t.Run("error message is cleared on key press", func(t *testing.T) {
		model := NewModel(nil, "test-project")
		model.errorMsg = "Some error message"
		
		// Simulate a key press
		updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
		
		// Error message should be cleared
		if updatedModel.errorMsg != "" {
			t.Errorf("Expected errorMsg to be cleared after key press, got %q", updatedModel.errorMsg)
		}
	})

	t.Run("error message is cleared on timer expiry", func(t *testing.T) {
		model := NewModel(nil, "test-project")
		displayTime := time.Now()
		model.errorMsg = "Some error message"
		model.errorDisplayTime = displayTime
		
		// Simulate timer expiry message
		updatedModel, _ := model.Update(errorClearMsg{displayTime: displayTime})
		
		// Error message should be cleared
		if updatedModel.errorMsg != "" {
			t.Errorf("Expected errorMsg to be cleared after timer expiry, got %q", updatedModel.errorMsg)
		}
	})

	t.Run("error message is NOT cleared if timer doesn't match", func(t *testing.T) {
		model := NewModel(nil, "test-project")
		displayTime1 := time.Now()
		displayTime2 := displayTime1.Add(1 * time.Second)
		
		model.errorMsg = "Some error message"
		model.errorDisplayTime = displayTime2
		
		// Simulate timer expiry with old timestamp
		updatedModel, _ := model.Update(errorClearMsg{displayTime: displayTime1})
		
		// Error message should NOT be cleared
		if updatedModel.errorMsg == "" {
			t.Error("Expected errorMsg to remain when timer timestamp doesn't match")
		}
	})

	t.Run("error is cleared on successful action", func(t *testing.T) {
		model := NewModel(nil, "test-project")
		model.errorMsg = "Some error message"
		
		// Simulate successful container action
		msg := containerActionResultMsg{
			success: true,
			message: "Container started successfully",
		}
		
		updatedModel, _ := model.Update(msg)
		
		// Error message should be cleared
		if updatedModel.errorMsg != "" {
			t.Errorf("Expected errorMsg to be cleared on success, got %q", updatedModel.errorMsg)
		}
	})

	t.Run("error is set on container list failure", func(t *testing.T) {
		model := NewModel(nil, "test-project")
		
		// Simulate container list failure
		msg := containerListMsg{
			containers: nil,
			err:        errors.New("Cannot connect to the Docker daemon"),
		}
		
		updatedModel, cmd := model.Update(msg)
		
		// Error message should be set
		if updatedModel.errorMsg == "" {
			t.Error("Expected errorMsg to be set after container list failure")
		}
		
		// Timer command should be returned
		if cmd == nil {
			t.Error("Expected timer command to be returned")
		}
	})

	t.Run("error is cleared on successful container list load", func(t *testing.T) {
		model := NewModel(nil, "test-project")
		model.errorMsg = "Some error message"
		
		// Simulate successful container list load
		msg := containerListMsg{
			containers: []docker.Container{
				{ID: "abc123", Service: "web", Status: docker.StatusRunning},
			},
			err: nil,
		}
		
		updatedModel, _ := model.Update(msg)
		
		// Error message should be cleared
		if updatedModel.errorMsg != "" {
			t.Errorf("Expected errorMsg to be cleared on successful load, got %q", updatedModel.errorMsg)
		}
	})
}

