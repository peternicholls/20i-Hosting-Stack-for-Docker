// Project: 20i Stack Manager TUI
// File: template.go
// Purpose: Template installation from demo-site-folder
// Version: 0.1.0
// Updated: 2025-12-28

package project

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// InstallTemplate copies the demo site template from demo-site-folder/public_html/
// into the current project's public_html/ directory.
//
// Template location detection:
//  1. If STACK_HOME environment variable is set, look for demo-site-folder there
//  2. Otherwise, locate demo-site-folder relative to the TUI executable
//  3. The template source is always demo-site-folder/public_html/
//
// The function creates the destination public_html/ directory if it doesn't exist,
// and recursively copies all files and subdirectories from the template.
//
// Parameters:
//   - projectRoot: The absolute path to the project root directory
//
// Returns:
//   - An error if the template cannot be found or copying fails
//
// Example:
//
//	err := InstallTemplate("/home/user/myproject")
//	// Creates /home/user/myproject/public_html/ with template contents
func InstallTemplate(projectRoot string) error {
	// Determine template source location
	templatePath, err := findTemplatePath()
	if err != nil {
		return err
	}

	// Verify template source exists
	if info, err := os.Stat(templatePath); err != nil || !info.IsDir() {
		return fmt.Errorf("template not found at %s", templatePath)
	}

	// Destination is projectRoot/public_html
	destPath := filepath.Join(projectRoot, "public_html")

	// Create destination directory
	if err := os.MkdirAll(destPath, 0755); err != nil {
		return fmt.Errorf("failed to create public_html directory: %w", err)
	}

	// Copy all files from template to destination
	return copyDir(templatePath, destPath)
}

// findTemplatePath locates the demo-site-folder/public_html template directory.
// It checks STACK_HOME first, then falls back to executable-relative path.
func findTemplatePath() (string, error) {
	// Check STACK_HOME environment variable
	if stackHome := os.Getenv("STACK_HOME"); stackHome != "" {
		templatePath := filepath.Join(stackHome, "demo-site-folder", "public_html")
		if info, err := os.Stat(templatePath); err == nil && info.IsDir() {
			return templatePath, nil
		}
	}

	// Fall back to executable-relative path
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}

	// The executable is in tui/bin/, so we go up to find the repo root
	// Expected structure: repo_root/tui/bin/tui-executable
	// We need to go up 2 levels to repo root, then to demo-site-folder
	execDir := filepath.Dir(execPath)
	repoRoot := filepath.Join(execDir, "..", "..")
	templatePath := filepath.Join(repoRoot, "demo-site-folder", "public_html")

	// Clean the path to resolve .. segments
	templatePath = filepath.Clean(templatePath)

	if info, err := os.Stat(templatePath); err == nil && info.IsDir() {
		return templatePath, nil
	}

	return "", fmt.Errorf("template not found: checked STACK_HOME and executable-relative paths")
}

// copyDir recursively copies all files and directories from src to dst.
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate the relative path from src
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if relPath == "." {
			return nil
		}

		// Calculate destination path
		destPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			// Create directory hierarchy with a sane default mode, then
			// explicitly set permissions to match the source directory.
			if err := os.MkdirAll(destPath, 0o755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", destPath, err)
			}
			if err := os.Chmod(destPath, info.Mode().Perm()); err != nil {
				return fmt.Errorf("failed to set permissions for directory %s: %w", destPath, err)
			}
			return nil
		}

		// Copy file
		return copyFile(path, destPath, info.Mode())
	})
}

// copyFile copies a single file from src to dst with the specified permissions.
func copyFile(src, dst string, mode os.FileMode) error {
	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %w", src, err)
	}
	defer srcFile.Close()

	// Create destination file
	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %w", dst, err)
	}
	defer dstFile.Close()

	// Copy contents
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	// Set permissions
	if err := os.Chmod(dst, mode); err != nil {
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	return nil
}
