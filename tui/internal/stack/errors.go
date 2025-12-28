// Project: 20i Stack Manager TUI
// File: errors.go
// Purpose: User-friendly error formatting for Docker/Compose failures
// Version: 0.1.0
// Updated: 2025-12-28

package stack

import (
	"fmt"
	"regexp"
	"strings"
)

// FormatUserError converts Docker/Compose error messages into user-friendly actionable messages.
// It recognizes common error patterns and provides specific guidance.
// For unknown errors, it returns a concise friendly message with a short Details excerpt.
func FormatUserError(err error) string {
	if err == nil {
		return ""
	}

	errStr := err.Error()
	errLower := strings.ToLower(errStr)

	// Port conflict errors
	if strings.Contains(errLower, "port is already allocated") ||
		strings.Contains(errLower, "address already in use") ||
		strings.Contains(errLower, "bind: address already in use") {
		port := extractPortNumber(errStr)
		if port != "" {
			return fmt.Sprintf("Port %s is already in use. Stop the conflicting service or use a different port.", port)
		}
		return "Port is already in use. Stop the conflicting service or use a different port."
	}

	// Docker daemon not running errors
	if strings.Contains(errLower, "cannot connect to the docker daemon") ||
		strings.Contains(errLower, "is the docker daemon running") ||
		strings.Contains(errLower, "docker daemon is not running") ||
		strings.Contains(errLower, "connection refused") && strings.Contains(errLower, "docker") ||
		strings.Contains(errLower, "open //./pipe/docker_engine") ||
		strings.Contains(errLower, "the system cannot find the file specified") && strings.Contains(errLower, "docker") {
		return "Docker daemon is not running. Start Docker Desktop and try again."
	}

	// Permission denied errors
	if strings.Contains(errLower, "permission denied") ||
		strings.Contains(errLower, "access is denied") {
		return "Permission denied. Run Docker as your user or check socket permissions."
	}

	// Timeout errors
	if strings.Contains(errLower, "timeout") ||
		strings.Contains(errLower, "context deadline exceeded") {
		return "Operation timed out. The service may be unresponsive or taking too long to respond."
	}

	// Not found errors
	if strings.Contains(errLower, "no such container") ||
		strings.Contains(errLower, "not found") {
		return "Container not found. It may have been removed. Try refreshing the view."
	}

	// Unknown error - provide concise excerpt
	return formatUnknownError(errStr)
}

// extractPortNumber attempts to extract a port number from an error message.
// It looks for common patterns like "port 80" or ":8080".
func extractPortNumber(errStr string) string {
	// Pattern 1: "port 80" or "Port 80"
	portPattern := regexp.MustCompile(`(?i)port\s+(\d+)`)
	if matches := portPattern.FindStringSubmatch(errStr); len(matches) > 1 {
		return matches[1]
	}

	// Pattern 2: ":8080" in bind address
	bindPattern := regexp.MustCompile(`:(\d+)`)
	if matches := bindPattern.FindStringSubmatch(errStr); len(matches) > 1 {
		return matches[1]
	}

	// Pattern 3: "0.0.0.0:80" or similar
	addrPattern := regexp.MustCompile(`\d+\.\d+\.\d+\.\d+:(\d+)`)
	if matches := addrPattern.FindStringSubmatch(errStr); len(matches) > 1 {
		return matches[1]
	}

	return ""
}

// formatUnknownError creates a user-friendly message for unknown errors.
// It provides a concise excerpt rather than dumping raw multi-line output.
func formatUnknownError(errStr string) string {
	// Limit the details to a single line, max 100 characters
	details := strings.Split(errStr, "\n")[0]
	if len(details) > 100 {
		details = details[:97] + "..."
	}

	return fmt.Sprintf("An error occurred. Details: %s", details)
}
