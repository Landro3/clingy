package config

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	inputs  []textinput.Model
	spinner spinner.Model
}

func InitialModel() Model {
	serverInput := serverInputModel()
	configInputs := []textinput.Model{serverInput}
	dotSpinner := spinner.New()
	dotSpinner.Spinner = spinner.Dot

	return Model{inputs: configInputs, spinner: dotSpinner}
}

func (m Model) Init() tea.Cmd {
	return nil
}
