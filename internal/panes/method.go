package panes

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/bugsy6/recurlit/internal/models"
	"github.com/bugsy6/recurlit/internal/styles"
)

type MethodPane struct {
	selected int
	active   bool
	insert   bool
	dirty    bool
	width    int
	height   int
}

func NewMethodPane() *MethodPane {
	return &MethodPane{selected: 0, dirty: true}
}

func (m *MethodPane) Method() models.HTTPMethod {
	return models.AllMethods[m.selected]
}

func (m *MethodPane) SetMethod(method models.HTTPMethod) {
	for i, me := range models.AllMethods {
		if me == method {
			m.selected = i
			m.dirty = true
			return
		}
	}
}

func (m *MethodPane) Update(msg tea.Msg) (Pane, tea.Cmd) {
	if !m.insert {
		return m, nil
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if m.selected < len(models.AllMethods)-1 {
				m.selected++
				m.dirty = true
			}
		case "k", "up":
			if m.selected > 0 {
				m.selected--
				m.dirty = true
			}
		}
	}
	return m, nil
}

func (m *MethodPane) View() string {
	inner := m.renderInner()
	style := styles.InactiveBorder
	if m.insert {
		style = styles.InsertBorder
	} else if m.active {
		style = styles.ActiveBorder
	}

	style = style.Width(m.width - 2).Height(m.height - 2)
	return style.Render(inner)
}

func (m *MethodPane) renderInner() string {
	method := models.AllMethods[m.selected]
	label := string(method)
	var content string
	if m.insert {
		content = methodColor(method).Bold(true).Render("↑ " + label + " ↓")
	} else {
		content = methodColor(method).Bold(true).Render(label)
	}
	title := styles.Dim.Render("METHOD")
	return title + "\n" + content
}

func (m *MethodPane) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *MethodPane) SetActive(active bool) {
	m.active = active
}

func (m *MethodPane) SetInsert(insert bool) {
	m.insert = insert
}

func (m *MethodPane) IsDirty() bool { return m.dirty }
func (m *MethodPane) ClearDirty()   { m.dirty = false }

func methodColor(method models.HTTPMethod) lipgloss.Style {
	switch method {
	case models.GET:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#10B981"))
	case models.POST:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#3B82F6"))
	case models.PUT:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#F59E0B"))
	case models.PATCH:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#8B5CF6"))
	case models.DELETE:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#EF4444"))
	case models.HEAD, models.OPTIONS:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280"))
	}
	return lipgloss.NewStyle()
}

// MethodBadge returns a colored method string for use in status bar etc.
func MethodBadge(method models.HTTPMethod) string {
	return methodColor(method).Bold(true).Render(strings.ToUpper(string(method)))
}
