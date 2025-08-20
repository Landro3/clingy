package config

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct{}

func InitialModel() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return nil
}
