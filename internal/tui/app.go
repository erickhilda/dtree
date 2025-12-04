package tui

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"dtree/internal/color"
	"dtree/internal/tree"
)

var (
	titleStyle        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	itemStyle         = lipgloss.NewStyle().PaddingLeft(2)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	helpStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginTop(1)
)

type model struct {
	treeRoot    *tree.Node
	currentPath string
	expanded    map[string]bool
	selected    string
	viewport    viewport.Model
	theme       *color.Theme
	searchMode  bool
	searchQuery string
	filtered    []*tree.Node
	height      int
	width       int
}

type item struct {
	node *tree.Node
}

func (i item) FilterValue() string {
	return i.node.Name
}

func (i item) Title() string {
	return i.node.Name
}

func (i item) Description() string {
	if i.node.IsDir {
		return fmt.Sprintf("Directory (%d items)", len(i.node.Children))
	}
	return fmt.Sprintf("File (%s)", formatSize(i.node.Size))
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func initialModel(root *tree.Node, theme *color.Theme) model {
	expanded := make(map[string]bool)
	expanded[root.Path] = true

	return model{
		treeRoot:    root,
		currentPath: root.Path,
		expanded:    expanded,
		selected:    root.Path,
		theme:       theme,
		searchMode:  false,
		searchQuery: "",
		filtered:    []*tree.Node{root},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 4
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			if m.searchMode {
				m.searchMode = false
				m.searchQuery = ""
				return m, nil
			}
			return m, tea.Quit

		case "/":
			m.searchMode = true
			m.searchQuery = ""
			return m, nil

		case "enter":
			if m.searchMode {
				m.searchMode = false
				m.filterTree()
				return m, nil
			}
			m.toggleExpand()
			return m, nil

		case "j", "down":
			m.moveSelection(1)
			return m, nil

		case "k", "up":
			m.moveSelection(-1)
			return m, nil

		case "gg":
			m.moveToTop()
			return m, nil

		case "G":
			m.moveToBottom()
			return m, nil

		default:
			if m.searchMode {
				if msg.String() == "backspace" {
					if len(m.searchQuery) > 0 {
						m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
					}
				} else if len(msg.String()) == 1 {
					m.searchQuery += msg.String()
				}
				m.filterTree()
				return m, nil
			}
		}
	}

	return m, nil
}

func (m *model) toggleExpand() {
	// Find selected node and toggle expansion
	for _, node := range m.getVisibleNodes() {
		if node.Path == m.selected {
			if node.IsDir {
				key := filepath.Join(node.Path, node.Name)
				m.expanded[key] = !m.expanded[key]
			}
			break
		}
	}
}

func (m *model) moveSelection(delta int) {
	nodes := m.getVisibleNodes()
	if len(nodes) == 0 {
		return
	}

	currentIdx := -1
	for i, node := range nodes {
		if node.Path == m.selected {
			currentIdx = i
			break
		}
	}

	if currentIdx == -1 {
		m.selected = nodes[0].Path
		return
	}

	newIdx := currentIdx + delta
	if newIdx < 0 {
		newIdx = 0
	} else if newIdx >= len(nodes) {
		newIdx = len(nodes) - 1
	}

	m.selected = nodes[newIdx].Path
}

func (m *model) moveToTop() {
	nodes := m.getVisibleNodes()
	if len(nodes) > 0 {
		m.selected = nodes[0].Path
	}
}

func (m *model) moveToBottom() {
	nodes := m.getVisibleNodes()
	if len(nodes) > 0 {
		m.selected = nodes[len(nodes)-1].Path
	}
}

func (m *model) getVisibleNodes() []*tree.Node {
	if len(m.filtered) > 0 {
		return m.filtered
	}
	return m.flattenTree(m.treeRoot, "")
}

func (m *model) flattenTree(node *tree.Node, prefix string) []*tree.Node {
	var result []*tree.Node
	result = append(result, node)

	if node.IsDir {
		key := filepath.Join(node.Path, node.Name)
		if m.expanded[key] {
			for _, child := range node.Children {
				result = append(result, m.flattenTree(child, prefix+"  ")...)
			}
		}
	}

	return result
}

func (m *model) filterTree() {
	if m.searchQuery == "" {
		m.filtered = nil
		return
	}

	query := strings.ToLower(m.searchQuery)
	m.filtered = []*tree.Node{}

	var search func(*tree.Node)
	search = func(node *tree.Node) {
		if strings.Contains(strings.ToLower(node.Name), query) {
			m.filtered = append(m.filtered, node)
		}
		for _, child := range node.Children {
			search(child)
		}
	}

	search(m.treeRoot)
}

func (m model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render("Directory Tree: " + m.treeRoot.Name))
	b.WriteString("\n\n")

	// Search bar
	if m.searchMode {
		b.WriteString(fmt.Sprintf("Search: %s_\n\n", m.searchQuery))
	}

	// Tree view
	nodes := m.getVisibleNodes()
	content := m.renderTree(nodes)
	m.viewport.SetContent(content)
	b.WriteString(m.viewport.View())

	// Help
	help := helpStyle.Render("↑/↓: Navigate  Enter: Expand/Collapse  /: Search  q: Quit")
	b.WriteString("\n" + help)

	return b.String()
}

func (m *model) renderTree(nodes []*tree.Node) string {
	var lines []string
	for _, node := range nodes {
		isSelected := node.Path == m.selected
		line := m.renderNode(node, "", true, true, isSelected)
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func (m *model) renderNode(node *tree.Node, prefix string, isLast bool, skipRoot bool, isSelected bool) string {
	var line string

	if !skipRoot {
		connector := tree.TreeLast
		if !isLast {
			connector = tree.TreeBranch
		}

		name := node.Name
		if m.theme != nil && m.theme.IsEnabled() {
			name = m.theme.Colorize(node.Name, node.IsDir, node.IsSymlink, node.Mode)
		}

		if isSelected {
			line = selectedItemStyle.Render(prefix + connector + name)
		} else {
			line = itemStyle.Render(prefix + connector + name)
		}
	}

	if node.IsDir {
		key := filepath.Join(node.Path, node.Name)
		if m.expanded[key] {
			children := node.Children
			for i, child := range children {
				isLastChild := i == len(children)-1
				childPrefix := prefix
				if !skipRoot {
					if isLast {
						childPrefix += tree.TreeSpace
					} else {
						childPrefix += tree.TreePipe
					}
				}
				childLine := m.renderNode(child, childPrefix, isLastChild, false, child.Path == m.selected)
				if line != "" {
					line += "\n" + childLine
				} else {
					line = childLine
				}
			}
		}
	}

	return line
}

// Run starts the TUI application
func Run(root *tree.Node, theme *color.Theme) error {
	m := initialModel(root, theme)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

