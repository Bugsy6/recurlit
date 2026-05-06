package panes

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/bugsy6/recurlit/internal/styles"
)

// kvRow is one editable key-value row
type kvRow struct {
	key   textinput.Model
	value textinput.Model
}

func newKVRow() kvRow {
	k := textinput.New()
	k.Placeholder = "key"
	k.CharLimit = 256
	k.Width = 20

	v := textinput.New()
	v.Placeholder = "value"
	v.CharLimit = 1024
	v.Width = 30

	return kvRow{key: k, value: v}
}

// KVPane manages a list of key-value pairs (headers or params)
type KVPane struct {
	rows    []kvRow
	cursor  int // row index
	col     int // 0=key, 1=value
	active  bool
	insert  bool
	dirty   bool
	width   int
	height  int
	title   string
}

func NewKVPane(title string) *KVPane {
	p := &KVPane{title: title, dirty: true}
	p.rows = []kvRow{newKVRow()}
	return p
}

type KVPair struct {
	Key     string
	Value   string
	Enabled bool
}

func (kv *KVPane) Pairs() []KVPair {
	var pairs []KVPair
	for _, row := range kv.rows {
		k := row.key.Value()
		if k == "" {
			continue
		}
		pairs = append(pairs, KVPair{Key: k, Value: row.value.Value(), Enabled: true})
	}
	return pairs
}

func (kv *KVPane) SetPairs(pairs []KVPair) {
	kv.rows = nil
	for _, p := range pairs {
		row := newKVRow()
		row.key.SetValue(p.Key)
		row.value.SetValue(p.Value)
		kv.rows = append(kv.rows, row)
	}
	if len(kv.rows) == 0 {
		kv.rows = []kvRow{newKVRow()}
	}
	kv.dirty = true
}

func (kv *KVPane) focusCurrent() {
	for i := range kv.rows {
		kv.rows[i].key.Blur()
		kv.rows[i].value.Blur()
	}
	if !kv.insert || kv.cursor >= len(kv.rows) {
		return
	}
	if kv.col == 0 {
		kv.rows[kv.cursor].key.Focus()
	} else {
		kv.rows[kv.cursor].value.Focus()
	}
}

func (kv *KVPane) Update(msg tea.Msg) (Pane, tea.Cmd) {
	if !kv.insert {
		return kv, nil
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			// advance col → next row → wrap
			kv.col++
			if kv.col > 1 {
				kv.col = 0
				kv.cursor++
				if kv.cursor >= len(kv.rows) {
					kv.rows = append(kv.rows, newKVRow())
				}
			}
			kv.focusCurrent()
			return kv, nil
		case "shift+tab":
			kv.col--
			if kv.col < 0 {
				kv.col = 1
				if kv.cursor > 0 {
					kv.cursor--
				}
			}
			kv.focusCurrent()
			return kv, nil
		case "enter":
			if kv.col == 0 {
				// key → value
				kv.col = 1
			} else {
				// value → new/next row
				kv.cursor++
				kv.col = 0
				if kv.cursor >= len(kv.rows) {
					kv.rows = append(kv.rows, newKVRow())
				}
			}
			kv.focusCurrent()
			kv.dirty = true
			return kv, nil
		case "ctrl+d":
			if len(kv.rows) > 1 {
				kv.rows = append(kv.rows[:kv.cursor], kv.rows[kv.cursor+1:]...)
				if kv.cursor >= len(kv.rows) {
					kv.cursor = len(kv.rows) - 1
				}
				kv.focusCurrent()
				kv.dirty = true
			}
			return kv, nil
		case "down":
			if kv.cursor < len(kv.rows)-1 {
				kv.cursor++
				kv.focusCurrent()
			}
			return kv, nil
		case "up":
			if kv.cursor > 0 {
				kv.cursor--
				kv.focusCurrent()
			}
			return kv, nil
		}
	}
	// Delegate to focused input
	if kv.cursor < len(kv.rows) {
		var cmd tea.Cmd
		if kv.col == 0 {
			prev := kv.rows[kv.cursor].key.Value()
			kv.rows[kv.cursor].key, cmd = kv.rows[kv.cursor].key.Update(msg)
			if kv.rows[kv.cursor].key.Value() != prev {
				kv.dirty = true
			}
		} else {
			prev := kv.rows[kv.cursor].value.Value()
			kv.rows[kv.cursor].value, cmd = kv.rows[kv.cursor].value.Update(msg)
			if kv.rows[kv.cursor].value.Value() != prev {
				kv.dirty = true
			}
		}
		return kv, cmd
	}
	return kv, nil
}

func (kv *KVPane) View() string {
	style := styles.InactiveBorder
	if kv.insert {
		style = styles.InsertBorder
	} else if kv.active {
		style = styles.ActiveBorder
	}
	style = style.Width(kv.width - 2).Height(kv.height - 2)

	inner := kv.renderInner()
	return style.Render(inner)
}

func (kv *KVPane) renderInner() string {
	keyW := (kv.width - 8) / 3
	valW := (kv.width - 8) * 2 / 3

	var lines []string
	title := styles.Title.Render(kv.title)
	lines = append(lines, title)
	header := lipgloss.NewStyle().Foreground(styles.ColorDim).Render(
		fmt.Sprintf("%-*s  %-*s", keyW, "KEY", valW, "VALUE"),
	)
	lines = append(lines, header)

	for i, row := range kv.rows {
		k := row.key.Value()
		v := row.value.Value()
		if k == "" {
			k = row.key.Placeholder
		}
		if v == "" {
			v = row.value.Placeholder
		}

		keyStr := truncate(k, keyW)
		valStr := truncate(v, valW)

		if kv.insert && i == kv.cursor {
			// show live textinput views
			kRow := lipgloss.NewStyle().Width(keyW).Render(row.key.View())
			vRow := lipgloss.NewStyle().Width(valW).Render(row.value.View())
			line := kRow + "  " + vRow
			lines = append(lines, line)
		} else if i == kv.cursor && kv.active {
			line := styles.Selected.Render(fmt.Sprintf("%-*s  %-*s", keyW, keyStr, valW, valStr))
			lines = append(lines, line)
		} else {
			line := fmt.Sprintf("%-*s  %-*s", keyW, keyStr, valW, valStr)
			lines = append(lines, line)
		}
	}

	if kv.insert {
		hint := styles.Dim.Render("enter: next/new row  ↑↓: move  ctrl+d: delete")
		lines = append(lines, hint)
	}

	return strings.Join(lines, "\n")
}

func (kv *KVPane) SetSize(width, height int) {
	kv.width = width
	kv.height = height
	keyW := (width - 8) / 3
	valW := (width - 8) * 2 / 3
	for i := range kv.rows {
		kv.rows[i].key.Width = keyW
		kv.rows[i].value.Width = valW
	}
}

func (kv *KVPane) SetActive(active bool) {
	kv.active = active
	if !active {
		kv.focusCurrent()
	}
}

func (kv *KVPane) SetInsert(insert bool) {
	kv.insert = insert
	kv.focusCurrent()
}

func (kv *KVPane) IsDirty() bool { return kv.dirty }
func (kv *KVPane) ClearDirty()   { kv.dirty = false }

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	if n <= 3 {
		return s[:n]
	}
	return s[:n-3] + "..."
}
