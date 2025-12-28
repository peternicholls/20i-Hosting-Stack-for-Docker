// Project: 20i Stack Manager TUI
// File: template_test.go
// Purpose: Unit tests for template installation
// Version: 0.1.0
// Updated: 2025-12-28

package project

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInstallTemplate(t *testing.T) {
	t.Run("successful template installation with STACK_HOME", func(t *testing.T) {
		// Create a temporary directory structure
		tempDir, err := os.MkdirTemp("", "template-test-*")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		// Create mock STACK_HOME with demo-site-folder
		stackHome := filepath.Join(tempDir, "stack")
		templateSource := filepath.Join(stackHome, "demo-site-folder", "public_html")
		if err := os.MkdirAll(templateSource, 0755); err != nil {
			t.Fatalf("failed to create template source: %v", err)
		}

		// Create some template files
		indexFile := filepath.Join(templateSource, "index.php")
		if err := os.WriteFile(indexFile, []byte("<?php echo 'Hello'; ?>"), 0644); err != nil {
			t.Fatalf("failed to create index.php: %v", err)
		}

		subDir := filepath.Join(templateSource, "assets")
		if err := os.MkdirAll(subDir, 0755); err != nil {
			t.Fatalf("failed to create assets dir: %v", err)
		}

		cssFile := filepath.Join(subDir, "style.css")
		if err := os.WriteFile(cssFile, []byte("body { margin: 0; }"), 0644); err != nil {
			t.Fatalf("failed to create style.css: %v", err)
		}

		// Set STACK_HOME environment variable
		oldStackHome := os.Getenv("STACK_HOME")
		os.Setenv("STACK_HOME", stackHome)
		defer os.Setenv("STACK_HOME", oldStackHome)

		// Create project directory
		projectDir := filepath.Join(tempDir, "myproject")
		if err := os.MkdirAll(projectDir, 0755); err != nil {
			t.Fatalf("failed to create project dir: %v", err)
		}

		// Install template
		if err := InstallTemplate(projectDir); err != nil {
			t.Fatalf("InstallTemplate() returned error: %v", err)
		}

		// Verify files were copied
		destIndex := filepath.Join(projectDir, "public_html", "index.php")
		if _, err := os.Stat(destIndex); os.IsNotExist(err) {
			t.Errorf("expected index.php to be copied, but it doesn't exist")
		}

		destCSS := filepath.Join(projectDir, "public_html", "assets", "style.css")
		if _, err := os.Stat(destCSS); os.IsNotExist(err) {
			t.Errorf("expected assets/style.css to be copied, but it doesn't exist")
		}

		// Verify file contents
		content, err := os.ReadFile(destIndex)
		if err != nil {
			t.Fatalf("failed to read copied index.php: %v", err)
		}
		if string(content) != "<?php echo 'Hello'; ?>" {
			t.Errorf("copied file has wrong content: %q", string(content))
		}
	})

	t.Run("template not found returns error", func(t *testing.T) {
		// Create a temporary project directory
		tempDir, err := os.MkdirTemp("", "template-test-*")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		// Clear STACK_HOME to force executable-relative lookup
		oldStackHome := os.Getenv("STACK_HOME")
		os.Setenv("STACK_HOME", "")
		defer os.Setenv("STACK_HOME", oldStackHome)

		projectDir := filepath.Join(tempDir, "myproject")
		if err := os.MkdirAll(projectDir, 0755); err != nil {
			t.Fatalf("failed to create project dir: %v", err)
		}

		// Try to install template (should fail because template doesn't exist)
		err = InstallTemplate(projectDir)
		if err == nil {
			t.Errorf("expected InstallTemplate to return error when template not found")
		}
	})

	t.Run("creates public_html if it doesn't exist", func(t *testing.T) {
		// Create a temporary directory structure
		tempDir, err := os.MkdirTemp("", "template-test-*")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		// Create mock STACK_HOME with demo-site-folder
		stackHome := filepath.Join(tempDir, "stack")
		templateSource := filepath.Join(stackHome, "demo-site-folder", "public_html")
		if err := os.MkdirAll(templateSource, 0755); err != nil {
			t.Fatalf("failed to create template source: %v", err)
		}

		// Create a template file
		indexFile := filepath.Join(templateSource, "index.html")
		if err := os.WriteFile(indexFile, []byte("<html></html>"), 0644); err != nil {
			t.Fatalf("failed to create index.html: %v", err)
		}

		// Set STACK_HOME
		oldStackHome := os.Getenv("STACK_HOME")
		os.Setenv("STACK_HOME", stackHome)
		defer os.Setenv("STACK_HOME", oldStackHome)

		// Create project directory WITHOUT public_html
		projectDir := filepath.Join(tempDir, "myproject")
		if err := os.MkdirAll(projectDir, 0755); err != nil {
			t.Fatalf("failed to create project dir: %v", err)
		}

		// Verify public_html doesn't exist yet
		publicHTMLDir := filepath.Join(projectDir, "public_html")
		if _, err := os.Stat(publicHTMLDir); !os.IsNotExist(err) {
			t.Fatalf("public_html should not exist yet")
		}

		// Install template
		if err := InstallTemplate(projectDir); err != nil {
			t.Fatalf("InstallTemplate() returned error: %v", err)
		}

		// Verify public_html was created
		if info, err := os.Stat(publicHTMLDir); os.IsNotExist(err) || !info.IsDir() {
			t.Errorf("expected public_html directory to be created")
		}

		// Verify file was copied
		destIndex := filepath.Join(publicHTMLDir, "index.html")
		if _, err := os.Stat(destIndex); os.IsNotExist(err) {
			t.Errorf("expected index.html to be copied")
		}
	})

	t.Run("preserves file permissions", func(t *testing.T) {
		// Create a temporary directory structure
		tempDir, err := os.MkdirTemp("", "template-test-*")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		// Create mock STACK_HOME with demo-site-folder
		stackHome := filepath.Join(tempDir, "stack")
		templateSource := filepath.Join(stackHome, "demo-site-folder", "public_html")
		if err := os.MkdirAll(templateSource, 0755); err != nil {
			t.Fatalf("failed to create template source: %v", err)
		}

		// Create a file with specific permissions
		scriptFile := filepath.Join(templateSource, "script.sh")
		if err := os.WriteFile(scriptFile, []byte("#!/bin/bash"), 0755); err != nil {
			t.Fatalf("failed to create script.sh: %v", err)
		}

		// Set STACK_HOME
		oldStackHome := os.Getenv("STACK_HOME")
		os.Setenv("STACK_HOME", stackHome)
		defer os.Setenv("STACK_HOME", oldStackHome)

		// Create project directory
		projectDir := filepath.Join(tempDir, "myproject")
		if err := os.MkdirAll(projectDir, 0755); err != nil {
			t.Fatalf("failed to create project dir: %v", err)
		}

		// Install template
		if err := InstallTemplate(projectDir); err != nil {
			t.Fatalf("InstallTemplate() returned error: %v", err)
		}

		// Verify permissions were preserved
		destScript := filepath.Join(projectDir, "public_html", "script.sh")
		info, err := os.Stat(destScript)
		if err != nil {
			t.Fatalf("failed to stat copied script.sh: %v", err)
		}

		if info.Mode().Perm() != 0755 {
			t.Errorf("expected file permissions 0755, got %o", info.Mode().Perm())
		}
	})
}

func TestFindTemplatePath(t *testing.T) {
	t.Run("finds template with STACK_HOME set", func(t *testing.T) {
		// Create a temporary STACK_HOME
		tempDir, err := os.MkdirTemp("", "template-test-*")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		templateSource := filepath.Join(tempDir, "demo-site-folder", "public_html")
		if err := os.MkdirAll(templateSource, 0755); err != nil {
			t.Fatalf("failed to create template source: %v", err)
		}

		// Set STACK_HOME
		oldStackHome := os.Getenv("STACK_HOME")
		os.Setenv("STACK_HOME", tempDir)
		defer os.Setenv("STACK_HOME", oldStackHome)

		// Find template path
		path, err := findTemplatePath()
		if err != nil {
			t.Fatalf("findTemplatePath() returned error: %v", err)
		}

		if path != templateSource {
			t.Errorf("expected path %q, got %q", templateSource, path)
		}
	})

	t.Run("returns error when template not found", func(t *testing.T) {
		// Clear STACK_HOME
		oldStackHome := os.Getenv("STACK_HOME")
		os.Setenv("STACK_HOME", "")
		defer os.Setenv("STACK_HOME", oldStackHome)

		// Try to find template (should fail in test environment)
		_, err := findTemplatePath()
		if err == nil {
			t.Errorf("expected findTemplatePath to return error when template not found")
		}
	})
}
