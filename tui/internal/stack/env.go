// Project: 20i Stack Manager TUI
// File: env.go
// Purpose: Stack environment detection (STACK_FILE and STACK_HOME)
// Version: 0.1.0
// Updated: 2025-12-28

// Package stack provides functionality for managing Docker Compose stacks.
// It handles stack environment detection, validation, and lifecycle operations
// (up, down, restart, destroy).
package stack

import (
	"fmt"
	"os"
	"path/filepath"
)

// StackEnv holds the detected stack file and home directory paths.
type StackEnv struct {
	StackFile string // Absolute path to docker-compose.yml
	StackHome string // Directory containing the stack file
}

// DetectStackEnv detects STACK_FILE and STACK_HOME from environment or executable location.
// This matches the 20i-gui logic exactly (lines 7-17 of 20i-gui script):
//
// 1. If STACK_FILE is set in environment, use it
// 2. If STACK_HOME is set, derive STACK_FILE from it
// 3. Otherwise, use executable directory as STACK_HOME
//
// Returns:
//   - StackEnv with StackFile and StackHome paths
//   - Error if paths cannot be determined
func DetectStackEnv() (*StackEnv, error) {
	var stackFile string
	var stackHome string

	// Check if STACK_FILE is explicitly set
	if envStackFile := os.Getenv("STACK_FILE"); envStackFile != "" {
		stackFile = envStackFile
		// Derive STACK_HOME from STACK_FILE if not set
		if envStackHome := os.Getenv("STACK_HOME"); envStackHome != "" {
			stackHome = envStackHome
		} else {
			// Get directory containing STACK_FILE
			dir := filepath.Dir(stackFile)
			absDir, err := filepath.Abs(dir)
			if err != nil {
				return nil, fmt.Errorf("failed to resolve STACK_HOME from STACK_FILE: %w", err)
			}
			stackHome = absDir
		}
	} else {
		// STACK_FILE not set, determine from executable or STACK_HOME
		if envStackHome := os.Getenv("STACK_HOME"); envStackHome != "" {
			stackHome = envStackHome
		} else {
			// Fall back to executable directory
			execPath, err := os.Executable()
			if err != nil {
				return nil, fmt.Errorf("failed to determine executable path: %w", err)
			}
			stackHome = filepath.Dir(execPath)
		}
		stackFile = filepath.Join(stackHome, "docker-compose.yml")
	}

	// Ensure paths are absolute
	absStackFile, err := filepath.Abs(stackFile)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve absolute path for STACK_FILE: %w", err)
	}

	absStackHome, err := filepath.Abs(stackHome)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve absolute path for STACK_HOME: %w", err)
	}

	return &StackEnv{
		StackFile: absStackFile,
		StackHome: absStackHome,
	}, nil
}

// ValidateStackFile verifies that the STACK_FILE exists and is readable.
// It returns a friendly error message if validation fails.
//
// Returns:
//   - nil if the stack file is valid
//   - error with helpful message if validation fails
func ValidateStackFile(stackFile string) error {
	if stackFile == "" {
		return fmt.Errorf("STACK_FILE not set and cannot be detected - please set STACK_FILE environment variable or run from stack directory")
	}

	info, err := os.Stat(stackFile)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("docker-compose file not found at %s - please verify STACK_FILE path", stackFile)
		}
		return fmt.Errorf("cannot access docker-compose file at %s: %w", stackFile, err)
	}

	if info.IsDir() {
		return fmt.Errorf("STACK_FILE points to a directory, not a file: %s", stackFile)
	}

	// Verify file is readable
	file, err := os.Open(stackFile)
	if err != nil {
		return fmt.Errorf("docker-compose file at %s is not readable: %w", stackFile, err)
	}
	file.Close()

	return nil
}
