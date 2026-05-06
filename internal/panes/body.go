package panes

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/bugsy6/recurlit/internal/styles"
)

type BodyPane struct {
	ta     textarea.Model
	active bool
	insert bool
	dirty  bool
	width  int
	height int
}

func NewBodyPane() *BodyPane {
	ta := textarea.New()
	ta.Placeholder = "Request body (JSON, XML, form data...)"
	ta.SetWidth(60)
	ta.SetHeight(6)
	return &BodyPane{ta: ta, dirty: true}
}

func (b *BodyPane) Body() string {
	return b.ta.Value()
}

func (b *BodyPane) SetBody(body string) {
	b.ta.SetValue(body)
	b.dirty = true
}

func (b *BodyPane) Update(msg tea.Msg) (Pane, tea.Cmd) {
	if !b.insert {
		return b, nil
	}
	var cmd tea.Cmd
	prev := b.ta.Value()
	b.ta, cmd = b.ta.Update(msg)
	if b.ta.Value() != prev {
		b.dirty = true
	}
	return b, cmd
}

func (b *BodyPane) View() string {
	style := styles.InactiveBorder
	if b.insert {
		style = styles.InsertBorder
	} else if b.active {
		style = styles.ActiveBorder
	}
	style = style.Width(b.width - 2).Height(b.height - 2)
	title := styles.Title.Render("Body")
	return style.Render(title + "\n" + b.ta.View())
}

func (b *BodyPane) SetSize(width, height int) {
	b.width = width
	b.height = height
	innerW := width - 4
	innerH := height - 5 // -1 for title line
	if innerW < 1 {
		innerW = 1
	}
	if innerH < 1 {
		innerH = 1
	}
	b.ta.SetWidth(innerW)
	b.ta.SetHeight(innerH)
}

func (b *BodyPane) SetActive(active bool) {
	b.active = active
	if !active {
		b.ta.Blur()
	}
}

func (b *BodyPane) SetInsert(insert bool) {
	b.insert = insert
	if insert {
		b.ta.Focus()
	} else {
		b.ta.Blur()
	}
}

func (b *BodyPane) IsDirty() bool { return b.dirty }
func (b *BodyPane) ClearDirty()   { b.dirty = false }
