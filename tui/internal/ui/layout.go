// Project: 20i Stack Manager TUI
// File: layout.go
// Purpose: Panel sizing helpers for responsive layouts.
// Version: 0.1.0
// Updated: 2025-12-28

package ui

import "github.com/charmbracelet/lipgloss"

const panelGap = 4

// calculatePanelWidths returns left and right panel widths for a 2-panel layout.
// It starts with a 30/70 split (respecting a global panelGap) and then expands
// panels to fit measured content using lipgloss.Width(). If content minima exceed
// the available space (totalWidth - panelGap), the function enforces content
// minimums and may return widths whose sum exceeds available space. Returned
// widths are clamped to >= 0.
func calculatePanelWidths(totalWidth int, leftContent, rightContent string) (int, int) {
	if totalWidth <= 0 {
		return 0, 0
	}

	available := totalWidth - panelGap
	if available < 0 {
		available = totalWidth
	}

	leftWidth := (available * 30) / 100
	rightWidth := available - leftWidth

	leftContentWidth := lipgloss.Width(leftContent)
	rightContentWidth := lipgloss.Width(rightContent)

	if leftWidth < leftContentWidth {
		leftWidth = leftContentWidth
		rightWidth = available - leftWidth
	}

	if rightWidth < rightContentWidth {
		rightWidth = rightContentWidth
		leftWidth = available - rightWidth
	}

	if leftWidth < leftContentWidth {
		leftWidth = leftContentWidth
	}

	if rightWidth < rightContentWidth {
		rightWidth = rightContentWidth
	}

	if leftWidth < 0 {
		leftWidth = 0
	}

	if rightWidth < 0 {
		rightWidth = 0
	}

	return leftWidth, rightWidth
}

// calculatePanelHeights returns header, content, and footer heights for a layout.
// Heights are derived from measured content so panels can be sized before rendering.
func calculatePanelHeights(totalHeight int, headerContent, footerContent string) (int, int, int) {
	if totalHeight <= 0 {
		return 0, 0, 0
	}

	headerHeight := lipgloss.Height(headerContent)
	footerHeight := lipgloss.Height(footerContent)
	contentHeight := totalHeight - headerHeight - footerHeight
	if contentHeight < 0 {
		contentHeight = 0
	}

	return headerHeight, contentHeight, footerHeight
}
