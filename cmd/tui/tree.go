package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"textedit/internal/editor"

	"github.com/charmbracelet/lipgloss"
)

type treeItem struct {
	name     string
	path     string
	isDir    bool
	depth    int
	expanded bool
}

type treeModel struct {
	items    []treeItem
	cursor   int
	offset   int
	rootPath string
	width    int
	height   int
}

func newTreeModel(rootPath string) treeModel {
	t := treeModel{
		rootPath: rootPath,
	}
	t.items = t.loadChildren(rootPath, 0)
	return t
}

func (t *treeModel) loadChildren(dirPath string, depth int) []treeItem {
	entries, err := editor.ListDirectory(dirPath)
	if err != nil {
		return nil
	}

	var items []treeItem
	for _, e := range entries {
		items = append(items, treeItem{
			name:  e.Name,
			path:  e.Path,
			isDir: e.IsDir,
			depth: depth,
		})
	}
	return items
}

func (t *treeModel) selected() *treeItem {
	if t.cursor < 0 || t.cursor >= len(t.items) {
		return nil
	}
	return &t.items[t.cursor]
}

func (t *treeModel) moveUp() {
	if t.cursor > 0 {
		t.cursor--
	}
	t.fixScroll()
}

func (t *treeModel) moveDown() {
	if t.cursor < len(t.items)-1 {
		t.cursor++
	}
	t.fixScroll()
}

func (t *treeModel) fixScroll() {
	if t.cursor < t.offset {
		t.offset = t.cursor
	}
	visible := t.height - 2 // account for border
	if visible < 1 {
		visible = 1
	}
	if t.cursor >= t.offset+visible {
		t.offset = t.cursor - visible + 1
	}
}

func (t *treeModel) toggle() {
	item := t.selected()
	if item == nil || !item.isDir {
		return
	}

	if item.expanded {
		// Collapse: remove all children below this item that have greater depth
		t.items[t.cursor].expanded = false
		start := t.cursor + 1
		end := start
		for end < len(t.items) && t.items[end].depth > item.depth {
			end++
		}
		t.items = append(t.items[:start], t.items[end:]...)
	} else {
		// Expand: load children and insert after current item
		t.items[t.cursor].expanded = true
		children := t.loadChildren(item.path, item.depth+1)
		if len(children) > 0 {
			after := append(children, t.items[t.cursor+1:]...)
			t.items = append(t.items[:t.cursor+1], after...)
		}
	}
}

func (t *treeModel) view(active bool) string {
	var borderStyle lipgloss.Style
	if active {
		borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#007acc")).
			Width(t.width).
			Height(t.height)
	} else {
		borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#3c3c3c")).
			Width(t.width).
			Height(t.height)
	}

	cursorStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#37373d")).
		Foreground(lipgloss.Color("#ffffff"))

	dirStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#dcb67a"))

	fileStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#cccccc"))

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#007acc"))

	visible := t.height
	if visible < 1 {
		visible = 1
	}

	var lines []string
	header := headerStyle.Render(fmt.Sprintf(" Files — %s", filepath.Base(t.rootPath)))
	lines = append(lines, header)

	end := t.offset + visible - 1 // -1 for header
	if end > len(t.items) {
		end = len(t.items)
	}

	for i := t.offset; i < end; i++ {
		item := t.items[i]
		indent := strings.Repeat("  ", item.depth)

		var icon string
		if item.isDir {
			if item.expanded {
				icon = "▾ "
			} else {
				icon = "▸ "
			}
		} else {
			icon = "  "
		}

		line := fmt.Sprintf(" %s%s%s", indent, icon, item.name)

		// Pad to width
		if len(line) < t.width {
			line += strings.Repeat(" ", t.width-len(line))
		} else if len(line) > t.width {
			line = line[:t.width]
		}

		if i == t.cursor {
			line = cursorStyle.Render(line)
		} else if item.isDir {
			line = dirStyle.Render(line)
		} else {
			line = fileStyle.Render(line)
		}

		lines = append(lines, line)
	}

	// Fill remaining lines
	for len(lines) < visible+1 {
		lines = append(lines, strings.Repeat(" ", t.width))
	}

	content := strings.Join(lines, "\n")
	return borderStyle.Render(content)
}
