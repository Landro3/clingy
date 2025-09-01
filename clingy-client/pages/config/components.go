package config

import (
	"clingy-client/shared"
	"fmt"

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

func usernameInputModel() textinput.Model {
	usernameInput := textinput.New()
	usernameInput.Placeholder = "Username"
	usernameInput.CharLimit = 32
	usernameInput.Width = 30
	usernameInput.Focus()

	return usernameInput
}

func (m Model) registerIDView() string {
	style := shared.ButtonStyle
	if m.focusIndex == len(m.inputs) {
		style = shared.FocusedButtonStyle
	}

	uniqueID := "N/A"
	if m.config.UniqueID != "" {
		uniqueID = m.config.UniqueID
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		style.Render("Register"),
		" ",
		shared.GrayText.Render(fmt.Sprintf("Current ID: %s", uniqueID)),
		// m.spinner.View(),
	)
}
