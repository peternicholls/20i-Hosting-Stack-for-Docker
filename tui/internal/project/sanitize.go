// Project: 20i Stack Manager TUI
// File: sanitize.go
// Purpose: Project name sanitization matching 20i-gui behavior
// Version: 0.1.0
// Updated: 2025-12-28

package project

import (
	"regexp"
	"strings"
)

var (
	// invalidCharsRegex matches any sequence of characters that are not lowercase letters or digits.
	// These will be replaced with a single hyphen.
	invalidCharsRegex = regexp.MustCompile(`[^a-z0-9]+`)

	// multiHyphenRegex matches consecutive hyphens or underscores.
	// These will be collapsed into a single hyphen.
	multiHyphenRegex = regexp.MustCompile(`[-_]+`)

	// leadingDigitsRegex matches one or more digits at the start of the string.
	// If the name starts with digits, we'll move them to the end.
	leadingDigitsRegex = regexp.MustCompile(`^(\d+)-?(.*)$`)
)

// SanitizeProjectName sanitizes a project name to match 20i-gui behavior.
// It performs the following transformations:
//   - Converts to lowercase
//   - Replaces sequences of invalid characters (non-alphanumeric) with hyphens
//   - Collapses consecutive hyphens/underscores into a single hyphen
//   - Trims leading and trailing hyphens
//   - Moves leading digits to the end (e.g., "123-test" -> "test-123")
//   - For names that are only digits, prepends "p" (e.g., "123" -> "p123")
//   - Returns "project" as fallback if the result is empty
//
// Examples:
//   - "My Project" -> "my-project"
//   - "test_123" -> "test-123"
//   - "123-test" -> "test-123"
//   - "ABC--xyz__123" -> "abc-xyz-123"
//   - "___" -> "project"
func SanitizeProjectName(name string) string {
	if name == "" {
		return "project"
	}

	// Step 1: Lowercase
	name = strings.ToLower(name)

	// Step 2: Replace invalid chars with hyphens
	name = invalidCharsRegex.ReplaceAllString(name, "-")

	// Step 3: Collapse consecutive hyphens/underscores
	name = multiHyphenRegex.ReplaceAllString(name, "-")

	// Step 4: Trim leading/trailing hyphens
	name = strings.Trim(name, "-")

	// Step 5: Handle leading digits by moving them to the end
	// e.g., "123-test" -> "test-123"
	if matches := leadingDigitsRegex.FindStringSubmatch(name); len(matches) == 3 {
		digits := matches[1]
		rest := matches[2]
		if rest != "" {
			// Move digits to end
			name = rest + "-" + digits
		} else {
			// Only digits, prepend "p"
			name = "p" + digits
		}
	}

	// Step 6: Ensure doesn't start with hyphen after digit transformation
	name = strings.TrimPrefix(name, "-")

	// Step 7: Fallback if empty
	if name == "" {
		return "project"
	}

	return name
}
