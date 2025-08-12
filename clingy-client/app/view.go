package app

import "clingy-client/shared"

func (m Model) View() string {
	switch m.currentPage {
	case shared.ChatPage:
		return m.chatModel.View()
	case shared.ContactPage:
		return m.contactModel.View()
	default:
		return "Unknown Page"
	}
}
