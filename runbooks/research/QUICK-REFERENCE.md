# Bubble Tea Quick Reference Card

**For agents building TUI components - keep this open while coding!**

---

## Bubble Tea Basics

```go
// Model interface (all components must implement)
type tea.Model interface {
    Init() tea.Cmd
    Update(tea.Msg) (tea.Model, tea.Cmd)
    View() string
}

// Start program
p := tea.NewProgram(initialModel())
p.Run()

// Common commands
tea.Quit                         // Exit
tea.Batch(cmd1, cmd2)            // Parallel
tea.Sequence(cmd1, cmd2)         // Sequential
tea.Tick(2*time.Second, func)    // Timer
```

## Essential Messages

```go
tea.KeyMsg          // Keyboard input
tea.MouseMsg        // Mouse events
tea.WindowSizeMsg   // Terminal resize
```

## Common Update Pattern

```go
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "q", "ctrl+c":
            return m, tea.Quit
        case "up", "k":
            if m.cursor > 0 {
                m.cursor--
            }
        case "down", "j":
            if m.cursor < len(m.items)-1 {
                m.cursor++
            }
        }
    
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
    }
    
    return m, nil
}
```

---

## Lipgloss Basics

```go
// Create style
var style = lipgloss.NewStyle().
    Foreground(lipgloss.Color("205")).
    Background(lipgloss.Color("235")).
    Bold(true).
    Padding(1, 2).
    Width(40).
    Border(lipgloss.RoundedBorder())

// Render
output := style.Render("Text")
```

## Quick Color Guide

```go
// Status
ColorRunning  = lipgloss.Color("10")   // Green
ColorStopped  = lipgloss.Color("8")    // Gray
ColorError    = lipgloss.Color("9")    // Red
ColorWarning  = lipgloss.Color("11")   // Yellow
ColorInfo     = lipgloss.Color("12")   // Blue

// UI
ColorPrimary  = lipgloss.Color("#7D56F4")  // Purple
ColorBorder   = lipgloss.Color("240")      // Dark gray
ColorText     = lipgloss.Color("252")      // Light gray
```

## Layout Functions

```go
// Horizontal
lipgloss.JoinHorizontal(lipgloss.Top, str1, str2)

// Vertical
lipgloss.JoinVertical(lipgloss.Left, str1, str2)

// Center
lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, str)

// Measure
w := lipgloss.Width(str)
h := lipgloss.Height(str)
```

---

## Standard Styles

```go
// Header
var HeaderStyle = lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("#FAFAFA")).
    Background(ColorPrimary).
    Padding(0, 1)

// Footer
var FooterStyle = lipgloss.NewStyle().
    Foreground(ColorMuted).
    Padding(0, 1)

// Panel
var PanelStyle = lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(ColorBorder).
    Padding(1, 2)

// Selected Row
var SelectedRowStyle = lipgloss.NewStyle().
    Background(ColorSelection).
    Bold(true)
```

---

## Bubbles Components

```go
import "github.com/charmbracelet/bubbles/..."

// List
list.New(items, delegate, width, height)

// Viewport (logs)
viewport.New(width, height)

// Text Input
textinput.New()

// Text Area
textarea.New()

// Table
table.New()

// Spinner
spinner.New()

// Progress
progress.New()

// Help
help.New()
```

---

## Unicode Icons

```go
"●"  // Running/active
"○"  // Stopped/inactive
"✓"  // Success
"✗"  // Error
"⚠"  // Warning
"ℹ"  // Info
"→"  // Arrow right
"↑"  // Arrow up
"↓"  // Arrow down
"⟳"  // Refresh
"▓"  // Filled block
"░"  // Empty block
```

---

## Common Patterns

### Service List

```go
for i, svc := range m.services {
    icon := "●"
    if svc.Status != "running" {
        icon = "○"
    }
    
    style := rowStyle
    if i == m.cursor {
        style = selectedRowStyle
    }
    
    row := fmt.Sprintf("%s %s", icon, svc.Name)
    rows = append(rows, style.Render(row))
}
```

### Status Badge

```go
func StatusBadge(status string) string {
    var style = lipgloss.NewStyle().
        Padding(0, 1).
        Bold(true)
    
    switch status {
    case "running":
        style = style.Background(ColorRunning)
    case "stopped":
        style = style.Background(ColorStopped)
    case "error":
        style = style.Background(ColorError)
    }
    
    return style.Render(status)
}
```

### Progress Bar

```go
func ProgressBar(percent float64, width int) string {
    filled := int(percent / 100.0 * float64(width))
    bar := strings.Repeat("▓", filled) + 
           strings.Repeat("░", width-filled)
    
    return lipgloss.NewStyle().
        Foreground(ColorInfo).
        Render(bar)
}
```

---

## Anti-Patterns

```go
// ❌ DON'T: Block in Update
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    m.data = fetchData() // BLOCKS UI!
    return m, nil
}

// ✓ DO: Use commands
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    return m, fetchDataCmd()
}

// ❌ DON'T: Create styles in View
func (m Model) View() string {
    style := lipgloss.NewStyle().Bold(true) // Wasteful
    return style.Render(m.text)
}

// ✓ DO: Define once
var boldStyle = lipgloss.NewStyle().Bold(true)
func (m Model) View() string {
    return boldStyle.Render(m.text)
}
```

---

## File Structure

```
components/myview/
  myview.go    # Model + Init/Update/View
  styles.go    # Lipgloss styles
  messages.go  # Custom messages
```

---

## Cheat Sheet URLs

- Bubble Tea: https://github.com/charmbracelet/bubbletea
- Lipgloss: https://github.com/charmbracelet/lipgloss
- Bubbles: https://github.com/charmbracelet/bubbles
- Examples: https://github.com/charmbracelet/bubbletea/tree/main/examples

---

**Keep this card handy! Print it, pin it, or keep it open in a split pane.**
