package styles

import "github.com/charmbracelet/lipgloss"

var (
	ColorActive   = lipgloss.Color("#10B981")
	ColorInsert   = lipgloss.Color("#7C3AED")
	ColorInactive = lipgloss.Color("#374151")
	ColorCurlBg   = lipgloss.Color("#1E40AF")
	ColorText     = lipgloss.Color("#F9FAFB")
	ColorDim      = lipgloss.Color("#6B7280")
	ColorSuccess  = lipgloss.Color("#10B981")
	ColorError    = lipgloss.Color("#EF4444")
	ColorWarning  = lipgloss.Color("#F59E0B")

	ActiveBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorActive)

	InsertBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorInsert)

	InactiveBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorInactive)

	CurlBar = lipgloss.NewStyle().
		Background(ColorCurlBg).
		Foreground(ColorText).
		Padding(0, 1)

	StatusBar = lipgloss.NewStyle().
			Background(lipgloss.Color("#111827")).
			Foreground(ColorDim).
			Padding(0, 1)

	StatusOK = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Bold(true)

	StatusErr = lipgloss.NewStyle().
			Foreground(ColorError).
			Bold(true)

	Title = lipgloss.NewStyle().
		Foreground(ColorText).
		Bold(true)

	Dim = lipgloss.NewStyle().
		Foreground(ColorDim)

	Selected = lipgloss.NewStyle().
			Foreground(ColorActive).
			Bold(true)
)
