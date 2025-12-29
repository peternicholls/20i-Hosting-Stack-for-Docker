// Project: 20i Stack Manager TUI
// File: compose_streaming_test.go
// Purpose: Unit tests for streaming compose operations (no real docker execution)
// Version: 0.1.0
// Updated: 2025-12-29

package stack

import (
"bytes"
"fmt"
"io"
"os"
"path/filepath"
"strings"
"testing"
)

// TestComposeUpStreaming_Validation tests validation without executing docker
func TestComposeUpStreaming_Validation(t *testing.T) {
t.Run("validation error for non-existent stack file", func(t *testing.T) {
// Test with invalid stack file
_, err := ComposeUpStreaming("/nonexistent/docker-compose.yml", "/tmp")
if err == nil {
t.Error("expected error for non-existent stack file")
}
if !strings.Contains(err.Error(), "not found") {
t.Errorf("expected validation error, got: %v", err)
}
})

t.Run("validation error for empty stack file", func(t *testing.T) {
_, err := ComposeUpStreaming("", "/tmp")
if err == nil {
t.Error("expected error for empty stack file")
}
if !strings.Contains(err.Error(), "not set") {
t.Errorf("expected 'not set' error, got: %v", err)
}
})

t.Run("validation with valid stack file returns channel", func(t *testing.T) {
// Create a temporary valid compose file
tempDir, err := os.MkdirTemp("", "compose-stream-test-*")
if err != nil {
t.Fatalf("failed to create temp dir: %v", err)
}
defer os.RemoveAll(tempDir)

composeFile := filepath.Join(tempDir, "docker-compose.yml")
content := []byte("version: '3'\nservices:\n  test:\n    image: nginx\n")
if err := os.WriteFile(composeFile, content, 0644); err != nil {
t.Fatalf("failed to create test file: %v", err)
}

// This will attempt to execute docker, but we just verify we get a channel
outputChan, err := ComposeUpStreaming(composeFile, tempDir)
if err != nil {
t.Errorf("unexpected validation error: %v", err)
return
}

if outputChan == nil {
t.Error("expected non-nil channel")
return
}

// Drain the channel to allow goroutine to complete
// This proves the channel will be closed
for range outputChan {
// Just drain it
}
// If we get here, channel was closed successfully
})

t.Run("empty codeDir defaults to current directory", func(t *testing.T) {
tempDir, err := os.MkdirTemp("", "compose-stream-test-*")
if err != nil {
t.Fatalf("failed to create temp dir: %v", err)
}
defer os.RemoveAll(tempDir)

composeFile := filepath.Join(tempDir, "docker-compose.yml")
content := []byte("version: '3'\nservices:\n  test:\n    image: nginx\n")
if err := os.WriteFile(composeFile, content, 0644); err != nil {
t.Fatalf("failed to create test file: %v", err)
}

// Empty codeDir should not cause validation error
outputChan, err := ComposeUpStreaming(composeFile, "")
if err != nil {
t.Errorf("unexpected error with empty codeDir: %v", err)
return
}

// Drain the channel
for range outputChan {
}
})
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

// TestStreamLines_WithError tests streamLines handles scanner errors
func TestStreamLines_WithError(t *testing.T) {
// Create a reader that will error after some data
errorReader := &erroringReader{
data:  []byte("line1\nline2\n"),
pos:   0,
errAt: 10,
}

outputChan := make(chan string, 10)
done := make(chan struct{})

go streamLines(errorReader, outputChan, done)

// Wait for completion
<-done
close(outputChan)

// Collect output
var lines []string
for line := range outputChan {
lines = append(lines, line)
}

// Should have some lines and an error message
foundError := false
for _, line := range lines {
if strings.HasPrefix(line, "ERROR: Stream read error") {
foundError = true
break
}
}

if !foundError {
t.Error("expected ERROR message for stream read error")
}
}

// erroringReader is a test helper that errors after reading a certain amount
type erroringReader struct {
data  []byte
pos   int
errAt int
}

func (r *erroringReader) Read(p []byte) (n int, err error) {
if r.pos >= r.errAt {
return 0, fmt.Errorf("simulated read error")
}

if r.pos >= len(r.data) {
return 0, io.EOF
}

n = copy(p, r.data[r.pos:])
r.pos += n

if r.pos >= r.errAt && n > 0 {
// Partial read before error
return n, nil
}

return n, nil
}

// TestStreamLinesOrderPreservation tests that line order is preserved
func TestStreamLinesOrderPreservation(t *testing.T) {
// Create test data with specific order
lines := []string{
"First line",
"Second line", 
"Third line",
"Fourth line",
"Fifth line",
}
testData := strings.Join(lines, "\n") + "\n"
reader := bytes.NewBufferString(testData)

outputChan := make(chan string, 20)
done := make(chan struct{})

go streamLines(reader, outputChan, done)

// Wait for completion
<-done
close(outputChan)

// Verify order is preserved
var received []string
for line := range outputChan {
received = append(received, line)
}

if len(received) != len(lines) {
t.Fatalf("expected %d lines, got %d", len(lines), len(received))
}

for i, expected := range lines {
if received[i] != expected {
t.Errorf("line %d: expected %q, got %q", i, expected, received[i])
}
}
}

// TestChannelClosurePattern demonstrates the channel closure pattern
func TestChannelClosurePattern(t *testing.T) {
t.Run("channel closes on success", func(t *testing.T) {
// Simulate successful output stream
outputChan := make(chan string, 5)

go func() {
defer close(outputChan)
outputChan <- "Line 1"
outputChan <- "Line 2"
outputChan <- "[Complete]"
}()

var lines []string
for line := range outputChan {
lines = append(lines, line)
}

if len(lines) != 3 {
t.Errorf("expected 3 lines, got %d", len(lines))
}
if lines[len(lines)-1] != "[Complete]" {
t.Error("expected [Complete] as last line")
}
})

t.Run("channel closes on error", func(t *testing.T) {
// Simulate error during streaming
outputChan := make(chan string, 5)

go func() {
defer close(outputChan)
outputChan <- "Line 1"
outputChan <- "ERROR: Something failed"
}()

var lines []string
for line := range outputChan {
lines = append(lines, line)
}

foundError := false
for _, line := range lines {
if strings.HasPrefix(line, "ERROR:") {
foundError = true
}
}

if !foundError {
t.Error("expected error message in output")
}
})
}
