// Project: 20i Stack Manager TUI
// File: client_test.go
// Purpose: Unit tests for Docker client wrapper (connection error mapping)
// Version: 0.1.0
// Updated: 2025-12-28

package docker

import (
	"context"
	"errors"
	"os"
	"testing"
)

func TestMapConnectionError(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want error
	}{
		{"permission denied", errors.New("permission denied"), ErrPermissionDenied},
		{"connection refused", errors.New("connect: connection refused"), ErrDaemonUnreachable},
		{"cannot connect to daemon", errors.New("Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?"), ErrDaemonUnreachable},
		{"windows pipe missing", errors.New("open //./pipe/docker_engine: The system cannot find the file specified."), ErrDaemonUnreachable},
		{"no such file", errors.New("no such file or directory"), ErrDaemonUnreachable},
		{"timeout", errors.New("i/o timeout while dialing"), ErrTimeout},
		{"deadline exceeded", context.DeadlineExceeded, ErrTimeout},
		{"not found", errors.New("not found"), ErrNotFound},
		{"port conflict", errors.New("address already in use"), ErrConflict},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := mapConnectionError(c.err)
			if !errors.Is(got, c.want) {
				t.Fatalf("mapConnectionError(%v) = %v, want %v", c.err, got, c.want)
			}
		})
	}
}

func TestNewClient_Integration(t *testing.T) {
	// Integration test: only run when explicitly requested
	if os.Getenv("DOCKER_INTEGRATION_TEST") != "1" {
		t.Skip("skipping integration test; set DOCKER_INTEGRATION_TEST=1 to run")
	}

	ctx := context.Background()
	c, err := NewClient(ctx)
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestNewClient_CanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := NewClient(ctx)
	if err == nil {
		t.Fatal("expected error when creating client with canceled context")
	}
}
