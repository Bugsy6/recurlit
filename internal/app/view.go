package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/bugsy6/recurlit/internal/styles"
)

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	var sb strings.Builder

	// Row 1: Method + URL
	row1 := lipgloss.JoinHorizontal(lipgloss.Top,
		m.method.View(),
		m.url.View(),
	)
	sb.WriteString(row1 + "\n")

	// Row 2: Headers + Params
	row2 := lipgloss.JoinHorizontal(lipgloss.Top,
		m.headers.View(),
		m.params.View(),
	)
	sb.WriteString(row2 + "\n")

	// Row 3: Body + Auth
	row3 := lipgloss.JoinHorizontal(lipgloss.Top,
		m.body.View(),
		m.auth.View(),
	)
	sb.WriteString(row3 + "\n")

	// Curl bar
	curlContent := m.curlStr
	if curlContent == "" {
		curlContent = "curl preview will appear here"
	}
	curlBar := styles.CurlBar.Width(m.width).Render(truncateLine(curlContent, m.width-2))
	sb.WriteString(curlBar + "\n")

	// Results
	sb.WriteString(m.results.View() + "\n")

	// Status bar
	sb.WriteString(m.renderStatusBar())

	out := sb.String()

	// Overlays
	switch m.mode {
	case ModeSave:
		out = renderOverlay(out, m.renderSaveModal(), m.width, m.height)
	case ModeLoad:
		out = renderOverlay(out, m.renderLoadModal(), m.width, m.height)
	}

	if m.showHelp {
		out = renderOverlay(out, m.renderHelp(), m.width, m.height)
	}

	return out
}

func (m Model) renderStatusBar() string {
	mode := "NORMAL"
	modeColor := lipgloss.Color("#10B981")
	if m.mode == ModeInsert {
		mode = "INSERT"
		modeColor = lipgloss.Color("#7C3AED")
	}
	modeStr := lipgloss.NewStyle().Foreground(modeColor).Bold(true).Render(mode)

	paneStr := styles.Dim.Render(" [" + m.activePane.String() + "]")

	var status string
	if m.executing {
		status = m.spinner.View() + " Running..."
	} else if m.lastStatus != "" {
		status = m.lastStatus
	}

	help := styles.Dim.Render("tab:pane  i:insert  r:run  s:save  o:open  ?:help  q:quit")

	left := modeStr + paneStr
	if status != "" {
		left += "  " + status
	}

	// Fill remaining space
	leftLen := lipgloss.Width(left)
	rightLen := lipgloss.Width(help)
	spaces := m.width - leftLen - rightLen - 2
	if spaces < 1 {
		spaces = 1
	}

	bar := left + strings.Repeat(" ", spaces) + help
	return styles.StatusBar.Width(m.width).Render(bar)
}

func (m Model) renderSaveModal() string {
	content := fmt.Sprintf(
		"Save Request\n\n%s\n\n%s",
		m.saveInput.View(),
		styles.Dim.Render("enter: save  esc: cancel"),
	)
	return renderModal(content, 40, 6)
}

func (m Model) renderLoadModal() string {
	var sb strings.Builder
	sb.WriteString("Open Request\n\n")
	if len(m.savedRequests) == 0 {
		sb.WriteString(styles.Dim.Render("No saved requests"))
	} else {
		for i, sr := range m.savedRequests {
			if i == m.loadCursor {
				sb.WriteString(styles.Selected.Render("> " + sr.Name))
			} else {
				sb.WriteString(styles.Dim.Render("  " + sr.Name))
			}
			sb.WriteString("\n")
		}
	}
	sb.WriteString("\n" + styles.Dim.Render("j/k: move  enter: load  esc: cancel"))
	return renderModal(sb.String(), 50, 20)
}

func (m Model) renderHelp() string {
	lines := []string{
		"recURLit — Help",
		"",
		"Normal Mode:",
		"  tab / l         next pane",
		"  shift+tab / h   prev pane",
		"  j / k           down/up (or scroll results)",
		"  i / enter       insert mode",
		"  r               execute request",
		"  s               save request",
		"  o               open saved request",
		"  ?               toggle help",
		"  q / ctrl+c      quit",
		"",
		"Insert Mode:",
		"  esc             back to normal",
		"  (keys routed to active pane)",
		"",
		"KV Panes (Headers/Params):",
		"  tab             next field",
		"  enter           next field / new row",
		"  ↑/↓             move between rows",
		"  ctrl+d          delete row",
		"",
		"Auth Pane:",
		"  ctrl+t          cycle auth type",
		"  tab             next field",
		"",
		"Results Pane:",
		"  j/k             scroll",
		"  ctrl+d/ctrl+u   half page",
		"  g/G             top/bottom",
		"",
		styles.Dim.Render("Press ? to close"),
	}
	return renderModal(strings.Join(lines, "\n"), 50, 32)
}

func renderModal(content string, width, height int) string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7C3AED")).
		Padding(1, 2).
		Width(width).
		Height(height).
		Background(lipgloss.Color("#111827"))
	return style.Render(content)
}

func renderOverlay(base, overlay string, totalW, totalH int) string {
	overlayW := lipgloss.Width(overlay)
	overlayH := lipgloss.Height(overlay)

	col := (totalW - overlayW) / 2
	row := (totalH - overlayH) / 2
	if col < 0 {
		col = 0
	}
	if row < 0 {
		row = 0
	}

	baseLines := strings.Split(base, "\n")
	overlayLines := strings.Split(overlay, "\n")

	for i, ol := range overlayLines {
		targetRow := row + i
		if targetRow >= len(baseLines) {
			break
		}
		bl := baseLines[targetRow]
		baseLines[targetRow] = overlayLine(bl, ol, col, totalW)
	}
	return strings.Join(baseLines, "\n")
}

func overlayLine(base, overlay string, col, totalW int) string {
	// Strip ANSI for width calculation, then place overlay at col
	// Simple approach: pad base to totalW, splice in overlay
	baseRunes := []rune(stripSimple(base))
	for len(baseRunes) < totalW {
		baseRunes = append(baseRunes, ' ')
	}
	overlayClean := stripSimple(overlay)
	for i, ch := range []rune(overlayClean) {
		pos := col + i
		if pos >= 0 && pos < len(baseRunes) {
			baseRunes[pos] = ch
		}
	}
	// Re-build: prepend prefix, overlay raw (with ANSI), append suffix
	prefix := string(baseRunes[:col])
	suffix := ""
	end := col + lipgloss.Width(overlay)
	if end < len(baseRunes) {
		suffix = string(baseRunes[end:])
	}
	return prefix + overlay + suffix
}

func stripSimple(s string) string {
	// Very basic ANSI stripper - just for width calculation
	var out strings.Builder
	inEscape := false
	for _, r := range s {
		if r == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
				inEscape = false
			}
			continue
		}
		out.WriteRune(r)
	}
	return out.String()
}

func truncateLine(s string, n int) string {
	if n <= 0 {
		return ""
	}
	if lipgloss.Width(s) <= n {
		return s
	}
	runes := []rune(s)
	if len(runes) > n-3 {
		return string(runes[:n-3]) + "..."
	}
	return s
}
