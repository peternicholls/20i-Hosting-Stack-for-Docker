// Project: 20i Stack Manager TUI
// File: compose.go
// Purpose: Docker Compose stack lifecycle operations
// Version: 0.1.0
// Updated: 2025-12-28

package stack

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/peternicholls/20i-stack/tui/internal/project"
)

// Default buffer size for compose output channel.
// Large enough to handle typical docker compose output without blocking.
const composeOutputBufferSize = 100

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

// ComposeUpStreaming starts the Docker Compose stack and streams output line-by-line.
// It executes `docker compose -f $STACK_FILE up -d` with stdout/stderr pipes
// and sends each output line through a buffered channel.
//
// The channel is ALWAYS closed when the function returns, on both success and error paths.
// The channel is buffered to prevent deadlocks if the consumer is slow.
//
// Parameters:
//   - stackFile: Path to the docker-compose.yml file
//   - codeDir: Project code directory (defaults to current directory if empty)
//
// Returns:
//   - <-chan string: Receive-only channel of output lines (closed on completion)
//   - error: Validation error (if any) - errors during execution are sent as output lines
func ComposeUpStreaming(stackFile, codeDir string) (<-chan string, error) {
	// Validate STACK_FILE before starting
	if err := ValidateStackFile(stackFile); err != nil {
		return nil, err
	}

	// Get effective code directory
	var err error
	codeDir, err = getEffectiveCodeDir(codeDir)
	if err != nil {
		return nil, err
	}

	// Create buffered channel to prevent deadlock if consumer is slow
	outputChan := make(chan string, composeOutputBufferSize)

	// Start goroutine to execute command and stream output
	go func() {
		defer close(outputChan) // ALWAYS close channel on completion

		// Build command
		cmd := exec.Command("docker", "compose", "-f", stackFile, "up", "-d")
		cmd.Env = buildComposeEnv(codeDir)

		// Create pipes for stdout and stderr
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			outputChan <- fmt.Sprintf("ERROR: Failed to create stdout pipe: %v", err)
			return
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			// Close stdout pipe to avoid resource leak if stderr pipe creation fails
			stdout.Close()
			outputChan <- fmt.Sprintf("ERROR: Failed to create stderr pipe: %v", err)
			return
		}

		// Start command
		if err := cmd.Start(); err != nil {
			outputChan <- fmt.Sprintf("ERROR: Failed to start command: %v", err)
			return
		}

		// Stream output from both stdout and stderr
		done := make(chan struct{})
		go streamLines(stdout, outputChan, done)
		go streamLines(stderr, outputChan, done)

		// Wait for both streams to complete
		<-done
		<-done

		// Wait for command to finish
		if err := cmd.Wait(); err != nil {
			outputChan <- fmt.Sprintf("ERROR: Command failed: %v", err)
		} else {
			outputChan <- "[Complete]"
		}
	}()

	return outputChan, nil
}

// ComposeDownStreaming stops the stack and streams output line-by-line.
// It executes `docker compose -f $STACK_FILE down` with stdout/stderr pipes
// and sends each output line through a buffered channel.
//
// Parameters:
//   - stackFile: Path to the docker-compose.yml file
//   - codeDir: Project code directory (defaults to current directory if empty)
//
// Returns:
//   - <-chan string: Receive-only channel of output lines (closed on completion)
//   - error: Validation error (if any) - errors during execution are sent as output lines
func ComposeDownStreaming(stackFile, codeDir string) (<-chan string, error) {
	return composeStreamingOperation(stackFile, codeDir, "down")
}

// ComposeRestartStreaming restarts the stack and streams output line-by-line.
// It executes `docker compose -f $STACK_FILE restart` with stdout/stderr pipes
// and sends each output line through a buffered channel.
//
// Parameters:
//   - stackFile: Path to the docker-compose.yml file
//   - codeDir: Project code directory (defaults to current directory if empty)
//
// Returns:
//   - <-chan string: Receive-only channel of output lines (closed on completion)
//   - error: Validation error (if any) - errors during execution are sent as output lines
func ComposeRestartStreaming(stackFile, codeDir string) (<-chan string, error) {
	return composeStreamingOperation(stackFile, codeDir, "restart")
}

// ComposeDestroyStreaming destroys the stack (including volumes) and streams output line-by-line.
// It executes `docker compose -f $STACK_FILE down -v` with stdout/stderr pipes
// and sends each output line through a buffered channel.
// WARNING: This will delete all data in volumes!
//
// Parameters:
//   - stackFile: Path to the docker-compose.yml file
//   - codeDir: Project code directory (defaults to current directory if empty)
//
// Returns:
//   - <-chan string: Receive-only channel of output lines (closed on completion)
//   - error: Validation error (if any) - errors during execution are sent as output lines
func ComposeDestroyStreaming(stackFile, codeDir string) (<-chan string, error) {
	return composeStreamingOperation(stackFile, codeDir, "destroy")
}

// composeStreamingOperation is a generic function for streaming compose operations.
// It handles the common pattern of running docker compose with streaming output.
//
// Parameters:
//   - stackFile: Path to the docker-compose.yml file
//   - codeDir: Project code directory (defaults to current directory if empty)
//   - operation: The operation to perform ("down", "restart", "destroy")
//
// Returns:
//   - <-chan string: Receive-only channel of output lines (closed on completion)
//   - error: Validation error (if any) - errors during execution are sent as output lines
func composeStreamingOperation(stackFile, codeDir, operation string) (<-chan string, error) {
	// Validate STACK_FILE before starting
	if err := ValidateStackFile(stackFile); err != nil {
		return nil, err
	}

	// Get effective code directory
	var err error
	codeDir, err = getEffectiveCodeDir(codeDir)
	if err != nil {
		return nil, err
	}

	// Create buffered channel to prevent deadlock if consumer is slow
	outputChan := make(chan string, composeOutputBufferSize)

	// Start goroutine to execute command and stream output
	go func() {
		defer close(outputChan) // ALWAYS close channel on completion

		// Build command based on operation
		var cmd *exec.Cmd
		switch operation {
		case "down":
			cmd = exec.Command("docker", "compose", "-f", stackFile, "down")
		case "restart":
			cmd = exec.Command("docker", "compose", "-f", stackFile, "restart")
		case "destroy":
			cmd = exec.Command("docker", "compose", "-f", stackFile, "down", "-v")
		default:
			outputChan <- fmt.Sprintf("ERROR: Unknown operation: %s", operation)
			return
		}

		cmd.Env = buildComposeEnv(codeDir)

		// Create pipes for stdout and stderr
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			outputChan <- fmt.Sprintf("ERROR: Failed to create stdout pipe: %v", err)
			return
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			// Close stdout pipe to avoid resource leak if stderr pipe creation fails
			stdout.Close()
			outputChan <- fmt.Sprintf("ERROR: Failed to create stderr pipe: %v", err)
			return
		}

		// Start command
		if err := cmd.Start(); err != nil {
			outputChan <- fmt.Sprintf("ERROR: Failed to start command: %v", err)
			return
		}

		// Stream output from both stdout and stderr
		done := make(chan struct{})
		go streamLines(stdout, outputChan, done)
		go streamLines(stderr, outputChan, done)

		// Wait for both streams to complete
		<-done
		<-done

		// Wait for command to finish
		if err := cmd.Wait(); err != nil {
			outputChan <- fmt.Sprintf("ERROR: Command failed: %v", err)
		} else {
			outputChan <- "[Complete]"
		}
	}()

	return outputChan, nil
}

// streamLines reads lines from a reader and sends them to the output channel.
// It signals completion on the done channel when the reader is exhausted.
func streamLines(r io.Reader, out chan<- string, done chan<- struct{}) {
	defer func() { done <- struct{}{} }()

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		out <- scanner.Text()
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		out <- fmt.Sprintf("ERROR: Stream read error: %v", err)
	}
}
