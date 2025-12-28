// Project: 20i Stack Manager TUI
// File: compose.go
// Purpose: Docker Compose stack lifecycle operations
// Version: 0.1.0
// Updated: 2025-12-28

package stack

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/peternicholls/20i-stack/tui/internal/project"
)

// ComposeResult holds the result of a compose operation.
type ComposeResult struct {
	Success bool
	Output  string
	Error   error
}

// getEffectiveCodeDir returns the effective code directory, defaulting to current directory if empty.
func getEffectiveCodeDir(codeDir string) (string, error) {
	if codeDir == "" {
		var err error
		codeDir, err = os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get current directory: %w", err)
		}
	}
	return codeDir, nil
}

// buildComposeProjectName builds the COMPOSE_PROJECT_NAME from the sanitized directory name.
func buildComposeProjectName(codeDir string) string {
	projectName := filepath.Base(codeDir)
	return project.SanitizeProjectName(projectName)
}

// buildComposeEnv builds the environment variables for docker compose commands.
func buildComposeEnv(codeDir string) []string {
	env := os.Environ()
	env = append(env, fmt.Sprintf("CODE_DIR=%s", codeDir))
	env = append(env, fmt.Sprintf("COMPOSE_PROJECT_NAME=%s", buildComposeProjectName(codeDir)))
	return env
}

// ComposeUp starts the Docker Compose stack in detached mode.
// It validates STACK_FILE, builds the environment with CODE_DIR and COMPOSE_PROJECT_NAME,
// and executes `docker compose -f $STACK_FILE up -d`.
//
// Parameters:
//   - stackFile: Path to the docker-compose.yml file
//   - codeDir: Project code directory (defaults to current directory if empty)
//
// Returns:
//   - ComposeResult with operation status and output
func ComposeUp(stackFile, codeDir string) *ComposeResult {
	// Validate STACK_FILE
	if err := ValidateStackFile(stackFile); err != nil {
		return &ComposeResult{
			Success: false,
			Error:   err,
		}
	}

	// Get effective code directory
	var err error
	codeDir, err = getEffectiveCodeDir(codeDir)
	if err != nil {
		return &ComposeResult{
			Success: false,
			Error:   err,
		}
	}

	// Build command
	cmd := exec.Command("docker", "compose", "-f", stackFile, "up", "-d")
	
	// Set environment variables
	cmd.Env = buildComposeEnv(codeDir)

	// Execute command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return &ComposeResult{
			Success: false,
			Output:  string(output),
			Error:   fmt.Errorf("docker compose up failed: %w", err),
		}
	}

	return &ComposeResult{
		Success: true,
		Output:  string(output),
	}
}

// ComposeDown stops and removes containers, networks, and default volumes.
// It executes `docker compose -f $STACK_FILE down`.
//
// Parameters:
//   - stackFile: Path to the docker-compose.yml file
//   - codeDir: Project code directory (for environment variables)
//
// Returns:
//   - ComposeResult with operation status and output
func ComposeDown(stackFile, codeDir string) *ComposeResult {
	// Validate STACK_FILE
	if err := ValidateStackFile(stackFile); err != nil {
		return &ComposeResult{
			Success: false,
			Error:   err,
		}
	}

	// Get effective code directory
	var err error
	codeDir, err = getEffectiveCodeDir(codeDir)
	if err != nil {
		return &ComposeResult{
			Success: false,
			Error:   err,
		}
	}

	// Build command
	cmd := exec.Command("docker", "compose", "-f", stackFile, "down")
	
	// Set environment variables
	cmd.Env = buildComposeEnv(codeDir)

	// Execute command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return &ComposeResult{
			Success: false,
			Output:  string(output),
			Error:   fmt.Errorf("docker compose down failed: %w", err),
		}
	}

	return &ComposeResult{
		Success: true,
		Output:  string(output),
	}
}

// ComposeRestart restarts all services in the stack.
// It executes `docker compose -f $STACK_FILE restart`.
//
// Parameters:
//   - stackFile: Path to the docker-compose.yml file
//   - codeDir: Project code directory (for environment variables)
//
// Returns:
//   - ComposeResult with operation status and output
func ComposeRestart(stackFile, codeDir string) *ComposeResult {
	// Validate STACK_FILE
	if err := ValidateStackFile(stackFile); err != nil {
		return &ComposeResult{
			Success: false,
			Error:   err,
		}
	}

	// Default codeDir to current directory if not provided
	var err error
	codeDir, err = getEffectiveCodeDir(codeDir)
	if err != nil {
		return &ComposeResult{
			Success: false,
			Error:   err,
		}
	}

	// Build command
	cmd := exec.Command("docker", "compose", "-f", stackFile, "restart")
	
	// Set environment variables
	cmd.Env = buildComposeEnv(codeDir)

	// Execute command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return &ComposeResult{
			Success: false,
			Output:  string(output),
			Error:   fmt.Errorf("docker compose restart failed: %w", err),
		}
	}

	return &ComposeResult{
		Success: true,
		Output:  string(output),
	}
}

// ComposeDestroy stops and removes containers, networks, and ALL volumes.
// It executes `docker compose -f $STACK_FILE down -v`.
// WARNING: This will delete all data in volumes!
//
// Parameters:
//   - stackFile: Path to the docker-compose.yml file
//   - codeDir: Project code directory (for environment variables)
//
// Returns:
//   - ComposeResult with operation status and output
func ComposeDestroy(stackFile, codeDir string) *ComposeResult {
	// Validate STACK_FILE
	if err := ValidateStackFile(stackFile); err != nil {
		return &ComposeResult{
			Success: false,
			Error:   err,
		}
	}

	// Default codeDir to current directory if not provided
	var err error
	codeDir, err = getEffectiveCodeDir(codeDir)
	if err != nil {
		return &ComposeResult{
			Success: false,
			Error:   err,
		}
	}

	// Build command
	cmd := exec.Command("docker", "compose", "-f", stackFile, "down", "-v")
	
	// Set environment variables
	cmd.Env = buildComposeEnv(codeDir)

	// Execute command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return &ComposeResult{
			Success: false,
			Output:  string(output),
			Error:   fmt.Errorf("docker compose down -v failed: %w", err),
		}
	}

	return &ComposeResult{
		Success: true,
		Output:  string(output),
	}
}
