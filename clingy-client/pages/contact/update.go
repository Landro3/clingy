package contact

import (
	"clingy-client/services"
	"clingy-client/shared"
	"clingy-client/util"

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

			if msg.String() == "enter" {
				selectedItem := m.list.SelectedItem()
				if selectedItem != nil {
					contact, ok := selectedItem.(services.ContactInfo)
					if ok {
						cmds = append(cmds, tea.Sequence(
							shared.NavigateCmd(shared.ChatPage),
							SelectChatCmd(contact),
						))
					} else {
						util.Log("Failed to cast selected item to ContactInfo")
					}
				}
			}
		} else {
			if msg.String() == "esc" {
				m.mode = None

				for i := range m.inputs {
					m.inputs[i].SetValue("")
				}

				return m, tea.Batch(cmds...)
			}

			if msg.String() == "enter" && m.mode == Add {
				contactInfo := services.NewContactInfo(m.inputs[0].Value(), m.inputs[1].Value())
				m.contact.AddContact(contactInfo)
				m.list.SetItems(m.contact.ToListItems())
				for i := range m.inputs {
					m.inputs[i].SetValue("")
				}
				m.mode = None
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
