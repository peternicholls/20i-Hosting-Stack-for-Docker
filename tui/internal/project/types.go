// Project: 20i Stack Manager TUI
// File: types.go
// Purpose: Project data structures for project detection and management
// Version: 0.1.0
// Updated: 2025-12-28

package project

// Project represents a detected web project in the current working directory.
// It contains information about the project's name, location, and readiness status.
type Project struct {
	// Name is the sanitized project name derived from the directory basename.
	// Sanitization follows 20i-gui rules: lowercase, hyphens for invalid chars.
	Name string

	// Path is the absolute path to the project root directory.
	Path string

	// HasPublicHTML indicates whether a public_html/ directory exists in the project root.
	// This is used for pre-flight validation before stack operations.
	HasPublicHTML bool
}
