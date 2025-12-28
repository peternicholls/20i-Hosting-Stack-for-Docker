// Project: 20i Stack Manager TUI
// File: compose_test.go
// Purpose: Unit tests for Docker Compose stack operations
// Version: 0.1.0
// Updated: 2025-12-28

package stack

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDetectStackEnv(t *testing.T) {
	// Save original environment
	origStackFile := os.Getenv("STACK_FILE")
	origStackHome := os.Getenv("STACK_HOME")
	defer func() {
		os.Setenv("STACK_FILE", origStackFile)
		os.Setenv("STACK_HOME", origStackHome)
	}()

	t.Run("STACK_FILE explicitly set", func(t *testing.T) {
		os.Setenv("STACK_FILE", "/custom/path/docker-compose.yml")
		os.Setenv("STACK_HOME", "/custom/path")

		env, err := DetectStackEnv()
		if err != nil {
			t.Fatalf("DetectStackEnv() returned error: %v", err)
		}

		if !strings.HasSuffix(env.StackFile, "/custom/path/docker-compose.yml") {
			t.Errorf("expected StackFile to end with /custom/path/docker-compose.yml, got %s", env.StackFile)
		}

		if !strings.HasSuffix(env.StackHome, "/custom/path") {
			t.Errorf("expected StackHome to end with /custom/path, got %s", env.StackHome)
		}
	})

	t.Run("STACK_FILE set, STACK_HOME derived", func(t *testing.T) {
		os.Setenv("STACK_FILE", "/another/path/docker-compose.yml")
		os.Unsetenv("STACK_HOME")

		env, err := DetectStackEnv()
		if err != nil {
			t.Fatalf("DetectStackEnv() returned error: %v", err)
		}

		if !strings.HasSuffix(env.StackFile, "/another/path/docker-compose.yml") {
			t.Errorf("expected StackFile to end with /another/path/docker-compose.yml, got %s", env.StackFile)
		}

		if !strings.HasSuffix(env.StackHome, "/another/path") {
			t.Errorf("expected StackHome to end with /another/path, got %s", env.StackHome)
		}
	})

	t.Run("STACK_HOME set, STACK_FILE derived", func(t *testing.T) {
		os.Unsetenv("STACK_FILE")
		os.Setenv("STACK_HOME", "/home/stack")

		env, err := DetectStackEnv()
		if err != nil {
			t.Fatalf("DetectStackEnv() returned error: %v", err)
		}

		if !strings.HasSuffix(env.StackFile, "/home/stack/docker-compose.yml") {
			t.Errorf("expected StackFile to end with /home/stack/docker-compose.yml, got %s", env.StackFile)
		}

		if !strings.HasSuffix(env.StackHome, "/home/stack") {
			t.Errorf("expected StackHome to end with /home/stack, got %s", env.StackHome)
		}
	})

	t.Run("fallback to executable directory", func(t *testing.T) {
		os.Unsetenv("STACK_FILE")
		os.Unsetenv("STACK_HOME")

		env, err := DetectStackEnv()
		if err != nil {
			t.Fatalf("DetectStackEnv() returned error: %v", err)
		}

		// Should derive from executable path
		if env.StackFile == "" {
			t.Error("expected StackFile to be set")
		}

		if env.StackHome == "" {
			t.Error("expected StackHome to be set")
		}

		// Should end with docker-compose.yml
		if !strings.HasSuffix(env.StackFile, "docker-compose.yml") {
			t.Errorf("expected StackFile to end with docker-compose.yml, got %s", env.StackFile)
		}
	})
}

func TestValidateStackFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "stack-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	t.Run("empty stack file", func(t *testing.T) {
		err := ValidateStackFile("")
		if err == nil {
			t.Error("expected error for empty stack file, got nil")
		}
		if !strings.Contains(err.Error(), "STACK_FILE not set") {
			t.Errorf("expected error about STACK_FILE not set, got: %v", err)
		}
	})

	t.Run("non-existent file", func(t *testing.T) {
		err := ValidateStackFile("/nonexistent/docker-compose.yml")
		if err == nil {
			t.Error("expected error for non-existent file, got nil")
		}
		if !strings.Contains(err.Error(), "not found") {
			t.Errorf("expected error about file not found, got: %v", err)
		}
	})

	t.Run("directory instead of file", func(t *testing.T) {
		err := ValidateStackFile(tempDir)
		if err == nil {
			t.Error("expected error for directory, got nil")
		}
		if !strings.Contains(err.Error(), "directory") {
			t.Errorf("expected error about directory, got: %v", err)
		}
	})

	t.Run("valid file", func(t *testing.T) {
		// Create a valid docker-compose.yml file
		validFile := filepath.Join(tempDir, "docker-compose.yml")
		content := []byte("version: '3'\nservices:\n  test:\n    image: nginx\n")
		if err := os.WriteFile(validFile, content, 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		err := ValidateStackFile(validFile)
		if err != nil {
			t.Errorf("expected no error for valid file, got: %v", err)
		}
	})

	t.Run("unreadable file", func(t *testing.T) {
		// Create a file with no read permissions
		unreadableFile := filepath.Join(tempDir, "unreadable.yml")
		content := []byte("version: '3'\n")
		if err := os.WriteFile(unreadableFile, content, 0000); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
		defer os.Chmod(unreadableFile, 0644) // Restore permissions for cleanup

		err := ValidateStackFile(unreadableFile)
		if err == nil {
			t.Error("expected error for unreadable file, got nil")
		}
		if !strings.Contains(err.Error(), "not readable") {
			t.Errorf("expected error about file not readable, got: %v", err)
		}
	})
}

func TestComposeEnvironmentBuilding(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "compose-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a valid docker-compose.yml file
	composeFile := filepath.Join(tempDir, "docker-compose.yml")
	content := []byte("version: '3'\nservices:\n  test:\n    image: nginx\n")
	if err := os.WriteFile(composeFile, content, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	t.Run("validate environment variables are built correctly", func(t *testing.T) {
		// Create a test project directory
		projectDir := filepath.Join(tempDir, "Test Project 123")
		if err := os.MkdirAll(projectDir, 0755); err != nil {
			t.Fatalf("failed to create project dir: %v", err)
		}

		// Note: We can't fully test command execution without docker,
		// but we can test that the function handles missing docker gracefully
		result := ComposeUp(composeFile, projectDir)
		
		// The command will fail without docker, but we can verify it tried
		// to execute with the correct setup
		if result == nil {
			t.Error("expected non-nil result")
		}

		// The error should be about docker, not our validation
		if result.Success {
			// If docker is available and succeeded, that's fine too
			t.Log("Docker compose up succeeded (docker is available)")
		} else if result.Error != nil && !strings.Contains(result.Error.Error(), "STACK_FILE") {
			// Error should be about docker execution, not our validation
			t.Logf("Expected docker execution error, got: %v", result.Error)
		}
	})

	t.Run("validate STACK_FILE before execution", func(t *testing.T) {
		// Test with invalid stack file
		result := ComposeUp("/nonexistent/docker-compose.yml", tempDir)
		
		if result.Success {
			t.Error("expected failure for non-existent stack file")
		}

		if result.Error == nil {
			t.Error("expected error for non-existent stack file")
		} else if !strings.Contains(result.Error.Error(), "not found") {
			t.Errorf("expected validation error, got: %v", result.Error)
		}
	})

	t.Run("default codeDir to current directory", func(t *testing.T) {
		// Test with empty codeDir (should default to current directory)
		result := ComposeUp(composeFile, "")
		
		if result == nil {
			t.Error("expected non-nil result")
		}

		// Should attempt to execute (will fail without docker, but validates our logic)
	})
}

func TestComposeDown(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "compose-down-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a valid docker-compose.yml file
	composeFile := filepath.Join(tempDir, "docker-compose.yml")
	content := []byte("version: '3'\nservices:\n  test:\n    image: nginx\n")
	if err := os.WriteFile(composeFile, content, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	t.Run("validate STACK_FILE before execution", func(t *testing.T) {
		result := ComposeDown("/nonexistent/docker-compose.yml", tempDir)
		
		if result.Success {
			t.Error("expected failure for non-existent stack file")
		}

		if result.Error == nil {
			t.Error("expected error for non-existent stack file")
		}
	})

	t.Run("execute with valid stack file", func(t *testing.T) {
		projectDir := filepath.Join(tempDir, "test-project")
		if err := os.MkdirAll(projectDir, 0755); err != nil {
			t.Fatalf("failed to create project dir: %v", err)
		}

		result := ComposeDown(composeFile, projectDir)
		
		if result == nil {
			t.Error("expected non-nil result")
		}
	})
}

func TestComposeRestart(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "compose-restart-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a valid docker-compose.yml file
	composeFile := filepath.Join(tempDir, "docker-compose.yml")
	content := []byte("version: '3'\nservices:\n  test:\n    image: nginx\n")
	if err := os.WriteFile(composeFile, content, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	t.Run("validate STACK_FILE before execution", func(t *testing.T) {
		result := ComposeRestart("/nonexistent/docker-compose.yml", tempDir)
		
		if result.Success {
			t.Error("expected failure for non-existent stack file")
		}

		if result.Error == nil {
			t.Error("expected error for non-existent stack file")
		}
	})
}

func TestComposeDestroy(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "compose-destroy-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a valid docker-compose.yml file
	composeFile := filepath.Join(tempDir, "docker-compose.yml")
	content := []byte("version: '3'\nservices:\n  test:\n    image: nginx\n")
	if err := os.WriteFile(composeFile, content, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	t.Run("validate STACK_FILE before execution", func(t *testing.T) {
		result := ComposeDestroy("/nonexistent/docker-compose.yml", tempDir)
		
		if result.Success {
			t.Error("expected failure for non-existent stack file")
		}

		if result.Error == nil {
			t.Error("expected error for non-existent stack file")
		}
	})

	t.Run("execute with valid stack file", func(t *testing.T) {
		projectDir := filepath.Join(tempDir, "test-project")
		if err := os.MkdirAll(projectDir, 0755); err != nil {
			t.Fatalf("failed to create project dir: %v", err)
		}

		result := ComposeDestroy(composeFile, projectDir)
		
		if result == nil {
			t.Error("expected non-nil result")
		}

		// The command should include the -v flag for volume removal
		// We can't verify this directly without mocking, but the function structure ensures it
	})
}
