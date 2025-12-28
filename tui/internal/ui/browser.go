// Project: 20i Stack Manager TUI
// File: browser.go
// Purpose: Platform-agnostic URL opening utility
// Version: 0.1.0
// Updated: 2025-12-28

package ui

import (
	"os/exec"
	"runtime"
)

// URLOpener is an interface for opening URLs in the default browser.
// This abstraction allows for easy mocking in tests and platform-independent URL handling.
type URLOpener interface {
	// OpenURL opens the given URL in the system's default browser.
	// Returns an error if the browser cannot be launched.
	OpenURL(url string) error
}

// DefaultURLOpener uses platform-specific commands to open URLs.
// It is the production implementation of URLOpener.
type DefaultURLOpener struct{}

// OpenURL opens the given URL in the default browser using platform-specific commands.
// On macOS it uses 'open', on Linux it uses 'xdg-open', and on Windows it uses 'start'.
//
// Parameters:
//   - url: The URL to open (e.g., "http://localhost:80")
//
// Returns:
//   - error if the command fails to start (but not if the browser fails to open the URL)
func (d *DefaultURLOpener) OpenURL(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}

	return cmd.Start()
}

// MockURLOpener is a test implementation that records opened URLs.
// It is used for testing URL click functionality without actually launching a browser.
type MockURLOpener struct {
	OpenedURLs    []string // List of URLs that were opened
	ErrorToReturn error    // Error to return from OpenURL (nil for success)
}

// OpenURL records the URL and returns the configured error.
// This allows tests to verify which URLs were opened without side effects.
//
// Parameters:
//   - url: The URL to "open" (actually just recorded)
//
// Returns:
//   - The configured ErrorToReturn value
func (m *MockURLOpener) OpenURL(url string) error {
	m.OpenedURLs = append(m.OpenedURLs, url)
	return m.ErrorToReturn
}
