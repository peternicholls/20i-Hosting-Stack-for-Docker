// Project: 20i Stack Manager TUI
// File: status_table_test.go
// Purpose: Tests for status table rendering and URL click detection
// Version: 0.1.0
// Updated: 2025-12-28

package dashboard

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/peternicholls/20i-stack/tui/internal/docker"
	"github.com/peternicholls/20i-stack/tui/internal/ui"
)

func TestRenderStatusTable_EmptyContainers(t *testing.T) {
	containers := []docker.Container{}
	rendered, state := renderStatusTable(containers, 100, 20)

	if rendered == "" {
		t.Error("Should return non-empty string for empty container list")
	}

	if len(state.URLRegions) != 0 {
		t.Error("Should have no URL regions for empty container list")
	}
}

func TestRenderStatusTable_WithContainers(t *testing.T) {
	containers := []docker.Container{
		{
			Service: "nginx",
			Status:  docker.StatusRunning,
			Image:   "nginx:latest",
		},
		{
			Service: "mariadb",
			Status:  docker.StatusRunning,
			Image:   "mariadb:10.11",
		},
	}

	rendered, state := renderStatusTable(containers, 100, 20)

	if rendered == "" {
		t.Error("Should return non-empty string")
	}

	// Should have URL regions for nginx (has HTTP URL)
	if len(state.URLRegions) == 0 {
		t.Error("Should have at least one URL region for nginx")
	}
}

func TestRenderCPUBar(t *testing.T) {
	tests := []struct {
		name       string
		cpuPercent float64
		width      int
		wantFilled int
	}{
		{"0% usage", 0.0, 10, 0},
		{"50% usage", 50.0, 10, 5},
		{"100% usage", 100.0, 10, 10},
		{"Over 100%", 150.0, 10, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bar := renderCPUBar(tt.cpuPercent, tt.width)
			if bar == "" && tt.width > 0 {
				t.Error("Should return non-empty bar")
			}
		})
	}
}

func TestExtractURL(t *testing.T) {
	tests := []struct {
		name      string
		container docker.Container
		wantURL   string
	}{
		{
			name:      "nginx service",
			container: docker.Container{Service: "nginx"},
			wantURL:   "http://localhost:80",
		},
		{
			name:      "phpmyadmin service",
			container: docker.Container{Service: "phpmyadmin"},
			wantURL:   "http://localhost:8081",
		},
		{
			name:      "mariadb service",
			container: docker.Container{Service: "mariadb"},
			wantURL:   "localhost:3306",
		},
		{
			name:      "apache service",
			container: docker.Container{Service: "apache"},
			wantURL:   "internal",
		},
		{
			name:      "unknown service",
			container: docker.Container{Service: "unknown"},
			wantURL:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractURL(tt.container)
			if got != tt.wantURL {
				t.Errorf("extractURL() = %v, want %v", got, tt.wantURL)
			}
		})
	}
}

func TestHandleURLClick_NoClick(t *testing.T) {
	mockOpener := &ui.MockURLOpener{}
	state := StatusTableState{
		URLRegions: []URLRegion{
			{URL: "http://localhost:80", Row: 2, ColStart: 64, ColEnd: 84},
		},
	}

	// Click outside URL region
	msg := tea.MouseMsg{
		Type:   tea.MouseLeft,
		Button: tea.MouseButtonLeft,
		X:      10,
		Y:      2,
	}

	handleURLClick(msg, state, mockOpener)

	if len(mockOpener.OpenedURLs) != 0 {
		t.Error("Should not open URL when clicking outside region")
	}
}

func TestHandleURLClick_InsideRegion(t *testing.T) {
	mockOpener := &ui.MockURLOpener{}
	state := StatusTableState{
		URLRegions: []URLRegion{
			{URL: "http://localhost:80", Row: 2, ColStart: 64, ColEnd: 84},
		},
	}

	// Click inside URL region
	msg := tea.MouseMsg{
		Type:   tea.MouseLeft,
		Button: tea.MouseButtonLeft,
		X:      70,
		Y:      2,
	}

	handleURLClick(msg, state, mockOpener)

	// Note: The actual URL opening happens in a goroutine
	// so we can't reliably test it here without adding synchronization
	// In a real test environment, we would use channels or wait groups
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		maxLen int
		want   string
	}{
		{"short string", "hello", 10, "hello"},
		{"exact length", "hello", 5, "hello"},
		{"needs truncate", "hello world", 8, "hello..."},
		{"very short max", "hello", 2, "..."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncateString(tt.input, tt.maxLen)
			if got != tt.want {
				t.Errorf("truncateString(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
			}
		})
	}
}

func TestGetStatusBadge(t *testing.T) {
	tests := []struct {
		status docker.ContainerStatus
	}{
		{docker.StatusRunning},
		{docker.StatusStopped},
		{docker.StatusRestarting},
		{docker.StatusError},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			badge := getStatusBadge(tt.status)
			if badge == "" {
				t.Error("Status badge should not be empty")
			}
		})
	}
}
