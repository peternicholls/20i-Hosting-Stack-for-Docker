// Project: 20i Stack Manager TUI
// File: main.go
// Purpose: TUI application entry point.
// Version: 0.1.0
// Updated: 2025-12-28
package main

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/peternicholls/20i-stack/tui/internal/app"
)

func main() {
	ctx := context.Background()

	// Create RootModel with Docker client
	rootModel, err := app.NewRootModel(ctx)
	if err != nil {
		fmt.Printf("Error initializing TUI: %v\n", err)
		os.Exit(1)
	}

	// Create and run the Bubble Tea program
	p := tea.NewProgram(rootModel)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
