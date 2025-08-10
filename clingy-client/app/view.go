package app

func (m Model) View() string {
	switch m.currentPage {
	case ChatPage:
		return m.chatModel.View()
	default:
		return "Unknown Page"
	}
}
