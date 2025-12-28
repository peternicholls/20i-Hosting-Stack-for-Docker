// Project: 20i Stack Manager TUI
// File: color_preview_test.go
// Purpose: Visual preview of UI color palette via go test output.
// Version: 0.1.0
// Updated: 2025-12-28
package ui_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/peternicholls/20i-stack/tui/internal/ui"
// ...existing code...
)

// TestColorPreview renders simple swatches for all defined palette entries.
// Run with: go test ./internal/ui -run TestColorPreview -v
func TestColorPreview(t *testing.T) {
	lipgloss.SetColorProfile(termenv.TrueColor)

	type entry struct {
		name  string
		color lipgloss.Color
	}

	palette := []entry{
		{"ColorRunning (#00ff00)", ui.ColorRunning},
		{"ColorStopped (#808080)", ui.ColorStopped},
		{"ColorError (#ff0000)", ui.ColorError},
		{"ColorWarning (#ffff00)", ui.ColorWarning},
		{"ColorInfo (#0000ff)", ui.ColorInfo},
		{"ColorPrimary (#7D56F4)", ui.ColorPrimary},
		{"ColorSecondary (#04B575)", ui.ColorSecondary},
		{"ColorAccent (#0000ff)", ui.ColorAccent},
		{"ColorBorder (#585858)", ui.ColorBorder},
		{"ColorText (#e4e4e4)", ui.ColorText},
		{"ColorMuted (#767676)", ui.ColorMuted},
		{"ColorSelection (#585858)", ui.ColorSelection},
		{"ColorHighlight (#ff00ff)", ui.ColorHighlight},
	}

	// Header
	fmt.Println("\n=== UI Color Palette Preview ===")
	fmt.Println("Each line shows a colored swatch and the name.")
	fmt.Println("(If you see no color, try: go test -v | cat or use a truecolor terminal)")
	if noColor := os.Getenv("NO_COLOR"); noColor != "" {
		fmt.Println("NO_COLOR is set; colors may be suppressed.")
	}

	// Render solid block swatch for each color, always show hex/rgb code
	for _, p := range palette {
		hex := string(p.color)

		swatch := lipgloss.NewStyle().
			Background(lipgloss.Color(hex)).
			Foreground(lipgloss.Color(hex)).
			Render("██████")

		label := lipgloss.NewStyle().
			Foreground(ui.ColorText).
			Padding(0, 1).
			Render(fmt.Sprintf("%s [%s]", p.name, hex))

		fmt.Printf("%s %s\n", swatch, label)
	}

	// Adaptive text demonstration on primary background
	adaptive := lipgloss.NewStyle().
		Background(ui.ColorPrimary).
		Foreground(lipgloss.Color(ui.AdaptiveText.Light)).
		Padding(0, 2).
		Render("Adaptive Light on Primary")

	adaptiveDark := lipgloss.NewStyle().
		Background(ui.ColorPrimary).
		Foreground(lipgloss.Color(ui.AdaptiveText.Dark)).
		Padding(0, 2).
		Render("Adaptive Dark on Primary")

	fmt.Printf("%s  %s\n", adaptive, adaptiveDark)
}
