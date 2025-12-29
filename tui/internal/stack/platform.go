// Project: 20i Stack Manager TUI
// File: platform.go
// Purpose: Platform-specific architecture detection and image selection
// Version: 0.1.0
// Updated: 2025-12-28

package stack

import (
	"os"
	"runtime"
)

// GetPhpMyAdminImage returns the appropriate phpMyAdmin Docker image based on system architecture.
// It supports environment variable override via PHPMYADMIN_IMAGE.
//
// Architecture mapping:
//   - ARM64 (Apple Silicon, ARM servers): arm64v8/phpmyadmin:latest
//   - x86_64 (Intel/AMD): phpmyadmin/phpmyadmin:latest
//
// Environment variable override:
//   - PHPMYADMIN_IMAGE: If set, this value is returned regardless of architecture
//
// Returns:
//   - string: Docker image name for phpMyAdmin
func GetPhpMyAdminImage() string {
	// Check for environment variable override first
	if envImage := os.Getenv("PHPMYADMIN_IMAGE"); envImage != "" {
		return envImage
	}

	// Detect architecture and return appropriate image
	switch runtime.GOARCH {
	case "arm64":
		return "arm64v8/phpmyadmin:latest"
	case "amd64":
		return "phpmyadmin/phpmyadmin:latest"
	default:
		// Default to cross-platform image for unknown architectures
		return "phpmyadmin/phpmyadmin:latest"
	}
}

// IsARM64 returns true if the current system is running on ARM64 architecture.
func IsARM64() bool {
	return runtime.GOARCH == "arm64"
}

// IsAMD64 returns true if the current system is running on x86_64/AMD64 architecture.
func IsAMD64() bool {
	return runtime.GOARCH == "amd64"
}

// GetArchitecture returns the current system architecture string.
func GetArchitecture() string {
	return runtime.GOARCH
}
