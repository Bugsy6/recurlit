package app

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/bugsy6/recurlit/internal/models"
	"github.com/bugsy6/recurlit/internal/panes"
	"github.com/bugsy6/recurlit/internal/storage"
)

type Mode int

const (
	ModeNormal Mode = iota
	ModeInsert
	ModeSave
	ModeLoad
	ModeHelp
)

type Model struct {
	// Panes (ordered for tab navigation)
	method  *panes.MethodPane
	url     *panes.URLPane
	headers *panes.HeadersPane
	params  *panes.ParamsPane
	body    *panes.BodyPane
	auth    *panes.AuthPane
	results *panes.ResultsPane

	activePane panes.PaneID
	mode       Mode

	// Curl preview
	curlStr string

	// Status bar
	lastStatus string
	executing  bool
	spinner    spinner.Model

	// Save modal
	saveInput textinput.Model

	// Load picker
	savedRequests []models.SavedRequest
	loadCursor    int

	// Help overlay
	showHelp bool

	// Terminal size
	width  int
	height int

	keymap KeyMap
}

func New() Model {
	sp := spinner.New()
	sp.Spinner = spinner.Dot

	si := textinput.New()
	si.Placeholder = "Request name..."
	si.CharLimit = 128

	m := Model{
		method:     panes.NewMethodPane(),
		url:        panes.NewURLPane(),
		headers:    panes.NewHeadersPane(),
		params:     panes.NewParamsPane(),
		body:       panes.NewBodyPane(),
		auth:       panes.NewAuthPane(),
		results:    panes.NewResultsPane(),
		activePane: panes.PaneURL,
		mode:       ModeNormal,
		spinner:    sp,
		saveInput:  si,
		keymap:     DefaultKeyMap,
	}

	// Load saved requests in background at startup
	saved, _ := storage.LoadAll()
	m.savedRequests = saved

	m.setActive(panes.PaneURL)
	return m
}

func (m *Model) setActive(id panes.PaneID) {
	// Deactivate all
	m.method.SetActive(false)
	m.url.SetActive(false)
	m.headers.SetActive(false)
	m.params.SetActive(false)
	m.body.SetActive(false)
	m.auth.SetActive(false)
	m.results.SetActive(false)

	m.activePane = id

	switch id {
	case panes.PaneMethod:
		m.method.SetActive(true)
	case panes.PaneURL:
		m.url.SetActive(true)
	case panes.PaneHeaders:
		m.headers.SetActive(true)
	case panes.PaneParams:
		m.params.SetActive(true)
	case panes.PaneBody:
		m.body.SetActive(true)
	case panes.PaneAuth:
		m.auth.SetActive(true)
	case panes.PaneResults:
		m.results.SetActive(true)
	}
}

func (m *Model) setInsert(insert bool) {
	m.method.SetInsert(false)
	m.url.SetInsert(false)
	m.headers.SetInsert(false)
	m.params.SetInsert(false)
	m.body.SetInsert(false)
	m.auth.SetInsert(false)
	m.results.SetInsert(false)

	if !insert {
		return
	}
	switch m.activePane {
	case panes.PaneMethod:
		m.method.SetInsert(true)
	case panes.PaneURL:
		m.url.SetInsert(true)
	case panes.PaneHeaders:
		m.headers.SetInsert(true)
	case panes.PaneParams:
		m.params.SetInsert(true)
	case panes.PaneBody:
		m.body.SetInsert(true)
	case panes.PaneAuth:
		m.auth.SetInsert(true)
	case panes.PaneResults:
		m.results.SetInsert(true)
	}
}

func (m *Model) currentRequest() models.Request {
	headers := m.headers.Headers()
	params := m.params.Params()
	return models.Request{
		Method:  m.method.Method(),
		URL:     m.url.URL(),
		Headers: headers,
		Params:  params,
		Body:    m.body.Body(),
		Auth:    m.auth.Auth(),
	}
}

func (m Model) Init() tea.Cmd {
	return m.spinner.Tick
}
