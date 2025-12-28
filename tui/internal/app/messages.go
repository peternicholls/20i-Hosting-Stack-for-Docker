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

// ContainerActionMsg sent when user requests action on container.
// Valid actions: "start", "stop", "restart".
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

// ComposeActionMsg sent when user requests action on entire stack.
// Valid actions: "stop", "restart", "down".
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

// ProjectDetectedMsg sent when project is successfully detected
type ProjectDetectedMsg struct {
	ProjectName   string
	ProjectPath   string
	HasPublicHTML bool
}

// TemplateInstalledMsg sent when template is successfully installed
type TemplateInstalledMsg struct {
	ProjectPath string
}

// StackStartMsg sent when user requests to start the stack
type StackStartMsg struct {
	CodeDir string // Project directory to start stack for
}

// StackStopMsg sent when user requests to stop the stack
type StackStopMsg struct {
	CodeDir string // Project directory to stop stack for
}

// StackRestartMsg sent when user requests to restart the stack
type StackRestartMsg struct {
	CodeDir string // Project directory to restart stack for
}

// StackDestroyMsg sent when user requests to destroy the stack (with volumes)
type StackDestroyMsg struct {
	CodeDir string // Project directory to destroy stack for
}

// StackOutputMsg sent when streaming compose command output
type StackOutputMsg struct {
	Line      string // Line of output from compose command
	IsError   bool   // True if from stderr
	// TODO: Populate Timestamp from compose output when streaming is implemented in compose.go.
	Timestamp int64 // Unix timestamp
}

// StackStatusMsg sent when stack operation completes
type StackStatusMsg struct {
	Operation string // "start", "stop", "restart", "destroy"
	Success   bool   // True if operation succeeded
	Message   string // Human-readable status message
	Error     error  // Error if operation failed
}

// StackContainersMsg sent when container status is refreshed (auto-refresh)
type StackContainersMsg struct {
	Containers []interface{} // List of containers with updated status
	Error      error         // Error if refresh failed
}
