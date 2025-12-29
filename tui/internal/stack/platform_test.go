// Project: 20i Stack Manager TUI
// File: platform_test.go
// Purpose: Unit tests for platform detection and image selection
// Version: 0.1.0
// Updated: 2025-12-28

package stack

import (
	"os"
	"runtime"
	"testing"
)

func TestGetPhpMyAdminImage(t *testing.T) {
	tests := []struct {
		name        string
		envValue    string
		expectedImg string
	}{
		{
			name:        "environment override",
			envValue:    "custom/phpmyadmin:v1",
			expectedImg: "custom/phpmyadmin:v1",
		},
		{
			name:     "no environment - uses architecture",
			envValue: "",
			// Expected value depends on runtime architecture
			expectedImg: "", // Will be set in test body
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable if provided
			if tt.envValue != "" {
				os.Setenv("PHPMYADMIN_IMAGE", tt.envValue)
				defer os.Unsetenv("PHPMYADMIN_IMAGE")
			} else {
				os.Unsetenv("PHPMYADMIN_IMAGE")
			}

			result := GetPhpMyAdminImage()

			if tt.envValue != "" {
				// If env var is set, it should return that value
				if result != tt.expectedImg {
					t.Errorf("GetPhpMyAdminImage() = %q, want %q", result, tt.expectedImg)
				}
			} else {
				// If env var is not set, check against architecture-specific value
				var expected string
				switch runtime.GOARCH {
				case "arm64":
					expected = "arm64v8/phpmyadmin:latest"
				default:
					expected = "phpmyadmin/phpmyadmin:latest"
				}
				if result != expected {
					t.Errorf("GetPhpMyAdminImage() = %q, want %q (for arch %s)", result, expected, runtime.GOARCH)
				}
			}
		})
	}
}

func TestIsARM64(t *testing.T) {
	result := IsARM64()
	expected := runtime.GOARCH == "arm64"
	if result != expected {
		t.Errorf("IsARM64() = %v, want %v (arch: %s)", result, expected, runtime.GOARCH)
	}
}

func TestIsAMD64(t *testing.T) {
	result := IsAMD64()
	expected := runtime.GOARCH == "amd64"
	if result != expected {
		t.Errorf("IsAMD64() = %v, want %v (arch: %s)", result, expected, runtime.GOARCH)
	}
}

func TestGetArchitecture(t *testing.T) {
	result := GetArchitecture()
	expected := runtime.GOARCH
	if result != expected {
		t.Errorf("GetArchitecture() = %q, want %q", result, expected)
	}
}

func TestGetPhpMyAdminImage_ArchitectureSpecific(t *testing.T) {
	// Clear environment to test architecture-based selection
	os.Unsetenv("PHPMYADMIN_IMAGE")

	result := GetPhpMyAdminImage()

	// Verify the result is one of the expected values
	validImages := map[string]bool{
		"arm64v8/phpmyadmin:latest":    true,
		"phpmyadmin/phpmyadmin:latest": true,
	}

	if !validImages[result] {
		t.Errorf("GetPhpMyAdminImage() = %q, want one of %v", result, validImages)
	}

	// Verify it matches the current architecture
	switch runtime.GOARCH {
	case "arm64":
		if result != "arm64v8/phpmyadmin:latest" {
			t.Errorf("On ARM64, expected arm64v8/phpmyadmin:latest, got %q", result)
		}
	case "amd64":
		if result != "phpmyadmin/phpmyadmin:latest" {
			t.Errorf("On AMD64, expected phpmyadmin/phpmyadmin:latest, got %q", result)
		}
	default:
		// For other architectures, should use cross-platform image
		if result != "phpmyadmin/phpmyadmin:latest" {
			t.Errorf("On %s, expected phpmyadmin/phpmyadmin:latest, got %q", runtime.GOARCH, result)
		}
	}
}
