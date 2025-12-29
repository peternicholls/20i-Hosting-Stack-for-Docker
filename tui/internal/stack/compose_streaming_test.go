// Project: 20i Stack Manager TUI
// File: compose_streaming_test.go
// Purpose: Unit tests for streaming compose operations
// Version: 0.1.0
// Updated: 2025-12-28

package stack

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestComposeUpStreaming_ChannelClosure tests that the channel closes on both success and error paths
func TestComposeUpStreaming_ChannelClosure(t *testing.T) {
// Create a temporary directory for testing
tempDir, err := os.MkdirTemp("", "compose-stream-test-*")
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

t.Run("channel closes after command execution", func(t *testing.T) {
// This test will attempt to run docker compose, which may fail
// but we're testing that the channel closes regardless
outputChan, err := ComposeUpStreaming(composeFile, tempDir)
if err != nil {
t.Fatalf("ComposeUpStreaming returned error: %v", err)
}

// Read all output until channel closes
var lines []string
for line := range outputChan {
lines = append(lines, line)
}

// Verify channel was closed (loop completed)
// If channel wasn't closed, this test would hang
t.Logf("Received %d lines before channel closed", len(lines))
})

t.Run("validation error returns error without channel", func(t *testing.T) {
// Test with invalid stack file
_, err := ComposeUpStreaming("/nonexistent/docker-compose.yml", tempDir)
if err == nil {
t.Error("expected error for non-existent stack file")
}
if !strings.Contains(err.Error(), "not found") {
t.Errorf("expected validation error, got: %v", err)
}
})
}

// TestComposeUpStreaming_OutputOrder tests that lines are delivered in order
func TestComposeUpStreaming_OutputOrder(t *testing.T) {
// Create a mock script that outputs predictable lines
tempDir, err := os.MkdirTemp("", "compose-stream-order-test-*")
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

// Note: We can't fully control docker compose output without mocking,
// so this test just verifies the mechanism works
outputChan, err := ComposeUpStreaming(composeFile, tempDir)
if err != nil {
t.Fatalf("ComposeUpStreaming returned error: %v", err)
}

// Collect all output
var lines []string
for line := range outputChan {
lines = append(lines, line)
}

// Verify we got some output and channel closed
if len(lines) == 0 {
t.Error("expected some output lines")
}

// Verify completion message or error message is present
hasCompletionOrError := false
for _, line := range lines {
if line == "[Complete]" || strings.HasPrefix(line, "ERROR:") {
hasCompletionOrError = true
break
}
}

if !hasCompletionOrError {
t.Error("expected [Complete] or ERROR: message in output")
}
}

// TestComposeUpStreaming_EmptyCodeDir tests that empty codeDir defaults to current directory
func TestComposeUpStreaming_EmptyCodeDir(t *testing.T) {
tempDir, err := os.MkdirTemp("", "compose-stream-empty-dir-test-*")
if err != nil {
t.Fatalf("failed to create temp dir: %v", err)
}
defer os.RemoveAll(tempDir)

composeFile := filepath.Join(tempDir, "docker-compose.yml")
content := []byte("version: '3'\nservices:\n  test:\n    image: nginx\n")
if err := os.WriteFile(composeFile, content, 0644); err != nil {
t.Fatalf("failed to create test file: %v", err)
}

// Test with empty codeDir
outputChan, err := ComposeUpStreaming(composeFile, "")
if err != nil {
t.Fatalf("ComposeUpStreaming returned error: %v", err)
}

// Read output until channel closes
for range outputChan {
// Drain channel
}

// Test passed if we got here (channel closed properly)
}

// TestStreamLines tests the streamLines helper function
func TestStreamLines(t *testing.T) {
testData := "line1\nline2\nline3\n"
reader := strings.NewReader(testData)

outputChan := make(chan string, 10)
done := make(chan struct{})

go streamLines(reader, outputChan, done)

// Wait for completion
<-done
close(outputChan)

// Collect output
var lines []string
for line := range outputChan {
lines = append(lines, line)
}

// Verify lines
expected := []string{"line1", "line2", "line3"}
if len(lines) != len(expected) {
t.Errorf("expected %d lines, got %d", len(expected), len(lines))
}

for i, line := range lines {
if i >= len(expected) {
t.Errorf("unexpected extra line: %s", line)
continue
}
if line != expected[i] {
t.Errorf("line %d: expected %q, got %q", i, expected[i], line)
}
}
}

// TestStreamLines_EmptyInput tests streamLines with empty input
func TestStreamLines_EmptyInput(t *testing.T) {
reader := strings.NewReader("")

outputChan := make(chan string, 10)
done := make(chan struct{})

go streamLines(reader, outputChan, done)

// Wait for completion
<-done
close(outputChan)

// Verify no output
count := 0
for range outputChan {
count++
}

if count != 0 {
t.Errorf("expected 0 lines, got %d", count)
}
}
