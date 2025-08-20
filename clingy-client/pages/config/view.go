package config

func (m Model) View() string {
	return "Config Page" + "\n\n" +
		"Clingy Server:" + "\n" +
		m.inputs[0].View() + "\n\n" +
		m.registerIdView()
}
