package config

import (
	"clingy-client/shared"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

func serverInputModel() textinput.Model {
	serverInput := textinput.New()
	serverInput.Placeholder = "clingy.server.com"
	serverInput.CharLimit = 32
	serverInput.Width = 30
	serverInput.Focus()

	return serverInput
}

func (m Model) registerIdView() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		shared.ButtonStyle.Render("Register"),
		" ",
		shared.GrayText.Render("Current ID: N/A"),
		m.spinner.View(),
	)
}
