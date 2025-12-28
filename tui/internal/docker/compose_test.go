// Project: 20i Stack Manager TUI
// File: compose_test.go
// Purpose: Unit tests for Docker compose operations
// Version: 0.1.0
// Updated: 2025-12-28

package docker

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

var (
	lastExecName   string
	lastExecArgs   []string
	helperExitCode int
	helperStderr   string
)

func fakeExecCommand(ctx context.Context, name string, args ...string) *exec.Cmd {
	lastExecName = name
	lastExecArgs = append([]string{}, args...)

	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=TestHelperProcess", "--")
	cmd.Env = append(os.Environ(),
		"GO_WANT_HELPER_PROCESS=1",
		fmt.Sprintf("HELPER_EXIT_CODE=%d", helperExitCode),
		fmt.Sprintf("HELPER_STDERR=%s", helperStderr),
	)
	return cmd
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	if msg := os.Getenv("HELPER_STDERR"); msg != "" {
		fmt.Fprint(os.Stderr, msg)
	}

	code, err := strconv.Atoi(os.Getenv("HELPER_EXIT_CODE"))
	if err != nil {
		code = 1
	}
	os.Exit(code)
}

func TestComposeStop_Success(t *testing.T) {
	originalExec := execCommand
	execCommand = fakeExecCommand
	defer func() { execCommand = originalExec }()

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "docker-compose.yml"), []byte("services: {}"), 0o644); err != nil {
		t.Fatalf("write compose file: %v", err)
	}

	resetHelper()
	client := &Client{ctx: context.Background()}
	if err := client.ComposeStop(dir); err != nil {
		t.Fatalf("ComposeStop returned error: %v", err)
	}

	assertComposeArgs(t, "stop", []string{"compose", "stop"})
}

func TestComposeRestart_Error(t *testing.T) {
	originalExec := execCommand
	execCommand = fakeExecCommand
	defer func() { execCommand = originalExec }()

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "docker-compose.yml"), []byte("services: {}"), 0o644); err != nil {
		t.Fatalf("write compose file: %v", err)
	}

	resetHelper()
	helperExitCode = 1
	helperStderr = "permission denied"

	client := &Client{ctx: context.Background()}
	err := client.ComposeRestart(dir)
	if err == nil {
		t.Fatal("expected error from ComposeRestart")
	}

	if !strings.Contains(err.Error(), "permission denied") {
		t.Fatalf("expected error to include stderr, got %v", err)
	}
}

func TestComposeDown_RemoveVolumes(t *testing.T) {
	originalExec := execCommand
	execCommand = fakeExecCommand
	defer func() { execCommand = originalExec }()

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "docker-compose.yml"), []byte("services: {}"), 0o644); err != nil {
		t.Fatalf("write compose file: %v", err)
	}

	resetHelper()
	client := &Client{ctx: context.Background()}
	if err := client.ComposeDown(dir, true); err != nil {
		t.Fatalf("ComposeDown returned error: %v", err)
	}

	assertComposeArgs(t, "down -v", []string{"compose", "down", "-v"})
}

func TestComposeCommand_MissingComposeFile(t *testing.T) {
	originalExec := execCommand
	execCommand = fakeExecCommand
	defer func() { execCommand = originalExec }()

	dir := t.TempDir()
	resetHelper()
	client := &Client{ctx: context.Background()}

	err := client.ComposeStop(dir)
	if err == nil {
		t.Fatal("expected error for missing compose file")
	}

	if lastExecName != "" {
		t.Fatal("expected compose command not to run when compose file missing")
	}
}

func resetHelper() {
	lastExecName = ""
	lastExecArgs = nil
	helperExitCode = 0
	helperStderr = ""
}

func assertComposeArgs(t *testing.T, label string, want []string) {
	t.Helper()
	if lastExecName != "docker" {
		t.Fatalf("expected docker command for %s, got %q", label, lastExecName)
	}
	if !equalStringSlice(lastExecArgs, want) {
		t.Fatalf("expected args %v for %s, got %v", want, label, lastExecArgs)
	}
}

func equalStringSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
