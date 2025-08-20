package config

import (
	"clingy-client/shared"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			cmds = append(cmds, shared.NavigateCmd(shared.ChatPage))
		}
		m.inputs[0], cmd = m.inputs[0].Update(msg)
		cmds = append(cmds, cmd)

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	if len(cmds) == 0 {
		cmds = append(cmds, m.spinner.Tick)
	}

	return m, tea.Batch(cmds...)
}
