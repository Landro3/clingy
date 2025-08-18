package contact

import "github.com/charmbracelet/lipgloss"

// a/e/r for add/edit/remove modal
// esc - back to chat?
// List component
// help along bottom

func (m Model) View() string {
	listView := docStyle.Render(m.list.View())

	var inputTitle string
	switch m.mode {
	case Add:
		inputTitle = "Create User:"
	case Edit:
		inputTitle = "Edit User:"
	case Remove:
		inputTitle = "Remove User? (y/n)"
	default:
		inputTitle = ""
	}

	// usernameLabel := labelStyle.Render("Username:")
	usernameView := inputStyle.Render(m.inputs[0].View())

	// uuidLabel := labelStyle.Render("UUID:")
	uuidView := inputStyle.Render(m.inputs[1].View())

	inputViews := lipgloss.JoinHorizontal(lipgloss.Left, usernameView, " ", uuidView)

	var helpText string
	switch m.mode {
	case None:
		helpText = "a: add • e: edit • r: remove • ↑↓: navigate • esc: back"
	default:
		helpText = "tab: switch • enter: save • esc: cancel"
	}

	if m.mode == None {
		return listView + "\n" + helpStyle.Render(helpText)
	}

	if m.mode == Remove {
		return listView + "\n" +
			labelStyle.Render(inputTitle) + "\n" +
			helpStyle.Render(helpText)
	}

	return listView + "\n" +
		labelStyle.Render(inputTitle) +
		inputViews + "\n" +
		helpStyle.Render(helpText)
}
