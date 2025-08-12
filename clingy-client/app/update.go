package app

import (
	"clingy-client/shared"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowHeight = msg.Height
		m.windowWidth = msg.Width
		switch m.currentPage {
		case shared.ChatPage:
			m.chatModel, _ = m.chatModel.Update(msg)
		case shared.ContactPage:
			m.contactModel, _ = m.contactModel.Update(msg)
		}
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}
	case shared.NavigateMsg:
		m.currentPage = msg.Page
		windowMsg := tea.WindowSizeMsg{Width: m.windowWidth, Height: m.windowHeight}
		switch m.currentPage {
		case shared.ChatPage:
			m.chatModel, _ = m.chatModel.Update(windowMsg)
		case shared.ContactPage:
			m.contactModel, _ = m.contactModel.Update(windowMsg)
		}
	}

	switch m.currentPage {
	case shared.ChatPage:
		m.chatModel, cmd = m.chatModel.Update(msg)
	case shared.ContactPage:
		m.contactModel, cmd = m.contactModel.Update(msg)
	}

	return m, cmd
}
