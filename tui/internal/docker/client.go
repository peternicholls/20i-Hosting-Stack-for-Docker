// Project: 20i Stack Manager TUI
// File: client.go
// Purpose: Docker client wrapper: Client struct and NewClient() initializer
// Version: 0.1.0
// Updated: 2025-12-28

package docker

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	dockerclient "github.com/docker/docker/client"
)

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

// Client is a thin wrapper around the official Docker SDK client
// that implements the minimal contract required by the TUI.
type Client struct {
	cli *dockerclient.Client
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
