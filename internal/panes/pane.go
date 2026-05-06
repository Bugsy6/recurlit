package panes

import tea "github.com/charmbracelet/bubbletea"

type PaneID int

const (
	PaneMethod PaneID = iota
	PaneURL
	PaneHeaders
	PaneParams
	PaneBody
	PaneAuth
	PaneResults
	PaneCount
)

func (p PaneID) String() string {
	switch p {
	case PaneMethod:
		return "Method"
	case PaneURL:
		return "URL"
	case PaneHeaders:
		return "Headers"
	case PaneParams:
		return "Params"
	case PaneBody:
		return "Body"
	case PaneAuth:
		return "Auth"
	case PaneResults:
		return "Results"
	}
	return "Unknown"
}

// Pane is implemented by each pane struct
type Pane interface {
	Update(msg tea.Msg) (Pane, tea.Cmd)
	View() string
	SetSize(width, height int)
	SetActive(active bool)
	SetInsert(insert bool)
	IsDirty() bool
	ClearDirty()
}
