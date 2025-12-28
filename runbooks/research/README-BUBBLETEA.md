# Bubble Tea & Lipgloss Research - Summary

**Research Completed**: 2025-12-28  
**Status**: ✓ Complete  
**Purpose**: Provide comprehensive guidance for agents building TUI applications with Bubble Tea and Lipgloss

---

## Documents Created

### 1. [Bubble Tea Component & View Guide](./bubbletea-component-guide.md)

**Comprehensive guide covering**:
- Core concepts (Elm Architecture, Model-Update-View)
- Component architecture and patterns
- Bubbles component library (list, viewport, textinput, table, etc.)
- Lipgloss styling system fundamentals
- Layout patterns (3-panel, header+content+footer, responsive)
- Component composition techniques
- Real-world examples (service lists, log viewers, modals)
- Anti-patterns to avoid

**Audience**: Agents implementing TUI features  
**Length**: ~500 lines with extensive code examples  
**Depth**: Deep technical reference with practical patterns

### 2. [Lipgloss Styling Reference](./lipgloss-styling-reference.md)

**Quick reference covering**:
- Standard color palette for 20i Stack Manager
- Common style patterns (headers, footers, panels, buttons)
- Layout patterns with code examples
- Border, alignment, and spacing references
- Unicode icon library
- Color system reference (ANSI 16/256, TrueColor)
- Common gotchas and solutions
- Testing guidance

**Audience**: Agents styling UI components  
**Length**: ~300 lines, focused on quick lookup  
**Depth**: Practical patterns with copy-paste examples

---

## Key Findings

### Architecture Best Practices

1. **Component Isolation**: Each view is a standalone Bubble Tea model with Init/Update/View
2. **Root Coordinator**: Top-level model delegates messages to active child components
3. **Message-Driven**: All state changes happen via messages, no direct mutation
4. **Commands for I/O**: Use `tea.Cmd` for async operations (API calls, timers)
5. **No Blocking**: Never block in Update() - use background goroutines

### Styling Best Practices

1. **Lipgloss Always**: Never use raw ANSI codes
2. **Define Once**: Create styles at package level, not in View()
3. **Color Palette**: Use consistent, semantic color names
4. **Responsive**: Handle `tea.WindowSizeMsg` and adapt layouts
5. **Measure First**: Use `lipgloss.Width/Height()` before laying out

### Component Library Recommendations

**Use Bubbles components for**:
- **Lists**: `bubbles/list` - Scrollable, filterable, with pagination
- **Logs**: `bubbles/viewport` - Efficient scrolling, auto-scroll support
- **Input**: `bubbles/textinput` - Single-line text entry with validation
- **Tables**: `bubbles/table` - Sortable columns, row selection
- **Loading**: `bubbles/spinner` - Multiple animation styles
- **Help**: `bubbles/help` - Auto-generated from keybindings

**Build custom for**:
- Complex multi-panel layouts
- Domain-specific visualizations (graphs, charts)
- Custom interaction patterns

---

## Quick Start for Agents

### Creating a New Component

1. Create package structure:
   ```
   components/containers/
     containers.go  # Model + Init/Update/View
     styles.go      # Lipgloss styles
     messages.go    # Custom message types
   ```

2. Define model:
   ```go
   type Model struct {
       items   []Item
       cursor  int
       width   int
       height  int
   }
   ```

3. Implement interface:
   ```go
   func (m Model) Init() tea.Cmd { ... }
   func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) { ... }
   func (m Model) View() string { ... }
   ```

4. Define styles in `styles.go`:
   ```go
   var panelStyle = lipgloss.NewStyle().
       Border(lipgloss.RoundedBorder()).
       BorderForeground(lipgloss.Color("240"))
   ```

### Styling a View

1. Use color palette from reference
2. Apply consistent borders and spacing
3. Handle responsive layout
4. Add unicode icons for clarity

```go
func (m Model) View() string {
    // Header
    header := headerStyle.Width(m.width).Render("Services")
    
    // Content with status badges
    rows := []string{}
    for i, svc := range m.services {
        icon := "●"  // Running
        if svc.Status != "running" {
            icon = "○"  // Stopped
        }
        
        style := rowStyle
        if i == m.cursor {
            style = selectedRowStyle
        }
        
        row := fmt.Sprintf("%s %s", icon, svc.Name)
        rows = append(rows, style.Render(row))
    }
    
    content := lipgloss.JoinVertical(lipgloss.Left, rows...)
    
    // Combine
    return lipgloss.JoinVertical(lipgloss.Left, header, content)
}
```

---

## Integration with 20i Stack Manager

### Recommended Component Structure

```
internal/tui/
  app.go                    # Root model (view coordinator)
  
  components/
    dashboard/
      dashboard.go          # Dashboard view
      styles.go             # Dashboard-specific styles
      
    services/
      services.go           # Service list component
      detail.go             # Service detail panel
      styles.go
      
    logs/
      logs.go               # Log viewer with auto-scroll
      styles.go
      
    help/
      help.go               # Help modal
      styles.go
  
  styles/
    colors.go               # Global color palette
    common.go               # Shared styles (headers, footers)
    
  docker/
    client.go               # Docker API wrapper
    messages.go             # Docker-related messages
```

### Color Palette for 20i

Already defined in styling reference:
- **Running**: Green (#10)
- **Stopped**: Gray (#8)
- **Error**: Red (#9)
- **Warning**: Yellow (#11)
- **Primary**: Purple (#7D56F4)
- **Border**: Dark Gray (#240)

---

## Common Patterns

### Pattern 1: Service List with Selection

See `bubbletea-component-guide.md` section "Real-World Examples → Service List with Status"

### Pattern 2: Log Viewer with Auto-Scroll

See `bubbletea-component-guide.md` section "Real-World Examples → Log Viewer with Auto-Scroll"

### Pattern 3: Modal Confirmation Dialog

See `bubbletea-component-guide.md` section "Real-World Examples → Modal Dialog"

### Pattern 4: 3-Panel Dashboard Layout

See `lipgloss-styling-reference.md` section "Layout Patterns → 3-Panel Layout"

---

## Testing Recommendations

1. **Unit test models**: Test Update() with various messages
2. **Test rendering**: Verify View() output format
3. **Test dimensions**: Ensure responsive behavior
4. **Integration test**: Test component composition

Example:
```go
func TestServiceListUpdate(t *testing.T) {
    m := NewServiceList()
    
    // Test cursor movement
    m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
    if m.cursor != 1 {
        t.Errorf("Expected cursor=1, got %d", m.cursor)
    }
}
```

---

## Next Steps for Implementation

1. ✓ Read both guides thoroughly
2. Start with simple component (e.g., service list)
3. Implement Model + Init/Update/View
4. Define styles using color palette
5. Test in isolation before integration
6. Compose with other components
7. Add real-time updates (Docker API)
8. Polish with proper error handling

---

## Resources

- **Bubble Tea**: https://github.com/charmbracelet/bubbletea
- **Lipgloss**: https://github.com/charmbracelet/lipgloss
- **Bubbles**: https://github.com/charmbracelet/bubbles
- **Examples**: https://github.com/charmbracelet/bubbletea/tree/main/examples
- **Video Tutorials**: https://charm.sh/yt

---

## Research Notes

This research builds on the earlier "TUI Excellence Study" (R01) which analyzed lazydocker, lazygit, k9s, and gh-dash. The key findings from that research:

- **3-panel layout wins** for Docker managers (navigation + content + detail)
- **Real-time updates** must not block UI (use background goroutines)
- **Keyboard-first navigation** with vim bindings as enhancement
- **Progressive disclosure** (summary in list, detail on selection)
- **Context-aware help** in footer with `?` for full help modal

These patterns are now codified in the component and styling guides for agents to follow.

---

**End of Summary**
