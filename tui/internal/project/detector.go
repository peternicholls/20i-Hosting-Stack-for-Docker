// Project: 20i Stack Manager TUI
// File: detector.go
// Purpose: Project detection from current working directory
// Version: 0.1.0
// Updated: 2025-12-28

// Package project provides functionality for detecting and managing web projects.
// It handles project detection from the current working directory, name sanitization,
// and template installation for new projects.
package project

import (
	"fmt"
	"os"
	"path/filepath"
)

// DetectProject detects the current working directory as a project.
// It reads the current working directory using os.Getwd(), derives the project name
// from the directory basename, sanitizes it, and checks for the existence of public_html/.
//
// Returns:
//   - A Project struct with Name (sanitized), Path (absolute), and HasPublicHTML flag
//   - An error if the current directory cannot be determined
//
// Example:
//
//	If pwd is "/home/user/My Project 123", the returned Project will have:
//	  Name: "my-project-123"
//	  Path: "/home/user/My Project 123"
//	  HasPublicHTML: true (if public_html/ exists)
func DetectProject() (*Project, error) {
	// Get current working directory
	projectRoot, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %w", err)
	}

	// Derive project name from directory basename
	rawName := filepath.Base(projectRoot)

	// Sanitize the project name
	sanitizedName := SanitizeProjectName(rawName)

	// Check if public_html/ exists
	publicHTMLPath := filepath.Join(projectRoot, "public_html")
	hasPublicHTML := false

	if info, err := os.Stat(publicHTMLPath); err == nil && info.IsDir() {
		hasPublicHTML = true
	}

	return &Project{
		Name:          sanitizedName,
		Path:          projectRoot,
		HasPublicHTML: hasPublicHTML,
	}, nil
}
