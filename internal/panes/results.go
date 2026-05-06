package panes

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/bugsy6/recurlit/internal/models"
	"github.com/bugsy6/recurlit/internal/styles"
)

type ResultsPane struct {
	vp       viewport.Model
	response *models.Response
	loading  bool
	active   bool
	insert   bool
	dirty    bool
	width    int
	height   int
}

func NewResultsPane() *ResultsPane {
	vp := viewport.New(80, 20)
	return &ResultsPane{vp: vp}
}

func (r *ResultsPane) SetResponse(resp *models.Response) {
	r.response = resp
	r.loading = false
	r.vp.SetContent(r.buildContent())
	r.vp.GotoTop()
	r.dirty = true
}

func (r *ResultsPane) SetLoading(loading bool) {
	r.loading = loading
	if loading {
		r.vp.SetContent(styles.Dim.Render("Executing request..."))
	}
	r.dirty = true
}

func (r *ResultsPane) buildContent() string {
	if r.response == nil {
		return styles.Dim.Render("Press r to execute request")
	}
	if r.response.Error != "" {
		return styles.StatusErr.Render("Error: " + r.response.Error)
	}

	var sb strings.Builder

	// Status line
	statusStyle := styles.StatusOK
	if r.response.StatusCode >= 400 {
		statusStyle = styles.StatusErr
	}
	sb.WriteString(statusStyle.Render(fmt.Sprintf("%s", r.response.Status)))
	sb.WriteString(styles.Dim.Render(fmt.Sprintf("  %dms", r.response.Duration.Milliseconds())))
	sb.WriteString("\n\n")

	// Response headers
	sb.WriteString(styles.Dim.Render("─── Headers ───\n"))
	for k, vals := range r.response.Headers {
		for _, v := range vals {
			sb.WriteString(fmt.Sprintf("%s: %s\n",
				styles.Dim.Render(k), v))
		}
	}
	sb.WriteString("\n")

	// Body
	sb.WriteString(styles.Dim.Render("─── Body ───\n"))
	body := r.response.Body
	if body == "" {
		sb.WriteString(styles.Dim.Render("(empty body)"))
	} else {
		// Try to pretty-print JSON
		var js interface{}
		if json.Unmarshal([]byte(body), &js) == nil {
			pretty, err := json.MarshalIndent(js, "", "  ")
			if err == nil {
				body = string(pretty)
			}
		}
		sb.WriteString(body)
	}

	return sb.String()
}

func (r *ResultsPane) Update(msg tea.Msg) (Pane, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "g":
			r.vp.GotoTop()
			return r, nil
		case "G":
			r.vp.GotoBottom()
			return r, nil
		}
	}
	var cmd tea.Cmd
	r.vp, cmd = r.vp.Update(msg)
	return r, cmd
}

func (r *ResultsPane) View() string {
	style := styles.InactiveBorder
	if r.insert {
		style = styles.InsertBorder
	} else if r.active {
		style = styles.ActiveBorder
	}
	style = style.Width(r.width - 2).Height(r.height - 2)
	title := styles.Title.Render("Response")
	return style.Render(title + "\n" + r.vp.View())
}

func (r *ResultsPane) SetSize(width, height int) {
	r.width = width
	r.height = height
	innerW := width - 4
	innerH := height - 5 // -1 for title line
	if innerW < 1 {
		innerW = 1
	}
	if innerH < 1 {
		innerH = 1
	}
	r.vp.Width = innerW
	r.vp.Height = innerH
	if r.response != nil {
		r.vp.SetContent(r.buildContent())
	}
}

func (r *ResultsPane) SetActive(active bool) { r.active = active }
func (r *ResultsPane) SetInsert(insert bool)  { r.insert = insert }
func (r *ResultsPane) IsDirty() bool          { return r.dirty }
func (r *ResultsPane) ClearDirty()            { r.dirty = false }

func (r *ResultsPane) Response() *models.Response { return r.response }
