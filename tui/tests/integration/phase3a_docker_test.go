// Project: 20i Stack Manager TUI
// File: phase3a_docker_test.go
// Purpose: Integration tests for Phase 3a with LIVE Docker (requires Docker daemon)
// Version: 0.1.0
// Updated: 2025-12-28

//go:build integration_docker
// +build integration_docker

package integration_test

import (
	"context"
	"testing"

	"github.com/peternicholls/20i-stack/tui/internal/app"
	"github.com/peternicholls/20i-stack/tui/internal/docker"
)

// TestPhase3aWithLiveDocker tests Phase 3a functionality with a live Docker daemon
// Run with: go test -tags=integration_docker ./tests/integration/...
func TestPhase3aWithLiveDocker(t *testing.T) {
	// Skip if Docker is not available
	ctx := context.Background()
	client, err := docker.NewClient(ctx)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	if client == nil {
		t.Skip("Docker client is nil")
	}
	
	// Verify we can ping Docker
	pingErr := client.Ping(ctx)
	if pingErr != nil {
		t.Skipf("Cannot ping Docker daemon: %v", pingErr)
	}
	
	t.Run("live_docker_initialization", func(t *testing.T) {
		model, err := app.NewRootModel(ctx)
		if err != nil {
			t.Fatalf("Failed to initialize RootModel with live Docker: %v", err)
		}
		
		if model == nil {
			t.Fatal("Expected non-nil model")
		}
		
		t.Log("✓ RootModel initialized successfully with live Docker")
	})
	
	t.Run("live_docker_container_list", func(t *testing.T) {
		// List containers (may be empty if none running)
		containers, err := client.ListContainers("")
		if err != nil {
			t.Fatalf("Failed to list containers: %v", err)
		}
		
		t.Logf("✓ Successfully listed %d containers", len(containers))
	})
	
	t.Log("✅ Live Docker integration tests complete")
	t.Log("ℹ️  These tests verify basic Docker connectivity")
	t.Log("ℹ️  For full stack testing, ensure Docker Compose stack is running")
}
