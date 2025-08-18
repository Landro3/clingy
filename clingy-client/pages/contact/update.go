package contact

import (
	"clingy-client/shared"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.mode == None && msg.String() == "esc" {
			cmds = append(cmds, shared.NavigateCmd(shared.ChatPage))
		}

		if m.mode == None {
			if msg.String() == "a" {
				m.mode = Add
				m.focusIndex = 0
			} else if msg.String() == "e" {
				m.mode = Edit
				m.focusIndex = 0
			} else if msg.String() == "r" {
				m.mode = Remove
			}
		} else {
			if msg.String() == "esc" {
				m.mode = None

				for i := range m.inputs {
					m.inputs[i].SetValue("")
				}

				return m, tea.Batch(cmds...)
			}

			for i := range m.inputs {
				m.inputs[i], cmd = m.inputs[i].Update(msg)
				cmds = append(cmds, cmd)
			}

			if msg.String() == "tab" {
				m.focusIndex += 1
				if m.focusIndex >= len(m.inputs) {
					m.focusIndex = 0
				}
			}
		}

		for i := range m.inputs {
			if i == m.focusIndex {
				m.inputs[i].Focus()
			} else {
				m.inputs[i].Blur()
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}
