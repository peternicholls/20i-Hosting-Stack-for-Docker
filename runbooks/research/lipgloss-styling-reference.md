# Lipgloss Styling Reference for Agents

**Quick reference for styling Bubble Tea TUI applications**

---

## Core Principles

1. **Always use Lipgloss** - Never use raw ANSI codes
2. **Define styles once** - At package level, not in View functions
3. **Use color palette** - Define project colors as constants
4. **Measure before placing** - Use `lipgloss.Width()` and `Height()` for layouts
5. **Responsive by default** - Handle `tea.WindowSizeMsg` and adapt to terminal size

---

## Standard Color Palette

```go
package styles

import "github.com/charmbracelet/lipgloss"

// Status Colors
var (
    ColorRunning   = lipgloss.Color("10")  // Green - ✓ success, running
    ColorStopped   = lipgloss.Color("8")   // Gray - stopped, disabled
    ColorError     = lipgloss.Color("9")   // Red - errors, critical
    ColorWarning   = lipgloss.Color("11")  // Yellow - warnings, attention
    ColorInfo      = lipgloss.Color("12")  // Blue - info, help
)

// UI Colors
var (
    ColorPrimary   = lipgloss.Color("#7D56F4")  // Purple - headers, focus
    ColorSecondary = lipgloss.Color("#04B575")  // Green - accents
    ColorBorder    = lipgloss.Color("240")      // Dark gray - borders
    ColorText      = lipgloss.Color("252")      // Light gray - text
    ColorMuted     = lipgloss.Color("243")      // Medium gray - hints
    ColorSelection = lipgloss.Color("240")      // Selection background
)

// Adaptive Colors (auto light/dark)
var AdaptiveText = lipgloss.AdaptiveColor{
    Light: "#1a1a1a",
    Dark:  "#f5f5f5",
}
```

---

## Common Style Patterns

### Header Style

```go
var HeaderStyle = lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("#FAFAFA")).
    Background(ColorPrimary).
    Padding(0, 1).
    Width(80)  // Set in Update based on terminal width
```

**Usage**:
```go
header := HeaderStyle.Render("20i Stack Manager v1.0")
```

### Footer/Help Style

```go
var FooterStyle = lipgloss.NewStyle().
    Foreground(ColorMuted).
    Padding(0, 1)
```

**Usage**:
```go
footer := FooterStyle.Render("Tab: switch | q: quit | ?: help")
```

### Panel/Box Style

```go
var PanelStyle = lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(ColorBorder).
    Padding(1, 2)
```

**Usage**:
```go
panel := PanelStyle.
    Width(40).
    Height(20).
    Render(content)
```

### List Row Styles

```go
var (
    RowStyle = lipgloss.NewStyle().
        Padding(0, 2)
    
    SelectedRowStyle = RowStyle.Copy().
        Background(ColorSelection).
        Bold(true)
    
    AlternateRowStyle = RowStyle.Copy().
        Foreground(ColorMuted)
)
```

**Usage**:
```go
for i, item := range items {
    style := RowStyle
    if i == cursor {
        style = SelectedRowStyle
    } else if i%2 == 1 {
        style = AlternateRowStyle
    }
    rows = append(rows, style.Render(item))
}
```

### Status Badge Styles

```go
func StatusBadge(status string) string {
    var style = lipgloss.NewStyle().
        Padding(0, 1).
        Bold(true)
    
    switch status {
    case "running":
        style = style.
            Foreground(lipgloss.Color("#000")).
            Background(ColorRunning)
    case "stopped":
        style = style.
            Foreground(ColorText).
            Background(ColorStopped)
    case "error":
        style = style.
            Foreground(lipgloss.Color("#FFF")).
            Background(ColorError)
    default:
        style = style.
            Foreground(ColorText).
            Background(ColorMuted)
    }
    
    return style.Render(strings.ToUpper(status))
}
```

### Error/Warning Message Styles

```go
var (
    ErrorStyle = lipgloss.NewStyle().
        Foreground(ColorError).
        Bold(true).
        Padding(0, 1)
    
    WarningStyle = lipgloss.NewStyle().
        Foreground(ColorWarning).
        Bold(true).
        Padding(0, 1)
    
    InfoStyle = lipgloss.NewStyle().
        Foreground(ColorInfo).
        Padding(0, 1)
)
```

**Usage**:
```go
// With icons
error := ErrorStyle.Render("✗ Failed to start container")
warning := WarningStyle.Render("⚠ Container using deprecated image")
info := InfoStyle.Render("ℹ Refreshing in 2s...")
```

### Button Styles

```go
var (
    ButtonStyle = lipgloss.NewStyle().
        Foreground(ColorText).
        Background(lipgloss.Color("238")).
        Padding(0, 2).
        Margin(0, 1)
    
    ActiveButtonStyle = ButtonStyle.Copy().
        Foreground(lipgloss.Color("#FFF")).
        Background(ColorPrimary).
        Bold(true)
)
```

**Usage**:
```go
yesBtn := ActiveButtonStyle.Render("Yes")
noBtn := ButtonStyle.Render("No")
buttons := lipgloss.JoinHorizontal(lipgloss.Left, yesBtn, noBtn)
```

---

## Layout Patterns

### 3-Panel Layout (Dashboard)

```go
func (m Model) View() string {
    // Calculate widths (with room for borders/padding)
    leftWidth := m.width / 5
    rightWidth := m.width - leftWidth - 4
    contentHeight := m.height - 5  // Header + footer
    
    // Left panel - navigation list
    leftPanel := PanelStyle.
        Width(leftWidth).
        Height(contentHeight).
        Render(m.renderServiceList())
    
    // Right panel - main content
    rightPanel := PanelStyle.
        Width(rightWidth).
        Height(contentHeight).
        Render(m.renderMainContent())
    
    // Join panels
    panels := lipgloss.JoinHorizontal(
        lipgloss.Top,
        leftPanel,
        rightPanel,
    )
    
    // Header
    header := HeaderStyle.
        Width(m.width).
        Render("20i Stack Manager")
    
    // Footer
    footer := FooterStyle.
        Width(m.width).
        Render("Tab: panels | q: quit")
    
    // Combine all
    return lipgloss.JoinVertical(
        lipgloss.Left,
        header,
        panels,
        footer,
    )
}
```

### Centered Modal Dialog

```go
func (m Model) renderDialog() string {
    // Dialog content
    title := lipgloss.NewStyle().
        Bold(true).
        Foreground(ColorPrimary).
        Render(m.dialog.Title)
    
    message := lipgloss.NewStyle().
        MarginTop(1).
        MarginBottom(1).
        Width(40).
        Render(m.dialog.Message)
    
    buttons := lipgloss.JoinHorizontal(
        lipgloss.Left,
        ActiveButtonStyle.Render("OK"),
        ButtonStyle.Render("Cancel"),
    )
    
    content := lipgloss.JoinVertical(
        lipgloss.Left,
        title,
        message,
        buttons,
    )
    
    // Wrap in box
    dialog := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(ColorPrimary).
        Padding(1, 2).
        Width(50).
        Render(content)
    
    // Center on screen
    return lipgloss.Place(
        m.width,
        m.height,
        lipgloss.Center,
        lipgloss.Center,
        dialog,
    )
}
```

### Responsive Table

```go
import "github.com/charmbracelet/lipgloss/table"

func (m Model) renderTable() string {
    // Define styles
    headerStyle := lipgloss.NewStyle().
        Bold(true).
        Foreground(ColorPrimary).
        Align(lipgloss.Center)
    
    cellStyle := lipgloss.NewStyle().
        Padding(0, 1)
    
    evenRowStyle := cellStyle.Foreground(ColorText)
    oddRowStyle := cellStyle.Foreground(ColorMuted)
    
    // Build table
    t := table.New().
        Border(lipgloss.NormalBorder()).
        BorderStyle(lipgloss.NewStyle().Foreground(ColorBorder)).
        Headers("SERVICE", "STATUS", "CPU", "MEMORY").
        Rows(m.rows...).
        Width(m.width - 4).  // Account for borders
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

### Progress Bar

```go
func RenderProgressBar(percent float64, width int, label string) string {
    filled := int(percent / 100.0 * float64(width))
    empty := width - filled
    
    // Bar segments
    filledBar := strings.Repeat("▓", filled)
    emptyBar := strings.Repeat("░", empty)
    
    // Color based on percentage
    var color lipgloss.Color
    switch {
    case percent >= 90:
        color = ColorError
    case percent >= 70:
        color = ColorWarning
    default:
        color = ColorInfo
    }
    
    // Styled bar
    bar := lipgloss.NewStyle().
        Foreground(color).
        Render(filledBar + emptyBar)
    
    // Label
    labelText := lipgloss.NewStyle().
        Width(8).
        Align(lipgloss.Right).
        Render(fmt.Sprintf("%.1f%%", percent))
    
    // Optional prefix
    prefix := ""
    if label != "" {
        prefix = lipgloss.NewStyle().
            Width(10).
            Render(label)
    }
    
    return lipgloss.JoinHorizontal(
        lipgloss.Left,
        prefix,
        bar,
        " ",
        labelText,
    )
}
```

**Usage**:
```go
cpuBar := RenderProgressBar(45.3, 20, "CPU:")
memBar := RenderProgressBar(67.8, 20, "Memory:")
```

---

## Border Reference

```go
// Built-in border styles
lipgloss.NormalBorder()     // ┌─┐│└─┘
lipgloss.RoundedBorder()    // ╭─╮│╰─╯
lipgloss.ThickBorder()      // ┏━┓┃┗━┛
lipgloss.DoubleBorder()     // ╔═╗║╚═╝
lipgloss.HiddenBorder()     // Invisible but spacing preserved

// Custom border
customBorder := lipgloss.Border{
    Top:         "─",
    Bottom:      "─",
    Left:        "│",
    Right:       "│",
    TopLeft:     "╭",
    TopRight:    "╮",
    BottomLeft:  "╰",
    BottomRight: "╯",
}

// Apply border
style := lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(ColorBorder)

// Selective borders (top and bottom only)
style := lipgloss.NewStyle().
    Border(lipgloss.NormalBorder(), true, false)
    
// Clockwise from top: top, right, bottom, left
style := lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder(), true, true, false, false)
```

---

## Alignment Reference

```go
// Horizontal alignment
lipgloss.Left
lipgloss.Center
lipgloss.Right

// Vertical alignment
lipgloss.Top
lipgloss.Center
lipgloss.Bottom

// Fractional positions (0.0 to 1.0)
0.0   // Top/Left
0.25  // Quarter
0.5   // Center
0.75  // Three quarters
1.0   // Bottom/Right

// Usage in styles
style := lipgloss.NewStyle().
    Width(40).
    Align(lipgloss.Center)  // Center text horizontally

// Usage in layout functions
result := lipgloss.JoinHorizontal(
    lipgloss.Bottom,  // Align to bottom edge
    col1,
    col2,
)

centered := lipgloss.Place(
    80, 24,
    lipgloss.Center,  // Horizontal center
    lipgloss.Center,  // Vertical center
    content,
)
```

---

## Spacing Reference

```go
// Padding (inside border)
.Padding(1)           // All sides: 1
.Padding(1, 2)        // Vertical: 1, Horizontal: 2
.Padding(1, 2, 3)     // Top: 1, Horizontal: 2, Bottom: 3
.Padding(1, 2, 3, 4)  // Top: 1, Right: 2, Bottom: 3, Left: 4

// Individual padding
.PaddingTop(1)
.PaddingRight(2)
.PaddingBottom(1)
.PaddingLeft(2)

// Margin (outside border)
.Margin(1)            // All sides: 1
.Margin(1, 2)         // Vertical: 1, Horizontal: 2
.Margin(1, 2, 3)      // Top: 1, Horizontal: 2, Bottom: 3
.Margin(1, 2, 3, 4)   // Top: 1, Right: 2, Bottom: 3, Left: 4

// Individual margin
.MarginTop(1)
.MarginRight(2)
.MarginBottom(1)
.MarginLeft(2)
```

---

## Unicode Icons Reference

Use these for visual indicators:

```go
// Status indicators
"●"  // Filled circle - running/active
"○"  // Empty circle - stopped/inactive
"◉"  // Circle with dot - partially active
"◐"  // Half circle - transitioning

// Arrows
"→"  // Right arrow
"←"  // Left arrow
"↑"  // Up arrow
"↓"  // Down arrow
"⇄"  // Bidirectional

// Symbols
"✓"  // Check mark - success
"✗"  // X mark - error
"⚠"  // Warning triangle
"ℹ"  // Info
"⟳"  // Refresh/reload
"⋯"  // Ellipsis - loading

// Progress
"▓"  // Filled block (for progress bars)
"░"  // Empty block (for progress bars)
"▁▂▃▄▅▆▇█"  // Vertical bars (for graphs)

// Brackets
"「」" // Japanese brackets
"『』" // Double Japanese brackets
"【】" // Lenticular brackets

// Box drawing
"─│┌┐└┘├┤┬┴┼"  // Normal
"═║╔╗╚╝╠╣╦╩╬"  // Double
"━┃┏┓┗┛┣┫┳┻╋"  // Thick
```

---

## Color Reference

### ANSI 16 Colors (4-bit)

```go
"0"   // Black
"1"   // Red
"2"   // Green
"3"   // Yellow
"4"   // Blue
"5"   // Magenta
"6"   // Cyan
"7"   // White
"8"   // Bright Black (Gray)
"9"   // Bright Red
"10"  // Bright Green
"11"  // Bright Yellow
"12"  // Bright Blue
"13"  // Bright Magenta
"14"  // Bright Cyan
"15"  // Bright White
```

### Common ANSI 256 Colors (8-bit)

```go
// Grayscale
"232" // Very dark gray
"238" // Dark gray
"240" // Medium-dark gray
"243" // Medium gray
"246" // Medium-light gray
"250" // Light gray
"255" // Almost white

// Vibrant colors
"196" // Red
"202" // Orange
"226" // Yellow
"46"  // Green
"51"  // Cyan
"201" // Magenta
"205" // Pink
"99"  // Purple
```

### Hex Colors (TrueColor)

```go
// Primary palette
"#7D56F4" // Purple
"#FF69B4" // Hot pink
"#04B575" // Green
"#F25D94" // Pink
"#FFA500" // Orange

// Background palette
"#1E1E1E" // Dark background
"#2D2D2D" // Panel background
"#3C3C3C" // Elevated surface
"#F5F5F5" // Light background
```

---

## Common Gotchas

### ❌ Creating Styles in View

```go
// BAD - Creates new styles every render (wasteful)
func (m Model) View() string {
    style := lipgloss.NewStyle().Bold(true)
    return style.Render(m.text)
}

// GOOD - Define once at package level
var boldStyle = lipgloss.NewStyle().Bold(true)

func (m Model) View() string {
    return boldStyle.Render(m.text)
}
```

### ❌ Hardcoding Dimensions

```go
// BAD - Won't adapt to terminal resize
style := lipgloss.NewStyle().Width(80).Height(24)

// GOOD - Use model dimensions
func (m Model) View() string {
    style := panelStyle.
        Width(m.width - 4).
        Height(m.height - 5)
    return style.Render(content)
}
```

### ❌ Forgetting to Measure

```go
// BAD - May overflow or misalign
content := lipgloss.JoinHorizontal(
    lipgloss.Left,
    leftPanel,  // Unknown width
    rightPanel, // Unknown width
)

// GOOD - Measure and calculate
leftWidth := lipgloss.Width(leftPanel)
rightWidth := m.width - leftWidth - 4
rightPanel := panelStyle.Width(rightWidth).Render(content)
```

### ❌ Not Handling WindowSizeMsg

```go
// BAD - Ignores terminal resize
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    // Never handles WindowSizeMsg
    return m, nil
}

// GOOD - Update dimensions
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        // Notify subcomponents
        m.list.SetSize(msg.Width, msg.Height-5)
    }
    return m, nil
}
```

---

## Testing Styles

```go
func TestStyles(t *testing.T) {
    // Create renderer with fixed width (for testing)
    r := lipgloss.NewRenderer(os.Stdout)
    
    style := r.NewStyle().
        Width(20).
        Padding(1)
    
    output := style.Render("Test")
    
    // Verify dimensions
    if lipgloss.Width(output) != 20 {
        t.Errorf("Expected width 20, got %d", lipgloss.Width(output))
    }
}
```

---

## Quick Checklist for Agents

When styling a TUI component:

- [ ] Define color palette at package level
- [ ] Define all styles outside View function
- [ ] Use semantic color names (ColorRunning, ColorError)
- [ ] Handle `tea.WindowSizeMsg` for responsive layout
- [ ] Measure content with `lipgloss.Width/Height()` before laying out
- [ ] Use `lipgloss.Place()` for centering modals
- [ ] Use `lipgloss.JoinHorizontal/Vertical()` for layouts
- [ ] Apply consistent borders (typically RoundedBorder or NormalBorder)
- [ ] Use unicode icons (✓✗⚠ℹ) for visual indicators
- [ ] Test with minimum terminal size (80x24)
- [ ] Consider both light and dark themes (use AdaptiveColor)

---

**End of Reference**
