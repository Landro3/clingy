package config

import (
	"clingy-client/services"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	inputs     []textinput.Model
	spinner    spinner.Model
	config     *services.Config
	focusIndex int
}

func InitialModel(config *services.Config) Model {
	serverInput := serverInputModel()
	serverInput.Focus()
	if config.ServerAddr != "" {
		serverInput.SetValue(config.ServerAddr)
	}

	usernameInput := usernameInputModel()
	usernameInput.Blur()
	if config.Username != "" {
		usernameInput.SetValue(config.Username)
	}

	configInputs := []textinput.Model{serverInput, usernameInput}
	dotSpinner := spinner.New()
	dotSpinner.Spinner = spinner.Dot

	return Model{inputs: configInputs, spinner: dotSpinner, config: config}
}

func (m Model) Init() tea.Cmd {
	return nil
}
