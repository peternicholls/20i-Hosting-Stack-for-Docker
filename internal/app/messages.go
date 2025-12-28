// Project: 20i Stack Manager TUI
// File: messages.go
// Purpose: Custom message types for Bubble Tea
// Version: 0.1.0
// Updated: 2025-12-28

package app

// Message types for application
type statsMsg struct{}

type containerListMsg struct{}

type logLineMsg struct {
	Line string
}

type containerActionMsg struct {
	Action      string
	ContainerID string
}

type containerActionResultMsg struct {
	Success bool
	Message string
}
