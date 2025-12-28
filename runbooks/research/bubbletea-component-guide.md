# Bubble Tea Component & View Guide

**Document Version**: 1.0.0  
**Created**: 2025-12-28  
**Last Updated**: 2025-12-28  
**Purpose**: Comprehensive guide for agents building Bubble Tea TUI applications with proper component architecture and Lipgloss styling

---

## Table of Contents

1. [Core Concepts](#core-concepts)
2. [Component Architecture](#component-architecture)
3. [Bubbles Component Library](#bubbles-component-library)
4. [Lipgloss Styling System](#lipgloss-styling-system)
5. [Layout Patterns](#layout-patterns)
6. [Component Composition](#component-composition)
7. [Real-World Examples](#real-world-examples)
8. [Anti-Patterns to Avoid](#anti-patterns-to-avoid)
9. [Quick Reference](#quick-reference)

---

## Core Concepts

### The Elm Architecture

Bubble Tea is based on The Elm Architecture with three core methods:

```go
type tea.Model interface {
    // Init returns the initial command to run
    Init() tea.Cmd
    
    // Update handles incoming messages and returns updated model + command
    Update(tea.Msg) (tea.Model, tea.Cmd)
    
    // View renders the model to a string for display
    View() string
}
```

**Key Principles**:
- **Models are immutable** - Update returns a new model, never mutates in place
- **No side effects in Update** - Use `tea.Cmd` for I/O, timers, HTTP requests
- **View is pure** - Derives UI from model state only, no side effects
- **Messages drive changes** - All state changes happen via Update receiving messages

### Model Types

**Root/App Model** - Top-level coordinator
```go
type AppModel struct {
    currentView string
    dashboard   DashboardModel
    containers  ContainersModel
    logs        LogsModel
    width       int
    height      int
}
```

**View/Component Model** - Self-contained feature
```go
type DashboardModel struct {
    services    []Service
    cursor      int
    loading     bool
    err         error
}
```

**Bubbles Components** - Pre-built UI components from library
```go
import "github.com/charmbracelet/bubbles/list"

type model struct {
    list list.Model
}
```

---

## Component Architecture

### Component Structure

Each component should be a standalone Go package with clear responsibilities:

```
internal/
  tui/
    app.go              # Root model that coordinates views
    
    components/
      dashboard/
        dashboard.go    # Model + Init/Update/View
        styles.go       # Lipgloss styles for this component
        messages.go     # Custom message types
        commands.go     # tea.Cmd factories
        
      containers/
        containers.go
        list.go         # Subcomponent: container list
        detail.go       # Subcomponent: container detail
        styles.go
        
      logs/
        logs.go
        viewport.go     # Uses Bubbles viewport
        styles.go
```

### Component Pattern

Every component follows this pattern:

```go
package dashboard

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

// Model holds component state
type Model struct {
    services []Service
    cursor   int
    width    int
    height   int
}

// Service represents a Docker service
type Service struct {
    Name   string
    Status string
    CPU    float64
    Memory string
}

// Init returns initial command
func (m Model) Init() tea.Cmd {
    return fetchServicesCmd()
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "up", "k":
            if m.cursor > 0 {
                m.cursor--
            }
        case "down", "j":
            if m.cursor < len(m.services)-1 {
                m.cursor++
            }
        }
    
    case servicesMsg:
        m.services = msg.services
        return m, nil
    
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
    }
    
    return m, nil
}

// View renders component
func (m Model) View() string {
    if len(m.services) == 0 {
        return emptyStateStyle.Render("No services running")
    }
    
    var rows []string
    for i, svc := range m.services {
        style := serviceRowStyle
        if i == m.cursor {
            style = selectedRowStyle
        }
        rows = append(rows, style.Render(svc.Name))
    }
    
    return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

// Custom message types
type servicesMsg struct {
    services []Service
}

// Command factory
func fetchServicesCmd() tea.Cmd {
    return func() tea.Msg {
        // This runs in background goroutine
        services := fetchServices()
        return servicesMsg{services: services}
    }
}
```

### Root Model Pattern

The root model delegates to child components:

```go
type AppModel struct {
    activeView  string
    dashboard   dashboard.Model
    containers  containers.Model
    width       int
    height      int
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Handle global shortcuts
        switch msg.String() {
        case "q", "ctrl+c":
            return m, tea.Quit
        case "tab":
            m.activeView = nextView(m.activeView)
            return m, nil
        }
    
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
    }
    
    // Delegate to active view
    switch m.activeView {
    case "dashboard":
        m.dashboard, cmd = m.dashboard.Update(msg)
    case "containers":
        m.containers, cmd = m.containers.Update(msg)
    }
    
    return m, cmd
}

func (m AppModel) View() string {
    // Render header
    header := headerStyle.Render("20i Stack Manager")
    
    // Render active view
    var view string
    switch m.activeView {
    case "dashboard":
        view = m.dashboard.View()
    case "containers":
        view = m.containers.View()
    }
    
    // Render footer
    footer := footerStyle.Render("Tab: switch | q: quit")
    
    // Combine
    return lipgloss.JoinVertical(
        lipgloss.Left,
        header,
        view,
        footer,
    )
}
```

---

## Bubbles Component Library

Bubbles provides pre-built components. Use these instead of building from scratch.

### List Component

Perfect for scrollable lists with selection:

```go
import "github.com/charmbracelet/bubbles/list"

// Define list item
type containerItem struct {
    name   string
    status string
}

func (i containerItem) Title() string       { return i.name }
func (i containerItem) Description() string { return i.status }
func (i containerItem) FilterValue() string { return i.name }

// Create list
func NewContainerList() list.Model {
    items := []list.Item{
        containerItem{name: "apache", status: "Running"},
        containerItem{name: "nginx", status: "Stopped"},
    }
    
    l := list.New(items, list.NewDefaultDelegate(), 20, 10)
    l.Title = "Containers"
    l.SetShowHelp(false)
    l.SetFilteringEnabled(true)
    
    return l
}

// In Update
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    var cmd tea.Cmd
    m.list, cmd = m.list.Update(msg)
    return m, cmd
}

// In View
func (m Model) View() string {
    return m.list.View()
}
```

**List Features**:
- Keyboard navigation (j/k, arrows, page up/down)
- Fuzzy filtering (/)
- Custom delegates for rendering
- Pagination built-in
- Status messages

### Viewport Component

For scrollable text (logs, diffs, markdown):

```go
import "github.com/charmbracelet/bubbles/viewport"

func NewLogViewer() viewport.Model {
    vp := viewport.New(80, 20)
    vp.SetContent("Log line 1\nLog line 2\n...")
    return vp
}

// Update
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    var cmd tea.Cmd
    m.viewport, cmd = m.viewport.Update(msg)
    return m, cmd
}

// Append logs
func (m Model) AppendLog(line string) {
    content := m.viewport.View() + "\n" + line
    m.viewport.SetContent(content)
    m.viewport.GotoBottom() // Auto-scroll
}
```

**Viewport Features**:
- Vertical scrolling
- Mouse wheel support
- High performance mode
- Page up/down, home/end keys
- Goto top/bottom

### Text Input Component

For single-line text entry:

```go
import "github.com/charmbracelet/bubbles/textinput"

func NewSearchInput() textinput.Model {
    ti := textinput.New()
    ti.Placeholder = "Search containers..."
    ti.CharLimit = 100
    ti.Width = 30
    ti.Focus()
    return ti
}

// Update
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    var cmd tea.Cmd
    m.input, cmd = m.input.Update(msg)
    
    // React to value changes
    if m.input.Value() != m.lastSearch {
        m.lastSearch = m.input.Value()
        return m, searchCmd(m.input.Value())
    }
    
    return m, cmd
}
```

**Text Input Features**:
- Unicode support
- Paste support
- Character limit
- Placeholder text
- Masked input (for passwords)
- Echo mode control

### Text Area Component

For multi-line text entry:

```go
import "github.com/charmbracelet/bubbles/textarea"

func NewCommentBox() textarea.Model {
    ta := textarea.New()
    ta.Placeholder = "Enter your comment..."
    ta.CharLimit = 500
    ta.SetWidth(60)
    ta.SetHeight(5)
    ta.Focus()
    return ta
}
```

### Table Component

For tabular data:

```go
import "github.com/charmbracelet/bubbles/table"

func NewServicesTable() table.Model {
    columns := []table.Column{
        {Title: "Name", Width: 15},
        {Title: "Status", Width: 10},
        {Title: "CPU", Width: 8},
        {Title: "Memory", Width: 10},
    }
    
    rows := []table.Row{
        {"apache", "Running", "45%", "128MB"},
        {"nginx", "Stopped", "0%", "0MB"},
    }
    
    t := table.New(
        table.WithColumns(columns),
        table.WithRows(rows),
        table.WithFocused(true),
        table.WithHeight(7),
    )
    
    return t
}
```

**Table Features**:
- Column headers
- Sortable columns
- Row selection
- Scrolling
- Custom styling per cell

### Spinner Component

For loading indicators:

```go
import "github.com/charmbracelet/bubbles/spinner"

func NewSpinner() spinner.Model {
    s := spinner.New()
    s.Spinner = spinner.Dot
    s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
    return s
}

// In Update
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    var cmd tea.Cmd
    m.spinner, cmd = m.spinner.Update(msg)
    return m, cmd
}

// In View
func (m Model) View() string {
    if m.loading {
        return m.spinner.View() + " Loading..."
    }
    return m.content
}
```

**Built-in Spinners**:
- `spinner.Line` - Simple line spinner
- `spinner.Dot` - Dot animation
- `spinner.MiniDot` - Small dot
- `spinner.Jump` - Jumping dot
- `spinner.Pulse` - Pulsing circle
- `spinner.Points` - Growing points
- `spinner.Globe` - Globe rotation
- `spinner.Meter` - Progress meter

### Progress Bar Component

For progress indication:

```go
import "github.com/charmbracelet/bubbles/progress"

func NewProgressBar() progress.Model {
    return progress.New(progress.WithDefaultGradient())
}

// Update progress
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    // Download progress example
    m.progress.SetPercent(m.downloadPercent)
    return m, nil
}
```

### Help Component

For keyboard shortcut help:

```go
import (
    "github.com/charmbracelet/bubbles/help"
    "github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
    Up   key.Binding
    Down key.Binding
    Quit key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
    return []key.Binding{k.Up, k.Down, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
    return [][]key.Binding{
        {k.Up, k.Down},
        {k.Quit},
    }
}

var keys = keyMap{
    Up: key.NewBinding(
        key.WithKeys("up", "k"),
        key.WithHelp("↑/k", "move up"),
    ),
    Down: key.NewBinding(
        key.WithKeys("down", "j"),
        key.WithHelp("↓/j", "move down"),
    ),
    Quit: key.NewBinding(
        key.WithKeys("q", "ctrl+c"),
        key.WithHelp("q", "quit"),
    ),
}

func NewModel() Model {
    return Model{
        help: help.New(),
        keys: keys,
    }
}

func (m Model) View() string {
    return m.help.View(m.keys)
}
```

---

## Lipgloss Styling System

Lipgloss is a CSS-like styling library for terminal UIs. All styling should use Lipgloss, never raw ANSI codes.

### Basic Style Creation

```go
import "github.com/charmbracelet/lipgloss"

// Create a style
var style = lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("#FAFAFA")).
    Background(lipgloss.Color("#7D56F4")).
    Padding(1, 2)

// Render text with style
output := style.Render("Hello, World!")
```

### Color System

**True Color (24-bit)** - Preferred for modern terminals:
```go
lipgloss.Color("#FF0000")  // Hex
lipgloss.Color("#F00")     // Short hex
```

**ANSI 256 (8-bit)** - Fallback for compatibility:
```go
lipgloss.Color("201")  // Magenta
lipgloss.Color("86")   // Cyan
```

**ANSI 16 (4-bit)** - Basic colors:
```go
lipgloss.Color("5")   // Magenta
lipgloss.Color("12")  // Bright blue
```

**Adaptive Colors** - Auto-detect light/dark theme:
```go
lipgloss.AdaptiveColor{
    Light: "#000000",  // Black on light background
    Dark:  "#FFFFFF",  // White on dark background
}
```

**Complete Colors** - Specify all profiles:
```go
lipgloss.CompleteColor{
    TrueColor: "#FF5F87",
    ANSI256:   "204",
    ANSI:      "13",
}
```

### Project Color Palette

Define consistent colors for your application:

```go
package styles

import "github.com/charmbracelet/lipgloss"

var (
    // Status colors
    ColorRunning   = lipgloss.Color("10")   // Green
    ColorStopped   = lipgloss.Color("8")    // Gray
    ColorError     = lipgloss.Color("9")    // Red
    ColorWarning   = lipgloss.Color("11")   // Yellow
    ColorInfo      = lipgloss.Color("12")   // Blue
    
    // UI colors
    ColorPrimary   = lipgloss.Color("#7D56F4")  // Purple
    ColorSecondary = lipgloss.Color("#04B575")  // Green
    ColorAccent    = lipgloss.Color("#FF69B4")  // Pink
    ColorBorder    = lipgloss.Color("240")      // Dark gray
    ColorText      = lipgloss.Color("252")      // Light gray
    ColorMuted     = lipgloss.Color("243")      // Medium gray
)
```

### Inline Formatting

Text formatting options:

```go
var style = lipgloss.NewStyle().
    Bold(true).           // Bold text
    Italic(true).         // Italic text
    Faint(true).          // Dim text
    Underline(true).      // Underline text
    Strikethrough(true).  // Strikethrough text
    Blink(true).          // Blinking text (use sparingly!)
    Reverse(true)         // Reverse fg/bg colors
```

### Block-Level Formatting

Layout and spacing:

```go
var style = lipgloss.NewStyle().
    Width(40).              // Fixed width
    Height(10).             // Fixed height
    MaxWidth(80).           // Maximum width
    MaxHeight(24).          // Maximum height
    
    Padding(1, 2).          // Top/bottom=1, left/right=2
    PaddingTop(2).          // Individual padding
    PaddingRight(4).
    PaddingBottom(2).
    PaddingLeft(4).
    
    Margin(1, 2).           // Top/bottom=1, left/right=2
    MarginTop(1).           // Individual margin
    MarginRight(2).
    MarginBottom(1).
    MarginLeft(2).
    
    Align(lipgloss.Center). // Text alignment
    AlignHorizontal(lipgloss.Right).
    AlignVertical(lipgloss.Bottom)
```

**Shorthand notation** (like CSS):
```go
// All sides
.Padding(2)           // 2 cells all around

// Vertical | Horizontal
.Margin(2, 4)         // 2 top/bottom, 4 left/right

// Top | Horizontal | Bottom
.Padding(1, 4, 2)     // 1 top, 4 sides, 2 bottom

// Top | Right | Bottom | Left (clockwise)
.Margin(1, 2, 3, 4)   // 1 top, 2 right, 3 bottom, 4 left
```

### Borders

**Border Styles**:
```go
lipgloss.NormalBorder()     // ┌─┐ │ └─┘
lipgloss.RoundedBorder()    // ╭─╮ │ ╰─╯
lipgloss.ThickBorder()      // ┏━┓ ┃ ┗━┛
lipgloss.DoubleBorder()     // ╔═╗ ║ ╚═╝
lipgloss.HiddenBorder()     // No border, but spacing
```

**Border Usage**:
```go
var style = lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color("63")).
    BorderBackground(lipgloss.Color("235")).
    BorderTop(true).
    BorderRight(true).
    BorderBottom(true).
    BorderLeft(true)

// Shorthand - all borders
.Border(lipgloss.NormalBorder(), true)

// Shorthand - top and bottom only
.Border(lipgloss.RoundedBorder(), true, false)

// Shorthand - clockwise from top
.Border(lipgloss.ThickBorder(), true, true, false, false)
```

**Custom Borders**:
```go
var customBorder = lipgloss.Border{
    Top:         "─",
    Bottom:      "─",
    Left:        "│",
    Right:       "│",
    TopLeft:     "┌",
    TopRight:    "┐",
    BottomLeft:  "└",
    BottomRight: "┘",
}

style := lipgloss.NewStyle().BorderStyle(customBorder)
```

### Layout Functions

**Horizontal Joining**:
```go
// Join strings horizontally
left := "Column 1"
right := "Column 2"

result := lipgloss.JoinHorizontal(
    lipgloss.Top,    // Align to top
    left,
    right,
)

// Position can be:
// lipgloss.Top, lipgloss.Center, lipgloss.Bottom
// or float64 between 0.0 (top) and 1.0 (bottom)
```

**Vertical Joining**:
```go
// Join strings vertically
header := "Header"
body := "Body content"
footer := "Footer"

result := lipgloss.JoinVertical(
    lipgloss.Left,   // Align to left
    header,
    body,
    footer,
)

// Position can be:
// lipgloss.Left, lipgloss.Center, lipgloss.Right
// or float64 between 0.0 (left) and 1.0 (right)
```

**Placing in Whitespace**:
```go
// Center text in a 80-cell wide space
content := "Hello"
centered := lipgloss.PlaceHorizontal(80, lipgloss.Center, content)

// Place at bottom of 24-line tall space
placed := lipgloss.PlaceVertical(24, lipgloss.Bottom, content)

// Place in 2D space (width, height, horizontal, vertical)
positioned := lipgloss.Place(
    80, 24,
    lipgloss.Center, lipgloss.Center,
    content,
)
```

**Measuring**:
```go
text := "Some rendered text"

width := lipgloss.Width(text)
height := lipgloss.Height(text)

// Or both at once
w, h := lipgloss.Size(text)
```

### Style Composition

**Inheritance**:
```go
baseStyle := lipgloss.NewStyle().
    Foreground(lipgloss.Color("252")).
    Padding(1)

// Inherit unset properties from baseStyle
derivedStyle := lipgloss.NewStyle().
    Background(lipgloss.Color("235")).
    Inherit(baseStyle)
```

**Copying**:
```go
original := lipgloss.NewStyle().Bold(true)

// Simple assignment creates a copy
copy := original

// Modify copy without affecting original
modified := original.Italic(true)
```

**Unsetting**:
```go
style := lipgloss.NewStyle().
    Bold(true).
    UnsetBold().           // Remove bold
    Background(lipgloss.Color("235")).
    UnsetBackground()      // Remove background
```

### Common Style Patterns

**Header Style**:
```go
var headerStyle = lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("#FAFAFA")).
    Background(lipgloss.Color("#7D56F4")).
    Padding(0, 1).
    Width(80)
```

**Panel Style**:
```go
var panelStyle = lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color("240")).
    Padding(1, 2).
    Width(40).
    Height(20)
```

**Selected Row Style**:
```go
var selectedRowStyle = lipgloss.NewStyle().
    Background(lipgloss.Color("240")).
    Foreground(lipgloss.Color("229")).
    Bold(true)
```

**Status Badge Styles**:
```go
func statusBadge(status string) string {
    style := lipgloss.NewStyle().
        Padding(0, 1).
        Bold(true)
    
    switch status {
    case "running":
        style = style.
            Foreground(lipgloss.Color("#000")).
            Background(lipgloss.Color("10"))
    case "stopped":
        style = style.
            Foreground(lipgloss.Color("#FFF")).
            Background(lipgloss.Color("8"))
    case "error":
        style = style.
            Foreground(lipgloss.Color("#FFF")).
            Background(lipgloss.Color("9"))
    }
    
    return style.Render(status)
}
```

---

## Layout Patterns

### 3-Panel Layout

Classic Docker manager layout:

```go
func (m Model) View() string {
    // Left panel - service list (20% width)
    leftWidth := m.width / 5
    leftPanel := panelStyle.
        Width(leftWidth).
        Height(m.height - 5).
        Render(m.renderServiceList())
    
    // Right panel - main content (80% width)
    rightWidth := m.width - leftWidth - 4 // Account for borders
    rightPanel := panelStyle.
        Width(rightWidth).
        Height(m.height - 5).
        Render(m.renderMainContent())
    
    // Join panels horizontally
    panels := lipgloss.JoinHorizontal(
        lipgloss.Top,
        leftPanel,
        rightPanel,
    )
    
    // Footer
    footer := footerStyle.
        Width(m.width).
        Render(m.renderFooter())
    
    // Join everything vertically
    return lipgloss.JoinVertical(
        lipgloss.Left,
        panels,
        footer,
    )
}
```

### Header + Content + Footer

Standard layout:

```go
func (m Model) View() string {
    header := headerStyle.
        Width(m.width).
        Render("20i Stack Manager v1.0")
    
    contentHeight := m.height - 3 // Header + footer
    content := contentStyle.
        Width(m.width).
        Height(contentHeight).
        Render(m.renderContent())
    
    footer := footerStyle.
        Width(m.width).
        Render("Tab: switch | q: quit | ?: help")
    
    return lipgloss.JoinVertical(
        lipgloss.Left,
        header,
        content,
        footer,
    )
}
```

### Responsive Layout

Adapt to terminal size:

```go
func (m Model) View() string {
    // Use different layouts based on width
    if m.width < 80 {
        // Narrow: Stack vertically
        return m.renderNarrowLayout()
    } else if m.width < 120 {
        // Medium: 2-panel
        return m.render2PanelLayout()
    } else {
        // Wide: 3-panel with sidebar
        return m.render3PanelLayout()
    }
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        // Notify subcomponents of resize
        m.list.SetSize(msg.Width, msg.Height-5)
    }
    return m, nil
}
```

### Table Layout

Using Lipgloss table rendering:

```go
import "github.com/charmbracelet/lipgloss/table"

func (m Model) renderServicesTable() string {
    rows := [][]string{
        {"apache", "Running", "45%", "128MB"},
        {"nginx", "Stopped", "0%", "0MB"},
        {"mariadb", "Running", "23%", "512MB"},
    }
    
    headerStyle := lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("205")).
        Align(lipgloss.Center)
    
    cellStyle := lipgloss.NewStyle().
        Padding(0, 1)
    
    oddRowStyle := cellStyle.Foreground(lipgloss.Color("252"))
    evenRowStyle := cellStyle.Foreground(lipgloss.Color("245"))
    
    t := table.New().
        Border(lipgloss.NormalBorder()).
        BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("240"))).
        Headers("SERVICE", "STATUS", "CPU", "MEMORY").
        Rows(rows...).
        StyleFunc(func(row, col int) lipgloss.Style {
            switch {
            case row == table.HeaderRow:
                return headerStyle
            case row%2 == 0:
                return evenRowStyle
            default:
                return oddRowStyle
            }
        })
    
    return t.Render()
}
```

---

## Component Composition

### Parent-Child Communication

**Parent sends message to child**:
```go
type RefreshMsg struct{}

func (parent ParentModel) Update(msg tea.Msg) (ParentModel, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "r" {
            // Refresh child component
            return parent, func() tea.Msg {
                return RefreshMsg{}
            }
        }
    }
    
    // Pass message to child
    var cmd tea.Cmd
    parent.child, cmd = parent.child.Update(msg)
    return parent, cmd
}
```

**Child sends message to parent**:
```go
// Child defines custom message
type SelectionChangedMsg struct {
    SelectedID string
}

// Child sends message
func (child ChildModel) Update(msg tea.Msg) (ChildModel, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "enter" {
            return child, func() tea.Msg {
                return SelectionChangedMsg{
                    SelectedID: child.items[child.cursor].ID,
                }
            }
        }
    }
    return child, nil
}

// Parent receives message
func (parent ParentModel) Update(msg tea.Msg) (ParentModel, tea.Cmd) {
    switch msg := msg.(type) {
    case SelectionChangedMsg:
        parent.selectedID = msg.SelectedID
        return parent, parent.loadDetails(msg.SelectedID)
    }
    
    var cmd tea.Cmd
    parent.child, cmd = parent.child.Update(msg)
    return parent, cmd
}
```

### Batch Commands

Run multiple commands at once:

```go
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "r" {
            // Refresh multiple things
            return m, tea.Batch(
                fetchServicesCmd(),
                fetchStatsCmd(),
                fetchLogsCmd(),
            )
        }
    }
    return m, nil
}
```

### Sequenced Commands

Run commands in order:

```go
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "d" {
            // Delete then refresh
            return m, tea.Sequence(
                deleteContainerCmd(m.selectedID),
                fetchServicesCmd(),
            )
        }
    }
    return m, nil
}
```

### Ticker for Updates

Periodic background updates:

```go
import "time"

func (m Model) Init() tea.Cmd {
    return tea.Batch(
        fetchInitialData(),
        tickCmd(), // Start ticker
    )
}

func tickCmd() tea.Cmd {
    return tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
        return tickMsg(t)
    })
}

type tickMsg time.Time

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tickMsg:
        // Update data
        return m, tea.Batch(
            fetchServicesCmd(),
            tickCmd(), // Schedule next tick
        )
    }
    return m, nil
}
```

---

## Real-World Examples

### Service List with Status

```go
package dashboard

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

type Model struct {
    services []Service
    cursor   int
    width    int
    height   int
}

type Service struct {
    Name   string
    Status string
    CPU    float64
    Memory string
}

var (
    serviceStyle = lipgloss.NewStyle().
        Padding(0, 2)
    
    selectedStyle = serviceStyle.Copy().
        Background(lipgloss.Color("240")).
        Bold(true)
    
    runningStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("10")).
        Bold(true)
    
    stoppedStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("8"))
)

func (m Model) View() string {
    var rows []string
    
    for i, svc := range m.services {
        // Status indicator
        var statusIcon string
        if svc.Status == "Running" {
            statusIcon = runningStyle.Render("●")
        } else {
            statusIcon = stoppedStyle.Render("○")
        }
        
        // Service name
        name := lipgloss.NewStyle().
            Width(20).
            Render(svc.Name)
        
        // CPU bar
        cpuBar := renderProgressBar(svc.CPU, 10)
        
        // Memory
        mem := lipgloss.NewStyle().
            Width(10).
            Render(svc.Memory)
        
        // Combine row
        row := lipgloss.JoinHorizontal(
            lipgloss.Left,
            statusIcon,
            " ",
            name,
            cpuBar,
            mem,
        )
        
        // Apply selection style
        style := serviceStyle
        if i == m.cursor {
            style = selectedStyle
        }
        
        rows = append(rows, style.Render(row))
    }
    
    return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func renderProgressBar(percent float64, width int) string {
    filled := int(percent / 100.0 * float64(width))
    empty := width - filled
    
    bar := strings.Repeat("▓", filled) + strings.Repeat("░", empty)
    
    return lipgloss.NewStyle().
        Foreground(lipgloss.Color("12")).
        Render(bar)
}
```

### Log Viewer with Auto-Scroll

```go
package logs

import (
    "github.com/charmbracelet/bubbles/viewport"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

type Model struct {
    viewport  viewport.Model
    following bool
    logs      []string
}

func New(width, height int) Model {
    vp := viewport.New(width, height)
    vp.MouseWheelEnabled = true
    
    return Model{
        viewport:  vp,
        following: true,
    }
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    var cmd tea.Cmd
    
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "f":
            // Toggle follow mode
            m.following = !m.following
            if m.following {
                m.viewport.GotoBottom()
            }
        case "g":
            // Go to top
            m.following = false
            m.viewport.GotoTop()
        case "G":
            // Go to bottom
            m.following = true
            m.viewport.GotoBottom()
        }
    
    case LogLineMsg:
        m.logs = append(m.logs, msg.Line)
        m.updateContent()
        if m.following {
            m.viewport.GotoBottom()
        }
    }
    
    m.viewport, cmd = m.viewport.Update(msg)
    return m, cmd
}

func (m *Model) updateContent() {
    content := lipgloss.JoinVertical(
        lipgloss.Left,
        m.logs...,
    )
    m.viewport.SetContent(content)
}

func (m Model) View() string {
    followIndicator := ""
    if m.following {
        followIndicator = lipgloss.NewStyle().
            Foreground(lipgloss.Color("10")).
            Render(" [FOLLOWING]")
    }
    
    header := lipgloss.NewStyle().
        Bold(true).
        Render("Logs") + followIndicator
    
    return lipgloss.JoinVertical(
        lipgloss.Left,
        header,
        m.viewport.View(),
    )
}

type LogLineMsg struct {
    Line string
}
```

### Modal Dialog

```go
package components

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

type ConfirmDialog struct {
    title   string
    message string
    visible bool
    result  bool
}

var (
    dialogBoxStyle = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("205")).
        Padding(1, 2).
        Width(50)
    
    buttonStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("252")).
        Background(lipgloss.Color("240")).
        Padding(0, 2).
        Margin(0, 1)
    
    activeButtonStyle = buttonStyle.Copy().
        Foreground(lipgloss.Color("16")).
        Background(lipgloss.Color("205")).
        Bold(true)
)

func NewConfirmDialog(title, message string) ConfirmDialog {
    return ConfirmDialog{
        title:   title,
        message: message,
        visible: true,
    }
}

func (d ConfirmDialog) Update(msg tea.Msg) (ConfirmDialog, tea.Cmd) {
    if !d.visible {
        return d, nil
    }
    
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "y", "enter":
            d.result = true
            d.visible = false
            return d, func() tea.Msg {
                return ConfirmResultMsg{Confirmed: true}
            }
        case "n", "esc":
            d.result = false
            d.visible = false
            return d, func() tea.Msg {
                return ConfirmResultMsg{Confirmed: false}
            }
        }
    }
    
    return d, nil
}

func (d ConfirmDialog) View() string {
    if !d.visible {
        return ""
    }
    
    title := lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("205")).
        Render(d.title)
    
    message := lipgloss.NewStyle().
        MarginTop(1).
        MarginBottom(1).
        Render(d.message)
    
    yesButton := activeButtonStyle.Render("Yes")
    noButton := buttonStyle.Render("No")
    buttons := lipgloss.JoinHorizontal(lipgloss.Left, yesButton, noButton)
    
    content := lipgloss.JoinVertical(
        lipgloss.Left,
        title,
        message,
        buttons,
    )
    
    dialog := dialogBoxStyle.Render(content)
    
    // Center in screen (assuming 80x24)
    return lipgloss.Place(
        80, 24,
        lipgloss.Center, lipgloss.Center,
        dialog,
    )
}

type ConfirmResultMsg struct {
    Confirmed bool
}
```

---

## Anti-Patterns to Avoid

### ❌ Don't: Mutate Model Without Returning

```go
// BAD
func (m *Model) Update(msg tea.Msg) {
    m.cursor++ // Bubble Tea won't detect change
}

// GOOD
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    m.cursor++
    return m, nil
}
```

### ❌ Don't: Block in Update

```go
// BAD - Freezes UI
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    m.data = fetchFromDocker() // HTTP call blocks
    return m, nil
}

// GOOD - Use tea.Cmd
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    return m, fetchDataCmd()
}

func fetchDataCmd() tea.Cmd {
    return func() tea.Msg {
        data := fetchFromDocker() // Runs in background
        return dataMsg{data: data}
    }
}
```

### ❌ Don't: Use Raw ANSI Codes

```go
// BAD
output := "\033[1;32mRunning\033[0m"

// GOOD - Use Lipgloss
output := lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("10")).
    Render("Running")
```

### ❌ Don't: Create Styles in View

```go
// BAD - Creates new styles every render
func (m Model) View() string {
    style := lipgloss.NewStyle().Bold(true) // Wasteful
    return style.Render("Text")
}

// GOOD - Define styles once
var boldStyle = lipgloss.NewStyle().Bold(true)

func (m Model) View() string {
    return boldStyle.Render("Text")
}
```

### ❌ Don't: Forget WindowSizeMsg

```go
// BAD - Ignores terminal resize
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    // Never handles tea.WindowSizeMsg
    return m, nil
}

// GOOD - Handle resize
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
    }
    return m, nil
}
```

---

## Quick Reference

### Bubble Tea Essentials

```go
// Create program
p := tea.NewProgram(initialModel())
if _, err := p.Run(); err != nil {
    log.Fatal(err)
}

// Model interface
type tea.Model interface {
    Init() tea.Cmd
    Update(tea.Msg) (tea.Model, tea.Cmd)
    View() string
}

// Common commands
tea.Quit                    // Exit program
tea.Batch(cmd1, cmd2)       // Run commands in parallel
tea.Sequence(cmd1, cmd2)    // Run commands in order
tea.Tick(duration, func)    // Periodic timer

// Common messages
tea.KeyMsg                  // Keyboard input
tea.MouseMsg                // Mouse input
tea.WindowSizeMsg           // Terminal resize
```

### Lipgloss Essentials

```go
// Create style
style := lipgloss.NewStyle().
    Foreground(lipgloss.Color("205")).
    Background(lipgloss.Color("235")).
    Bold(true).
    Padding(1, 2).
    Margin(1).
    Border(lipgloss.RoundedBorder()).
    Width(40).
    Align(lipgloss.Center)

// Render
output := style.Render("Text")

// Layout
lipgloss.JoinHorizontal(pos, str1, str2)
lipgloss.JoinVertical(pos, str1, str2)
lipgloss.Place(w, h, posX, posY, str)

// Measure
w := lipgloss.Width(str)
h := lipgloss.Height(str)
```

### Bubbles Components

```go
import "github.com/charmbracelet/bubbles/..."

list.New()          // Scrollable list
viewport.New()      // Scrollable text
textinput.New()     // Single-line input
textarea.New()      // Multi-line input
table.New()         // Table with columns
spinner.New()       // Loading spinner
progress.New()      // Progress bar
help.New()          // Help view
```

### File Structure Template

```
internal/tui/
  app.go           # Root model
  components/
    dashboard/
      dashboard.go   # Component model
      styles.go      # Lipgloss styles
      messages.go    # Custom messages
    containers/
      containers.go
      list.go
      styles.go
  styles/
    colors.go        # Color palette
    common.go        # Shared styles
```

---

## Additional Resources

- **Bubble Tea Docs**: https://github.com/charmbracelet/bubbletea
- **Lipgloss Docs**: https://github.com/charmbracelet/lipgloss
- **Bubbles Library**: https://github.com/charmbracelet/bubbles
- **Examples**: https://github.com/charmbracelet/bubbletea/tree/main/examples
- **Tutorial**: https://github.com/charmbracelet/bubbletea/tree/main/tutorials

---

**End of Guide**
