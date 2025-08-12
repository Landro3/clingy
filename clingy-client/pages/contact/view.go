package contact

// a/e/d for add/edit/delete modal
// esc - back to chat?
// List component
// help along bottom

func (m Model) View() string {
	return docStyle.Render(m.list.View())
}
