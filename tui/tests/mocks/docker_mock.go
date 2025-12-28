// Project: 20i Stack Manager TUI
// File: docker_mock.go
// Purpose: Mock Docker client interfaces and test data structures.
// Version: 0.1.0
// Updated: 2025-12-28
package mocks

import (
	"context"
	"io"
)

// MockDockerClient provides a mock implementation of the Docker client interface
// for testing purposes. This matches the docker-api.md contract.
type MockDockerClient interface {
	// Container lifecycle operations
	ListContainers(ctx context.Context, projectName string) ([]Container, error)
	StartContainer(ctx context.Context, containerID string) error
	StopContainer(ctx context.Context, containerID string, timeout int) error
	RestartContainer(ctx context.Context, containerID string, timeout int) error

	// Container inspection
	InspectContainer(ctx context.Context, containerID string) (Container, error)

	// Stats and monitoring
	StreamStats(ctx context.Context, containerID string) (<-chan Stats, error)

	// Logs
	StreamLogs(ctx context.Context, containerID string, tail int, follow bool) (io.ReadCloser, error)

	// Docker Compose operations
	ComposeStop(ctx context.Context, projectPath string) error
	ComposeRestart(ctx context.Context, projectPath string) error
	ComposeDown(ctx context.Context, projectPath string, removeVolumes bool) error

	// Project discovery
	DiscoverProjects(basePath string) ([]Project, error)

	// Connection check
	Ping(ctx context.Context) error
}

// Container represents container metadata
type Container struct {
	ID        string
	Name      string
	Service   string
	Image     string
	Status    string
	State     string
	Ports     []PortMapping
	CreatedAt int64
	StartedAt int64
}

// PortMapping represents container port configuration
type PortMapping struct {
	ContainerPort int
	HostPort      int
	Protocol      string
}

// Stats represents container resource usage statistics
type Stats struct {
	CPUPercent     float64
	MemoryUsed     uint64
	MemoryLimit    uint64
	MemoryPercent  float64
	NetworkRxBytes uint64
	NetworkTxBytes uint64
	Timestamp      int64
}

// Project represents a Docker Compose project
type Project struct {
	Name     string
	Path     string
	IsActive bool
}
