package contact

import (
	"clingy-client/shared"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

var (
	docStyle = lipgloss.NewStyle().Margin(1, 2)

	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(shared.Purple))

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(shared.Fg)).
			Bold(true)

	inputStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(shared.BgHighlight)).
			Padding(0, 1).
			MarginTop(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(shared.Fg))
)

type user struct {
	name, id string
}

func (u user) Title() string       { return u.name }
func (u user) FilterValue() string { return u.name }
func (u user) Description() string { return u.id }

func usernameInputModel() textinput.Model {
	usernameInput := textinput.New()
	usernameInput.Cursor.Style = cursorStyle
	usernameInput.PromptStyle = cursorStyle
	usernameInput.TextStyle = cursorStyle
	usernameInput.Placeholder = "Username"
	usernameInput.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(shared.Comment))
	usernameInput.CharLimit = 32
	usernameInput.Width = 30
	usernameInput.Focus()

	return usernameInput
}

func uuidInputModel() textinput.Model {
	uuidInput := textinput.New()
	uuidInput.Cursor.Style = cursorStyle
	uuidInput.PromptStyle = cursorStyle
	uuidInput.TextStyle = cursorStyle
	uuidInput.Placeholder = "UUID"
	uuidInput.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(shared.Comment))
	uuidInput.CharLimit = 32
	uuidInput.Width = 30 // Add this to fix placeholder display

	return uuidInput
}
