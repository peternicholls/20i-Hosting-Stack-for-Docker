// Project: 20i Stack Manager TUI
// File: client.go
// Purpose: Docker client wrapper for container lifecycle and compose operations
// Version: 0.1.0
// Updated: 2025-12-28

package docker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	dockerclient "github.com/docker/docker/client"
)

var execCommand = exec.CommandContext

// Sentinel errors returned by the Docker client wrapper.
var (
	// ErrDaemonUnreachable means we couldn't reach the Docker daemon (not running or socket missing)
	ErrDaemonUnreachable = errors.New("docker daemon unreachable")
	// ErrPermissionDenied means we don't have permission to access the Docker socket
	ErrPermissionDenied = errors.New("permission denied to access docker")
	// ErrTimeout indicates the operation timed out
	ErrTimeout = errors.New("docker operation timed out")
	// ErrNotFound indicates requested resource was not found
	ErrNotFound = errors.New("not found")
	// ErrConflict indicates a conflict such as port allocation
	ErrConflict = errors.New("conflict")
	// ErrUnknown is a catch-all for unclassified errors
	ErrUnknown = errors.New("unknown docker error")
)

// ContainerStatus represents the normalized status of a container for UI rendering.
type ContainerStatus string

const (
	// StatusRunning indicates a container is actively running.
	StatusRunning ContainerStatus = "running"
	// StatusStopped indicates a container exists but is not running.
	StatusStopped ContainerStatus = "stopped"
	// StatusRestarting indicates a container is restarting.
	StatusRestarting ContainerStatus = "restarting"
	// StatusError indicates a container is in a failed or unknown state.
	StatusError ContainerStatus = "error"
)

// Container represents a Docker container for Phase 3 lifecycle operations.
// Extended in Phase 5 with Ports, CreatedAt, StartedAt.
type Container struct {
	ID         string
	Name       string
	Service    string
	Image      string
	Status     ContainerStatus
	State      string
	CPUPercent float64 // CPU usage percentage (0-100)
}

type dockerAPI interface {
	Ping(ctx context.Context) (types.Ping, error)
	ContainerList(ctx context.Context, options container.ListOptions) ([]container.Summary, error)
	ContainerStart(ctx context.Context, containerID string, options container.StartOptions) error
	ContainerStop(ctx context.Context, containerID string, options container.StopOptions) error
	ContainerRestart(ctx context.Context, containerID string, options container.StopOptions) error
	ContainerStats(ctx context.Context, containerID string, stream bool) (container.StatsResponseReader, error)
}

// Client is a thin wrapper around the official Docker SDK client
// that implements the minimal contract required by the TUI.
type Client struct {
	cli dockerAPI
	ctx context.Context
}

// NewClient creates a new Docker client, negotiates the API version and
// verifies connectivity to the daemon. It returns wrapped, user-friendly
// sentinel errors defined in this package (e.g., ErrDaemonUnreachable).
func NewClient(ctx context.Context) (*Client, error) {
	cli, err := dockerclient.NewClientWithOpts(dockerclient.FromEnv, dockerclient.WithAPIVersionNegotiation())
	if err != nil {
		// If creating the SDK client fails, return an ErrUnknown wrapped error
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}

	// Use a short timeout when checking connectivity so startup doesn't hang
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err = cli.Ping(pingCtx); err != nil {
		return nil, mapConnectionError(err)
	}

	return &Client{cli: cli, ctx: ctx}, nil
}

// Ping checks connectivity with the Docker daemon using a short timeout and
// maps SDK errors to package sentinel errors.
func (c *Client) Ping(ctx context.Context) error {
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if _, err := c.cli.Ping(pingCtx); err != nil {
		return mapConnectionError(err)
	}
	return nil
}

// ListContainers returns all containers for the provided Docker Compose project.
func (c *Client) ListContainers(projectName string) ([]Container, error) {
	listCtx, cancel := context.WithTimeout(c.ctx, 5*time.Second)
	defer cancel()

	args := filters.NewArgs()
	if strings.TrimSpace(projectName) != "" {
		args.Add("label", fmt.Sprintf("com.docker.compose.project=%s", projectName))
	}

	containers, err := c.cli.ContainerList(listCtx, container.ListOptions{
		All:     true,
		Filters: args,
	})
	if err != nil {
		return nil, mapOperationError(err)
	}

	results := make([]Container, 0, len(containers))
	for _, summary := range containers {
		results = append(results, Container{
			ID:         summary.ID,
			Name:       containerName(summary.Names),
			Service:    containerService(summary),
			Image:      summary.Image,
			Status:     mapDockerState(string(summary.State)),
			State:      summary.Status,
			CPUPercent: 0, // CPU not fetched in basic listing
		})
	}

	return results, nil
}

// ListContainersWithStats returns all containers with CPU statistics for the provided Docker Compose project.
// This method is used for auto-refresh to get updated status and CPU usage.
func (c *Client) ListContainersWithStats(projectName string) ([]Container, error) {
	// First, get the basic container list
	containers, err := c.ListContainers(projectName)
	if err != nil {
		return nil, err
	}

	// Then fetch CPU stats for running containers
	for i := range containers {
		if containers[i].Status == StatusRunning {
			cpuPercent, err := c.GetContainerStats(containers[i].ID)
			if err == nil {
				containers[i].CPUPercent = cpuPercent
			}
			// Ignore errors for individual stats - continue with 0% CPU
		}
	}

	return containers, nil
}

// StartContainer starts a stopped container by ID.
func (c *Client) StartContainer(containerID string) error {
	startCtx, cancel := context.WithTimeout(c.ctx, 5*time.Second)
	defer cancel()

	if err := c.cli.ContainerStart(startCtx, containerID, container.StartOptions{}); err != nil {
		return mapOperationError(err)
	}

	return nil
}

// StopContainer stops a running container, waiting up to timeout seconds.
func (c *Client) StopContainer(containerID string, timeout int) error {
	timeoutSeconds := normalizeTimeout(timeout)
	stopCtx, cancel := context.WithTimeout(c.ctx, time.Duration(timeoutSeconds+5)*time.Second)
	defer cancel()

	if err := c.cli.ContainerStop(stopCtx, containerID, container.StopOptions{Timeout: &timeoutSeconds}); err != nil {
		return mapOperationError(err)
	}

	return nil
}

// RestartContainer restarts a container, waiting up to timeout seconds for stop.
func (c *Client) RestartContainer(containerID string, timeout int) error {
	timeoutSeconds := normalizeTimeout(timeout)
	restartCtx, cancel := context.WithTimeout(c.ctx, time.Duration(timeoutSeconds+5)*time.Second)
	defer cancel()

	if err := c.cli.ContainerRestart(restartCtx, containerID, container.StopOptions{Timeout: &timeoutSeconds}); err != nil {
		return mapOperationError(err)
	}

	return nil
}

// GetContainerStats retrieves CPU usage for a container using one-shot stats API.
// Returns CPU percentage (0-100) or 0 if stats unavailable.
func (c *Client) GetContainerStats(containerID string) (float64, error) {
	statsCtx, cancel := context.WithTimeout(c.ctx, 3*time.Second)
	defer cancel()

	resp, err := c.cli.ContainerStats(statsCtx, containerID, false) // stream=false for one-shot
	if err != nil {
		return 0, mapOperationError(err)
	}
	defer resp.Body.Close()

	var stats container.StatsResponse
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return 0, fmt.Errorf("failed to decode stats: %w", err)
	}

	// Calculate CPU percentage
	cpuPercent := calculateCPUPercent(&stats)
	return cpuPercent, nil
}

// ComposeStop stops all containers in a Docker Compose project.
func (c *Client) ComposeStop(projectPath string) error {
	return c.runComposeCommand(projectPath, 30*time.Second, "stop")
}

// ComposeRestart restarts all containers in a Docker Compose project.
func (c *Client) ComposeRestart(projectPath string) error {
	return c.runComposeCommand(projectPath, 60*time.Second, "restart")
}

// ComposeDown stops and removes containers, networks, and optionally volumes.
func (c *Client) ComposeDown(projectPath string, removeVolumes bool) error {
	args := []string{"down"}
	if removeVolumes {
		args = append(args, "-v")
	}
	return c.runComposeCommand(projectPath, 60*time.Second, args...)
}

// mapConnectionError maps common errors returned by the Docker SDK into
// the sentinel error values defined in this package so callers can react
// in a user-friendly way according to the contract.
func mapConnectionError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return ErrTimeout
	}
	if errors.Is(err, context.Canceled) {
		return ErrTimeout
	}
	msg := strings.ToLower(err.Error())

	switch {
	case strings.Contains(msg, "permission denied"):
		return ErrPermissionDenied
	case strings.Contains(msg, "cannot connect to the docker daemon"):
		return ErrDaemonUnreachable
	case strings.Contains(msg, "open //./pipe/docker_engine"):
		return ErrDaemonUnreachable
	case strings.Contains(msg, "the system cannot find the file specified"):
		return ErrDaemonUnreachable
	case strings.Contains(msg, "connect: connection refused"):
		return ErrDaemonUnreachable
	case strings.Contains(msg, "no such file or directory"):
		// common when socket not present on unix systems
		return ErrDaemonUnreachable
	case strings.Contains(msg, "i/o timeout") || strings.Contains(msg, "timeout"):
		return ErrTimeout
	case strings.Contains(msg, "not found"):
		return ErrNotFound
	case strings.Contains(msg, "port is already allocated") || strings.Contains(msg, "address already in use"):
		return ErrConflict
	default:
		// Preserve original error text while signalling unknown error category
		return fmt.Errorf("%w: %s", ErrUnknown, err)
	}
}

func mapOperationError(err error) error {
	if err == nil {
		return nil
	}

	mapped := mapConnectionError(err)
	if errors.Is(mapped, ErrDaemonUnreachable) || errors.Is(mapped, ErrPermissionDenied) || errors.Is(mapped, ErrTimeout) {
		return mapped
	}
	return err
}

func normalizeTimeout(timeout int) int {
	if timeout <= 0 {
		return 10
	}
	return timeout
}

func (c *Client) runComposeCommand(projectPath string, timeout time.Duration, args ...string) error {
	composePath, err := resolveComposePath(projectPath)
	if err != nil {
		return err
	}

	cmdArgs := append([]string{"compose"}, args...)
	ctx, cancel := context.WithTimeout(c.ctx, timeout)
	defer cancel()

	cmd := execCommand(ctx, "docker", cmdArgs...)
	cmd.Dir = filepath.Dir(composePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		message := strings.TrimSpace(string(output))
		if message != "" {
			return fmt.Errorf("docker compose %s failed: %w: %s", strings.Join(args, " "), err, message)
		}
		return fmt.Errorf("docker compose %s failed: %w", strings.Join(args, " "), err)
	}

	return nil
}

func resolveComposePath(projectPath string) (string, error) {
	if strings.TrimSpace(projectPath) == "" {
		return "", fmt.Errorf("project path is required")
	}

	composePath := filepath.Join(projectPath, "docker-compose.yml")
	if _, err := os.Stat(composePath); err != nil {
		return "", fmt.Errorf("compose file not found at %s: %w", composePath, err)
	}

	return composePath, nil
}

func mapDockerState(state string) ContainerStatus {
	switch strings.ToLower(strings.TrimSpace(state)) {
	case "running":
		return StatusRunning
	case "restarting":
		return StatusRestarting
	case "exited", "created", "paused", "dead":
		return StatusStopped
	default:
		return StatusError
	}
}

func containerName(names []string) string {
	for _, name := range names {
		trimmed := strings.TrimPrefix(name, "/")
		if trimmed != "" {
			return trimmed
		}
	}
	return ""
}

func containerService(summary container.Summary) string {
	if summary.Labels != nil {
		if service, ok := summary.Labels["com.docker.compose.service"]; ok && service != "" {
			return service
		}
	}

	name := containerName(summary.Names)
	if name != "" {
		return name
	}

	return summary.ID
}

// calculateCPUPercent computes CPU usage percentage from Docker stats.
// Formula: ((CPUDelta / SystemCPUDelta) * NumCPUs) * 100
func calculateCPUPercent(stats *container.StatsResponse) float64 {
	cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage) - float64(stats.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(stats.CPUStats.SystemUsage) - float64(stats.PreCPUStats.SystemUsage)

	if systemDelta > 0 && cpuDelta > 0 {
		numCPUs := float64(len(stats.CPUStats.CPUUsage.PercpuUsage))
		if numCPUs == 0 {
			numCPUs = 1 // Fallback to 1 CPU if not available
		}
		return (cpuDelta / systemDelta) * numCPUs * 100.0
	}

	return 0
}
