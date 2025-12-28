// Project: 20i Stack Manager TUI
// File: status_test.go
// Purpose: Unit tests for stack status retrieval with URLs and CPU%
// Version: 0.1.0
// Updated: 2025-12-28

package stack

import (
	"os"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

func TestExtractContainerName(t *testing.T) {
	tests := []struct {
		name     string
		names    []string
		expected string
	}{
		{
			name:     "single name with slash",
			names:    []string{"/my-container"},
			expected: "my-container",
		},
		{
			name:     "single name without slash",
			names:    []string{"my-container"},
			expected: "my-container",
		},
		{
			name:     "multiple names",
			names:    []string{"/my-container", "/alias"},
			expected: "my-container",
		},
		{
			name:     "empty names",
			names:    []string{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractContainerName(tt.names)
			if result != tt.expected {
				t.Errorf("extractContainerName(%v) = %q, want %q", tt.names, result, tt.expected)
			}
		})
	}
}

func TestExtractServiceName(t *testing.T) {
	tests := []struct {
		name     string
		summary  container.Summary
		expected string
	}{
		{
			name: "service label present",
			summary: container.Summary{
				Names: []string{"/project-nginx-1"},
				Labels: map[string]string{
					"com.docker.compose.service": "nginx",
				},
			},
			expected: "nginx",
		},
		{
			name: "service label missing",
			summary: container.Summary{
				Names:  []string{"/my-container"},
				Labels: map[string]string{},
			},
			expected: "my-container",
		},
		{
			name: "no labels",
			summary: container.Summary{
				Names:  []string{"/my-container"},
				Labels: nil,
			},
			expected: "my-container",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractServiceName(tt.summary)
			if result != tt.expected {
				t.Errorf("extractServiceName() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestGenerateURL(t *testing.T) {
	tests := []struct {
		name     string
		summary  container.Summary
		expected string
	}{
		{
			name: "nginx with default port",
			summary: container.Summary{
				Names: []string{"/nginx"},
				Labels: map[string]string{
					"com.docker.compose.service": "nginx",
				},
				Ports: []types.Port{
					{PrivatePort: 80, PublicPort: 80, Type: "tcp"},
				},
			},
			expected: "http://localhost:80",
		},
		{
			name: "nginx with custom port",
			summary: container.Summary{
				Names: []string{"/nginx"},
				Labels: map[string]string{
					"com.docker.compose.service": "nginx",
				},
				Ports: []types.Port{
					{PrivatePort: 80, PublicPort: 8080, Type: "tcp"},
				},
			},
			expected: "http://localhost:8080",
		},
		{
			name: "phpmyadmin with default port",
			summary: container.Summary{
				Names: []string{"/phpmyadmin"},
				Labels: map[string]string{
					"com.docker.compose.service": "phpmyadmin",
				},
				Ports: []types.Port{
					{PrivatePort: 80, PublicPort: 8081, Type: "tcp"},
				},
			},
			expected: "http://localhost:8081",
		},
		{
			name: "mariadb with default port",
			summary: container.Summary{
				Names: []string{"/mariadb"},
				Labels: map[string]string{
					"com.docker.compose.service": "mariadb",
				},
				Ports: []types.Port{
					{PrivatePort: 3306, PublicPort: 3306, Type: "tcp"},
				},
			},
			expected: "localhost:3306",
		},
		{
			name: "apache internal",
			summary: container.Summary{
				Names: []string{"/apache"},
				Labels: map[string]string{
					"com.docker.compose.service": "apache",
				},
				Ports: []types.Port{},
			},
			expected: "internal",
		},
		{
			name: "unknown service",
			summary: container.Summary{
				Names: []string{"/unknown"},
				Labels: map[string]string{
					"com.docker.compose.service": "unknown",
				},
				Ports: []types.Port{},
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateURL(tt.summary)
			if result != tt.expected {
				t.Errorf("generateURL() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestResolvePort(t *testing.T) {
	tests := []struct {
		name          string
		ports         []types.Port
		containerPort uint16
		envVar        string
		envValue      string
		defaultPort   string
		expected      string
	}{
		{
			name: "port from Docker binding",
			ports: []types.Port{
				{PrivatePort: 80, PublicPort: 8080, Type: "tcp"},
			},
			containerPort: 80,
			envVar:        "HOST_PORT",
			envValue:      "",
			defaultPort:   "80",
			expected:      "8080",
		},
		{
			name:          "port from environment variable",
			ports:         []types.Port{},
			containerPort: 80,
			envVar:        "HOST_PORT",
			envValue:      "9090",
			defaultPort:   "80",
			expected:      "9090",
		},
		{
			name:          "port from default",
			ports:         []types.Port{},
			containerPort: 80,
			envVar:        "HOST_PORT",
			envValue:      "",
			defaultPort:   "80",
			expected:      "80",
		},
		{
			name: "Docker binding takes precedence over env",
			ports: []types.Port{
				{PrivatePort: 80, PublicPort: 7070, Type: "tcp"},
			},
			containerPort: 80,
			envVar:        "HOST_PORT",
			envValue:      "9090",
			defaultPort:   "80",
			expected:      "7070",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable if provided
			if tt.envValue != "" {
				os.Setenv(tt.envVar, tt.envValue)
				defer os.Unsetenv(tt.envVar)
			} else {
				os.Unsetenv(tt.envVar)
			}

			result := resolvePort(tt.ports, tt.containerPort, tt.envVar, tt.defaultPort)
			if result != tt.expected {
				t.Errorf("resolvePort() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestCalculateCPUPercent(t *testing.T) {
	tests := []struct {
		name     string
		stats    *container.StatsResponse
		expected float64
	}{
		{
			name:     "nil stats",
			stats:    nil,
			expected: 0.0,
		},
		{
			name: "zero delta",
			stats: &container.StatsResponse{
				CPUStats: container.CPUStats{
					CPUUsage: container.CPUUsage{
						TotalUsage:  1000000,
						PercpuUsage: []uint64{500000, 500000},
					},
					SystemUsage: 10000000,
				},
				PreCPUStats: container.CPUStats{
					CPUUsage: container.CPUUsage{
						TotalUsage: 1000000,
					},
					SystemUsage: 10000000,
				},
			},
			expected: 0.0,
		},
		{
			name: "valid CPU usage",
			stats: &container.StatsResponse{
				CPUStats: container.CPUStats{
					CPUUsage: container.CPUUsage{
						TotalUsage:  2000000,
						PercpuUsage: []uint64{1000000, 1000000},
					},
					SystemUsage: 20000000,
				},
				PreCPUStats: container.CPUStats{
					CPUUsage: container.CPUUsage{
						TotalUsage: 1000000,
					},
					SystemUsage: 10000000,
				},
			},
			expected: 20.0, // (2000000 - 1000000) / (20000000 - 10000000) * len(PercpuUsage) * 100 = 20%
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateCPUPercent(tt.stats)
			if result != tt.expected {
				t.Errorf("calculateCPUPercent() = %f, want %f", result, tt.expected)
			}
		})
	}
}

func TestContainerInfo_Fields(t *testing.T) {
	info := ContainerInfo{
		ID:      "abc123",
		Name:    "my-container",
		Service: "nginx",
		Image:   "nginx:latest",
		Status:  "Up 5 minutes",
		URL:     "http://localhost:80",
		CPUPerc: 12.5,
	}

	if info.ID != "abc123" {
		t.Errorf("ID = %q, want %q", info.ID, "abc123")
	}
	if info.Name != "my-container" {
		t.Errorf("Name = %q, want %q", info.Name, "my-container")
	}
	if info.Service != "nginx" {
		t.Errorf("Service = %q, want %q", info.Service, "nginx")
	}
	if info.Image != "nginx:latest" {
		t.Errorf("Image = %q, want %q", info.Image, "nginx:latest")
	}
	if info.Status != "Up 5 minutes" {
		t.Errorf("Status = %q, want %q", info.Status, "Up 5 minutes")
	}
	if info.URL != "http://localhost:80" {
		t.Errorf("URL = %q, want %q", info.URL, "http://localhost:80")
	}
	if info.CPUPerc != 12.5 {
		t.Errorf("CPUPerc = %f, want %f", info.CPUPerc, 12.5)
	}
}

// TestGetStackStatus_Integration is an integration test that requires a running Docker daemon
func TestGetStackStatus_Integration(t *testing.T) {
	if os.Getenv("DOCKER_INTEGRATION_TEST") != "1" {
		t.Skip("skipping integration test; set DOCKER_INTEGRATION_TEST=1 to run")
	}

	// Test with empty project name (should list all containers)
	containers, err := GetStackStatus("")
	if err != nil {
		t.Fatalf("GetStackStatus returned error: %v", err)
	}

	// Just verify we got a result (even if empty)
	if containers == nil {
		t.Fatal("GetStackStatus returned nil containers slice")
	}
}
