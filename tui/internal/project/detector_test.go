// Project: 20i Stack Manager TUI
// File: detector_test.go
// Purpose: Unit tests for project detection and sanitization
// Version: 0.1.0
// Updated: 2025-12-28

package project

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSanitizeProjectName(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		// Required spec examples
		{"spec: My Website!", "My Website!", "my-website"},
		{"spec: my_site!", "my_site!", "my-site"},
		{"spec: MYSITE", "MYSITE", "mysite"},
		{"spec: 123-test", "123-test", "test-123"},
		{"spec: empty string", "", "project"},
		{"spec: only hyphens", "---", "project"},

		// Additional edge cases
		{"simple lowercase", "myproject", "myproject"},
		{"uppercase to lowercase", "MyProject", "myproject"},
		{"spaces to hyphens", "My Project", "my-project"},
		{"underscores to hyphens", "test_project", "test-project"},
		{"mixed separators", "test__project--name", "test-project-name"},
		{"special chars", "project@#$%name", "project-name"},
		{"only digits", "123", "p123"},
		{"leading trailing hyphens", "---test---", "test"},
		{"complex case", "ABC--xyz__123", "abc-xyz-123"},
		{"multiple spaces", "my   project   name", "my-project-name"},
		{"only special chars", "___", "project"},
		{"digits in middle", "test123name", "test123name"},
		{"hyphen after digits", "456-project", "project-456"},
		{"digits and text", "2024website", "website-2024"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := SanitizeProjectName(c.input)
			if result != c.expected {
				t.Errorf("SanitizeProjectName(%q) = %q, want %q", c.input, result, c.expected)
			}
		})
	}
}

func TestDetectProject(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "project-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Save current directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current dir: %v", err)
	}
	defer os.Chdir(originalDir)

	t.Run("project without public_html", func(t *testing.T) {
		projectDir := filepath.Join(tempDir, "Test Project 123")
		if err := os.MkdirAll(projectDir, 0755); err != nil {
			t.Fatalf("failed to create project dir: %v", err)
		}

		if err := os.Chdir(projectDir); err != nil {
			t.Fatalf("failed to change to project dir: %v", err)
		}

		project, err := DetectProject()
		if err != nil {
			t.Fatalf("DetectProject() returned error: %v", err)
		}

		if project.Name != "test-project-123" {
			t.Errorf("expected Name to be 'test-project-123', got %q", project.Name)
		}

		if project.Path != projectDir {
			t.Errorf("expected Path to be %q, got %q", projectDir, project.Path)
		}

		if project.HasPublicHTML {
			t.Errorf("expected HasPublicHTML to be false, got true")
		}
	})

	t.Run("project with public_html", func(t *testing.T) {
		projectDir := filepath.Join(tempDir, "my-website")
		publicHTMLDir := filepath.Join(projectDir, "public_html")
		if err := os.MkdirAll(publicHTMLDir, 0755); err != nil {
			t.Fatalf("failed to create public_html dir: %v", err)
		}

		if err := os.Chdir(projectDir); err != nil {
			t.Fatalf("failed to change to project dir: %v", err)
		}

		project, err := DetectProject()
		if err != nil {
			t.Fatalf("DetectProject() returned error: %v", err)
		}

		if project.Name != "my-website" {
			t.Errorf("expected Name to be 'my-website', got %q", project.Name)
		}

		if project.Path != projectDir {
			t.Errorf("expected Path to be %q, got %q", projectDir, project.Path)
		}

		if !project.HasPublicHTML {
			t.Errorf("expected HasPublicHTML to be true, got false")
		}
	})

	t.Run("project with public_html as file", func(t *testing.T) {
		projectDir := filepath.Join(tempDir, "bad-project")
		if err := os.MkdirAll(projectDir, 0755); err != nil {
			t.Fatalf("failed to create project dir: %v", err)
		}

		// Create public_html as a file instead of directory
		publicHTMLFile := filepath.Join(projectDir, "public_html")
		if err := os.WriteFile(publicHTMLFile, []byte("not a directory"), 0644); err != nil {
			t.Fatalf("failed to create public_html file: %v", err)
		}

		if err := os.Chdir(projectDir); err != nil {
			t.Fatalf("failed to change to project dir: %v", err)
		}

		project, err := DetectProject()
		if err != nil {
			t.Fatalf("DetectProject() returned error: %v", err)
		}

		// Should be false because public_html is not a directory
		if project.HasPublicHTML {
			t.Errorf("expected HasPublicHTML to be false when public_html is a file, got true")
		}
	})

	t.Run("project with leading digits in name", func(t *testing.T) {
		projectDir := filepath.Join(tempDir, "2024-website")
		if err := os.MkdirAll(projectDir, 0755); err != nil {
			t.Fatalf("failed to create project dir: %v", err)
		}

		if err := os.Chdir(projectDir); err != nil {
			t.Fatalf("failed to change to project dir: %v", err)
		}

		project, err := DetectProject()
		if err != nil {
			t.Fatalf("DetectProject() returned error: %v", err)
		}

		// "2024-website" should become "website-2024" after sanitization
		if project.Name != "website-2024" {
			t.Errorf("expected Name to be 'website-2024', got %q", project.Name)
		}
	})
}
