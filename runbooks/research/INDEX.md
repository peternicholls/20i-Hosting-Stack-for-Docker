# Bubble Tea & Lipgloss Research Index

**Research Area**: TUI Framework & Styling  
**Completed**: 2025-12-28  
**Status**: ✓ Complete and Ready for Implementation

---

## 📚 Documentation Structure

### For Quick Reference

**[QUICK-REFERENCE.md](./QUICK-REFERENCE.md)** ⚡  
**Best for**: Agents actively coding - keep this open while implementing  
**Contains**: Copy-paste patterns, common styles, anti-patterns, essential APIs  
**Size**: 1-page cheat sheet

### For Learning & Implementation

**[README-BUBBLETEA.md](./README-BUBBLETEA.md)** 📖  
**Best for**: Overview and getting started  
**Contains**: Research summary, document index, integration guide, next steps  
**Size**: 5-minute read

**[bubbletea-component-guide.md](./bubbletea-component-guide.md)** 🔧  
**Best for**: Deep technical reference when building components  
**Contains**: 
- Core concepts (Elm Architecture)
- Component patterns and structure
- Bubbles component library usage
- Real-world examples with full code
- Layout patterns
- Component composition techniques
**Size**: Comprehensive technical guide (~500 lines)

**[lipgloss-styling-reference.md](./lipgloss-styling-reference.md)** 🎨  
**Best for**: Styling components and creating consistent UI  
**Contains**:
- Standard color palette
- Common style patterns (ready to copy)
- Layout patterns with examples
- Border/alignment/spacing reference
- Unicode icon library
- Color system guide
**Size**: Practical styling guide (~300 lines)

---

## 🎯 Usage Guide for Agents

### Scenario 1: "I need to create a new TUI view"

1. Read [README-BUBBLETEA.md](./README-BUBBLETEA.md) - Quick Start section
2. Open [bubbletea-component-guide.md](./bubbletea-component-guide.md) - Component Architecture section
3. Keep [QUICK-REFERENCE.md](./QUICK-REFERENCE.md) open while coding
4. Reference [lipgloss-styling-reference.md](./lipgloss-styling-reference.md) for styling

### Scenario 2: "How do I style this button/panel/header?"

1. Check [lipgloss-styling-reference.md](./lipgloss-styling-reference.md) - Common Style Patterns
2. Copy the pattern that matches your need
3. Customize colors from the Standard Color Palette section

### Scenario 3: "I need a scrollable list component"

1. Check [bubbletea-component-guide.md](./bubbletea-component-guide.md) - Bubbles Component Library → List Component
2. Copy the example code
3. Customize styles from [lipgloss-styling-reference.md](./lipgloss-styling-reference.md)

### Scenario 4: "How do I lay out 3 panels side-by-side?"

1. Check [lipgloss-styling-reference.md](./lipgloss-styling-reference.md) - Layout Patterns → 3-Panel Layout
2. Or check [bubbletea-component-guide.md](./bubbletea-component-guide.md) - Layout Patterns → 3-Panel Layout
3. Both have complete examples

### Scenario 5: "Quick syntax check - how do I center text?"

1. Check [QUICK-REFERENCE.md](./QUICK-REFERENCE.md) - Layout Functions section
2. Use: `lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, str)`

---

## 📊 Document Comparison

| Document | Size | Purpose | Best For | Depth |
|----------|------|---------|----------|-------|
| QUICK-REFERENCE | 1 page | Cheat sheet | Active coding | Quick lookup |
| README-BUBBLETEA | 5 min read | Overview | Getting started | Summary |
| bubbletea-component-guide | 500 lines | Technical guide | Building components | Deep dive |
| lipgloss-styling-reference | 300 lines | Styling guide | UI consistency | Practical |

---

## 🔑 Key Concepts Summary

### Bubble Tea (Framework)

**Architecture**: Elm Architecture (Model-Update-View)
- **Model**: Application state (struct)
- **Update**: Handle messages, return new model + command
- **View**: Render model to string (pure function)

**Key Principles**:
1. Models are immutable
2. No side effects in Update (use `tea.Cmd`)
3. Messages drive all changes
4. Components are composable

### Lipgloss (Styling)

**Philosophy**: CSS-like styling for terminals

**Key Principles**:
1. Never use raw ANSI codes
2. Define styles once at package level
3. Use semantic color names
4. Measure before laying out
5. Handle responsive design

### Bubbles (Components)

**Pre-built components** for common UI patterns:
- `list` - Scrollable, filterable lists
- `viewport` - Scrollable text (logs)
- `textinput` - Single-line input
- `textarea` - Multi-line input
- `table` - Tables with sorting
- `spinner` - Loading indicators
- `progress` - Progress bars
- `help` - Auto-generated help

---

## 🎨 Standard Color Palette

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

---

## ✅ Implementation Checklist

Before starting TUI development:

- [ ] Read README-BUBBLETEA.md for overview
- [ ] Bookmark QUICK-REFERENCE.md in editor
- [ ] Review component guide for architecture patterns
- [ ] Copy color palette to `internal/tui/styles/colors.go`
- [ ] Set up component directory structure
- [ ] Handle `tea.WindowSizeMsg` in all components
- [ ] Define styles at package level (not in View)
- [ ] Use Bubbles components for common patterns
- [ ] Test with minimum terminal size (80x24)

---

## 🔗 External Resources

- **Bubble Tea GitHub**: https://github.com/charmbracelet/bubbletea
- **Lipgloss GitHub**: https://github.com/charmbracelet/lipgloss
- **Bubbles GitHub**: https://github.com/charmbracelet/bubbles
- **Examples Directory**: https://github.com/charmbracelet/bubbletea/tree/main/examples
- **Video Tutorials**: https://charm.sh/yt
- **Community Templates**: https://github.com/charm-and-friends/additional-bubbles

---

## 📝 Related Research

This research builds on:

**R01: TUI Excellence Study** ([01-tui-excellence/findings.md](./01-tui-excellence/findings.md))
- Analyzed lazydocker, lazygit, k9s, gh-dash
- Identified 3-panel layout as optimal for Docker managers
- Documented keyboard navigation patterns
- Real-time update strategies
- Help system best practices

Key finding: **Great TUIs use simple patterns consistently**

---

## 🚀 Next Steps

1. **Create `internal/tui/styles/colors.go`** with standard palette
2. **Implement first component** (e.g., service list)
3. **Test in isolation** before integration
4. **Add Docker API integration**
5. **Iterate based on user feedback**

---

## 📞 Questions?

If you're an agent and need clarification:

1. Check [QUICK-REFERENCE.md](./QUICK-REFERENCE.md) first
2. Then check the relevant detailed guide
3. Look for similar examples in the guides
4. Check official Bubble Tea examples

All common patterns are documented with copy-paste examples!

---

**Happy Building! 🎉**
