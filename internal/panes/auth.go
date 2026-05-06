package panes

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/bugsy6/recurlit/internal/models"
	"github.com/bugsy6/recurlit/internal/styles"
)

type AuthPane struct {
	authType int // 0=none,1=bearer,2=basic
	token    textinput.Model
	username textinput.Model
	password textinput.Model
	focusIdx int
	active   bool
	insert   bool
	dirty    bool
	width    int
	height   int
}

var authTypes = []models.AuthType{models.AuthNone, models.AuthBearer, models.AuthBasic}
var authTypeLabels = []string{"None", "Bearer", "Basic"}

func NewAuthPane() *AuthPane {
	token := textinput.New()
	token.Placeholder = "Bearer token"
	token.CharLimit = 512

	user := textinput.New()
	user.Placeholder = "Username"
	user.CharLimit = 256

	pass := textinput.New()
	pass.Placeholder = "Password"
	pass.EchoMode = textinput.EchoPassword
	pass.EchoCharacter = '•'
	pass.CharLimit = 256

	return &AuthPane{
		token:    token,
		username: user,
		password: pass,
		dirty:    true,
	}
}

func (a *AuthPane) Auth() models.Auth {
	return models.Auth{
		Type:     authTypes[a.authType],
		Token:    a.token.Value(),
		Username: a.username.Value(),
		Password: a.password.Value(),
	}
}

func (a *AuthPane) SetAuth(auth models.Auth) {
	for i, t := range authTypes {
		if t == auth.Type {
			a.authType = i
			break
		}
	}
	a.token.SetValue(auth.Token)
	a.username.SetValue(auth.Username)
	a.password.SetValue(auth.Password)
	a.dirty = true
}

func (a *AuthPane) inputs() []*textinput.Model {
	switch authTypes[a.authType] {
	case models.AuthBearer:
		return []*textinput.Model{&a.token}
	case models.AuthBasic:
		return []*textinput.Model{&a.username, &a.password}
	}
	return nil
}

func (a *AuthPane) focusCurrent() {
	a.token.Blur()
	a.username.Blur()
	a.password.Blur()
	if !a.insert {
		return
	}
	inputs := a.inputs()
	if a.focusIdx < len(inputs) {
		inputs[a.focusIdx].Focus()
	}
}

func (a *AuthPane) Update(msg tea.Msg) (Pane, tea.Cmd) {
	if !a.insert {
		return a, nil
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			inputs := a.inputs()
			if len(inputs) == 0 {
				// cycle auth type
				a.authType = (a.authType + 1) % len(authTypes)
				a.dirty = true
				return a, nil
			}
			a.focusIdx++
			if a.focusIdx >= len(inputs) {
				a.focusIdx = 0
				// also cycle type
				a.authType = (a.authType + 1) % len(authTypes)
				a.dirty = true
			}
			a.focusCurrent()
			return a, nil
		case "ctrl+t":
			a.authType = (a.authType + 1) % len(authTypes)
			a.focusIdx = 0
			a.focusCurrent()
			a.dirty = true
			return a, nil
		}
	}
	// delegate to active input
	inputs := a.inputs()
	if a.focusIdx < len(inputs) {
		var cmd tea.Cmd
		prev := inputs[a.focusIdx].Value()
		*inputs[a.focusIdx], cmd = inputs[a.focusIdx].Update(msg)
		if inputs[a.focusIdx].Value() != prev {
			a.dirty = true
		}
		return a, cmd
	}
	return a, nil
}

func (a *AuthPane) View() string {
	style := styles.InactiveBorder
	if a.insert {
		style = styles.InsertBorder
	} else if a.active {
		style = styles.ActiveBorder
	}
	style = style.Width(a.width - 2).Height(a.height - 2)
	return style.Render(a.renderInner())
}

func (a *AuthPane) renderInner() string {
	var sb strings.Builder

	sb.WriteString(styles.Title.Render("Auth") + "\n\n")

	// Type selector
	typeLine := "Type: "
	for i, label := range authTypeLabels {
		if i == a.authType {
			typeLine += styles.Selected.Render("[" + label + "]")
		} else {
			typeLine += styles.Dim.Render(" " + label + " ")
		}
		if i < len(authTypeLabels)-1 {
			typeLine += " "
		}
	}
	sb.WriteString(typeLine + "\n\n")

	switch authTypes[a.authType] {
	case models.AuthBearer:
		sb.WriteString(fmt.Sprintf("Token:\n%s\n", a.token.View()))
	case models.AuthBasic:
		sb.WriteString(fmt.Sprintf("Username:\n%s\n\nPassword:\n%s\n",
			a.username.View(), a.password.View()))
	default:
		sb.WriteString(styles.Dim.Render("No authentication"))
	}

	if a.insert {
		sb.WriteString("\n" + styles.Dim.Render("ctrl+t: change type  tab: next"))
	}

	return sb.String()
}

func (a *AuthPane) SetSize(width, height int) {
	a.width = width
	a.height = height
	inner := width - 6
	if inner < 10 {
		inner = 10
	}
	a.token.Width = inner
	a.username.Width = inner
	a.password.Width = inner
}

func (a *AuthPane) SetActive(active bool) {
	a.active = active
	if !active {
		a.focusCurrent()
	}
}

func (a *AuthPane) SetInsert(insert bool) {
	a.insert = insert
	a.focusCurrent()
}

func (a *AuthPane) IsDirty() bool { return a.dirty }
func (a *AuthPane) ClearDirty()   { a.dirty = false }
