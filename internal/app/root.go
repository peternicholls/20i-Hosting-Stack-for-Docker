// Project: 20i Stack Manager TUI
// File: root.go
// Purpose: Root model for managing application state
// Version: 0.1.0
// Updated: 2025-12-28

package app

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/peternicholls/20i-stack/internal/docker"
)

// RootModel struct to manage application state
type RootModel struct {
	activeView   string
	dockerClient *docker.Client
	width        int
	height       int
}

// Init initializes the RootModel and Docker client
func (m *RootModel) Init() bubbletea.Cmd {
	var err error
	m.dockerClient, err = docker.NewClient()
	if err != nil {
		return nil
	}
	return nil
}

// Update handles incoming messages and updates the model
func (m *RootModel) Update(msg bubbletea.Msg) (bubbletea.Model, bubbletea.Cmd) {
	switch msg := msg.(type) {
	case bubbletea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, bubbletea.Quit
		case "?":
			// Show help
		case "p":
			// Show projects
		}
	}
	return m, nil
}
