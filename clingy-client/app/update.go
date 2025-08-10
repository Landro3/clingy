package app

import tea "github.com/charmbracelet/bubbletea"

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}
	}

	switch m.currentPage {
	case ChatPage:
		m.chatModel, cmd = m.chatModel.Update(msg)
	}

	return m, cmd
}
