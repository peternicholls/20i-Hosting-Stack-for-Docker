// Project: 20i Stack Manager TUI
// File: errors_test.go
// Purpose: Unit tests for user-friendly error formatting
// Version: 0.1.0
// Updated: 2025-12-28

package stack

import (
	"errors"
	"strings"
	"testing"
)

func TestFormatUserError(t *testing.T) {
	tests := []struct {
		name          string
		err           error
		expectedMsg   string
		shouldContain []string // Multiple strings that should be present
	}{
		{
			name:          "nil error",
			err:           nil,
			expectedMsg:   "",
			shouldContain: nil,
		},
		{
			name:          "port conflict with port number (lowercase)",
			err:           errors.New("Error response from daemon: driver failed programming external connectivity on endpoint web (abc123): Bind for 0.0.0.0:80 failed: port is already allocated"),
			expectedMsg:   "",
			shouldContain: []string{"Port 80 is already in use", "Stop the conflicting service"},
		},
		{
			name:          "port conflict with port number (uppercase)",
			err:           errors.New("Error: Bind for 0.0.0.0:8080 failed: Port is already allocated"),
			expectedMsg:   "",
			shouldContain: []string{"Port 8080 is already in use", "Stop the conflicting service"},
		},
		{
			name:          "port conflict - address already in use",
			err:           errors.New("Error starting userland proxy: listen tcp4 0.0.0.0:443: bind: address already in use"),
			expectedMsg:   "",
			shouldContain: []string{"Port 443 is already in use", "Stop the conflicting service"},
		},
		{
			name:          "port conflict without detectable port number",
			err:           errors.New("address already in use"),
			expectedMsg:   "",
			shouldContain: []string{"Port is already in use", "Stop the conflicting service"},
		},
		{
			name:          "docker daemon not running - cannot connect",
			err:           errors.New("Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?"),
			expectedMsg:   "Docker daemon is not running. Start Docker Desktop and try again.",
			shouldContain: nil,
		},
		{
			name:          "docker daemon not running - connection refused",
			err:           errors.New("Error: failed to connect to docker: connection refused"),
			expectedMsg:   "Docker daemon is not running. Start Docker Desktop and try again.",
			shouldContain: nil,
		},
		{
			name:          "docker daemon not running - Windows named pipe",
			err:           errors.New("error during connect: This error may indicate that the docker daemon is not running.: open //./pipe/docker_engine: The system cannot find the file specified."),
			expectedMsg:   "Docker daemon is not running. Start Docker Desktop and try again.",
			shouldContain: nil,
		},
		{
			name:          "permission denied - unix socket",
			err:           errors.New("Got permission denied while trying to connect to the Docker daemon socket at unix:///var/run/docker.sock"),
			expectedMsg:   "Permission denied. Run Docker as your user or check socket permissions.",
			shouldContain: nil,
		},
		{
			name:          "permission denied - generic",
			err:           errors.New("permission denied"),
			expectedMsg:   "Permission denied. Run Docker as your user or check socket permissions.",
			shouldContain: nil,
		},
		{
			name:          "timeout error",
			err:           errors.New("context deadline exceeded"),
			expectedMsg:   "Operation timed out. The service may be unresponsive or taking too long to respond.",
			shouldContain: nil,
		},
		{
			name:          "timeout - explicit timeout message",
			err:           errors.New("Error: timeout waiting for container to start"),
			expectedMsg:   "Operation timed out. The service may be unresponsive or taking too long to respond.",
			shouldContain: nil,
		},
		{
			name:          "container not found",
			err:           errors.New("Error: No such container: abc123"),
			expectedMsg:   "Container not found. It may have been removed. Try refreshing the view.",
			shouldContain: nil,
		},
		{
			name:          "not found - generic",
			err:           errors.New("not found"),
			expectedMsg:   "Container not found. It may have been removed. Try refreshing the view.",
			shouldContain: nil,
		},
		{
			name:          "unknown error - short message",
			err:           errors.New("unexpected error occurred"),
			expectedMsg:   "",
			shouldContain: []string{"An error occurred", "Details: unexpected error occurred"},
		},
		{
			name:          "unknown error - long message truncated",
			err:           errors.New("This is a very long error message that should be truncated to avoid overwhelming the user interface with too much text in a single error display line which could make it hard to read"),
			expectedMsg:   "",
			shouldContain: []string{"An error occurred", "Details:", "..."},
		},
		{
			name:          "unknown error - multiline (only first line shown)",
			err:           errors.New("First line of error\nSecond line of error\nThird line of error"),
			expectedMsg:   "",
			shouldContain: []string{"An error occurred", "Details: First line of error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatUserError(tt.err)

			// If expectedMsg is set, check exact match
			if tt.expectedMsg != "" {
				if result != tt.expectedMsg {
					t.Errorf("FormatUserError() = %q, want %q", result, tt.expectedMsg)
				}
				return
			}

			// If shouldContain is set, check all substrings are present
			if tt.shouldContain != nil {
				for _, substr := range tt.shouldContain {
					if !strings.Contains(result, substr) {
						t.Errorf("FormatUserError() = %q, should contain %q", result, substr)
					}
				}
			}
		})
	}
}

func TestExtractPortNumber(t *testing.T) {
	tests := []struct {
		name     string
		errStr   string
		expected string
	}{
		{
			name:     "port with 'port' keyword",
			errStr:   "Bind for 0.0.0.0:80 failed: port is already allocated",
			expected: "80",
		},
		{
			name:     "port with uppercase 'Port'",
			errStr:   "Port 8080 is already in use",
			expected: "8080",
		},
		{
			name:     "port from bind address only",
			errStr:   "listen tcp4 0.0.0.0:443: bind: address already in use",
			expected: "443",
		},
		{
			name:     "port from colon notation",
			errStr:   "error on :3306 connection",
			expected: "3306",
		},
		{
			name:     "no port number present",
			errStr:   "address already in use",
			expected: "",
		},
		{
			name:     "multiple ports - returns first match",
			errStr:   "port 80 conflicts with port 8080",
			expected: "80",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractPortNumber(tt.errStr)
			if result != tt.expected {
				t.Errorf("extractPortNumber() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFormatUnknownError(t *testing.T) {
	tests := []struct {
		name     string
		errStr   string
		maxLen   int
		hasEllip bool
	}{
		{
			name:     "short error",
			errStr:   "short error",
			maxLen:   100,
			hasEllip: false,
		},
		{
			name:     "exact boundary (100 chars)",
			errStr:   strings.Repeat("a", 100),
			maxLen:   100,
			hasEllip: false,
		},
		{
			name:     "long error truncated",
			errStr:   strings.Repeat("a", 150),
			maxLen:   100,
			hasEllip: true,
		},
		{
			name:     "multiline - only first line",
			errStr:   "first line\nsecond line\nthird line",
			maxLen:   100,
			hasEllip: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatUnknownError(tt.errStr)

			// Should always start with "An error occurred. Details: "
			if !strings.HasPrefix(result, "An error occurred. Details: ") {
				t.Errorf("formatUnknownError() should start with 'An error occurred. Details: ', got %q", result)
			}

			// Check for ellipsis if expected
			if tt.hasEllip && !strings.Contains(result, "...") {
				t.Errorf("formatUnknownError() should contain '...' for long errors, got %q", result)
			}

			// Check that multiline errors only show first line
			if strings.Contains(tt.errStr, "\n") && strings.Contains(result, "\n") {
				t.Errorf("formatUnknownError() should not contain newlines, got %q", result)
			}
		})
	}
}

// TestRealWorldDockerErrors tests formatting with realistic Docker/Compose error messages
func TestRealWorldDockerErrors(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		wantContain string
	}{
		{
			name:        "docker compose up - port conflict",
			err:         errors.New("Error response from daemon: driver failed programming external connectivity on endpoint myproject-web-1 (1234567890ab): Bind for 0.0.0.0:80 failed: port is already allocated"),
			wantContain: "Port 80 is already in use",
		},
		{
			name:        "docker start - daemon not running",
			err:         errors.New("Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?"),
			wantContain: "Docker daemon is not running",
		},
		{
			name:        "docker compose - permission denied on Linux",
			err:         errors.New("Got permission denied while trying to connect to the Docker daemon socket at unix:///var/run/docker.sock: Get \"http://%2Fvar%2Frun%2Fdocker.sock/v1.43/containers/json\": dial unix /var/run/docker.sock: connect: permission denied"),
			wantContain: "Permission denied",
		},
		{
			name:        "compose up - network conflict",
			err:         errors.New("Error response from daemon: network with name myproject_default already exists"),
			wantContain: "An error occurred",
		},
		{
			name:        "docker stop - container not found",
			err:         errors.New("Error response from daemon: No such container: mycontainer"),
			wantContain: "Container not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatUserError(tt.err)
			if !strings.Contains(result, tt.wantContain) {
				t.Errorf("FormatUserError() = %q, want to contain %q", result, tt.wantContain)
			}
		})
	}
}
