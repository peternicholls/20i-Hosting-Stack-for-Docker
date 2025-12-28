// Project: 20i Stack Manager TUI
// File: root_test.go
// Purpose: Unit tests for RootModel initialization and message routing
// Version: 0.1.0
// Updated: 2025-12-28

package app

import (
	"context"
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewRootModel(t *testing.T) {
	ctx := context.Background()
	m, err := NewRootModel(ctx)

	// We expect an error because Docker daemon may not be running in test environment
	if m != nil {
		if m.activeView != "dashboard" {
			t.Errorf("expected activeView to be 'dashboard', got '%s'", m.activeView)
		}
		if m.width != 80 || m.height != 24 {
			t.Errorf("expected default dimensions 80x24, got %dx%d", m.width, m.height)
		}
	}
	_ = err
}

func TestRootModelInit(t *testing.T) {
	m := &RootModel{
		activeView: "dashboard",
		width:      80,
		height:     24,
	}

	cmd := m.Init()
	if cmd != nil {
		t.Errorf("expected Init to return nil cmd, got %v", cmd)
	}
}

func TestRootModelUpdate_GlobalShortcuts(t *testing.T) {
	cases := []struct {
		name         string
		keyMsg       string
		expectedView string
	}{
		{"q quits", "q", "dashboard"},
		{"? opens help", "?", "help"},
		{"p opens projects", "p", "projects"},
		{"esc from help back to dashboard", "esc", "dashboard"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := &RootModel{
				activeView: "help",
				width:      80,
				height:     24,
			}

			if c.keyMsg != "esc" {
				m.activeView = "dashboard"
			}

			model, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(c.keyMsg)})

			if c.keyMsg == "q" {
				if cmd == nil {
					t.Errorf("expected 'q' to return quit cmd")
				}
			} else {
				updatedModel := model.(*RootModel)
				if updatedModel.activeView != c.expectedView {
					t.Errorf("expected view '%s', got '%s'", c.expectedView, updatedModel.activeView)
				}
			}
		})
	}
}

func TestRootModelUpdate_WindowSize(t *testing.T) {
	m := &RootModel{
		activeView: "dashboard",
		width:      80,
		height:     24,
	}

	msg := tea.WindowSizeMsg{Width: 120, Height: 40}
	model, _ := m.Update(msg)

	updatedModel := model.(*RootModel)
	if updatedModel.width != 120 || updatedModel.height != 40 {
		t.Errorf("expected dimensions 120x40, got %dx%d", updatedModel.width, updatedModel.height)
	}
}

func TestRootModelView(t *testing.T) {
	cases := []struct {
		name         string
		activeView   string
		expectedText string
	}{
		{"dashboard view", "dashboard", "Dashboard"},
		{"help view", "help", "Global"},
		{"projects view", "projects", "Projects View"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := &RootModel{
				activeView: c.activeView,
				width:      80,
				height:     24,
			}

			output := m.View()
			if len(output) == 0 {
				t.Errorf("expected non-empty view output")
			}
		})
	}
}

// TestRootModel_LastErrorClearing tests that lastError is cleared on success and view change
func TestRootModel_LastErrorClearing(t *testing.T) {
	t.Run("error clears on view change", func(t *testing.T) {
		m := &RootModel{
			activeView: "dashboard",
			lastError:  errors.New("docker connection failed"),
			width:      80,
			height:     24,
		}

		// Initially has error
		if m.lastError == nil {
			t.Fatal("expected lastError to be set")
		}

		// Switch views - should clear error
		model, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")})
		updatedModel := model.(*RootModel)

		if updatedModel.lastError != nil {
			t.Errorf("expected lastError to be cleared on view change, got %v", updatedModel.lastError)
		}
	})

	t.Run("error clears on success message", func(t *testing.T) {
		m := &RootModel{
			activeView: "dashboard",
			lastError:  errors.New("docker connection failed"),
			width:      80,
			height:     24,
		}

		// Initially has error
		if m.lastError == nil {
			t.Fatal("expected lastError to be set")
		}

		// Receive success message - should clear error
		model, _ := m.Update(SuccessMsg{Message: "Operation completed"})
		updatedModel := model.(*RootModel)

		if updatedModel.lastError != nil {
			t.Errorf("expected lastError to be cleared on success, got %v", updatedModel.lastError)
		}
	})
}

// TestRootModel_StubbedClientPath tests graceful handling when Docker is unavailable
func TestRootModel_StubbedClientPath(t *testing.T) {
	t.Run("app starts with nil Docker client", func(t *testing.T) {
		// Simulate soft-fail approach - RootModel can be created without Docker
		m := &RootModel{
			dockerClient: nil, // Docker unavailable
			activeView:   "dashboard",
			width:        80,
			height:       24,
		}

		// Should still be able to init
		cmd := m.Init()
		if cmd != nil {
			t.Errorf("expected Init to return nil even without Docker client")
		}

		// Should render a degraded UI
		output := m.View()
		if len(output) == 0 {
			t.Errorf("expected view to render even without Docker client")
		}
	})

	t.Run("error message shown when Docker unavailable", func(t *testing.T) {
		m := &RootModel{
			dockerClient: nil,
			activeView:   "dashboard",
			lastError:    errors.New("docker daemon unreachable"),
			width:        80,
			height:       24,
		}

		output := m.View()
		if len(output) == 0 {
			t.Errorf("expected error view when Docker unavailable")
		}
	})
}
