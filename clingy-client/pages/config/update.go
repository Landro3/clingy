package config

import (
	"clingy-client/shared"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			cmds = append(cmds, shared.NavigateCmd(shared.ChatPage))
		}
	}
	return m, tea.Batch(cmds...)
}
