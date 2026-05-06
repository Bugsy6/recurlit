package panes

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/bugsy6/recurlit/internal/styles"
)

type URLPane struct {
	input  textinput.Model
	active bool
	insert bool
	dirty  bool
	width  int
	height int
}

func NewURLPane() *URLPane {
	ti := textinput.New()
	ti.Placeholder = "https://api.example.com/endpoint"
	ti.CharLimit = 2048
	return &URLPane{input: ti, dirty: true}
}

func (u *URLPane) URL() string {
	return u.input.Value()
}

func (u *URLPane) SetURL(url string) {
	u.input.SetValue(url)
	u.dirty = true
}

func (u *URLPane) Update(msg tea.Msg) (Pane, tea.Cmd) {
	if !u.insert {
		return u, nil
	}
	var cmd tea.Cmd
	prev := u.input.Value()
	u.input, cmd = u.input.Update(msg)
	if u.input.Value() != prev {
		u.dirty = true
	}
	return u, cmd
}

func (u *URLPane) View() string {
	style := styles.InactiveBorder
	if u.insert {
		style = styles.InsertBorder
	} else if u.active {
		style = styles.ActiveBorder
	}
	style = style.Width(u.width - 2).Height(u.height - 2)
	title := styles.Dim.Render("URL")
	return style.Render(title + "\n" + u.input.View())
}

func (u *URLPane) SetSize(width, height int) {
	u.width = width
	u.height = height
	u.input.Width = width - 4
}

func (u *URLPane) SetActive(active bool) {
	u.active = active
	if active && u.insert {
		u.input.Focus()
	} else if !active {
		u.input.Blur()
	}
}

func (u *URLPane) SetInsert(insert bool) {
	u.insert = insert
	if insert {
		u.input.Focus()
	} else {
		u.input.Blur()
	}
}

func (u *URLPane) IsDirty() bool { return u.dirty }
func (u *URLPane) ClearDirty()   { u.dirty = false }
