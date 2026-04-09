package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"textedit/internal/editor"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type panel int

const (
	treePanel panel = iota
	editorPanel
)

type model struct {
	tree      treeModel
	textarea  textarea.Model
	textinput textinput.Model
	active    panel
	width     int
	height    int

	currentFile  string
	currentName  string
	savedContent string
	hasCRLF      bool // true if the file originally used \r\n line endings
	dirty        bool
	statusMsg    string
	quitting     bool
	creating     bool // true when the new-file prompt is visible

	showWhitespace bool   // toggle for visible whitespace characters
	realContent    string // actual content while whitespace view is active
}

func newModel(dir string, openFile string) model {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		absDir = dir
	}

	ta := textarea.New()
	ta.Placeholder = "Open a file to start editing..."
	ta.ShowLineNumbers = true
	ta.CharLimit = 0
	ta.Blur()

	ti := textinput.New()
	ti.Placeholder = "filename.txt"
	ti.CharLimit = 255

	m := model{
		tree:      newTreeModel(absDir),
		textarea:  ta,
		textinput: ti,
		active:    treePanel,
	}

	if openFile != "" {
		absFile, err := filepath.Abs(openFile)
		if err == nil {
			m.openFile(absFile)
		}
	}

	return m
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateSizes()
		return m, nil

	case tea.KeyMsg:
		// Handle new-file prompt input
		if m.creating {
			switch msg.String() {
			case "enter":
				name := m.textinput.Value()
				if name != "" {
					m.createFile(name)
				}
				m.creating = false
				m.textinput.Blur()
				m.textinput.SetValue("")
				return m, nil
			case "esc":
				m.creating = false
				m.textinput.Blur()
				m.textinput.SetValue("")
				m.statusMsg = ""
				return m, nil
			}
			var cmd tea.Cmd
			m.textinput, cmd = m.textinput.Update(msg)
			return m, cmd
		}

		switch msg.String() {
		case "ctrl+q":
			m.quitting = true
			return m, tea.Quit
		case "ctrl+s":
			if m.showWhitespace {
				m.statusMsg = "Cannot save in whitespace view"
				return m, nil
			}
			m.save()
			return m, nil
		case "ctrl+n":
			m.creating = true
			m.textinput.Focus()
			m.textinput.SetValue("")
			m.statusMsg = "New file name:"
			return m, textinput.Blink
		case "ctrl+h":
			if m.currentFile == "" {
				return m, nil
			}
			if m.showWhitespace {
				m.showWhitespace = false
				m.textarea.SetValue(m.realContent)
				m.dirty = m.realContent != m.savedContent
				m.realContent = ""
			} else {
				m.showWhitespace = true
				m.realContent = m.textarea.Value()
				m.textarea.SetValue(annotateWhitespace(m.realContent))
			}
			return m, nil
		case "tab":
			if m.active == treePanel {
				m.active = editorPanel
				m.textarea.Focus()
			} else {
				m.active = treePanel
				m.textarea.Blur()
			}
			return m, nil
		}

		if m.active == treePanel {
			switch msg.String() {
			case "up", "k":
				m.tree.moveUp()
			case "down", "j":
				m.tree.moveDown()
			case "enter":
				item := m.tree.selected()
				if item != nil {
					if item.isDir {
						m.tree.toggle()
					} else {
						m.openFile(item.path)
					}
				}
			case "n":
				m.creating = true
				m.textinput.Focus()
				m.textinput.SetValue("")
				m.statusMsg = "New file name:"
				return m, textinput.Blink
			case "esc", "q":
				m.quitting = true
				return m, tea.Quit
			}
			return m, nil
		}

		// Editor panel
		if msg.String() == "esc" {
			m.active = treePanel
			m.textarea.Blur()
			return m, nil
		}
	}

	if m.active == editorPanel && !m.showWhitespace {
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		m.dirty = m.textarea.Value() != m.savedContent
		return m, cmd
	}

	return m, nil
}

func (m *model) updateSizes() {
	treeWidth := 30
	if m.width > 120 {
		treeWidth = 40
	}
	// Subtract borders (2 per panel) and 1 for gap
	editorWidth := m.width - treeWidth - 5
	if editorWidth < 10 {
		editorWidth = 10
	}
	contentHeight := m.height - 3 // status bar + borders
	if contentHeight < 3 {
		contentHeight = 3
	}

	m.tree.width = treeWidth
	m.tree.height = contentHeight
	m.textarea.SetWidth(editorWidth)
	m.textarea.SetHeight(contentHeight)
}

func (m *model) openFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		m.statusMsg = fmt.Sprintf("Error: %v", err)
		return
	}
	if editor.IsBinary(data) {
		m.statusMsg = "Cannot open binary file"
		return
	}
	raw := string(data)
	m.hasCRLF = strings.Contains(raw, "\r\n")
	content := strings.ReplaceAll(raw, "\r\n", "\n")
	m.currentFile = path
	m.currentName = filepath.Base(path)
	m.savedContent = content
	m.textarea.SetValue(content)
	m.dirty = false
	m.statusMsg = ""
	m.active = editorPanel
	m.textarea.Focus()
}

func (m *model) createFile(name string) {
	// Determine target directory: selected dir, or parent of selected file
	dir := m.tree.rootPath
	item := m.tree.selected()
	if item != nil {
		if item.isDir {
			dir = item.path
		} else {
			dir = filepath.Dir(item.path)
		}
	}

	fullPath := filepath.Join(dir, name)
	f, err := os.Create(fullPath)
	if err != nil {
		m.statusMsg = fmt.Sprintf("Error: %v", err)
		return
	}
	f.Close()

	// Refresh the tree, preserving dimensions
	w, h := m.tree.width, m.tree.height
	m.tree = newTreeModel(m.tree.rootPath)
	m.tree.width = w
	m.tree.height = h
	m.openFile(fullPath)
	m.statusMsg = fmt.Sprintf("Created %s", name)
}

func (m *model) save() {
	if m.currentFile == "" {
		m.statusMsg = "No file open"
		return
	}
	content := m.textarea.Value()
	output := content
	if m.hasCRLF {
		output = strings.ReplaceAll(output, "\n", "\r\n")
	}
	err := os.WriteFile(m.currentFile, []byte(output), 0644)
	if err != nil {
		m.statusMsg = fmt.Sprintf("Error saving: %v", err)
		return
	}
	m.savedContent = content
	m.dirty = false
	m.statusMsg = "Saved"
}

func (m model) View() string {
	if m.quitting {
		return ""
	}
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	// Tree panel
	treeView := m.tree.view(m.active == treePanel)

	// Editor panel
	editorWidth := m.width - m.tree.width - 5
	if editorWidth < 10 {
		editorWidth = 10
	}
	contentHeight := m.tree.height

	var editorBorderStyle lipgloss.Style
	if m.active == editorPanel {
		editorBorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#007acc")).
			Width(editorWidth).
			Height(contentHeight)
	} else {
		editorBorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#3c3c3c")).
			Width(editorWidth).
			Height(contentHeight)
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#007acc"))

	title := "Untitled"
	if m.currentName != "" {
		title = m.currentName
	}
	if m.dirty {
		title = "● " + title
	}

	editorHeader := titleStyle.Render(" " + title)
	editorContent := editorHeader + "\n" + m.textarea.View()
	editorView := editorBorderStyle.Render(editorContent)

	// Join panels horizontally
	panels := lipgloss.JoinHorizontal(lipgloss.Top, treeView, editorView)

	// Status bar
	statusStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#007acc")).
		Foreground(lipgloss.Color("#ffffff")).
		Width(m.width)

	var statusLeft string
	if m.creating {
		statusLeft = " New file: " + m.textinput.View()
	} else if m.statusMsg != "" {
		statusLeft = " " + m.statusMsg
	} else if m.currentFile != "" {
		statusLeft = " " + m.currentFile
		if m.dirty {
			statusLeft += " (modified)"
		}
	} else {
		statusLeft = " TextEdit TUI"
	}

	var indicators string
	if m.currentFile != "" {
		if m.hasCRLF {
			indicators = "CRLF "
		} else {
			indicators = "LF "
		}
	}
	if m.showWhitespace {
		indicators += "[WS] "
	}
	statusRight := indicators + "Ctrl+H: whitespace  Ctrl+S: save  Ctrl+Q: quit "
	padding := m.width - len(statusLeft) - len(statusRight)
	if padding < 1 {
		padding = 1
	}
	statusBar := statusStyle.Render(statusLeft + strings.Repeat(" ", padding) + statusRight)

	return panels + "\n" + statusBar
}

func annotateWhitespace(s string) string {
	var b strings.Builder
	for _, line := range strings.Split(s, "\n") {
		for _, ch := range line {
			switch ch {
			case ' ':
				b.WriteRune('·')
			case '\t':
				b.WriteString("→   ")
			default:
				b.WriteRune(ch)
			}
		}
		b.WriteString("↵\n")
	}
	// Remove the trailing ↵\n added after the last split element
	result := b.String()
	if len(result) >= 3 {
		result = result[:len(result)-len("↵\n")]
	}
	return result
}

func main() {
	dir := "."
	openFile := ""

	if len(os.Args) > 1 {
		arg := os.Args[1]
		info, err := os.Stat(arg)
		if err == nil && !info.IsDir() {
			openFile = arg
			dir = filepath.Dir(arg)
		} else {
			dir = arg
		}
	}

	p := tea.NewProgram(
		newModel(dir, openFile),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
