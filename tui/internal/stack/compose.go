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

	// Default codeDir to current directory if not provided
	if codeDir == "" {
		var err error
		codeDir, err = os.Getwd()
		if err != nil {
			return &ComposeResult{
				Success: false,
				Error:   fmt.Errorf("failed to get current directory: %w", err),
			}
		}
	}

	// Build COMPOSE_PROJECT_NAME from sanitized directory name
	projectName := filepath.Base(codeDir)
	sanitizedName := project.SanitizeProjectName(projectName)

	// Build command
	cmd := exec.Command("docker", "compose", "-f", stackFile, "up", "-d")
	
	// Set environment variables
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("CODE_DIR=%s", codeDir))
	cmd.Env = append(cmd.Env, fmt.Sprintf("COMPOSE_PROJECT_NAME=%s", sanitizedName))

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

	// Default codeDir to current directory if not provided
	if codeDir == "" {
		var err error
		codeDir, err = os.Getwd()
		if err != nil {
			return &ComposeResult{
				Success: false,
				Error:   fmt.Errorf("failed to get current directory: %w", err),
			}
		}
	}

	// Build COMPOSE_PROJECT_NAME from sanitized directory name
	projectName := filepath.Base(codeDir)
	sanitizedName := project.SanitizeProjectName(projectName)

	// Build command
	cmd := exec.Command("docker", "compose", "-f", stackFile, "down")
	
	// Set environment variables
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("CODE_DIR=%s", codeDir))
	cmd.Env = append(cmd.Env, fmt.Sprintf("COMPOSE_PROJECT_NAME=%s", sanitizedName))

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
	if codeDir == "" {
		var err error
		codeDir, err = os.Getwd()
		if err != nil {
			return &ComposeResult{
				Success: false,
				Error:   fmt.Errorf("failed to get current directory: %w", err),
			}
		}
	}

	// Build COMPOSE_PROJECT_NAME from sanitized directory name
	projectName := filepath.Base(codeDir)
	sanitizedName := project.SanitizeProjectName(projectName)

	// Build command
	cmd := exec.Command("docker", "compose", "-f", stackFile, "restart")
	
	// Set environment variables
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("CODE_DIR=%s", codeDir))
	cmd.Env = append(cmd.Env, fmt.Sprintf("COMPOSE_PROJECT_NAME=%s", sanitizedName))

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
	if codeDir == "" {
		var err error
		codeDir, err = os.Getwd()
		if err != nil {
			return &ComposeResult{
				Success: false,
				Error:   fmt.Errorf("failed to get current directory: %w", err),
			}
		}
	}

	// Build COMPOSE_PROJECT_NAME from sanitized directory name
	projectName := filepath.Base(codeDir)
	sanitizedName := project.SanitizeProjectName(projectName)

	// Build command
	cmd := exec.Command("docker", "compose", "-f", stackFile, "down", "-v")
	
	// Set environment variables
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("CODE_DIR=%s", codeDir))
	cmd.Env = append(cmd.Env, fmt.Sprintf("COMPOSE_PROJECT_NAME=%s", sanitizedName))

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
