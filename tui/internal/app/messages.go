// Project: 20i Stack Manager TUI
// File: messages.go
// Purpose: Custom message types for Bubble Tea event-driven architecture
// Version: 0.1.0
// Updated: 2025-12-28

package app

// StatsMsg sent when container stats are updated
type StatsMsg struct {
	Timestamp int64
	Stats     map[string]interface{}
}

// ContainerListMsg sent when container list is refreshed
type ContainerListMsg struct {
	Error error
	Data  interface{}
}

// LogLineMsg sent when a new log line arrives
type LogLineMsg struct {
	ContainerID string
	Timestamp   int64
	Source      string // "stdout" or "stderr"
	Line        string
}

// ContainerActionMsg sent when user requests action on container
type ContainerActionMsg struct {
	Action      string // "start", "stop", "restart"
	ContainerID string
}

// ContainerActionResultMsg sent when container action completes
type ContainerActionResultMsg struct {
	Success bool
	Message string
	Error   error
}

// ComposeActionMsg sent when user requests action on entire stack
type ComposeActionMsg struct {
	Action string // "stop", "restart", "down"
}

// ComposeActionResultMsg sent when compose action completes
type ComposeActionResultMsg struct {
	Success bool
	Message string
	Error   error
}

// ProjectSwitchMsg sent when user switches projects
type ProjectSwitchMsg struct {
	ProjectName string
	ProjectPath string
}

// ErrorMsg sent when an error occurs
type ErrorMsg struct {
	Title   string
	Message string
	Err     error
}

// SuccessMsg sent when operation succeeds
type SuccessMsg struct {
	Message string
}
