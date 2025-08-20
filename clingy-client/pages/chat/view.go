package chat

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	inputAndButtons := lipgloss.PlaceHorizontal(
		m.viewport.Width,
		lipgloss.Center,
		lipgloss.JoinHorizontal(lipgloss.Center, m.chatInputView(), m.contactButtonView(), m.configButtonView()),
	)
	return fmt.Sprintf("%s\n%s\n%s\n%s", m.headerView(), m.viewport.View(), inputAndButtons, m.footerView())
}
