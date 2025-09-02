package config

func (m Model) View() string {
	return "Config Page" + "\n\n" +
		"Clingy Server" +
		m.inputs[0].View() + "\n\n" +
		"Username" +
		m.inputs[1].View() + "\n\n" +
		m.registerIDView()
}
