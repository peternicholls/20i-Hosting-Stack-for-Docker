// Project: 20i Stack Manager TUI
// File: status.go
// Purpose: Retrieve stack status with container info, URLs, and CPU%
// Version: 0.1.0
// Updated: 2025-12-28

package stack

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	dockerclient "github.com/docker/docker/client"
)

// ContainerInfo holds runtime information about a container in the stack.
type ContainerInfo struct {
	ID      string  // Container ID
	Name    string  // Container name
	Service string  // Service name from compose label
	Image   string  // Image name
	Status  string  // Container status (running, stopped, etc.)
	URL     string  // Access URL (if applicable)
	CPUPerc float64 // CPU usage percentage
}

// GetStackStatus retrieves the status of all containers in a Docker Compose stack.
// It lists containers filtered by the com.docker.compose.project label,
// generates access URLs based on port bindings, and collects CPU statistics.
//
// Parameters:
//   - projectName: The Docker Compose project name (from COMPOSE_PROJECT_NAME)
//
// Returns:
//   - []ContainerInfo: List of container information
//   - error: Any error encountered during retrieval
func GetStackStatus(projectName string) ([]ContainerInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create Docker client
	cli, err := dockerclient.NewClientWithOpts(dockerclient.FromEnv, dockerclient.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}

	// Build filter for project
	args := filters.NewArgs()
	if strings.TrimSpace(projectName) != "" {
		args.Add("label", fmt.Sprintf("com.docker.compose.project=%s", projectName))
	}

	// List containers
	containers, err := cli.ContainerList(ctx, container.ListOptions{
		All:     true,
		Filters: args,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	// Build container info list
	results := make([]ContainerInfo, 0, len(containers))
	for _, c := range containers {
		info := ContainerInfo{
			ID:      c.ID,
			Name:    extractContainerName(c.Names),
			Service: extractServiceName(c),
			Image:   c.Image,
			Status:  c.Status,
			URL:     generateURL(c),
			CPUPerc: 0.0, // Will be populated by stats
		}

		// Get CPU stats if container is running
		if c.State == "running" {
			info.CPUPerc = getCPUPercentage(ctx, cli, c.ID)
		}

		results = append(results, info)
	}

	return results, nil
}

// extractContainerName extracts the container name from the names list.
func extractContainerName(names []string) string {
	for _, name := range names {
		trimmed := strings.TrimPrefix(name, "/")
		if trimmed != "" {
			return trimmed
		}
	}
	return ""
}

// extractServiceName extracts the service name from container labels.
func extractServiceName(c container.Summary) string {
	if c.Labels != nil {
		if service, ok := c.Labels["com.docker.compose.service"]; ok && service != "" {
			return service
		}
	}
	return extractContainerName(c.Names)
}

// generateURL generates the access URL for a container based on its service type and port bindings.
// Ports are resolved at runtime with precedence: Docker bindings > env vars > defaults.
//
// URL rules:
//   - nginx: http://localhost:{HOST_PORT} (default 80)
//   - phpmyadmin: http://localhost:{PMA_PORT} (default 8081)
//   - mariadb: localhost:{MYSQL_PORT} (no http, default 3306)
//   - apache: "internal" (proxied via nginx)
func generateURL(c container.Summary) string {
	service := extractServiceName(c)

	switch service {
	case "nginx":
		port := resolvePort(c.Ports, 80, "HOST_PORT", "80")
		return fmt.Sprintf("http://localhost:%s", port)

	case "phpmyadmin":
		port := resolvePort(c.Ports, 80, "PMA_PORT", "8081")
		return fmt.Sprintf("http://localhost:%s", port)

	case "mariadb":
		port := resolvePort(c.Ports, 3306, "MYSQL_PORT", "3306")
		return fmt.Sprintf("localhost:%s", port)

	case "apache":
		return "internal"

	default:
		return ""
	}
}

// resolvePort resolves the host port for a container with precedence:
// 1. Docker port bindings (from actual container state)
// 2. Environment variable
// 3. Default value
func resolvePort(ports []types.Port, containerPort uint16, envVar, defaultPort string) string {
	// First, check Docker port bindings
	for _, port := range ports {
		if port.PrivatePort == containerPort && port.PublicPort > 0 {
			return strconv.Itoa(int(port.PublicPort))
		}
	}

	// Second, check environment variable
	if envVal := os.Getenv(envVar); envVal != "" {
		return envVal
	}

	// Third, use default
	return defaultPort
}

// getCPUPercentage retrieves the current CPU usage percentage for a container.
func getCPUPercentage(ctx context.Context, cli *dockerclient.Client, containerID string) float64 {
	// Use a short timeout for stats collection
	statsCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	stats, err := cli.ContainerStats(statsCtx, containerID, false)
	if err != nil {
		return 0.0
	}
	defer stats.Body.Close()

	// Decode stats
	var v container.StatsResponse
	decoder := json.NewDecoder(stats.Body)
	if err := decoder.Decode(&v); err != nil {
		return 0.0
	}

	// Calculate CPU percentage
	return calculateCPUPercent(&v)
}

// calculateCPUPercent calculates CPU usage percentage from Docker stats.
// This uses the same formula as `docker stats` command.
func calculateCPUPercent(stats *container.StatsResponse) float64 {
	if stats == nil {
		return 0.0
	}

	// Calculate the change in CPU usage
	cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage - stats.PreCPUStats.CPUUsage.TotalUsage)

	// Calculate the change in system CPU usage
	systemDelta := float64(stats.CPUStats.SystemUsage - stats.PreCPUStats.SystemUsage)

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		// Calculate percentage and multiply by number of CPUs
		cpuPercent := (cpuDelta / systemDelta) * float64(len(stats.CPUStats.CPUUsage.PercpuUsage)) * 100.0
		return cpuPercent
	}

	return 0.0
}
