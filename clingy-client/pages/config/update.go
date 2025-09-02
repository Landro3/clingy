package config

import (
	"clingy-client/http3"
	"clingy-client/shared"
	"clingy-client/util"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	validForm := m.inputs[0].Value() != "" && m.inputs[1].Value() != ""

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" {
			cmds = append(cmds, shared.NavigateCmd(shared.ChatPage))
		}

		if msg.String() == "enter" && validForm && m.focusIndex == len(m.inputs) {
			m.config.ServerAddr = m.inputs[0].Value()
			m.config.Username = m.inputs[1].Value()
			m.config.UniqueID = util.GenerateRandomString(32)
			http3.RegisterInServer(m.config)
		}

		if msg.String() == "tab" {
			m.focusIndex += 1
			if !validForm && m.focusIndex >= len(m.inputs) {
				m.focusIndex = 0
			} else if validForm && m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			}

			for i := range m.inputs {
				if i == m.focusIndex {
					m.inputs[i].Focus()
				} else {
					m.inputs[i].Blur()
				}
			}
		}

		for i := range m.inputs {
			m.inputs[i], cmd = m.inputs[i].Update(msg)
			cmds = append(cmds, cmd)
		}

		cmds = append(cmds, cmd)

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	// if len(cmds) == 0 {
	// 	cmds = append(cmds, m.spinner.Tick)
	// }

	return m, tea.Batch(cmds...)
}
