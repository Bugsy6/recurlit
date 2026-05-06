package app

import (
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/bugsy6/recurlit/internal/curl"
	"github.com/bugsy6/recurlit/internal/executor"
	"github.com/bugsy6/recurlit/internal/models"
	"github.com/bugsy6/recurlit/internal/panes"
	"github.com/bugsy6/recurlit/internal/storage"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.recalcSizes()
		return m, nil

	case executor.ResponseMsg:
		m.executing = false
		m.results.SetResponse(msg.Response)
		if msg.Response.Error != "" {
			m.lastStatus = "Error: " + msg.Response.Error
		} else {
			m.lastStatus = msg.Response.Status + " " +
				msg.Response.Duration.Round(time.Millisecond).String()
		}
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	// Modal modes take priority
	switch m.mode {
	case ModeSave:
		return m.updateSaveMode(msg)
	case ModeLoad:
		return m.updateLoadMode(msg)
	}

	// Key events
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		// Esc always exits insert
		if m.mode == ModeInsert && key.Matches(keyMsg, m.keymap.Normal) {
			m.mode = ModeNormal
			m.setInsert(false)
			return m, nil
		}

		if m.mode == ModeInsert {
			return m.updateInsertMode(msg)
		}

		// Normal mode
		return m.updateNormalMode(keyMsg)
	}

	// Pass non-key msgs to active pane (e.g. blink)
	if m.mode == ModeInsert {
		return m.updateInsertMode(msg)
	}

	return m, nil
}

func (m Model) updateNormalMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	km := m.keymap

	switch {
	case key.Matches(msg, km.Quit):
		return m, tea.Quit

	case key.Matches(msg, km.Help):
		m.showHelp = !m.showHelp
		return m, nil

	case key.Matches(msg, km.NextPane):
		next := (int(m.activePane) + 1) % int(panes.PaneCount)
		m.setActive(panes.PaneID(next))
		return m, nil

	case key.Matches(msg, km.PrevPane):
		prev := (int(m.activePane) - 1 + int(panes.PaneCount)) % int(panes.PaneCount)
		m.setActive(panes.PaneID(prev))
		return m, nil

	case key.Matches(msg, km.Down):
		if m.activePane == panes.PaneResults {
			var cmd tea.Cmd
			p, c := m.results.Update(msg)
			m.results = p.(*panes.ResultsPane)
			cmd = c
			return m, cmd
		}
		next := (int(m.activePane) + 1) % int(panes.PaneCount)
		m.setActive(panes.PaneID(next))
		return m, nil

	case key.Matches(msg, km.Up):
		if m.activePane == panes.PaneResults {
			var cmd tea.Cmd
			p, c := m.results.Update(msg)
			m.results = p.(*panes.ResultsPane)
			cmd = c
			return m, cmd
		}
		prev := (int(m.activePane) - 1 + int(panes.PaneCount)) % int(panes.PaneCount)
		m.setActive(panes.PaneID(prev))
		return m, nil

	case key.Matches(msg, km.Insert):
		m.mode = ModeInsert
		m.setInsert(true)
		return m, nil

	case key.Matches(msg, km.Run):
		if m.executing {
			return m, nil
		}
		m.executing = true
		m.lastStatus = "Running..."
		m.results.SetLoading(true)
		req := m.currentRequest()
		return m, tea.Batch(executor.Execute(req), m.spinner.Tick)

	case key.Matches(msg, km.Save):
		m.mode = ModeSave
		m.saveInput.SetValue("")
		m.saveInput.Focus()
		return m, nil

	case key.Matches(msg, km.Open):
		saved, _ := storage.LoadAll()
		m.savedRequests = saved
		m.loadCursor = 0
		m.mode = ModeLoad
		return m, nil
	}

	return m, nil
}

func (m Model) updateInsertMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Enter on method pane confirms selection and exits insert mode
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "enter" {
		if m.activePane == panes.PaneMethod {
			m.mode = ModeNormal
			m.setInsert(false)
			return m, nil
		}
	}

	var cmd tea.Cmd
	switch m.activePane {
	case panes.PaneMethod:
		p, c := m.method.Update(msg)
		m.method = p.(*panes.MethodPane)
		cmd = c
	case panes.PaneURL:
		p, c := m.url.Update(msg)
		m.url = p.(*panes.URLPane)
		cmd = c
	case panes.PaneHeaders:
		p, c := m.headers.Update(msg)
		m.headers = p.(*panes.HeadersPane)
		cmd = c
	case panes.PaneParams:
		p, c := m.params.Update(msg)
		m.params = p.(*panes.ParamsPane)
		cmd = c
	case panes.PaneBody:
		p, c := m.body.Update(msg)
		m.body = p.(*panes.BodyPane)
		cmd = c
	case panes.PaneAuth:
		p, c := m.auth.Update(msg)
		m.auth = p.(*panes.AuthPane)
		cmd = c
	case panes.PaneResults:
		p, c := m.results.Update(msg)
		m.results = p.(*panes.ResultsPane)
		cmd = c
	}

	// Rebuild curl if any pane is dirty
	m.rebuildCurlIfNeeded()
	return m, cmd
}

func (m Model) updateSaveMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.mode = ModeNormal
			m.saveInput.Blur()
			return m, nil
		case "enter":
			name := m.saveInput.Value()
			if name != "" {
				now := time.Now().Unix()
				sr := models.SavedRequest{
					Name:      name,
					CreatedAt: now,
					UpdatedAt: now,
					Request:   m.currentRequest(),
				}
				if err := storage.Save(sr); err != nil {
					m.lastStatus = "Save error: " + err.Error()
				} else {
					m.lastStatus = "Saved: " + name
					m.savedRequests = append(m.savedRequests, sr)
				}
			}
			m.mode = ModeNormal
			m.saveInput.Blur()
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.saveInput, cmd = m.saveInput.Update(msg)
	return m, cmd
}

func (m Model) updateLoadMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.mode = ModeNormal
			return m, nil
		case "j", "down":
			if m.loadCursor < len(m.savedRequests)-1 {
				m.loadCursor++
			}
			return m, nil
		case "k", "up":
			if m.loadCursor > 0 {
				m.loadCursor--
			}
			return m, nil
		case "enter":
			if len(m.savedRequests) > 0 {
				sr := m.savedRequests[m.loadCursor]
				m.loadRequest(sr.Request)
				m.lastStatus = "Loaded: " + sr.Name
			}
			m.mode = ModeNormal
			return m, nil
		}
	}
	return m, nil
}

func (m *Model) loadRequest(req models.Request) {
	m.method.SetMethod(req.Method)
	m.url.SetURL(req.URL)
	m.headers.SetHeaders(req.Headers)
	m.params.SetParams(req.Params)
	m.body.SetBody(req.Body)
	m.auth.SetAuth(req.Auth)
	m.rebuildCurl()
}

func (m *Model) rebuildCurlIfNeeded() {
	if m.method.IsDirty() || m.url.IsDirty() || m.headers.IsDirty() ||
		m.params.IsDirty() || m.body.IsDirty() || m.auth.IsDirty() {
		m.rebuildCurl()
	}
}

func (m *Model) rebuildCurl() {
	req := m.currentRequest()
	m.curlStr = curl.Build(req)
	m.method.ClearDirty()
	m.url.ClearDirty()
	m.headers.ClearDirty()
	m.params.ClearDirty()
	m.body.ClearDirty()
	m.auth.ClearDirty()
}

func (m *Model) recalcSizes() {
	if m.width == 0 || m.height == 0 {
		return
	}

	// Layout:
	// row1: method(12) + url(rest)          height=3
	// row2: headers(half) + params(half)    height=minRow2H+
	// row3: body(60%) + auth(40%)           height=minRow3H+
	// curlbar: height=1
	// results: remaining, capped at maxResultsH
	// statusbar: height=1

	const (
		row1H       = 4
		minRow2H    = 9
		minRow3H    = 11
		curlH       = 1
		statusH     = 1
		maxResultsH = 14
		methodW     = 14
	)

	// Cap results and distribute leftover to input rows
	baseH := row1H + minRow2H + minRow3H + curlH + statusH
	resultsH := m.height - baseH
	if resultsH > maxResultsH {
		resultsH = maxResultsH
	}
	if resultsH < 4 {
		resultsH = 4
	}

	// Give any extra vertical space to the input rows
	leftover := m.height - row1H - minRow2H - minRow3H - curlH - statusH - resultsH
	if leftover < 0 {
		leftover = 0
	}
	row2H := minRow2H + leftover/2
	row3H := minRow3H + leftover - leftover/2

	urlW := m.width - methodW

	headersW := m.width / 2
	paramsW := m.width - headersW

	bodyW := m.width * 6 / 10
	authW := m.width - bodyW

	m.method.SetSize(methodW, row1H)
	m.url.SetSize(urlW, row1H)
	m.headers.SetSize(headersW, row2H)
	m.params.SetSize(paramsW, row2H)
	m.body.SetSize(bodyW, row3H)
	m.auth.SetSize(authW, row3H)
	m.results.SetSize(m.width, resultsH)
}
