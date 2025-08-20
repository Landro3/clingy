package shared

import "github.com/charmbracelet/lipgloss"

var (
	ButtonStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(Fg))

	FocusedButtonStyle = ButtonStyle.BorderForeground(lipgloss.Color(Purple))

	GrayText = lipgloss.NewStyle().Foreground(lipgloss.Color(Comment))
)
