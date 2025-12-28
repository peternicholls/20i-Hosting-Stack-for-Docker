package docker
// Project: 20i Stack Manager TUI
// File: client.go
// Purpose: Docker client for managing containers
// Version: 0.1.0
// Updated: 2025-12-28

package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// Client struct for Docker operations
type Client struct {
	dockerClient *client.Client
}

// NewClient initializes a new Docker client
func NewClient() (*Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.41"))
	if err != nil {
		return nil, err
	}
	return &Client{dockerClient: cli}, nil
}